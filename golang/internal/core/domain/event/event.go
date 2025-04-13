package event

type Event interface {
	Topic() string
	Type() string
	Subject() string
	EventSubject() string
	TenantID() string
	MakeData() []byte
}
