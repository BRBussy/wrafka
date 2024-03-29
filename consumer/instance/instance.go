package instance

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	consumerInstanceException "gitlab.com/iotTracker/messaging/consumer/instance/exception"
	"gitlab.com/iotTracker/messaging/log"
	messageHandler "gitlab.com/iotTracker/messaging/message/handler"
	messagingWrappedMessage "gitlab.com/iotTracker/messaging/message/wrapped"
	"os"
	"os/signal"
	"strings"
)

type instance struct {
	brokers  []string
	topic    string
	handlers []messageHandler.Handler
}

func New(
	brokers []string,
	topic string,
	handlers []messageHandler.Handler,
) *instance {
	return &instance{
		brokers:  brokers,
		topic:    topic,
		handlers: handlers,
	}
}

func (i *instance) Start() error {
	log.Info(fmt.Sprintf(
		"Starting a Consumer Instance for Topic: %s, using Brokers: %s",
		i.topic,
		strings.Join(i.brokers, ", ")),
	)

	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(i.brokers, config)
	if err != nil {
		return consumerInstanceException.Starting{Reasons: []string{"failed to create kafka client", err.Error()}}
	}
	defer func() { _ = client.Close() }()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return consumerInstanceException.Starting{Reasons: []string{"failed to get new consumer from client", err.Error()}}
	}
	// close consumer on termination
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Error(consumerInstanceException.Termination{Reasons: []string{"closing consumer", err.Error()}}.Error())
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(i.topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}
	// close partition consumer on termination
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Error(consumerInstanceException.Termination{Reasons: []string{"closing partition consumer", err.Error()}}.Error())
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

ConsumerLoop:
	for {
		select {
		case message := <-partitionConsumer.Messages():
			var wrappedMessage messagingWrappedMessage.Wrapped
			if err := json.Unmarshal(message.Value, &wrappedMessage); err != nil {
				log.Error(consumerInstanceException.Consumption{Reasons: []string{"unmarshalling wrapped message", err.Error()}}.Error())
				continue ConsumerLoop
			}
			for _, handler := range i.handlers {
				if handler.WantsMessage(wrappedMessage.Message) {
					if err := handler.HandleMessage(wrappedMessage.Message); err != nil {
						log.Error(consumerInstanceException.MessageHandling{Reasons: []string{err.Error()}}.Error())
					}
				}
			}

		case <-signals:
			break ConsumerLoop
		}
	}
	return nil
}
