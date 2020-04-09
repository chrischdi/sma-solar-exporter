package plant

import "github.com/chrischdi/sma-solar-exporter/pkg/inverter"

type Initializer interface {
	Initialize() error
	Detect() ([]inverter.ScrapeableInverter, error)
	Shutdown()
}
