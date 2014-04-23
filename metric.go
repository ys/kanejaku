package kanejaku

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

var db *sql.DB

func AddMetric(key string, value float32, timestamp time.Time) {
	if timestamp.IsZero() {
		timestamp = time.Now().UTC()
	}
	sStmt := "insert into metrics(key, value, timestamp) values ($1, $2, $3)"
	stmt, err := db.Prepare(sStmt)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(key, value, timestamp)
	if err != nil || res == nil {
		log.Fatal(err)
	}
}

func main() {
	url := os.Getenv("DATABASE_URL")
	var err error
	db, err = sql.Open("postgres", url)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	AddMetric("cool.gang", 23.0, time.Time{})
}
