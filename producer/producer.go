package producer

import "gitlab.com/iotTracker/messaging/message"

type Producer interface {
	Start() error
	Produce(message message.Message) error
}
