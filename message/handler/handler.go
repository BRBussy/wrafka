package handler

import "gitlab.com/iotTracker/messaging/message"

type Handler interface {
	WantsMessage(message message.Message) bool
	HandleMessage(message message.Message)
}
