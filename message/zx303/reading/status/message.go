package status

import (
	zx303StatusReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/status"
	"gitlab.com/iotTracker/messaging/message"
)

type Message struct {
	Reading zx303StatusReading.Reading `json:"reading"`
}

func (m Message) Type() message.Type {
	return message.ZX303StatusReading
}
