package hub

import (
	"gitlab.com/iotTracker/messaging/client"
	"gitlab.com/iotTracker/messaging/message"
)

type Hub interface {
	Broadcast(message message.Message) error
	SendToClient(identifier client.Identifier, message message.Message) error
}
