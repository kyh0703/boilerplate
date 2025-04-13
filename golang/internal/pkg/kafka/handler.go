package kafka

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	cloud "github.com/cloudevents/sdk-go/v2"
)

const (
	timeout  = time.Second * 2
	MaxEvent = 30
)

var ErrTimeout = errors.New("event handler timeout")

type DoneFn = func()

type EventFunc = func(ctx context.Context, e cloud.Event, done DoneFn) error

type Handler struct {
	id      uuid.UUID
	timer   *time.Timer
	bus     chan cloud.Event
	expects []string
	doneCh  chan struct{}
	isDone  bool

	onCallTerminate  EventFunc
	onCallDisconnect EventFunc
	onSipAlert       EventFunc
	onSipConnect     EventFunc
	onSipRelease     EventFunc
	onSipTransfer    EventFunc
	onSipJoin        EventFunc
	onSipSwitch      EventFunc
	onSipHeld        EventFunc
	onSipRetrieve    EventFunc
}

func NewHandler() *Handler {
	uuid, _ := uuid.NewV4()
	handler := &Handler{
		id:     uuid,
		bus:    make(chan cloud.Event, MaxEvent),
		doneCh: make(chan struct{}, 1), // Buffered to avoid locking up the event feed
	}
	return handler
}

func (h *Handler) ID() uuid.UUID {
	return h.id
}

func (h *Handler) SetOnCallTerminate(f EventFunc)  { h.onCallTerminate = f }
func (h *Handler) SetOnCallDisconnect(f EventFunc) { h.onCallDisconnect = f }
func (h *Handler) SetOnSipAlert(f EventFunc)       { h.onSipAlert = f }
func (h *Handler) SetOnSipConnect(f EventFunc)     { h.onSipConnect = f }
func (h *Handler) SetOnSipRelease(f EventFunc)     { h.onSipRelease = f }
func (h *Handler) SetOnSipTransfer(f EventFunc)    { h.onSipTransfer = f }
func (h *Handler) SetOnSipJoin(f EventFunc)        { h.onSipJoin = f }
func (h *Handler) SetOnSipSwitch(f EventFunc)      { h.onSipSwitch = f }
func (h *Handler) SetOnSipHeld(f EventFunc)        { h.onSipHeld = f }
func (h *Handler) SetOnSipRetrieve(f EventFunc)    { h.onSipRetrieve = f }

func (h *Handler) SetTimer(dur time.Duration) {
	if h.timer == nil {
		h.timer = time.NewTimer(dur)
	} else {
		h.timer.Reset(dur)
	}
}

func (h *Handler) OnEvent(event cloud.Event) {
	h.bus <- event
}

func (h *Handler) IfInExpects(id string) bool {
	for _, expect := range h.expects {
		if id == expect {
			return true
		}
	}
	return false
}

func (h *Handler) Expects() []string {
	return h.expects
}

func (h *Handler) Subscribe(expects ...string) {
	h.expects = expects
}

func (h *Handler) Close() {
	if h.bus != nil {
		close(h.bus)
	}
	if h.timer != nil {
		h.timer.Stop()
	}
	if h.doneCh != nil {
		close(h.doneCh)
		h.doneCh = nil
	}
}

func (h *Handler) done() {
	h.isDone = true
	if h.doneCh != nil {
		h.doneCh <- struct{}{}
	}
}

func (h *Handler) ProcEvent(ctx context.Context, event cloud.Event) error {
	if h.isDone {
		return nil
	}
	switch event.Type() {
	default:
		return nil
	}
}

func (h *Handler) Poll(ctx context.Context) error {
	if h.timer == nil {
		h.timer = time.NewTimer(timeout)
	}
	for {
		select {
		case <-h.doneCh:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-h.timer.C:
			return ErrTimeout
		case event := <-h.bus:
			if err := h.ProcEvent(ctx, event); err != nil {
				return err
			}
		}
	}
}
