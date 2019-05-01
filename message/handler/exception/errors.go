package exception

import "strings"

type Handling struct {
	Reasons []string
}

func (e Handling) Error() string {
	return "error handling message: " + strings.Join(e.Reasons, "; ")
}
