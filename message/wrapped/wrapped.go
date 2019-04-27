package wrapped

import (
	"encoding/json"
	"gitlab.com/iotTracker/messaging/message"
	"gitlab.com/iotTracker/messaging/message/gpsLocation"
	wrappedMessageException "gitlab.com/iotTracker/messaging/message/wrapped/exception"
)

type Wrapped struct {
	Type    message.Type    `json:"type"`
	Value   json.RawMessage `json:"value"`
	Message message.Message `json:"-"`
}

func Wrap(msg message.Message) (*Wrapped, error) {
	if msg == nil {
		return nil, wrappedMessageException.Wrapping{Reasons: []string{"message to wrap is nil"}}
	}
	value, err := json.Marshal(msg)
	if err != nil {
		return nil, wrappedMessageException.Wrapping{Reasons: []string{"json marshalling", err.Error()}}
	}

	return &Wrapped{
		Type:  msg.Type(),
		Value: value,
	}, nil
}

func (m *Wrapped) UnmarshalJSON(data []byte) error {
	type Alias Wrapped
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return wrappedMessageException.Unwrapping{Reasons: []string{"json unmarshalling wrapped", err.Error()}}
	}

	switch aux.Type {
	case message.GPSLocation:
		var unmarshalledMessage gpsLocation.Message
		if err := json.Unmarshal(m.Value, &unmarshalledMessage); err != nil {
			return wrappedMessageException.Unwrapping{Reasons: []string{"json unmarshalling value", err.Error()}}
		}
		m.Message = unmarshalledMessage

	default:
		return wrappedMessageException.Unwrapping{Reasons: []string{"invalid type", string(aux.Type)}}
	}

	if m.Message == nil {
		return wrappedMessageException.Unwrapping{Reasons: []string{"message is still nil"}}
	}

	return nil
}
