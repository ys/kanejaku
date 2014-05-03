package metric

import (
	"fmt"
	"log"
	"time"
)

type Metric struct {
	Key       string  `json:"key",db:"key"`
	Value     float32 `json:"value",db:"value"`
	Timestamp int64   `json:"timestamp",db:"timestamp"`
}

func (s *DefaultStore) AddMany(metrics []Metric) {
	for _, m := range metrics {
		s.Add(m)
	}
}

func (s *DefaultStore) Add(m Metric) {
	if m.Timestamp == 0 {
		m.Timestamp = time.Now().UTC().Unix()
	}
	fmt.Println(time.Unix(m.Timestamp, 0))
	sStmt := "insert into metrics(key, value, timestamp) values ($1, $2, $3)"
	stmt, err := s.DB.Prepare(sStmt)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(m.Key, m.Value, m.Timestamp)
	if err != nil || res == nil {
		log.Fatal(err)
	}
}

func (s *DefaultStore) Get(key string, function string, resolution int) []Metric {
	result := []Metric{}
	if function == "" {
		function = "avg"
	}
	if resolution == 0 {
		resolution = 30
	}
	rows, err := s.DB.Query(`SELECT MAX(key) AS key,
                                        (CASE $2
                                         WHEN 'avg'   THEN avg(value)
                                         WHEN 'sum'   THEN sum(value)
                                         WHEN 'count' THEN count(value)
                                         ELSE avg(value)
                                         END
                                        ) AS value,
                                        (ROUND(timestamp / $3) * $3)::bigint as timestamp
                               FROM metrics
                               WHERE key = $1
                               AND extract(epoch from (now() - interval '1 hour')) < timestamp
                               GROUP BY (ROUND(timestamp / $3) * $3)
                               ORDER BY (ROUND(timestamp / $3) * $3) DESC`, key, function, resolution)
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

func (s *DefaultStore) GetKeys() []string {
	keys := []string{}
	rows, err := s.DB.Query("SELECT DISTINCT key FROM metrics ORDER BY key ASC")
	if err != nil {
		log.Println(err)
		return nil
	}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			log.Fatal(err)
		}
		keys = append(keys, key)
	}
	return keys
}
