package main

import (
	"math"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/siro20/p1p2decoder/pkg/p1p2"
)

func runPrometheusServer(sys p1p2.System) {
	for i := range sys.Temperatures {
		gauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: "P1P2",
			Subsystem: "Temperature",
			Name:      sys.Temperatures[i].Name(),
			Help:      sys.Temperatures[i].Description(),
		}, func() float64 {
			f, ok := sys.Temperatures[i].Value().(float32)
			if !ok {
				return math.NaN()
			}
			return float64(f)
		})

		prometheus.MustRegister(gauge)
	}
	for i := range sys.Status {
		gauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: "P1P2",
			Subsystem: "State",
			Name:      sys.Status[i].Name(),
			Help:      sys.Status[i].Description(),
		}, func() float64 {
			f, ok := sys.Status[i].Value().(bool)
			if !ok || !f {
				return math.NaN()
			}
			return 1.0
		})

		prometheus.MustRegister(gauge)
	}
	for i := range sys.Pumps {
		gauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: "P1P2",
			Subsystem: "Pumps",
			Name:      sys.Pumps[i].Name(),
			Help:      sys.Pumps[i].Description(),
		}, func() float64 {
			f, ok := sys.Pumps[i].Value().(bool)
			if !ok || !f {
				return math.NaN()
			}
			return 1.0
		})

		prometheus.MustRegister(gauge)
	}
	for i := range sys.Valves {
		gauge := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: "P1P2",
			Subsystem: "Valves",
			Name:      sys.Valves[i].Name(),
			Help:      sys.Valves[i].Description(),
		}, func() float64 {
			f, ok := sys.Valves[i].Value().(bool)
			if !ok || !f {
				return math.NaN()
			}
			return 1.0
		})

		prometheus.MustRegister(gauge)
	}

	defaultPort := ":2112"
	http.Handle("/metrics", promhttp.Handler())
	if *prometheusAddr == "" {
		prometheusAddr = &defaultPort
	}
	http.ListenAndServe(*prometheusAddr, nil)
}
