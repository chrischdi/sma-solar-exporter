package main

import (
	"flag"
	"net/http"
	"strings"

	"github.com/chrischdi/sma-solar-exporter/pkg/plant"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog"
)

var (
	metricsAddr    = flag.String("metrics-addr", ":8080", "The address the metric endpoint binds to.")
	metricChannels = flag.String("metric-channels", "A.Ms.Amp,A.Ms.Vol,A.Ms.Watt,E-Total,GridMs.PhV.phsA,Mt.TotOpTmh,Mt.TotTmh,Pac", "A comma separated list of channel names which should get exposed as metric.")

	yasdiConfig  = flag.String("yasdi-config", "/etc/yasdi/yasdi.ini", "The path to the yasdi config file.")
	yasdiDriver  = flag.Int("yasdi-driver", 0, "The driver reference in yasdi-config to use.")
	yasdiDevices = flag.Int("yasdi-devices", 2, "The number of inverters expected to get detected.")
)

func main() {
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	flag.Parse()

	channels := strings.Split(*metricChannels, ",")
	if len(channels) == 0 || (len(channels) == 1 && channels[0] == "") {
		klog.Fatalf("need at least one channel in metric-channels parameter")
	}

	init := plant.YasdiInitializer{
		YasdiFile: *yasdiConfig,
		DriverID:  *yasdiDriver,
		Devices:   *yasdiDevices,
	}

	plant.RegisterMetrics()

	if err := plant.Run(&init, channels); err != nil {
		klog.Fatalf("error starting metrics loop: %v", err)
	}
	klog.Info("metrics loop  was started")

	klog.Infof("running metrics handler, metrics will be available at %s/metrics", *metricsAddr)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(*metricsAddr, nil); err != nil {
		klog.Fatalf("error starting metrics handler:  %v", err)
	}
}
