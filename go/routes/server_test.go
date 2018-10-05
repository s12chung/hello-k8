package routes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func DefaultServer() *Server {
	return NewServer()
}

func NewTestServer(server *Server) (*httptest.Server, func()) {
	testServer := httptest.NewServer(server.mux)
	return testServer, testServer.Close
}

func NewRoutedTestServer() (*httptest.Server, func()) {
	server := DefaultServer()
	server.setRoutes()
	return NewTestServer(server)
}

func StringBody(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), response.Body.Close()
}

func TestServer_get(t *testing.T) {
	server := DefaultServer()
	responseBody := `{ "cpu_used": 100 }`

	server.get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(responseBody))
	})

	testServer, clean := NewTestServer(server)
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
