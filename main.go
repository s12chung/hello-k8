package main

import (
	"github.com/s12chung/hello-k8/go/database"
	"github.com/s12chung/hello-k8/go/routes"
	"log"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := database.DefaultDataBase()
	if err != nil {
		return err
	}
	return routes.NewRouter(db).Run()
}
