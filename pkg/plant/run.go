package plant

import (
	"time"

	"github.com/chrischdi/sma-solar-exporter/pkg/inverter"
	"github.com/chrischdi/sma-solar-exporter/pkg/yasdi"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
)

var (
	// detectionInterval is the default interval for the detection retries
	detectionInterval = time.Second * 5
	// detectionSteps are the number of retries for device detection - 0 for unlimited
	detectionSteps = 0
)

func Run(init Initializer, channels []string) (err error) {
	// We want to execute shutdown on errors
	defer func() {
		if err != nil {
			klog.Errorf("error %v detected. executing init.Shutdown()", err)
			init.Shutdown()
		}
	}()

	klog.Info("initializing plant")
	if err := init.Initialize(); err != nil {
		return err
	}

	klog.Info("detecting inverters")
	var inverters []inverter.ScrapeableInverter

	if err := wait.PollImmediateInfinite(detectionInterval, func() (done bool, err error) {
		var innerErr error
		inverters, innerErr = init.Detect()
		if innerErr != nil {
			if errors.Cause(innerErr) == yasdi.ErrNotAllDevicesFound {
				klog.Warningf("detection returned retryable error: %v", innerErr)
				return false, nil
			}
			return false, innerErr
		}
		return true, nil
	}); err != nil {
		return err
	}

	klog.Info("creating plant")
	p, err := NewPlant(init, inverters, channels)
	if err != nil {
		return err
	}

	klog.Info("starting metrics loop")
	StartMetricsLoop(p, channels)
	return nil
}

func StartMetricsLoop(p *Plant, channels []string) {
	go func() {
		// Ensure shutdown is getting executed on fatal errors
		defer func() {
			if recover() != nil {
				p.Initializer.Shutdown()
			}
		}()

		for {
			start := time.Now()
			p.Scrape()
			time.Sleep(5*time.Second - time.Since(start))
		}
	}()
}
