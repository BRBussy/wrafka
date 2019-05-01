package basic

import (
	"gitlab.com/iotTracker/messaging/client"
	messagingHub "gitlab.com/iotTracker/messaging/hub"
	hubException "gitlab.com/iotTracker/messaging/hub/exception"
	"gitlab.com/iotTracker/messaging/message"
)

type hub struct {
	Clients []client.Client
}

func New() messagingHub.Hub {
	return &hub{}
}

func (h *hub) Broadcast(message message.Message) error {
	sendErrors := make([]string, 0)
	for _, c := range h.Clients {
		if err := c.Send(message); err != nil {
			sendErrors = append(sendErrors, hubException.SendToClient{
				ClientId: c.Identifier(),
				Reasons:  []string{err.Error()},
			}.Error())
		}
	}

	if len(sendErrors) > 0 {
		return hubException.Broadcast{Reasons: sendErrors}
	}
	return nil
}

func (h *hub) SendToClient(identifier client.Identifier, message message.Message) error {
	for _, c := range h.Clients {
		if c.IdentifiedBy(identifier) {
			if err := c.Send(message); err != nil {
				return hubException.SendToClient{
					ClientId: c.Identifier(),
					Reasons:  []string{err.Error()},
				}
			}
			return nil
		}
	}
	return hubException.SendToClient{
		ClientId: identifier,
		Reasons:  []string{"no such client in hub"},
	}
}
