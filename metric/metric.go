package metric

import (
	"log"
	"time"
)

type Metric struct {
	Key          string     `json:"key"`
	Value        float32    `json:"value"`
	Timestamp    *time.Time `json:"-"`
	TimestampInt int64      `json:"timestamp"`
}

func (s *DefaultStore) AddMany(metrics []Metric) {
	for _, m := range metrics {
		s.Add(m)
	}
}

func (s *DefaultStore) Add(m Metric) {
	if m.TimestampInt == 0 {
		if m.Timestamp == nil || m.Timestamp.IsZero() {
			t := time.Now().UTC()
			m.Timestamp = &t
		}
		m.TimestampInt = m.Timestamp.Unix()
	}
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
                                         WHEN 'avg'       THEN avg(value)
                                         WHEN 'sum'       THEN sum(value)
                                         WHEN 'count'     THEN count(value)
                                         WHEN 'median'    THEN median(value)
                                         WHEN 'max'       THEN max(value)
                                         WHEN 'min'       THEN min(value)
                                         WHEN 'perc90'    THEN percentile_cont(array_agg(value), 0.90)
                                         WHEN 'perc95'    THEN percentile_cont(array_agg(value), 0.95)
                                         WHEN 'perc99'    THEN percentile_cont(array_agg(value), 0.99)
                                         ELSE avg(value)
                                         END
                                        ) AS value,
                                        round_timestamp(timestamp, $3) as timestamp
                               FROM metrics
                               WHERE key = $1
                               AND (now() - interval '1 hour') < timestamp
                               GROUP BY round_timestamp(timestamp, $3)
                               ORDER BY timestamp DESC`, key, function, resolution)
	if err != nil {
		log.Println(err)
		return nil
	}
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.Key, &m.Value, &m.Timestamp); err != nil {
			log.Fatal(err)
		}
		m.TimestampInt = m.Timestamp.Unix()
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
