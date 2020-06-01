package inverter

import (
	"fmt"

	"github.com/chrischdi/sma-solar-exporter/pkg/yasdi"
)

// Plant implements the scraper interface to make its metrics available via prometheus
type Plant struct {
	Inverters []*Inverter
	Channels  []string
	yasdiFile string
	driverID  int
	devices   int
}

func NewPlant(yasdiFile string, driverID, devices int, channels []string) (*Plant, error) {
	return &Plant{
		Inverters: nil,
		Channels:  channels,
		yasdiFile: yasdiFile,
		driverID:  driverID,
		devices:   devices,
	}, nil
}

func (p *Plant) Initialize() error {
	if p.Inverters != nil {
		return fmt.Errorf("plant has already assigned inverters, execute Shutdown before using Initialize again")
	}
	if _, err := yasdi.YasdiMasterInitialize(p.yasdiFile); err != nil {
		return err
	}

	return yasdi.YasdiMasterSetDriverOnline(p.driverID)
}

func (p *Plant) Detect() error {
	if err := yasdi.DoStartDeviceDetection(p.devices, true); err != nil {
		return err
	}

	deviceHandles, err := yasdi.GetDeviceHandles(p.devices)
	if err != nil {
		return err
	}

	inverters := []*Inverter{}
	for _, deviceHandle := range deviceHandles {
		i, err := NewInverter(deviceHandle)
		if err != nil {
			return fmt.Errorf("error detecting channels for device %d: %v", deviceHandle, err)
		}
		inverters = append(inverters, i)
	}
	p.Inverters = inverters

	return nil
}

func (p *Plant) Shutdown() {
	p.Inverters = nil
	yasdi.YasdiMasterShutdown()
}
