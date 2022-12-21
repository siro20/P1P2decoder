package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/siro20/p1p2decoder/pkg/p1p2"
)

type PrometheusConfig struct {
	Enable        bool   `yaml:"enable"`
	DefaultPort   int    `yaml:"port"`
	ListenAddress string `yaml:"listen_address"`
}

func runPrometheusServer(cfg PrometheusConfig) {
	if !cfg.Enable {
		return
	}
	for i := range p1p2.Sys.Temperatures {
		gauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: "P1P2",
			Subsystem: "Temperature",
			Name:      p1p2.Sys.Temperatures[i].Name(),
			Help:      p1p2.Sys.Temperatures[i].Description(),
		}, func() float64 {
			f, ok := p1p2.Sys.Temperatures[i].Value().(float32)
			if !ok {
				return math.NaN()
			}
			return float64(f)
		})

		prometheus.MustRegister(gauge)
	}
	for i := range p1p2.Sys.Status {
		gauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: "P1P2",
			Subsystem: "State",
			Name:      p1p2.Sys.Status[i].Name(),
			Help:      p1p2.Sys.Status[i].Description(),
		}, func() float64 {
			f, ok := p1p2.Sys.Status[i].Value().(bool)
			if !ok || !f {
				return math.NaN()
			}
			return 1.0
		})

		prometheus.MustRegister(gauge)
	}
	for i := range p1p2.Sys.Pumps {
		gauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: "P1P2",
			Subsystem: "Pumps",
			Name:      p1p2.Sys.Pumps[i].Name(),
			Help:      p1p2.Sys.Pumps[i].Description(),
		}, func() float64 {
			f, ok := p1p2.Sys.Pumps[i].Value().(bool)
			if !ok || !f {
				return math.NaN()
			}
			return 1.0
		})

		prometheus.MustRegister(gauge)
	}
	for i := range p1p2.Sys.Valves {
		gauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: "P1P2",
			Subsystem: "Valves",
			Name:      p1p2.Sys.Valves[i].Name(),
			Help:      p1p2.Sys.Valves[i].Description(),
		}, func() float64 {
			f, ok := p1p2.Sys.Valves[i].Value().(bool)
			if !ok || !f {
				return math.NaN()
			}
			return 1.0
		})

		prometheus.MustRegister(gauge)
	}

	if cfg.DefaultPort == 0 {
		cfg.DefaultPort = 2112
	}
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.ListenAddress, cfg.DefaultPort), nil)
}
