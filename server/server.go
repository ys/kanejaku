package server

import (
	"github.com/daneharrigan/bourbon"
	"github.com/ys/kanejaku/metric"
)

func Run() {
	metric.InitDB()
	metrics().Run()
}

func metrics() bourbon.Bourbon {
	b := bourbon.New()
	b.Get("/metrics", func() (int, bourbon.Encodeable) {
		return 418, "TEAPOT"
	})
	b.Post("/metrics", func(m metric.Metric) (int, bourbon.Encodeable) {
		println("LOL")
		return 201, "tto"
	})
	return b
}
