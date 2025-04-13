package eventhandler

import (
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gofrs/uuid"
)

type Manager interface {
	Deregister(id string)
}

type Subject interface {
	ID() string
	Register(Listener)
	Deregister(Listener)
	Notify(cloudevents.Event)
	ListenerCount() int
	CreateAt() time.Time
}

type Listener interface {
	ID() uuid.UUID
	Close()
	OnEvent(cloudevents.Event)
}

type Client interface {
	Topic() string
	Connect(brokers []string) error
	Close()
}
