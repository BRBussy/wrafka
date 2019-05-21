package client

import (
	"fmt"
	"gitlab.com/iotTracker/messaging/message"
)

type Client interface {
	Send(message message.Message) error
	IdentifiedBy(identifier Identifier) bool
	Identifier() Identifier
}

type Type string

const ZX303 Type = "ZX303"

type Identifier struct {
	Type Type
	Id   string
}

func (i Identifier) String() string {
	return fmt.Sprintf("type: %s, id: %s", i.Type, i.Id)
}
