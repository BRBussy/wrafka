package group

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	consumerGroupException "gitlab.com/iotTracker/messaging/consumer/group/exception"
	"gitlab.com/iotTracker/messaging/log"
	messageHandler "gitlab.com/iotTracker/messaging/message/handler"
	messagingWrappedMessage "gitlab.com/iotTracker/messaging/message/wrapped"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// consumer represents a Sarama consumer group consumer
type consumer struct {
	ready    chan bool
	handlers []messageHandler.Handler
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	c.ready <- true
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		var wrappedMessage messagingWrappedMessage.Wrapped
		if err := json.Unmarshal(message.Value, &wrappedMessage); err != nil {
			log.Error(consumerGroupException.Consumption{Reasons: []string{"unmarshalling wrapped message", err.Error()}}.Error())
			continue
		}

		for _, handler := range c.handlers {
			if handler.WantsMessage(wrappedMessage.Message) {
				if err := handler.HandleMessage(wrappedMessage.Message); err != nil {
					log.Error(consumerGroupException.MessageHandling{Reasons: []string{err.Error()}}.Error())
				}
			}
		}

		session.MarkMessage(message, "")
	}

	return nil
}

type group struct {
	brokers   []string
	topics    []string
	handlers  []messageHandler.Handler
	groupName string
}

func New(
	brokers []string,
	topics []string,
	groupName string,
	handlers []messageHandler.Handler,
) *group {
	return &group{
		brokers:   brokers,
		topics:    topics,
		groupName: groupName,
		handlers:  handlers,
	}
}

func (g *group) Start() error {
	log.Info(fmt.Sprintf(
		"Starting a Consumer Group %s for Topics: %s, using Brokers: %s",
		g.groupName,
		strings.Join(g.topics, ", "),
		strings.Join(g.brokers, ", ")),
	)

	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(g.brokers, config)
	if err != nil {
		return consumerGroupException.Starting{Reasons: []string{"failed to create kafka client", err.Error()}}
	}
	// close client on termination
	defer func() {
		if err = client.Close(); err != nil {
			log.Error(consumerGroupException.Termination{Reasons: []string{"closing client", err.Error()}}.Error())
		}
	}()

	consumer := consumer{
		ready:    make(chan bool, 0),
		handlers: g.handlers,
	}

	ctx := context.Background()
	consumerGroup, err := sarama.NewConsumerGroupFromClient(g.groupName, client)
	if err != nil {
		return consumerGroupException.GroupCreation{GroupName: g.groupName, Reasons: []string{err.Error()}}
	}
	// close group on termination
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Error(consumerGroupException.Termination{Reasons: []string{"closing consumer group", err.Error()}}.Error())
		}
	}()

	go func() {
		for {
			err := consumerGroup.Consume(ctx, g.topics, &consumer)
			if err != nil {
				log.Fatal(consumerGroupException.Consumption{Reasons: []string{err.Error()}}.Error())
			}
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Info(fmt.Sprintf("Consumer Group %s up and running", g.groupName))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm

	return nil
}
