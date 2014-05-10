package metric

import (
	"github.com/ys/influxdb-go"
	"log"
	"strconv"
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
	series := []*influxdb.Series{
		&influxdb.Series{
			Name:    m.Key,
			Columns: []string{"time", "value"},
			Points: [][]interface{}{
				[]interface{}{m.Timestamp.Unix(), m.Value},
			},
		},
	}
	DB.WriteSeries(series)
	return m
}

func (s *DefaultStore) Get(key string, function string, resolution int) []Metric {
	result := []Metric{}
	if function == "" {
		function = "MEAN"
	}
	if resolution == 0 {
		resolution = 30
	}
	// moment := time.Now().Add(-1 * time.Hour).Format("2006-01-02 03:04:00")
	query := "SELECT " + toFunction(function) + " FROM " + key + " group by time(" + strconv.Itoa(resolution) + "s)"
	series, err := s.DB.Query(query)
	if err != nil {
		log.Println(err)
		return result
	}
	for _, serie := range series {
		for _, point := range serie.Points {
			timestamp := time.Unix(int64(point[0].(float64)), 0).UTC()
			m := Metric{
				Key:          serie.Name,
				Value:        float32(point[1].(float64)),
				TimestampInt: int64(point[0].(float64)),
				Timestamp:    &timestamp,
			}
			result = append(result, m)
		}
	}
	return result
}

func toFunction(function string) string {
	var queryFunction string
	switch function {
	case "mean":
		queryFunction = "mean(value)"
	case "sum":
		queryFunction = "sum(value)"
	case "count":
		queryFunction = "count(value)"
	case "median":
		queryFunction = "median(value)"
	case "max":
		queryFunction = "max(value)"
	case "min":
		queryFunction = "min(value)"
	case "perc90":
		queryFunction = "percentile(value, 90)"
	case "perc95":
		queryFunction = "percentile(value, 95)"
	case "perc99":
		queryFunction = "percentile(value, 99)"
	default:
		queryFunction = "mean(value)"
	}

	return queryFunction
}
