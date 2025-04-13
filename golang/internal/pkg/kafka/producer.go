package kafka

import (
	"context"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	"github.com/pkg/errors"

	cloud "github.com/cloudevents/sdk-go/v2"
)

const MaxProduceEvent = 300

type Producer struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	config *sarama.Config
	sender *kafka_sarama.Sender
	topic  string
	client cloud.Client
	postCh chan cloud.Event
}

func NewProducer(ctx context.Context, topic string) *Producer {
	ctx, cancel := context.WithCancel(ctx)
	saramaConfig := sarama.NewConfig()

	return &Producer{
		ctx:    ctx,
		cancel: cancel,
		config: saramaConfig,
		topic:  topic,
		postCh: make(chan cloud.Event, MaxProduceEvent),
	}
}

func (p *Producer) Topic() string {
	return p.topic
}

func (p *Producer) Connect(brokers []string) error {
	sender, err := kafka_sarama.NewSender(brokers, p.config, p.topic)
	if err != nil {
		return errors.Wrap(err, "failed to new sender")
	}
	client, err := cloud.NewClient(sender, cloud.WithTimeNow(), cloud.WithUUIDs())
	if err != nil {
		return errors.Wrap(err, "failed to new client")
	}
	p.sender = sender
	p.client = client
	go p.run()
	return nil
}

func (p *Producer) Close() {
	p.cancel()
	if p.sender != nil {
		p.sender.Close(p.ctx)
		p.sender = nil
	}
	p.wg.Wait()
}

func (p *Producer) Send(event cloud.Event) {
	p.postCh <- event
}

func (p *Producer) ExtractCallID(event cloud.Event) (string, bool) {
	subject := event.Subject()
	return strings.CutPrefix(subject, "call/")
}

func (p *Producer) setMessageKey(ctx context.Context, callID string) context.Context {
	return kafka_sarama.WithMessageKey(ctx, sarama.StringEncoder(callID))
}

func (p *Producer) run() {
	for {
		select {
		case <-p.ctx.Done():
			p.wg.Done()
			p.Close()
			return
		case event := <-p.postCh:
			callID, ok := p.ExtractCallID(event)
			if !ok {
				continue
			}
			ctx := p.setMessageKey(p.ctx, callID)
			result := p.client.Send(ctx, event)
			if cloud.IsUndelivered(result) {
			}
		}
	}
}
