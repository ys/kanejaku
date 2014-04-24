package server

import (
	"github.com/daneharrigan/bourbon"
	"github.com/ys/kanejaku/metric"
)

func Run() {
	metric.InitDB()
	bourbon.Run(metrics())

}

func metrics() bourbon.Bourbon {
	b := bourbon.New()
	b.Post("/metrics", func(m metric.Metric) (int, bourbon.Encodeable) {
		return 201, m
	})
	return b
}
