package gps

import (
	zx303GPSReading "gitlab.com/iotTracker/brain/tracker/zx303/reading/gps"
	"gitlab.com/iotTracker/messaging/message"
)

type Message struct {
	Reading zx303GPSReading.Reading `json:"reading"`
}

func (m Message) Type() message.Type {
	return message.ZX303GPSReading
}
