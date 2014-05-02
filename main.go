package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/ys/kanejaku/metric"
	"github.com/ys/kanejaku/server"
)

func main() {
	store := metric.NewDefaultStore()
	defer store.Close()
	s := &server.DefaultServer{Store: store}
	s.Run()
}
