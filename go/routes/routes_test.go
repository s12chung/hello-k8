package routes

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/s12chung/gostatic/go/test"

	"github.com/s12chung/hello-k8/go/models"
)

func testUnmarkshallJSONBody(t *testing.T, body io.ReadCloser, v interface{}) (err error) {
	return _unmarkshallJSONBody(body, v, func(b []byte) {
		t.Log(string(b))
	})
}

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

func Test_postNodeMetric(t *testing.T) {
	testServer, clean := NewRoutedServer(t)
	defer clean()

	var err error
	reqMetric := &metricRequest{
		10,
		20,
		30,
	}

	body, err := marshallJSON(reqMetric)
	if err != nil {
		t.Error(err)
	}
	var response *http.Response
	response, err = http.Post(testServer.URL+"/nodes/blah/metrics", jsonContentType, bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}
	test.AssertLabel(t, "response.StatusCode", response.StatusCode, http.StatusOK)

	gotMetric := &models.Metric{}
	err = testUnmarkshallJSONBody(t, response.Body, gotMetric)
	if err != nil {
		t.Error(err)
	}

	expMetric := &models.Metric{
		Time:     models.RoundSecond(testClock.Now()),
		NodeName: "blah",
		CPUUsed:  reqMetric.CPUUSed,
		MemUsed:  reqMetric.MemUsed,
	}

	if cmp.Equal(gotMetric, expMetric) {
		t.Error(test.AssertLabelString("metric", gotMetric, expMetric))
	}
}
