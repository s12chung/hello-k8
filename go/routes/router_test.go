package routes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/s12chung/hello-k8/go/database"
)

func DefaultRouter(t *testing.T) *Router {
	db, err := database.DefaultDataBase()
	if err != nil {
		t.Error(err)
	}
	return NewRouter(db)
}

func NewServer(router *Router) (*httptest.Server, func()) {
	testServer := httptest.NewServer(router.mux)
	return testServer, testServer.Close
}

func NewRoutedServer(t *testing.T) (*httptest.Server, func()) {
	server := DefaultRouter(t)
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
	router := DefaultRouter(t)
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

func Test_setDefaultRoutes(t *testing.T) {
	testServer, clean := NewRoutedServer(t)
	defer clean()

	response, err := http.Get(testServer.URL)
	if err != nil {
		t.Error(err)
	}

	got := response.StatusCode
	if err != nil {
		t.Error(err)
	}
	exp := http.StatusNotFound
	if got != exp {
		t.Errorf("got: %v, exp: %v\n", got, exp)
	}
}
