package inverter

import (
	"fmt"

	"github.com/chrischdi/sma-solar-exporter/pkg/yasdi"
)

// Initialize initializes yasdi and detects devices.
// If Initialize returns without an error, then you will need to execute Shutdown
// after usage.
func Initialize(yasdiConfigFile string, driver, devices int) (handles []int, err error) {
	// if we exit with error we run Shutdown
	defer func() {
		if err != nil {
			Shutdown()
		}
	}()

	// initialize yasdi
	drivers, err := yasdi.YasdiMasterInitialize(yasdiConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error initializing yasdi using %q", yasdiConfigFile)
	}

	if driver >= drivers {
		return nil, fmt.Errorf("driver %d was not found (got %d)", driver, drivers)
	}

	// activate the driver
	if err := yasdi.YasdiMasterSetDriverOnline(driver); err != nil {
		return nil, fmt.Errorf("error setting yasdi driver %d online", driver)
	}

	// detect the devices
	if err := yasdi.DoStartDeviceDetection(devices, true); err != nil {
		return nil, fmt.Errorf("error during device detection: %v", err)
	}

	// get the device handles
	deviceHandles, err := yasdi.GetDeviceHandles(devices)
	if err != nil {
		return nil, fmt.Errorf("error getting device handles: %v", err)
	}

	return deviceHandles, nil
}

// Shutdown needs to be run after successful initialization and usage.
func Shutdown() {
	yasdi.YasdiMasterShutdown()
}
