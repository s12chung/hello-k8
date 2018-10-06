package models

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/s12chung/gostatic/go/test"

	"github.com/s12chung/hello-k8/go/database"
)

func DefaultDb(t *testing.T) *sql.DB {
	db, err := database.DefaultDataBase()
	if err != nil {
		t.Error(err)
	}

	err = DeleteAllMetrics(db)
	if err != nil {
		t.Error(err)
	}

	return db
}

func DefaultTx(t *testing.T) (*sql.DB, *sql.Tx) {
	db := DefaultDb(t)
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
	}
	return db, tx
}

func TestMetric_Create(t *testing.T) {
	db, tx := DefaultTx(t)
	metric := &Metric{
		time.Now(),
		"my_node",
		"",
		10,
		20,
	}
	t.Log(metric.createString())
	err := metric.Create(tx)
	if err != nil {
		t.Error(err)
	}

	err = tx.Commit()
	if err != nil {
		t.Error(err)
	}

	metrics, err := AllMetrics(db)
	if err != nil {
		t.Error(err)
	}

	test.AssertLabel(t, "len", len(metrics), 1)

	gotMetric := metrics[0]
	metric.Time = RoundSecond(metric.Time)
	if !cmp.Equal(gotMetric, metric) {
		t.Error(test.AssertLabelString("metric", gotMetric, metric))
	}
}
