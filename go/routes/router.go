package routes

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

// Router is a custom server to set routes on
type Router struct {
	mux *chi.Mux
	db  *sql.DB
}

// NewRouter returns a new Router
func NewRouter(db *sql.DB) *Router {
	return &Router{
		mux: chi.NewRouter(),
		db:  db,
	}
}

// Run sets the routes defined in the router and runs/serves the routes
func (router *Router) Run() error {
	router.setRoutes()

	port := "8080"
	log.Printf("Router listening on port %s\n", port)
	return http.ListenAndServe(":"+port, router.mux)
}
