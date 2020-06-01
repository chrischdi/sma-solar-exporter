package inverter

import (
	"log"

	"github.com/chrischdi/sma-solar-exporter/pkg/yasdi"
)

type Inverter struct {
	Handle       int
	SerialNumber uint
	Name         string
	Channels     map[string]*Channel
}

func NewInverter(deviceHandle int) (*Inverter, error) {
	serial, err := yasdi.GetDeviceSN(deviceHandle)
	if err != nil {
		return nil, err
	}

	name, err := yasdi.GetDeviceName(deviceHandle)
	if err != nil {
		return nil, err
	}

	channelHandles, err := yasdi.GetChannelHandlesEx(deviceHandle, yasdi.MaxDeviceChannelHandles, yasdi.ALLCHANNELS)
	if err != nil {
		log.Printf("warning: error getting channel handles for yasdi.ALLCHANNELS=%d", yasdi.ALLCHANNELS)
	}

	channels := map[string]*Channel{}
	for _, channelHandle := range channelHandles {
		channel, err := NewChannel(deviceHandle, channelHandle)
		if err != nil {
			return nil, err
		}
		channels[channel.Name] = channel
	}

	inv := Inverter{
		Handle:       deviceHandle,
		SerialNumber: serial,
		Name:         name,
		Channels:     channels,
	}

	return &inv, nil
}
