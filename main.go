package main

import (
	"log"
	"os"

	"github.com/s12chung/hello-k8/go/routes"
)

func main() {
	err := routes.NewRouter().Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
