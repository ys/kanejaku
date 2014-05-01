package server

import (
	"github.com/daneharrigan/bourbon"
	"github.com/ys/kanejaku/metric"
)

func Run() {
	db := metric.InitDB()
	defer db.Close()
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
	b.Get("/metrics/{key}", func(params bourbon.Params) (int, bourbon.Encodeable) {
		return 200, metric.Get(params["key"])
	})
	return b
}
