package metric

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"math"
	"os"
	"testing"
	"time"
)

var S *DefaultStore
var DB *sql.DB

func TestAdd(t *testing.T) {
	setup(t)
	S.Add(Metric{Key: "key", Value: 1, Timestamp: 1398978882})
	m := &Metric{}
	err := DB.QueryRow("SELECT key, value, timestamp FROM metrics LIMIT 1").Scan(&m.Key, &m.Value, &m.Timestamp)
	if err != nil {
		t.Errorf("Err '%s' while getting row", err)
	}
	if m.Key != "key" || m.Value != 1 || m.Timestamp != 1398978882 {
		t.Errorf("Metric not correct: '%v'", m)
	}
	teardown(t)
}

func TestGet(t *testing.T) {
	setup(t)
	timestamp := time.Now().Unix()
	insertMetric("key", 1, timestamp)
	metrics := S.Get("key", "", 0)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m := metrics[0]
	if m.Key != "key" || m.Value != 1 || m.Timestamp != int64(math.Floor(float64(timestamp/30))*30) {
		t.Errorf("Metric not correct: '%v'", m)
	}
	teardown(t)
}

func TestGetAggregated(t *testing.T) {
	setup(t)
	timestamp := time.Now().Unix()
	timestamp = int64(math.Floor(float64(timestamp/60)) * 60)
	insertMetric("key", 1, timestamp)
	insertMetric("key", 1, timestamp+40)
	metrics := S.Get("key", "", 0)
	if len(metrics) != 2 {
		t.Error("Error was expecting two rows")
	}
	metrics = S.Get("key", "", 60)
	if len(metrics) != 1 {
		t.Error("Error was expecting one row, it should have been aggregated")
	}
	teardown(t)
}

func TestGetAverage(t *testing.T) {
	setup(t)
	timestamp := time.Now().Unix()
	insertMetric("key", 1, timestamp)
	insertMetric("key", 1, timestamp+1)
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
	timestamp := time.Now().Unix()
	insertMetric("key", 1, timestamp)
	insertMetric("key", 1, timestamp+1)
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
	timestamp := time.Now().Unix()
	insertMetric("key", 10, timestamp)
	insertMetric("key", 1, timestamp+1)
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
	timestamp := time.Now().Unix()
	insertMetric("key", 10, timestamp)
	insertMetric("key", 1, timestamp+1)
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
	timestamp := time.Now().Unix()
	insertMetric("key", 10, timestamp)
	insertMetric("key", 1, timestamp+1)
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
	timestamp := time.Now().Unix()
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

func TestGetKeys(t *testing.T) {
	setup(t)
	timestamp := time.Now().Unix()
	insertMetric("key", 1, timestamp)
	insertMetric("key1", 1, timestamp)
	keys := S.GetKeys()
	if len(keys) != 2 {
		t.Error("Must have 2 keys")
	}
	if keys[0] != "key" || keys[1] != "key1" {
		t.Error("Keys are ordered")
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
	S.InitDB()
	DB = S.DB
	DB.Exec("TRUNCATE TABLE metrics")
}

func teardown(t *testing.T) {
	DB.Exec("TRUNCATE TABLE metrics")
	if err := DB.Close(); err != nil {
		t.Errorf("Error '%s' was not expected while closing the database", err)
	}
}

func insertMetric(key string, value float32, timestamp int64) {
	DB.Exec("INSERT INTO metrics(key, value, timestamp) VALUES ($1, $2, $3)", key, value, timestamp)
}
