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
		t.Errorf("Error '%s' was not expected while closing the database", err)
	}
	if m.Key != "key" || m.Value != 1 || m.Timestamp != 1398978882 {
		t.Error("Metric not correct")
	}
	teardown(t)
}

func TestGet(t *testing.T) {
	setup(t)
	timestamp := time.Now().Unix()
	insertMetric("key", 1, timestamp)
	metrics := S.Get("key")
	if len(metrics) != 1 {
		t.Error("Error was expecting one row")
	}
	m := metrics[0]
	if m.Key != "key" || m.Value != 1 || m.Timestamp != int64(math.Floor(float64(timestamp/30))*30) {
		t.Error("Metric not correct")
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
