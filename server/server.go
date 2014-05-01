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
	b.Post("/metrics", func(metrics []metric.Metric) (int, bourbon.Encodeable) {
		metric.AddMany(metrics)
		return 201, metrics
	})
	})
	return b
}
