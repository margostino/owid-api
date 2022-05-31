package main

import (
	handler "github.com/margostino/owid-api/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/query", handler.Query)
	http.HandleFunc("/playground", handler.Playground)
	log.Println("Starting OWID-API Server in :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
