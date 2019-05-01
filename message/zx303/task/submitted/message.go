package submitted

import (
	zx303Task "gitlab.com/iotTracker/brain/tracker/zx303/task"
	"gitlab.com/iotTracker/messaging/message"
)

type Message struct {
	Task zx303Task.Task `json:"task"`
}

func (m Message) Type() message.Type {
	return message.ZX303TaskSubmitted
}
