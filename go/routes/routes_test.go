package routes

import (
	"bytes"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/s12chung/gostatic/go/test"

	"github.com/s12chung/hello-k8/go/models"
)

func testMarshallJSON(t *testing.T, v interface{}) []byte {
	b, err := marshallJSON(v)
	if err != nil {
		t.Error(err)
	}
	return b
}

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
	routedServer := func(t *testing.T) (*httptest.Server, func()) {
		testServer, _, clean := NewRoutedServer(t)
		return testServer, clean
	}

	testCases := []struct {
		serverFunc func(t *testing.T) (*httptest.Server, func())
		success    bool
	}{
		{routedServer, true},
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
	testServer, _, clean := NewRoutedServer(t)
	defer clean()

	var err error
	reqMetric := &metricRequestResponse{
		20,
		30,
	}

	var response *http.Response
	response, err = http.Post(testServer.URL+"/nodes/blah/metrics", jsonContentType, bytes.NewBuffer(testMarshallJSON(t, reqMetric)))
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
		CPUUsed:  reqMetric.CPUUsed,
		MemUsed:  reqMetric.MemUsed,
	}

	if cmp.Equal(gotMetric, expMetric) {
		t.Error(test.AssertLabelString("metric", gotMetric, expMetric))
	}
}

func Test_postNodeMetricBadReqBody(t *testing.T) {
	testServer, _, clean := NewRoutedServer(t)
	defer clean()

	response, err := http.Post(testServer.URL+"/nodes/blah/metrics", jsonContentType, bytes.NewBuffer([]byte(`not_json`)))
	if err != nil {
		t.Error(err)
	}
	test.AssertLabel(t, "response.StatusCode", response.StatusCode, http.StatusBadRequest)
}

func Test_getNodeMetricsAverage(t *testing.T) {
	testServer, router, clean := NewRoutedServer(t)
	defer clean()

	ts := []time.Duration{0, 10, 30, 60}
	useds := []int{0, 20, 60, 0}

	var metrics []*models.Metric
	for i := 0; i < len(ts); i++ {
		metrics = append(metrics, &models.Metric{
			Time:     testClock.Now().Add(time.Second * ts[i]),
			NodeName: "blah",
			CPUUsed:  useds[i],
			MemUsed:  100 - useds[i],
		})
	}

	err := models.CreateMetrics(router.db, metrics)
	if err != nil {
		t.Error(err)
	}

	var response *http.Response
	response, err = http.Get(testServer.URL + "/nodes/blah/metrics/average")
	if err != nil {
		t.Error(err)
	}
	test.AssertLabel(t, "response.StatusCode", response.StatusCode, http.StatusOK)

	got := &metricRequestResponse{}
	err = testUnmarkshallJSONBody(t, response.Body, got)
	if err != nil {
		t.Error(err)
	}
	exp := &metricRequestResponse{
		30,
		70,
	}

	if !cmp.Equal(got, exp) {
		t.Error(test.AssertLabelString("metricRequestResponse", got, exp))
	}
}
