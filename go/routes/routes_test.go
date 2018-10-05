package routes

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
)

func brokenDbServer(t *testing.T) (*httptest.Server, func()) {
	db, err := sql.Open("postgres", "user=the_user password=pw dbname=the_db host=192.168.0.1 port=5000 sslmode=disable")
	if err != nil {
		t.Error(err)
	}
	router := NewRouter(db)
	router.setRoutes()
	return NewServer(router)
}

func Test_getPing(t *testing.T) {
	testCases := []struct {
		serverFunc func(t *testing.T) (*httptest.Server, func())
		success    bool
	}{
		{NewRoutedServer, true},
		{brokenDbServer, false},
	}

	for _, tc := range testCases {
		func() {
			testServer, clean := tc.serverFunc(t)
			defer clean()

			response, err := http.Get(testServer.URL + "/ping")
			if err != nil {
				t.Error(err)
			}

			got, err := StringBody(response)
			if err != nil {
				t.Error(err)
			}
			exp := `{ "success": false }`
			if tc.success {
				exp = `{ "success": true }`
			}

			if got != exp {
				t.Errorf("got: %v, exp: %v\n", got, exp)
			}
		}()
	}
}
