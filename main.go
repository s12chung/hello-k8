package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	err := runServer()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func runServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.Write([]byte(`{ "cpu_used": 100 }`))
	})

	port := "8080"
	log.Printf("Server listening on port %s\n", port)
	return http.ListenAndServe(":"+port, mux)
}
