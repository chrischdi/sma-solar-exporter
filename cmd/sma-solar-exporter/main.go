package main

import (
	"flag"
	"net/http"
	"os"
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

	influxUrl    = flag.String("influx-url", "", "URL of InfluxDB server")
	influxToken  = flag.String("influx-token", "", "Token (influx 2.x) or username:password (influx 1.8) to authenticate with. INFLUX_TOKEN environment variable will be used if present instead of this parameter.")
	influxOrg    = flag.String("influx-org", "", "Organization name (optional for influx 1.8)")
	influxBucket = flag.String("influx-bucket", "", "Bucket name (influx 2.x) or database/retention-policy (influx 1.8) or database (influx 1.8)")
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

	var influxConfig *plant.InfluxClient
	if *influxUrl != "" {
		influxEnvToken := os.Getenv("INFLUX_TOKEN")
		if influxEnvToken != "" {
			influxToken = &influxEnvToken
		}
		influxConfig = &plant.InfluxClient{
			Url:    *influxUrl,
			Token:  *influxToken,
			Org:    *influxOrg,
			Bucket: *influxBucket,
		}
	}

	if err := plant.Run(&init, channels, influxConfig); err != nil {
		klog.Fatalf("error starting metrics loop: %v", err)
	}
	klog.Info("metrics loop  was started")

	klog.Infof("running metrics handler, metrics will be available at %s/metrics", *metricsAddr)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(*metricsAddr, nil); err != nil {
		klog.Fatalf("error starting metrics handler:  %v", err)
	}
}
