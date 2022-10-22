package plant

import (
	"fmt"
	"testing"
	"time"

	"github.com/chrischdi/sma-solar-exporter/pkg/inverter"
	"github.com/chrischdi/sma-solar-exporter/pkg/yasdi"
)

type DummyInitializer struct {
	initErr    error
	detectErrs []error
	inverters  []inverter.ScrapeableInverter
}

func (i *DummyInitializer) Initialize() error {
	return i.initErr
}
func (i *DummyInitializer) Detect() ([]inverter.ScrapeableInverter, error) {
	if len(i.detectErrs) == 0 {
		return i.inverters, nil
	}
	err := i.detectErrs[0]
	i.detectErrs = i.detectErrs[1:]
	return nil, err
}
func (i *DummyInitializer) Shutdown() {}

func Test_Run(t *testing.T) {
	detectionInterval = time.Millisecond
	detectionSteps = 3
	tests := []struct {
		name     string
		init     Initializer
		channels []string
		wantErr  bool
	}{
		{
			"Initialize error",
			&DummyInitializer{
				initErr: fmt.Errorf("error initializing"),
			},
			[]string{"A.Ms.Amp", "A.Ms.Vol", "A.Ms.Watt", "E-Total", "GridMs.PhV.phsA", "Mt.TotOpTmh", "Mt.TotTmh", "Pac"},
			true,
		},
		{
			"non-retryable detection",
			&DummyInitializer{
				detectErrs: []error{yasdi.ErrDetectionAlreadyRunning},
			},
			[]string{"A.Ms.Amp", "A.Ms.Vol", "A.Ms.Watt", "E-Total", "GridMs.PhV.phsA", "Mt.TotOpTmh", "Mt.TotTmh", "Pac"},
			true,
		},
		{
			"Happy Path",
			&DummyInitializer{
				detectErrs: []error{},
				inverters:  []inverter.ScrapeableInverter{},
			},
			[]string{"A.Ms.Amp", "A.Ms.Vol", "A.Ms.Watt", "E-Total", "GridMs.PhV.phsA", "Mt.TotOpTmh", "Mt.TotTmh", "Pac"},
			false,
		},
		{
			"Happy Path after multiple detections",
			&DummyInitializer{
				detectErrs: []error{yasdi.ErrNotAllDevicesFound, yasdi.ErrNotAllDevicesFound},
				inverters:  []inverter.ScrapeableInverter{},
			},
			[]string{"A.Ms.Amp", "A.Ms.Vol", "A.Ms.Watt", "E-Total", "GridMs.PhV.phsA", "Mt.TotOpTmh", "Mt.TotTmh", "Pac"},
			false,
		},
		{
			"Detection failure after multiple detections",
			&DummyInitializer{
				detectErrs: []error{yasdi.ErrNotAllDevicesFound, yasdi.ErrNotAllDevicesFound, yasdi.ErrNotAllDevicesFound, yasdi.ErrDetectionAlreadyRunning},
				inverters:  []inverter.ScrapeableInverter{},
			},
			[]string{"A.Ms.Amp", "A.Ms.Vol", "A.Ms.Watt", "E-Total", "GridMs.PhV.phsA", "Mt.TotOpTmh", "Mt.TotTmh", "Pac"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.init, tt.channels, nil); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
