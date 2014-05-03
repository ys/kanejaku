package server

import (
	"github.com/daneharrigan/bourbon"
	"github.com/ys/kanejaku/metric"
	"net/http"
	"strconv"
)

type Server interface {
	Run()
}

type DefaultServer struct {
	Store metric.Store
}

func (s *DefaultServer) Run() {
	s.metrics().Run()
}

func (s *DefaultServer) metrics() bourbon.Bourbon {
	store := s.Store
	b := bourbon.New()
	b.Get("/metrics", func() (int, bourbon.Encodeable) {
		return 200, store.GetKeys()
	})
	b.Post("/metrics", func(metrics []metric.Metric) (int, bourbon.Encodeable) {
		store.AddMany(metrics)
		return 201, metrics
	})
	b.Get("/metrics/{key}", func(req *http.Request, params bourbon.Params) (int, bourbon.Encodeable) {
		queryParams := req.URL.Query()
		resolution, _ := strconv.Atoi(queryParams.Get("resolution"))
		return 200, store.Get(params["key"], queryParams.Get("function"), resolution)
	})
	return b
}
