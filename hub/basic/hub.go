package basic

import (
	messagingClient "gitlab.com/iotTracker/messaging/client"
	messagingHub "gitlab.com/iotTracker/messaging/hub"
	hubException "gitlab.com/iotTracker/messaging/hub/exception"
	"gitlab.com/iotTracker/messaging/message"
)

type hub struct {
	clients map[messagingClient.Identifier]messagingClient.Client
}

func New() messagingHub.Hub {
	return &hub{
		clients: make(map[messagingClient.Identifier]messagingClient.Client),
	}
}

func (h *hub) GetClient(identifier messagingClient.Identifier) (messagingClient.Client, error) {
	client, clientRegistered := h.clients[identifier]
	if !clientRegistered {
		return nil, hubException.GetClient{
			ClientId: identifier,
			Reasons:  []string{"no such client in hub"},
		}
	}

	return client, nil
}

func (h *hub) ClientCount() int {
	return len(h.clients)
}

func (h *hub) RegisterClient(client messagingClient.Client) error {
	// check if the client identifier is blank, cannot be registered
	if client.Identifier().Id == "" || client.Identifier().Type == "" {
		return hubException.ClientRegistration{Reasons: []string{"identifier is blank", client.Identifier().String()}}
	}

	// check if the client is already registered
	if _, clientRegistered := h.clients[client.Identifier()]; clientRegistered {
		return hubException.ClientAlreadyRegistered{ClientId: client.Identifier()}
	}
	// if not register the client
	h.clients[client.Identifier()] = client
	return nil
}

func (h *hub) DeRegisterClient(client messagingClient.Client) error {
	// check if the client identifier is blank, cannot be registered
	if client.Identifier().Id == "" || client.Identifier().Type == "" {
		return hubException.ClientDeRegistration{Reasons: []string{"identifier is blank", client.Identifier().String()}}
	}

	// check if the client is registered on this hub
	if _, clientRegistered := h.clients[client.Identifier()]; !clientRegistered {
		return hubException.ClientDeRegistration{Reasons: []string{"client not in hub", client.Identifier().String()}}
	}

	// remove client from hub
	delete(h.clients, client.Identifier())
	return nil
}

func (h *hub) ReRegisterClient(client messagingClient.Client) error {
	// check if the client identifier is blank, cannot be registered
	if client.Identifier().Id == "" || client.Identifier().Type == "" {
		return hubException.ClientDeRegistration{Reasons: []string{"identifier is blank", client.Identifier().String()}}
	}

	// check if the client is registered on this hub
	if _, clientRegistered := h.clients[client.Identifier()]; !clientRegistered {
		return hubException.ClientDeRegistration{Reasons: []string{"client not in hub", client.Identifier().String()}}
	}

	// remove old client from hub
	delete(h.clients, client.Identifier())

	// add new client back to hub
	h.clients[client.Identifier()] = client

	return nil
}

func (h *hub) Broadcast(message message.Message) error {
	sendErrors := make([]string, 0)
	for _, client := range h.clients {
		if err := client.Send(message); err != nil {
			sendErrors = append(sendErrors, hubException.SendToClient{
				ClientId: client.Identifier(),
				Reasons:  []string{err.Error()},
			}.Error())
		}
	}

	if len(sendErrors) > 0 {
		return hubException.Broadcast{Reasons: sendErrors}
	}
	return nil
}

func (h *hub) SendToClient(identifier messagingClient.Identifier, message message.Message) error {
	// get client from map of registered clients
	client, clientRegistered := h.clients[identifier]
	if !clientRegistered {
		return hubException.SendToClient{
			ClientId: identifier,
			Reasons:  []string{"no such client in hub"},
		}
	}
	// try and send the message to the client
	if err := client.Send(message); err != nil {
		return hubException.SendToClient{
			ClientId: client.Identifier(),
			Reasons:  []string{err.Error()},
		}
	}

	return nil
}
