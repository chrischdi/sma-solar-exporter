package inverter

import "fmt"

type Simulation struct {
	Serial       uint
	Name         string
	Channels     map[string]SimulatedChannel
	ChannelUnits map[string]string
}

func NewSimulation(serial uint, name string) *Simulation {
	return &Simulation{
		Serial: serial,
		Name:   name,
		Channels: map[string]SimulatedChannel{
			"A.Ms.Amp":        SimulatedChannel{0.2, 1.4, 0.1, 0.4},
			"A.Ms.Vol":        SimulatedChannel{200, 400, 20.2, 220},
			"A.Ms.Watt":       SimulatedChannel{0, 500, 11.4, 220},
			"B.Ms.Amp":        SimulatedChannel{0.3, 1.5, 0.1, 0.6},
			"B.Ms.Vol":        SimulatedChannel{200, 400, 20.2, 260},
			"B.Ms.Watt":       SimulatedChannel{0, 500, 11.4, 320},
			"E-Total":         SimulatedChannel{0, -1, 0.5, 10000},
			"GridMs.PhV.phsA": SimulatedChannel{229.750, 230.250, 0.01, 229.750},
			"GridMs.PhV.phsB": SimulatedChannel{0, -1, 0, 0},
			"GridMs.PhV.phsC": SimulatedChannel{0, -1, 0, 0},
			"Mt.TotOpTmh":     SimulatedChannel{0, -1, 0.0014, 3000},
			"Mt.TotTmh":       SimulatedChannel{0, -1, 0.0014, 2500},
			"Pac":             SimulatedChannel{0, 500, 10, 0},
		},
		ChannelUnits: map[string]string{
			"A.Ms.Amp":        "A",
			"A.Ms.Vol":        "V",
			"A.Ms.Watt":       "W",
			"B.Ms.Amp":        "A",
			"B.Ms.Vol":        "V",
			"B.Ms.Watt":       "W",
			"E-Total":         "kWh",
			"GridMs.PhV.phsA": "V",
			"GridMs.PhV.phsB": "V",
			"GridMs.PhV.phsC": "V",
			"Mt.TotOpTmh":     "h",
			"Mt.TotTmh":       "h",
			"Pac":             "W",
		},
	}
}

func (i *Simulation) GetChannelValue(channel string) (float64, error) {
	c, ok := i.Channels[channel]
	if !ok {
		return 0, fmt.Errorf("channel %s not found", channel)
	}
	val := c.iterate()
	i.Channels[channel] = c
	return val, nil
}

func (i *Simulation) GetName() string {
	return i.Name
}

func (i *Simulation) GetChannelUnit(channelName string) (string, error) {
	unit, ok := i.ChannelUnits[channelName]
	if !ok {
		return "", fmt.Errorf("channel %s not found", channelName)
	}
	return unit, nil
}

func (i *Simulation) GetSerialNumber() uint {
	return i.Serial
}

type SimulatedChannel struct {
	Min   float64
	Max   float64
	Adder float64
	Value float64
}

func (c *SimulatedChannel) iterate() float64 {
	current := c.Value
	c.Value = c.Value + c.Adder
	if c.Value > c.Max && c.Max != -1 {
		c.Value = c.Min
	}
	return current
}
