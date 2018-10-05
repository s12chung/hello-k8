package routes

import (
	"database/sql"
	"log"
	"net/http"
)

// Router is a custom server to set routes on
type Router struct {
	mux *http.ServeMux
	db  *sql.DB
}

// NewRouter returns a new Router
func NewRouter(db *sql.DB) *Router {
	return &Router{
		mux: http.NewServeMux(),
		db:  db,
	}
}

// Run sets the routes defined in the router and runs/serves the routes
func (router *Router) Run() error {
	router.setDefaultRoutes()
	router.setRoutes()

	port := "8080"
	log.Printf("Router listening on port %s\n", port)
	return http.ListenAndServe(":"+port, router.mux)
}

func (router *Router) setDefaultRoutes() {
	router.mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.Error(writer, "Not Found", http.StatusNotFound)
	})
}

func (router *Router) get(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	router.mux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if request.Method == http.MethodGet {
			handler(writer, request)
		}
	})
}
