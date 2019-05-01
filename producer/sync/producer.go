package sync

import (
	"encoding/json"
	"fmt"
	"gitlab.com/iotTracker/messaging/log"
	"gitlab.com/iotTracker/messaging/message"
	wrappedMessage "gitlab.com/iotTracker/messaging/message/wrapped"
	messagingProducer "gitlab.com/iotTracker/messaging/producer"
	producerException "gitlab.com/iotTracker/messaging/producer/exception"
	"gopkg.in/Shopify/sarama.v1"
	"strings"
)

type producer struct {
	producer sarama.SyncProducer
	brokers  []string
	topic    string
}

func New(
	brokers []string,
	topic string,
) messagingProducer.Producer {
	return &producer{
		brokers: brokers,
		topic:   topic,
	}
}

func (p *producer) Start() error {
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(p.brokers, config)
	if err != nil {
		return producerException.Start{Reasons: []string{"failed to connect new producer", err.Error()}}
	}

	log.Info(fmt.Sprintf("Started Producer for Topic: %s, Using Brokers: %s", p.topic, strings.Join(p.brokers, ", ")))

	p.producer = producer

	return nil
}

func (p *producer) Produce(message message.Message) error {
	// We are not setting a message key, which means that all messages will
	// be distributed randomly over the different partitions.

	wrappedMessage, err := wrappedMessage.Wrap(message)
	if err != nil {
		return producerException.Produce{Reasons: []string{"wrapping", err.Error()}}
	}
	messageData, err := json.Marshal(wrappedMessage)
	if err != nil {
		return producerException.Produce{Reasons: []string{"marshalling wrapped", err.Error()}}
	}

	//partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(messageData),
	})
	if err != nil {
		return producerException.Produce{Reasons: []string{err.Error()}}
	}
	return nil
}
