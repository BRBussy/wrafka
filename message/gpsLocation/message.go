package gpsLocation

import (
	"encoding/json"
	"gitlab.com/iotTracker/brain/search/identifier"
	wrappedIdentifier "gitlab.com/iotTracker/brain/search/identifier/wrapped"
	"gitlab.com/iotTracker/brain/tracker/device"
	"gitlab.com/iotTracker/messaging/log"
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

	log.Info("wrapping!!!!!!")

	type Alias Message
	return json.Marshal(&struct {
		Alias
		DeviceId wrappedIdentifier.Wrapped `json:"deviceId"`
	}{
		DeviceId: *wrappedDeviceId,
		Alias:    (Alias)(m),
	})
}
