package main

import (
	"github.com/s12chung/hello-k8/go/database"
	"log"
	"os"

	"github.com/s12chung/hello-k8/go/routes"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	db, err := database.DefaultDataBase()
	if err != nil {
		return err
	}
	return routes.NewRouter(db).Run()
}
