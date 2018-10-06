package models

import (
	"database/sql"
	"fmt"
	"github.com/s12chung/hello-k8/go/database"
	"log"
	"time"
)

// Metric is a metric at a certain Time
type Metric struct {
	Time        time.Time `json:"time"`
	NodeName    string    `json:"node_name"`
	ProcessName string    `json:"process_name"`
	CPUUsed     float32   `json:"cpu_used"`
	MemUsed     float32   `json:"mem_used"`
}

// DeleteAllMetrics deletes all the metrics in the db (used for testing)
func DeleteAllMetrics(db *sql.DB) error {
	_, err := db.Exec(fmt.Sprintf(`TRUNCATE %v`, database.TableName("metrics")))
	return err
}

// AllMetrics returns all the metrics in the db (used for testing)
func AllMetrics(db *sql.DB) (metrics []*Metric, err error) {
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM %v;`, database.TableName("metrics")))
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := rows.Close()
		if err == nil {
			err = cerr
		}
	}()

	for rows.Next() {
		metric := &Metric{}
		err := rows.Scan(&metric.Time, &metric.NodeName, &metric.ProcessName, &metric.CPUUsed, &metric.MemUsed)
		if err != nil {
			log.Fatal(err)
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

// Create creates a metric in the DB
func (metric *Metric) Create(tx *sql.Tx) error {
	_, err := tx.Exec(metric.createString())
	return err
}

func (metric *Metric) createString() string {
	return fmt.Sprintf(
		`INSERT INTO %v VALUES ('%v', '%v', '%v', %v, %v);`,
		database.TableName("metrics"),
		metric.Time.Format(time.RFC3339),
		metric.NodeName,
		metric.ProcessName,
		metric.CPUUsed,
		metric.MemUsed,
	)
}

// RoundSecond returns the same time rounded by the second
func RoundSecond(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}
