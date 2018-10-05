package routes

import (
	"net/http"
	"testing"
)

func Test_getHome(t *testing.T) {
	testServer, clean := NewRoutedTestServer()
	defer clean()

	response, err := http.Get(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	got, err := StringBody(response)
	if err != nil {
		t.Error(err)
	}
	exp := `{ "cpu_used": 100 }`
	if got != exp {
		t.Errorf("got: %v, exp: %v\n", got, exp)
	}
}
