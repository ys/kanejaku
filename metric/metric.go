package metric

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

type Metric struct {
	Key       string
	Value     string
	Timestamp string
}

// type Metric struct {
// 	Key       string     `json:"key"`
// 	Value     float32    `json:"value"`
// 	Timestamp *time.Time `json:"timestamp"`
// }

var Db *sql.DB

func Add(key string, value float32, timestamp time.Time) {
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
	}
	sStmt := "insert into metrics(key, value, timestamp) values ($1, $2, $3)"
	stmt, err := Db.Prepare(sStmt)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(key, value, timestamp)
	if err != nil || res == nil {
		log.Fatal(err)
	}
}

func InitDB() {
	url := os.Getenv("DATABASE_URL")
	var err error
	Db, err = sql.Open("postgres", url)
	defer Db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
