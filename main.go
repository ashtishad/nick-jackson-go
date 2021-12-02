package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ashtishad/go-microservice/handlers"
)

const Port = ":8080"

func main() {

	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	hh := handlers.NewHello(l)
	gh := handlers.Newbye(l)
	mux := http.NewServeMux()

	mux.Handle("/", hh)
	mux.Handle("/bye", gh)

	http.ListenAndServe(Port, mux)
}
