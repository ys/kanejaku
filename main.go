package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/ys/kanejaku/metric"
	"github.com/ys/kanejaku/server"
)

func main() {
	metric.DB = metric.InitDB()
	defer metric.DB.Close()
	server.Run()
}
