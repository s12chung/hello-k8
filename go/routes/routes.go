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

func (router *Router) setRoutes() {
	router.mux.Get("/ping", router.getPing)
	router.mux.Post("/nodes/{nodeName}/metrics", router.postNodeMetric)
	router.mux.Get("/nodes/{nodeName}/metrics/average", router.getNodeMetricsAverage)
}

type metricRequestResponse struct {
	CPUUsed int `json:"cpu_used"`
	MemUsed int `json:"mem_used"`
}

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
	for _, metric := range metrics {
		if previousMetric != nil {
			seconds := int64(models.RoundSecond(metric.Time).Sub(models.RoundSecond(previousMetric.Time)).Seconds())
			totalSeconds += seconds
			weightedCPUSum += seconds * int64((metric.CPUUsed+previousMetric.CPUUsed)/2)
			weightedMemSum += seconds * int64((metric.MemUsed+previousMetric.MemUsed)/2)
		}
		previousMetric = metric
	}

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
