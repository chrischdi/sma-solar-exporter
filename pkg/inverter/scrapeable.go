package inverter

import (
	"fmt"
)

type ScrapeableInverter interface {
	GetChannelValue(channel string) (float64, error)
	GetName() string
	GetChannelUnit(channelName string) (string, error)
	GetSerialNumber() uint
}

func (i *Inverter) GetChannelValue(channelName string) (float64, error) {
	c, ok := i.Channels[channelName]
	if !ok {
		return 0, fmt.Errorf("channel %q was not found for inverter %d", channelName, i.SerialNumber)
	}

	val, _, err := c.GetValues()
	if err != nil {
		// TODO check if alternative on errors.Is(err, yasdi.ErrTimeout) is to return 0,nil
		return 0, err
	}
	return val, nil
}

func (i *Inverter) GetName() string {
	return i.Name
}

func (i *Inverter) GetChannelUnit(channelName string) (string, error) {
	c, ok := i.Channels[channelName]
	if !ok {
		return "", fmt.Errorf("channel %q was not found for inverter %d", channelName, i.SerialNumber)
	}

	return c.Unit, nil
}

func (i *Inverter) GetSerialNumber() uint {
	return i.SerialNumber
}
