package exception

import (
	"fmt"
	"gitlab.com/iotTracker/messaging/client"
	"strings"
)

type ClientRegistration struct {
	Reasons []string
}

func (e ClientRegistration) Error() string {
	return "client registration error: " + strings.Join(e.Reasons, "; ")
}

type ClientDeRegistration struct {
	Reasons []string
}

func (e ClientDeRegistration) Error() string {
	return "client deRegistration error: " + strings.Join(e.Reasons, "; ")
}

type ClientAlreadyRegistered struct {
	ClientId client.Identifier
}

func (e ClientAlreadyRegistered) Error() string {
	return fmt.Sprintf("client already registered: %s", e.ClientId.String())
}

type Broadcast struct {
	Reasons []string
}

func (e Broadcast) Error() string {
	return "broadcast error: " + strings.Join(e.Reasons, "; ")
}

type SendToClient struct {
	ClientId client.Identifier
	Reasons  []string
}

func (e SendToClient) Error() string {
	return fmt.Sprintf("error sending to client: %s, %s", e.ClientId.String(), strings.Join(e.Reasons, "; "))
}

type GetClient struct {
	ClientId client.Identifier
	Reasons  []string
}

func (e GetClient) Error() string {
	return fmt.Sprintf("error getting client: %s, %s", e.ClientId.String(), strings.Join(e.Reasons, "; "))
}

type StopClient struct {
	ClientId client.Identifier
	Reasons  []string
}

func (e StopClient) Error() string {
	return fmt.Sprintf("error stopping client: %s, %s", e.ClientId.String(), strings.Join(e.Reasons, "; "))
}
