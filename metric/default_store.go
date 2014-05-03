package metric

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Store interface {
	AddMany(metrics []Metric)
	Add(m Metric)
	Get(key string, function string, resolution int) []Metric
	GetKeys() []string
}

type DefaultStore struct {
	DB *sql.DB
}

func NewDefaultStore() *DefaultStore {
	store := &DefaultStore{}
	store.InitDB()
	return store
}

func (s *DefaultStore) InitDB() {
	url := os.Getenv("DATABASE_URL")
	var err error
	s.DB, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *DefaultStore) Close() {
	s.DB.Close()
}
