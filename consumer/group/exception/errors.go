package exception

import (
	"fmt"
	"strings"
)

type GroupCreation struct {
	GroupName string
	Reasons   []string
}

func (e GroupCreation) Error() string {
	return fmt.Sprintf("error creating consumer group %s: %s", e.GroupName, strings.Join(e.Reasons, "; "))
}

type Consumption struct {
	Reasons []string
}

func (e Consumption) Error() string {
	return "error consuming: " + strings.Join(e.Reasons, "; ")
}

type Starting struct {
	Reasons []string
}

func (e Starting) Error() string {
	return "error starting consumer group: " + strings.Join(e.Reasons, "; ")
}

type MessageHandling struct {
	Reasons []string
}

func (e MessageHandling) Error() string {
	return "error handling message: " + strings.Join(e.Reasons, "; ")
}

type Termination struct {
	Reasons []string
}

func (e Termination) Error() string {
	return "error terminating consumer group: " + strings.Join(e.Reasons, "; ")
}
