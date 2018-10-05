package routes

import (
	"log"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (server *Server) Run() error {
	server.setRoutes()

	port := "8080"
	log.Printf("Server listening on port %s\n", port)
	return http.ListenAndServe(":"+port, server.mux)
}

func (server *Server) get(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	server.mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if request.Method == http.MethodGet {
			handler(writer, request)
		}
	})
}
