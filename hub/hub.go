package hub

import (
	"gitlab.com/iotTracker/messaging/client"
	"gitlab.com/iotTracker/messaging/message"
)

type Hub interface {
	RegisterClient(client client.Client) error
	DeRegisterClient(client client.Client) error
	ReRegisterClient(client client.Client) error
	Broadcast(message message.Message) error
	SendToClient(identifier client.Identifier, message message.Message) error
	GetClient(identifier client.Identifier) (client.Client, error)
	ClientCount() int
}
