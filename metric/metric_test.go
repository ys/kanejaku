package metric

import (
	"github.com/joho/godotenv"
	"github.com/ys/influxdb-go"
	"log"
	"os"
	"testing"
	"time"
)

var S *DefaultStore
var DB *influxdb.Client

func TestGet(t *testing.T) {
	setup(t)
	timestamp := time.Now().UTC().Round(1 * time.Millisecond)
	insertMetric("key", 1, timestamp)
	metrics := S.Get("key", "", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m := metrics[0]
	if m.Key != "key" || m.Value != 1 || !m.Timestamp.UTC().Equal(timestamp.Round(30*time.Second)) {
		t.Errorf("Metric not correct: '%v', %v", m, timestamp.Round(30*time.Second))
	}
	teardown(t)
}

func TestGetAggregated(t *testing.T) {
	setup(t)
	timestamp := time.Now().Round(1 * time.Minute)
	insertMetric("key", 1, timestamp)
	// Rounding first is going to the minute and this to the 30s when rounding by 30
	insertMetric("key", 1, timestamp.Add(25*time.Second))
	metrics := S.Get("key", "", 0)
	if len(metrics) != 2 {
		t.Error("Error was expecting two rows")
	}
	// both are rounded to the minute
	metrics = S.Get("key", "", 60)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row, it should have been aggregated")
	}
	teardown(t)
}

func TestGetAverage(t *testing.T) {
	setup(t)
	timestamp := time.Now()
	insertMetric("key", 1, timestamp)
	insertMetric("key", 1, timestamp)
	metrics := S.Get("key", "avg", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m := metrics[0]
	if m.Value != 1 {
		t.Errorf("Metric not correct: '%v'", m)
	}
	teardown(t)
}

func TestGetSum(t *testing.T) {
	setup(t)
	timestamp := time.Now()
	insertMetric("key", 1, timestamp)
	insertMetric("key", 1, timestamp)
	metrics := S.Get("key", "sum", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m := metrics[0]
	if m.Value != 2 {
		t.Errorf("Metric not correct: '%v'", m)
	}
	teardown(t)
}

func TestGetCount(t *testing.T) {
	setup(t)
	timestamp := time.Now()
	insertMetric("key", 10, timestamp)
	insertMetric("key", 1, timestamp)
	metrics := S.Get("key", "count", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m := metrics[0]
	if m.Value != 2 {
		t.Errorf("Metric not correct: '%v'", m)
	}
	teardown(t)
}

func TestGetMax(t *testing.T) {
	setup(t)
	timestamp := time.Now()
	insertMetric("key", 10, timestamp)
	insertMetric("key", 1, timestamp)
	metrics := S.Get("key", "max", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m := metrics[0]
	if m.Value != 10 {
		t.Errorf("Metric not correct: '%v'", m)
	}
	teardown(t)
}

func TestGetMin(t *testing.T) {
	setup(t)
	timestamp := time.Now()
	insertMetric("key", 10, timestamp)
	insertMetric("key", 1, timestamp)
	metrics := S.Get("key", "min", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m := metrics[0]
	if m.Value != 1 {
		t.Errorf("Metric not correct: '%v'", m)
	}
	teardown(t)
}

func TestGetPercAndMedian(t *testing.T) {
	setup(t)
	timestamp := time.Now()
	for i := 1; i <= 100; i = i + 1 {
		insertMetric("key", float32(i), timestamp)
	}
	metrics := S.Get("key", "median", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m := metrics[0]
	if int(m.Value) != 50 {
		t.Error("Metric should have a 90 value")
	}
	metrics = S.Get("key", "perc90", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m = metrics[0]
	if int(m.Value) != 90 {
		t.Error("Metric should have a 90 value")
	}
	metrics = S.Get("key", "perc95", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m = metrics[0]
	if int(m.Value) != 95 {
		t.Error("Metric should have a 95 value")
	}
	metrics = S.Get("key", "perc99", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m = metrics[0]
	if int(m.Value) != 99 {
		t.Error("Metric should have a 99 value")
	}
	teardown(t)
}

func setup(t *testing.T) {
	env := ".env.test"
	if os.Getenv("TRAVIS") == "true" {
		env = ".env.travis"
	}
	err := godotenv.Load(env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	S = &DefaultStore{}
	S.InitClient()
	DB = S.DB
	log.Println(DB.DeleteSerie("key"))
	DB.DeleteSerie("key1")
}

func teardown(t *testing.T) {
	DB.DeleteSerie("key")
	DB.DeleteSerie("key1")
}

func insertMetric(key string, value float32, timestamp time.Time) {
	series := []*influxdb.Series{
		&influxdb.Series{
			Name:    key,
			Columns: []string{"time", "value"},
			Points: [][]interface{}{
				[]interface{}{timestamp.Unix(), value},
			},
		},
	}
	DB.WriteSeries(series)
}
