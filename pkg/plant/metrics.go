package plant

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricItems struct {
	PrometheusGauge *prometheus.GaugeVec
	InfluxPoint *write.Point
}

var (
	// The metrics, for descriptions see RegisterMetrics

	MetricDCCurrent          *MetricItems
	MetricDCVoltage          *MetricItems
	MetricDCPower            *MetricItems
	MetricETotal             *MetricItems
	MetricHoursTotal         *MetricItems
	MetricInverterHoursTotal *MetricItems
	MetricGridVoltage        *MetricItems
	MetricActivePower        *MetricItems
	MetricScrapeDuration     *prometheus.HistogramVec

	// metricsMinimumValue is used during scrape to skip too low values.
	// There are some metrics which only increase over time but as soon as the
	// inverters go offline they would drop to 0.
	metricsMinimumValue = map[string]float64{
		"E-Total":     1,
		"Mt.TotOpTmh": 1,
		"Mt.TotTmh":   1,
	}

	// channelToMetricMapping contains a mapping from channels to their corresponding metric
	channelToMetricMapping map[string]*MetricItems

	metricLabels = []string{"serial", "inverter_name", "channel_name", "unit"}
)

func RegisterMetrics() {
	MetricDCCurrent = &MetricItems{
		PrometheusGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sma",
			Subsystem: "dc",
			Name:      "current",
			Help:      "DC current in A",
		}, metricLabels),
		InfluxPoint: influxdb2.NewPointWithMeasurement("sma.dc.current"),
	}
	MetricDCVoltage = &MetricItems{
		PrometheusGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sma",
			Subsystem: "dc",
			Name:      "voltage",
			Help:      "DC voltage in V",
		}, metricLabels),
		InfluxPoint: influxdb2.NewPointWithMeasurement("sma.dc.voltage"),
	}
	MetricDCPower = &MetricItems{
		PrometheusGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sma",
			Subsystem: "dc",
			Name:      "power",
			Help:      "DC power",
		}, metricLabels),
		InfluxPoint: influxdb2.NewPointWithMeasurement("sma.dc.power"),
	}

	MetricETotal = &MetricItems{
		PrometheusGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sma",
			Subsystem: "feeding",
			Name:      "energy_total",
			Help:      "Total amount of feeding-in energy in kWh",
		}, metricLabels),
		InfluxPoint: influxdb2.NewPointWithMeasurement("sma.feeding.energy_total"),
	}
	MetricHoursTotal = &MetricItems{
		PrometheusGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sma",
			Subsystem: "feeding",
			Name:      "hours_total",
			Help:      "Total number of grid-feeding operational hours",
		}, metricLabels),
		InfluxPoint: influxdb2.NewPointWithMeasurement("sma.feeding.hours_total"),
	}
	MetricInverterHoursTotal = &MetricItems{
		PrometheusGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sma",
			Subsystem: "operating",
			Name:      "hours_total",
			Help:      "Total number of operating hours of inverter",
		}, metricLabels),
		InfluxPoint: influxdb2.NewPointWithMeasurement("sma.operating.hours_total"),
	}

	MetricGridVoltage = &MetricItems{
		PrometheusGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sma",
			Subsystem: "grid",
			Name:      "phase_voltage",
			Help:      "Grid voltage on phase in V",
		}, metricLabels),
		InfluxPoint: influxdb2.NewPointWithMeasurement("sma.grid.phase_voltage"),
	}
	MetricActivePower = &MetricItems{
		PrometheusGauge: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sma",
			// Subsystem: "grid",
			Name: "actual_power",
			Help: "Delivered active power in W",
		}, metricLabels),
		InfluxPoint: influxdb2.NewPointWithMeasurement("sma.grid.actual_power"),
	}

	MetricScrapeDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "sma",
		Subsystem: "scrape",
		Name:      "duration_seconds",
		Help:      "The time needed to scrape the plant",
	}, nil)

	channelToMetricMapping = map[string]*MetricItems {
		"A.Ms.Amp":        MetricDCCurrent,
		"B.Ms.Amp":        MetricDCCurrent,
		"A.Ms.Vol":        MetricDCVoltage,
		"B.Ms.Vol":        MetricDCVoltage,
		"A.Ms.Watt":       MetricDCPower,
		"B.Ms.Watt":       MetricDCPower,
		"E-Total":         MetricETotal,
		"GridMs.PhV.phsA": MetricGridVoltage,
		"GridMs.PhV.phsB": MetricGridVoltage,
		"GridMs.PhV.phsC": MetricGridVoltage,
		"Mt.TotOpTmh":     MetricHoursTotal,
		"Mt.TotTmh":       MetricInverterHoursTotal,
		"Pac":             MetricActivePower,
	}
}
