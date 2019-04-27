package exception

import "strings"

type InvalidMessage struct {
	Reasons []string
}

func (e InvalidMessage) Error() string {
	return "invalid message: " + strings.Join(e.Reasons, "; ")
}
