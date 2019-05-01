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
