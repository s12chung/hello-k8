package routes

import "net/http"

func (server *Server) setRoutes() {
	server.get("/", getHome)
}

func getHome(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(`{ "cpu_used": 100 }`))
}
