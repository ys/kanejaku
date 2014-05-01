package metric

import (
	"database/sql"
	"fmt"
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

var DB *sql.DB

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
	stmt, err := DB.Prepare(sStmt)
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
	rows, err := DB.Query(`SELECT MAX(key) AS key,
                                      AVG(value) AS value,
                                      (ROUND(timestamp / 30) * 30)::bigint as timestamp
                               FROM metrics
                               WHERE key = $1
                               AND extract(epoch from (now() - interval '1 hour')) < timestamp
                               GROUP BY timestamp
                               ORDER BY timestamp DESC`, key)
	if err != nil {
		log.Println(err)
		return nil
	}
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.Key, &m.Value, &m.Timestamp); err != nil {
			log.Fatal(err)
		}
		result = append(result, m)
	}
	return result
}

func InitDB() *sql.DB {
	url := os.Getenv("DATABASE_URL")
	var err error
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
