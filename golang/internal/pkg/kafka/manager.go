package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kyh0703/template/internal/core/domain/event"
	"github.com/kyh0703/template/internal/pkg/eventhandler"

	cloud "github.com/cloudevents/sdk-go/v2"
)

type Manager struct {
	lock        sync.RWMutex
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	subjectMap  map[string]eventhandler.Subject
	producerMap map[string]*Producer
	consumerMap map[string]*Consumer
}

func ProvideManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	manager := &Manager{
		ctx:         ctx,
		cancel:      cancel,
		subjectMap:  make(map[string]eventhandler.Subject),
		producerMap: make(map[string]*Producer),
		consumerMap: make(map[string]*Consumer),
	}
	manager.InitLoop()
	go manager.poll()
	return manager
}

func (m *Manager) InitLoop() {
reconnect_producer:
	if err := m.initProducer(); err != nil {
		time.Sleep(time.Second)
		goto reconnect_producer
	}
reconnect_consumer:
	if err := m.initConsumer(); err != nil {
		time.Sleep(time.Second)
		goto reconnect_consumer
	}
}

func (m *Manager) initConsumer() error {
	// call := NewConsumer(m.ctx, ievent.CloudTopicCall, m.Consume)
	// if err := call.Connect(configs.Env.KafkaBroker); err != nil {
	// 	call.Close()
	// 	return err
	// }
	// media := NewConsumer(m.ctx, ievent.CloudTopicMedia, m.Consume)
	// if err := media.Connect(configs.Env.KafkaBroker); err != nil {
	// 	media.Close()
	// 	return err
	// }
	// sip := NewConsumer(m.ctx, ievent.CloudTypeSip, m.Consume)
	// if err := sip.Connect(configs.Env.KafkaBroker); err != nil {
	// 	sip.Close()
	// 	return err
	// }

	// m.registerConsumer(call)
	// m.registerConsumer(media)
	// m.registerConsumer(sip)
	return nil
}

func (m *Manager) initProducer() error {
	return nil
}

func (m *Manager) registerProducer(p *Producer) {
	m.producerMap[p.Topic()] = p
}

func (m *Manager) registerConsumer(c *Consumer) {
	m.consumerMap[c.Topic()] = c
}

// Close Event Manager.
func (m *Manager) Close() {
	m.cancel()
	m.wg.Wait()
}

func (m *Manager) Register(id string, l eventhandler.Listener) {
	m.lock.Lock()
	defer m.lock.Unlock()

	subject, ok := m.subjectMap[id]
	if !ok {
		subject = eventhandler.NewDispatcher(id)
	}
	subject.Register(l)
	m.subjectMap[id] = subject
}

// Deregister is delete subject.
func (m *Manager) Deregister(id string, l eventhandler.Listener) {
	m.lock.Lock()
	defer m.lock.Unlock()

	defer l.Close()

	// get dispatcher
	subject, ok := m.subjectMap[id]
	if !ok {
		return
	}

	// detach listener
	subject.Deregister(l)

	// if listener count is zero, delete subject
	if subject.ListenerCount() == 0 {
		delete(m.subjectMap, id)
	}
}

func (m *Manager) Consume(_ context.Context, ce cloud.Event) {
	var callID string
	switch ce.Type() {
	default:
		return
	}

	m.lock.RLock()
	defer m.lock.RUnlock()

	subject, ok := m.subjectMap[callID]
	if !ok {
		return
	}
	subject.Notify(ce)
}

func (m *Manager) Produce(events ...event.Event) {
	for _, e := range events {
		ce := cloud.NewEvent()
		ce.SetID(uuid.New().String())
		ce.SetType(e.Type())
		ce.SetSubject(e.Subject())
		ce.SetExtension("eventsubject", e.EventSubject())
		ce.SetExtension("tenantid", e.TenantID())
		ce.SetTime(time.Now())
		if err := ce.SetData(cloud.ApplicationJSON, e.MakeData()); err != nil {
			continue
		}
		producer, ok := m.producerMap[e.Topic()]
		if !ok {
			continue
		}
		producer.Send(ce)
	}
}

func (m *Manager) print() {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if len(m.subjectMap) == 0 {
		return
	}

	type garbageLog struct {
		CallId        string    `json:"call_id"`
		ListenerCount int       `json:"listener_cnt"`
		CreateAt      time.Time `json:"create_at"`
	}
	garbage := make([]garbageLog, 0)

	for _, subject := range m.subjectMap {
		garbage = append(garbage, garbageLog{
			subject.ID(),
			subject.ListenerCount(),
			subject.CreateAt(),
		})
	}

	b, err := json.Marshal(garbage)
	if err != nil {
		return
	}
	log.Print(string(b))
}

func (m *Manager) poll() {
	m.wg.Add(1)
	defer m.wg.Done()
	for {
		select {
		case <-time.After(30 * time.Second):
			m.print()
		case <-m.ctx.Done():
			return
		}
	}
}
