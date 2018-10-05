/*
Package routes is where the routes are Set, also has a custom Router
*/
package routes

import "net/http"

func (router *Router) setRoutes() {
	router.get("/", getHome)
}

func getHome(writer http.ResponseWriter, request *http.Request) {
	_, err := writer.Write([]byte(`{ "cpu_used": 100 }`))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
