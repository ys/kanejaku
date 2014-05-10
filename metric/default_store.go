package metric

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/ys/influxdb-go"
	"log"
	"net/url"
	"os"
)

type Store interface {
	AddMany(metrics []Metric) []Metric
	Add(m Metric) Metric
	Get(key string, function string, resolution int) []Metric
}

type DefaultStore struct {
	DB *influxdb.Client
}

func NewDefaultStore() *DefaultStore {
	store := &DefaultStore{}
	store.InitClient()
	return store
}

func (s *DefaultStore) InitClient() {
	influxUrl := os.Getenv("DATABASE_URL")
	u, err := url.Parse(influxUrl)
	if err != nil {
		log.Fatal(err)
	}
	var username string
	var password string
	if u.User != nil {
		username = u.User.Username()
		var is_set bool
		password, is_set = u.User.Password()
		if !is_set {
			password = ""
		}
	} else {
		username = ""
		password = ""
	}
	config := &influxdb.ClientConfig{
		Host:     u.Host,
		Database: u.Path[1:],
		Username: username,
		Password: password,
	}
	s.DB, err = influxdb.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *DefaultStore) Close() {
}
