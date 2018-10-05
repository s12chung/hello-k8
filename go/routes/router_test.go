package routes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func DefaultRouter() *Router {
	return NewRouter()
}

func NewServer(router *Router) (*httptest.Server, func()) {
	testServer := httptest.NewServer(router.mux)
	return testServer, testServer.Close
}

func NewRoutedServer() (*httptest.Server, func()) {
	server := DefaultRouter()
	server.setRoutes()
	return NewServer(server)
}

func StringBody(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), response.Body.Close()
}

func TestRouter_get(t *testing.T) {
	router := DefaultRouter()
	responseBody := `{ "cpu_used": 100 }`

	router.get("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte(responseBody))
		if err != nil {
			t.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})

	testServer, clean := NewServer(router)
	defer clean()

	response, err := http.Get(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	got := response.Header.Get("Content-Type")
	exp := "application/json"
	if got != exp {
		t.Errorf("got: %v, exp: %v\n", got, exp)
	}

	got, err = StringBody(response)
	if err != nil {
		t.Error(err)
	}
	exp = responseBody
	if got != exp {
		t.Errorf("got: %v, exp: %v\n", got, exp)
	}
}
