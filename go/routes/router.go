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
	c   clock
}

// NewRouter returns a new Router
func NewRouter(db *sql.DB) *Router {
	return &Router{
		mux: chi.NewRouter(),
		db:  db,
		c:   &realClock{},
	}
}

// Run sets the routes defined in the router and runs/serves the routes
func (router *Router) Run() error {
	server := router.server()
	log.Printf("Router listening on port %s\n", server.Addr)
	return server.ListenAndServe()
}

func (router *Router) server() *http.Server {
	router.mux = chi.NewRouter()
	router.setRoutes()
	port := "8080"
	return &http.Server{Addr: ":" + port, Handler: router.mux}
}

func (router *Router) withTx(writer http.ResponseWriter, callback func(tx *sql.Tx) bool) {
	tx, err := router.db.Begin()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	success := callback(tx)
	if !success {
		cerr := tx.Rollback()
		if cerr != nil {
			log.Println(cerr)
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
