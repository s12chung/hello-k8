package routes

import (
	"net/http"
	"testing"
)

func Test_getPing(t *testing.T) {
	testServer, clean := NewRoutedServer(t)
	defer clean()

	response, err := http.Get(testServer.URL + "/ping")
	if err != nil {
		t.Error(err)
	}

	got, err := StringBody(response)
	if err != nil {
		t.Error(err)
	}
	exp := `{ "success": true }`
	if got != exp {
		t.Errorf("got: %v, exp: %v\n", got, exp)
	}
}
