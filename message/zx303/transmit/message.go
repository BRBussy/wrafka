package transmit

import (
	"gitlab.com/iotTracker/messaging/message"
	nerveServerMessage "gitlab.com/iotTracker/nerve/server/message"
)

type Message struct {
	Message nerveServerMessage.Message `json:"message"`
}

func (m Message) Type() message.Type {
	return message.ZX303Transmit
}
