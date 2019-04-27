package exception

import "strings"

type Wrapping struct {
	Reasons []string
}

func (e Wrapping) Error() string {
	return "wrapping error: " + strings.Join(e.Reasons, "; ")
}

type Unwrapping struct {
	Reasons []string
}

func (e Unwrapping) Error() string {
	return "unwrapping error: " + strings.Join(e.Reasons, "; ")
}
