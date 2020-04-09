package plant

import (
	"fmt"

	"github.com/chrischdi/sma-solar-exporter/pkg/inverter"
	"github.com/chrischdi/sma-solar-exporter/pkg/yasdi"
)

type YasdiInitializer struct {
	YasdiFile string
	DriverID  int
	Devices   int
}

func (i *YasdiInitializer) Initialize() error {
	if _, err := yasdi.YasdiMasterInitialize(i.YasdiFile); err != nil {
		return err
	}
	return yasdi.YasdiMasterSetDriverOnline(i.DriverID)
}

func (i *YasdiInitializer) Shutdown() {
	yasdi.YasdiMasterShutdown()
}

func (i *YasdiInitializer) Detect() ([]inverter.ScrapeableInverter, error) {
	if err := yasdi.DoStartDeviceDetection(i.Devices, true); err != nil {
		return nil, err
	}

	deviceHandles, err := yasdi.GetDeviceHandles(i.Devices)
	if err != nil {
		return nil, err
	}

	inverters := []inverter.ScrapeableInverter{}
	for _, deviceHandle := range deviceHandles {
		i, err := inverter.NewInverter(deviceHandle)
		if err != nil {
			return nil, fmt.Errorf("error detecting channels for device %d: %v", deviceHandle, err)
		}
		inverters = append(inverters, i)
	}
	return inverters, nil
}
