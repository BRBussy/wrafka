package message

type Type string

type Message interface {
	Type() Type
}
