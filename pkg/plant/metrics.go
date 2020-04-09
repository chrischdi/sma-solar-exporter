package plant

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// The metrics, for descriptions see RegisterMetrics

	MetricDCCurrent          *prometheus.GaugeVec
	MetricDCVoltage          *prometheus.GaugeVec
	MetricDCPower            *prometheus.GaugeVec
	MetricETotal             *prometheus.GaugeVec
	MetricHoursTotal         *prometheus.GaugeVec
	MetricInverterHoursTotal *prometheus.GaugeVec
	MetricGridVoltage        *prometheus.GaugeVec
	MetricActivePower        *prometheus.GaugeVec
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
	channelToMetricMapping map[string]*prometheus.GaugeVec

	metricLabels = []string{"serial", "inverter_name", "channel_name", "unit"}
)

func RegisterMetrics() {
	MetricDCCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sma",
		Subsystem: "dc",
		Name:      "current",
		Help:      "DC current in A",
	}, metricLabels)
	MetricDCVoltage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sma",
		Subsystem: "dc",
		Name:      "voltage",
		Help:      "DC voltage in V",
	}, metricLabels)
	MetricDCPower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sma",
		Subsystem: "dc",
		Name:      "power",
		Help:      "DC power",
	}, metricLabels)

	MetricETotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sma",
		Subsystem: "feeding",
		Name:      "energy_total",
		Help:      "Total amount of feeding-in energy in kWh",
	}, metricLabels)
	MetricHoursTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sma",
		Subsystem: "feeding",
		Name:      "hours_total",
		Help:      "Total number of grid-feeding operational hours",
	}, metricLabels)
	MetricInverterHoursTotal = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sma",
		Subsystem: "operating",
		Name:      "hours_total",
		Help:      "Total number of operating hours of inverter",
	}, metricLabels)

	MetricGridVoltage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sma",
		Subsystem: "grid",
		Name:      "phase_voltage",
		Help:      "Grid voltage on phase in V",
	}, metricLabels)
	MetricActivePower = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "sma",
		// Subsystem: "grid",
		Name: "actual_power",
		Help: "Delivered active power in W",
	}, metricLabels)

	MetricScrapeDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "sma",
		Subsystem: "scrape",
		Name:      "duration_seconds",
		Help:      "The time needed to scrape the plant",
	}, nil)

	channelToMetricMapping = map[string]*prometheus.GaugeVec{
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
