/*
Package routes is where the routes are Set, also has a custom Router
*/
package routes

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/s12chung/hello-k8/go/models"
	"net/http"
)

// setRoutes defines the routes for the router
//
// I changed the URL patterns from the docs because I wanted the API to look REST-ful
// "/nodes/{nodeName}/metrics" says create a metric (the last word) for the node with nodeName
// "/nodes/{nodeName}/metrics/average" says give me the average of the metric, of the node with nodeName
//
// also removed the v1 to make it cleaner.
func (router *Router) setRoutes() {
	router.mux.Get("/ping", router.getPing)
	router.mux.Post("/nodes/{nodeName}/metrics", router.postNodeMetric)
	router.mux.Get("/nodes/{nodeName}/metrics/average", router.getNodeMetricsAverage)
}

type metricRequestResponse struct {
	CPUUsed int `json:"cpu_used"`
	MemUsed int `json:"mem_used"`
}

// postNodeMetric creates a new Metric.
// Originally planned to create 2 metrics. One metric for each of the two ends of the timeslice,
// and removing any metrics within the timeslice. That way, the algorithm in getNodeMetricsAverage
// would still work.
func (router *Router) postNodeMetric(writer http.ResponseWriter, request *http.Request) {
	router.withTx(writer, func(tx *sql.Tx) bool {
		rMetric := &metricRequestResponse{}
		err := unmarkshallJSONBody(request.Body, rMetric)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return false
		}

		metric := models.Metric{
			Time:     router.c.Now(),
			NodeName: chi.URLParam(request, "nodeName"),
			CPUUsed:  rMetric.CPUUsed,
			MemUsed:  rMetric.MemUsed,
		}
		err = metric.Create(tx)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return false
		}

		writeJSON(writer, metric)
		return true
	})
}

// getNodeMetricsAverage returns the metricsAverage for all the metrics
// each metric is like a point on a line graph
// metricsAverage = areaOfLineGraph/totalTimeOfLineGraph
func (router *Router) getNodeMetricsAverage(writer http.ResponseWriter, request *http.Request) {
	metrics, err := models.AllMetrics(router.db)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var previousMetric *models.Metric
	var weightedCPUSum int64
	var weightedMemSum int64
	var totalSeconds int64
	// for every pair of metrics, in ASC metric.Time order
	for _, metric := range metrics {
		if previousMetric != nil {
			// find the time between the two points in seconds
			seconds := int64(models.RoundSecond(metric.Time).Sub(models.RoundSecond(previousMetric.Time)).Seconds())
			totalSeconds += seconds

			// seconds * (the average of the two metrics)
			weightedCPUSum += seconds * int64((metric.CPUUsed+previousMetric.CPUUsed)/2)
			weightedMemSum += seconds * int64((metric.MemUsed+previousMetric.MemUsed)/2)
		}
		previousMetric = metric
	}

	// areaOfLineGraph/totalTimeOfLineGraph
	writeJSON(writer, &metricRequestResponse{
		int(weightedCPUSum / totalSeconds),
		int(weightedMemSum / totalSeconds),
	})
}

func (router *Router) getPing(writer http.ResponseWriter, request *http.Request) {
	err := router.db.Ping()

	success := true
	if err != nil {
		success = false
	}
	_, err = writer.Write([]byte(fmt.Sprintf(`{ "success": %v }`, success)))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
