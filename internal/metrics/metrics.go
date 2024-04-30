package metrics

import (
	"fmt"
	"log/slog"
	"net/http"

	"carbonaut.cloud/internal/connector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	// MetricsPort is the port on which the metrics server will listen
	MetricsPort int `yaml:"metricsPort" json:"metricsPort" validate:"required" default:"8080"`
	// CollectorName is the name of the prometheus collector
	CollectorName string `yaml:"collectorName" json:"collectorName" validate:"required" default:"carbonaut"`
}

func Serve(cfg Config, providers []connector.IProvider) {
	prometheus.Register(version.NewCollector(cfg.CollectorName))
	prometheus.Register(&Collector{
		providers: providers,
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{
			EnableOpenMetrics: false,
		})
		h.ServeHTTP(w, r)
	})
	prometheus.Unregister(collectors.NewGoCollector())

	slog.Info("Starting http server", slog.Int("metrics port", cfg.MetricsPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.MetricsPort), nil); err != nil {
		slog.Error("Failed to start http server", slog.String("error", err.Error()), slog.Int("port", cfg.MetricsPort))
	}
}

type Collector struct {
	providers []connector.IProvider
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for i := range c.providers {
		c.providers[i].Describe(ch)
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	for i := range c.providers {
		// TODO: check if this should be called in a goroutine
		c.providers[i].Collect(ch)
	}
}
