package plant

import (
	"fmt"

	"github.com/chrischdi/sma-solar-exporter/pkg/inverter"
)

type SimulationInitializer struct {
	Devices int
}

func (i *SimulationInitializer) Initialize() error {
	return nil
}

func (i *SimulationInitializer) Shutdown() {
}

func (i *SimulationInitializer) Detect() ([]inverter.ScrapeableInverter, error) {
	inverters := []inverter.ScrapeableInverter{}
	for j := 0; j < i.Devices; j++ {
		inverters = append(inverters, inverter.NewSimulation(uint(j), fmt.Sprintf("inverter %d", j)))
	}
	return inverters, nil
}
