package plant

import (
	"fmt"
	"time"

	"github.com/chrischdi/sma-solar-exporter/pkg/inverter"
	"github.com/chrischdi/sma-solar-exporter/pkg/yasdi"
	"github.com/pkg/errors"
	"k8s.io/klog"
)

// Plant implements the scraper interface to make its metrics available via prometheus
type Plant struct {
	Inverters   []inverter.ScrapeableInverter
	Channels    []string
	Initializer Initializer
}

func NewPlant(init Initializer, inverters []inverter.ScrapeableInverter, channels []string) (*Plant, error) {
	// check if there are channels given we cannot handle
	for _, channel := range channels {
		if _, ok := channelToMetricMapping[channel]; !ok {
			return nil, fmt.Errorf("channel %q is not supported to get exported as metric", channel)
		}
	}

	return &Plant{
		Inverters:   inverters,
		Channels:    channels,
		Initializer: init,
	}, nil
}

func (p *Plant) Scrape() {
	start := time.Now()

	for _, i := range p.Inverters {
		for _, channelName := range p.Channels {
			metric, ok := channelToMetricMapping[channelName]
			if !ok {
				klog.Warningf("channel %q has no corresponding metric to expose", channelName)
				continue
			}

			metricMinimum := metricsMinimumValue[channelName]

			v, err := i.GetChannelValue(channelName)
			if err != nil && errors.Cause(err) != yasdi.ErrTimeout {
				klog.Warningf("error getting values for channel %q for inverter %d: %v", channelName, i.GetSerialNumber(), err)
				continue
			}

			if metricMinimum > v {
				continue
			}

			unit, err := i.GetChannelUnit(channelName)
			if err != nil {
				klog.Warningf("error getting unit for channel %q for inverter %d: %v", channelName, i.GetSerialNumber(), err)
			}

			metric.WithLabelValues(fmt.Sprintf("%d", i.GetSerialNumber()), i.GetName(), channelName, unit).Set(v)
		}
	}
	MetricScrapeDuration.WithLabelValues().Observe(time.Since(start).Seconds())
}
