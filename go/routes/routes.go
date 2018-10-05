/*
Package routes is where the routes are Set, also has a custom Router
*/
package routes

import (
	"fmt"
	"net/http"
)

func (router *Router) setRoutes() {
	router.get("/ping", router.getPing)
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
