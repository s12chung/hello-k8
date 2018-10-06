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

type metricRequest struct {
	Timeslice int     `json:"timeslice"`
	CPUUSed   float32 `json:"cpu_used"`
	MemUsed   float32 `json:"mem_used"`
}

func (router *Router) postNodeMetric(writer http.ResponseWriter, request *http.Request) {
	router.withTx(writer, func(tx *sql.Tx) bool {
		rMetric := &metricRequest{}
		err := unmarkshallJSONBody(request.Body, rMetric)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return false
		}

		metric := models.Metric{
			Time:        router.c.Now(),
			NodeName:    chi.URLParam(request, "nodeName"),
			ProcessName: "",
			CPUUsed:     rMetric.CPUUSed,
			MemUsed:     rMetric.MemUsed,
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
