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
	Key       string  `json:"key",db:"key"`
	Value     float32 `json:"value",db:"value"`
	Timestamp int64   `json:"timestamp",db:"timestamp"`
}

// type Metric struct {
// 	Key       string     `json:"key"`
// 	Value     float32    `json:"value"`
// 	Timestamp *time.Time `json:"timestamp"`
// }

var Db *sql.DB
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

func InitDB() {
	url := os.Getenv("DATABASE_URL")
	var err error
	Db, err = sql.Open("postgres", url)
	defer Db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
