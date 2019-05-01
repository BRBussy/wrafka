package wrapped

import (
	"encoding/json"
	"gitlab.com/iotTracker/messaging/message"
	wrappedMessageException "gitlab.com/iotTracker/messaging/message/wrapped/exception"
	zx303GPSReadingMessage "gitlab.com/iotTracker/messaging/message/zx303/reading/gps"
	zx303StatusReadingMessage "gitlab.com/iotTracker/messaging/message/zx303/reading/status"
	zx303TaskSubmittedMessage "gitlab.com/iotTracker/messaging/message/zx303/task/submitted"
	zx303TaskTransitionedMessage "gitlab.com/iotTracker/messaging/message/zx303/task/transitioned"
	zx303TransmitMessage "gitlab.com/iotTracker/messaging/message/zx303/transmit"
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
	case message.ZX303GPSReading:
		var unmarshalledMessage zx303GPSReadingMessage.Message
		if err := json.Unmarshal(m.Value, &unmarshalledMessage); err != nil {
			return wrappedMessageException.Unwrapping{Reasons: []string{"json unmarshalling value", err.Error()}}
		}
		m.Message = unmarshalledMessage

	case message.ZX303StatusReading:
		var unmarshalledMessage zx303StatusReadingMessage.Message
		if err := json.Unmarshal(m.Value, &unmarshalledMessage); err != nil {
			return wrappedMessageException.Unwrapping{Reasons: []string{"json unmarshalling value", err.Error()}}
		}
		m.Message = unmarshalledMessage

	case message.ZX303Transmit:
		var unmarshalledMessage zx303TransmitMessage.Message
		if err := json.Unmarshal(m.Value, &unmarshalledMessage); err != nil {
			return wrappedMessageException.Unwrapping{Reasons: []string{"json unmarshalling value", err.Error()}}
		}
		m.Message = unmarshalledMessage

	case message.ZX303TaskSubmitted:
		var unmarshalledMessage zx303TaskSubmittedMessage.Message
		if err := json.Unmarshal(m.Value, &unmarshalledMessage); err != nil {
			return wrappedMessageException.Unwrapping{Reasons: []string{"json unmarshalling value", err.Error()}}
		}
		m.Message = unmarshalledMessage

	case message.ZX303TaskTransitioned:
		var unmarshalledMessage zx303TaskTransitionedMessage.Message
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
