package routes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/s12chung/hello-k8/go/database"
	"github.com/s12chung/hello-k8/go/models"
)

func DefaultRouter(t *testing.T) *Router {
	db, err := database.DefaultDataBase()
	if err != nil {
		t.Error(err)
	}
	router := NewRouter(db)

	err = models.DeleteAllMetrics(db)
	if err != nil {
		t.Error(err)
	}

	router.c = testClock
	return router
}

func NewServer(router *Router) (*httptest.Server, func()) {
	testServer := httptest.NewServer(router.mux)
	return testServer, testServer.Close
}

func NewRoutedServer(t *testing.T) (*httptest.Server, func()) {
	router := DefaultRouter(t)
	router.setRoutes()
	return NewServer(router)
}

func StringBody(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), response.Body.Close()
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
