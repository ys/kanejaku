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

func (s *DefaultStore) AddMany(metrics []Metric) []Metric {
	results := []Metric{}
	for _, m := range metrics {
		m := s.Add(m)
		results = append(results, m)
	}
	return results
}

func (s *DefaultStore) Add(m Metric) Metric {
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
		log.Println(err)
		return Metric{}
	}
	return m
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
                                        `+toSQLFunction(function)+`AS value,
                                        round_timestamp(timestamp, $2) as timestamp
                               FROM metrics
                               WHERE key = $1
                               AND (now() - interval '1 hour') < timestamp
                               GROUP BY round_timestamp(timestamp, $2)
                               ORDER BY timestamp DESC`, key, resolution)
	if err != nil {
		log.Println(err)
		return result
	}
	for rows.Next() {
		var m Metric
		if err := rows.Scan(&m.Key, &m.Value, &m.Timestamp); err != nil {
			log.Println(err)
			return result
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
		return keys
	}
	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			log.Println(err)
			return keys
		}
		keys = append(keys, key)
	}
	return keys
}

func toSQLFunction(function string) string {
	var sqlFunction string
	switch function {
	case "avg":
		sqlFunction = "avg(value)"
	case "sum":
		sqlFunction = "sum(value)"
	case "count":
		sqlFunction = "count(value)"
	case "median":
		sqlFunction = "median(value)"
	case "max":
		sqlFunction = "max(value)"
	case "min":
		sqlFunction = "min(value)"
	case "perc90":
		sqlFunction = "percentile_cont(array_agg(value), 0.90)"
	case "perc95":
		sqlFunction = "percentile_cont(array_agg(value), 0.95)"
	case "perc99":
		sqlFunction = "percentile_cont(array_agg(value), 0.99)"
	default:
		sqlFunction = "avg(value)"
	}

	return sqlFunction
}
