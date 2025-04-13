package eventhandler

import (
	"slices"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Dispatcher struct {
	id        string
	listeners []Listener
	createAt  time.Time
}

func NewDispatcher(id string) Subject {
	return &Dispatcher{
		id:       id,
		createAt: time.Now(),
	}
}

// ID implements Subject.
func (d *Dispatcher) ID() string {
	return d.id
}

func (d *Dispatcher) Register(l Listener) {
	d.listeners = append(d.listeners, l)
}

func (d *Dispatcher) Deregister(l Listener) {
	d.RemoveListener(l)
}

func (d *Dispatcher) Notify(ce cloudevents.Event) {
	for i := range d.listeners {
		d.listeners[i].OnEvent(ce)
	}
}

func (d *Dispatcher) ListenerCount() int {
	return len(d.listeners)
}

func (d *Dispatcher) RemoveListener(l Listener) {
	for i, listener := range d.listeners {
		if l.ID() == listener.ID() {
			d.listeners = slices.Delete(d.listeners, i, i+1)
			return
		}
	}
}

// CreateAt implements Subject.
func (d *Dispatcher) CreateAt() time.Time {
	return d.createAt
}
