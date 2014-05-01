package metric

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

type Metric struct {
	Key       string  `json:"key",db:"key"`
	Value     float32 `json:"value",db:"value"`
	Timestamp int64   `json:"timestamp",db:"timestamp"`
}

var Db *sqlx.DB

func AddMany(metrics []Metric) {
	for _, m := range metrics {
		Add(m)
	}
}

func Add(m Metric) {
	if m.Timestamp == 0 {
		m.Timestamp = time.Now().UTC().Unix()
	}
	fmt.Println(time.Unix(m.Timestamp, 0))
	sStmt := "insert into metrics(key, value, timestamp) values ($1, $2, $3)"
	stmt, err := Db.Prepare(sStmt)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(m.Key, m.Value, m.Timestamp)
	if err != nil || res == nil {
		log.Fatal(err)
	}
}

func Get(key string) []Metric {
	result := []Metric{}
	err := Db.Select(&result, "SELECT MAX(key) AS key, avg(value) AS value, (ROUND(timestamp / 30) * 30)::bigint as timestamp FROM metrics WHERE key = $1 GROUP BY timestamp ORDER BY timestamp DESC", key)
	if err != nil {
		log.Println(err)
		return nil
	}
	return result
}

func InitDB() *sqlx.DB {
	url := os.Getenv("DATABASE_URL")
	var err error
	Db, err = sqlx.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}
	return Db
}
