package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	"github.com/pkg/errors"

	cloud "github.com/cloudevents/sdk-go/v2"
)

type ReceiveFunc func(context.Context, cloud.Event)

type Consumer struct {
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	config   *sarama.Config
	receiver *kafka_sarama.Consumer
	topic    string
	client   cloud.Client
	recvFn   ReceiveFunc
}

func NewConsumer(ctx context.Context, topic string, recvFn ReceiveFunc) *Consumer {
	ctx, cancel := context.WithCancel(ctx)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0
	saramaConfig.Consumer.Offsets.Retention = time.Minute
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	return &Consumer{
		ctx:    ctx,
		cancel: cancel,
		config: saramaConfig,
		topic:  topic,
		recvFn: recvFn,
	}
}

func (c *Consumer) Close() {
	c.cancel()
	if c.receiver != nil {
		c.receiver.Close(c.ctx)
		c.receiver = nil
	}
	c.wg.Wait()
}

func (c *Consumer) Topic() string {
	return c.topic
}

func (c *Consumer) Connect(brokers []string) error {
	receiver, err := kafka_sarama.NewConsumer(
		brokers,
		c.config,
		"my_app",
		c.topic,
	)
	if err != nil {
		return err
	}

	client, err := cloud.NewClient(receiver, cloud.WithTimeNow(), cloud.WithUUIDs())
	if err != nil {
		return errors.Wrap(err, "failed to create client")
	}

	c.receiver = receiver
	c.client = client

	c.wg.Add(1)
	go c.run()

	return nil
}

func (c *Consumer) run() {
	for {
		select {
		case <-c.ctx.Done():
			c.wg.Done()
			c.Close()
			return
		default:
			if err := c.client.StartReceiver(c.ctx, c.recvFn); err != nil {
				// reconnect receiver, after 5 second
				time.Sleep(time.Second * 5)
			}
		}
	}
}
