package gpsLocation

import (
	"encoding/json"
	"errors"
	"gitlab.com/iotTracker/brain/search/identifier"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	"gitlab.com/iotTracker/brain/tracker/device"
	"gitlab.com/iotTracker/messaging/message"
)

type Message struct {
	// Device Details
	DeviceId   identifier.Identifier `json:"deviceId"`
	DeviceType device.Type           `json:"deviceType"`

	// Reading Details
	TimeStamp int64   `json:"timeStamp"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func (m Message) Type() message.Type {
	return message.GPSLocation
}

func (m Message) MarshalJSON() ([]byte, error) {

	wrappedDeviceId, err := wrappedIdentifier.Wrap(m.DeviceId)
	if err != nil {
		return nil, err
	}

	type Alias Message
	return json.Marshal(&struct {
		Alias
		DeviceId wrappedIdentifier.Wrapped `json:"deviceId"`
	}{
		DeviceId: *wrappedDeviceId,
		Alias:    (Alias)(m),
	})
}

func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		*Alias
		DeviceId wrappedIdentifier.Wrapped `json:"deviceId"`
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return errors.New("unmarshalling: " + err.Error())
	}

	m.DeviceId = aux.DeviceId.Identifier
	m.DeviceType = aux.DeviceType
	m.TimeStamp = aux.TimeStamp
	m.Latitude = aux.Latitude
	m.Longitude = aux.Longitude

	return nil
}
