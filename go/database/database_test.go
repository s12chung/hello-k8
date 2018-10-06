package database

import (
	"github.com/s12chung/gostatic/go/test"
	"os"
	"testing"
)

func SetEnv(t *testing.T, key, value string) string {
	old := os.Getenv(key)
	err := os.Setenv(key, value)
	if err != nil {
		t.Error(err)
	}
	return old
}

func TestDefaultDataBaseConfig(t *testing.T) {
	// docker env has ENV set
	oldUser := SetEnv(t, "POSTGRES_USER", "")
	oldDb := SetEnv(t, "POSTGRES_DB", "")
	oldPassword := SetEnv(t, "POSTGRES_PASSWORD", "")
	oldHost := SetEnv(t, "POSTGRES_SERVICE_HOST", "")
	oldPort := SetEnv(t, "POSTGRES_SERVICE_PORT", "")

	got := DefaultDataBaseConfig()

	exp := "user=postgres password= dbname=postgres host= port= sslmode=disable"
	if got != exp {
		t.Errorf("got: %v, exp: %v\n", got, exp)
	}

	SetEnv(t, "POSTGRES_USER", "the_user")
	got = DefaultDataBaseConfig()
	exp = "user=the_user password= dbname=the_user host= port= sslmode=disable"
	if got != exp {
		t.Errorf("got: %v, exp: %v\n", got, exp)
	}

	SetEnv(t, "POSTGRES_USER", "the_user")
	SetEnv(t, "POSTGRES_DB", "the_db")
	SetEnv(t, "POSTGRES_PASSWORD", "pw")
	SetEnv(t, "POSTGRES_SERVICE_HOST", "192.168.0.1")
	SetEnv(t, "POSTGRES_SERVICE_PORT", "5000")

	got = DefaultDataBaseConfig()
	exp = "user=the_user password=pw dbname=the_db host=192.168.0.1 port=5000 sslmode=disable"
	if got != exp {
		t.Errorf("got: %v, exp: %v\n", got, exp)
	}

	SetEnv(t, "POSTGRES_USER", oldUser)
	SetEnv(t, "POSTGRES_DB", oldDb)
	SetEnv(t, "POSTGRES_PASSWORD", oldPassword)
	SetEnv(t, "POSTGRES_SERVICE_HOST", oldHost)
	SetEnv(t, "POSTGRES_SERVICE_PORT", oldPort)
}

func TestTableName(t *testing.T) {
	test.AssertLabel(t, "result", TableName("blah"), "blah_test")
}
