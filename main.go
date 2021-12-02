package main

import (
	"log"
	"net/http"
	"os"
	"time"

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

	s := &http.Server{
		Addr:         Port,
		Handler:      mux,
		IdleTimeout:  20 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	s.ListenAndServe()
}
