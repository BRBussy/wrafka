package hub

import (
	"gitlab.com/iotTracker/messaging/client"
	"gitlab.com/iotTracker/messaging/message"
)

type Hub interface {
	RegisterClient(client client.Client) error
	Broadcast(message message.Message) error
	SendToClient(identifier client.Identifier, message message.Message) error
}
