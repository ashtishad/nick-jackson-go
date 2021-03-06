package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/ashtishad/go-microservice/internal/handlers"
)

const Port = ":8080"

func main() {

	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProducts(l)

	// create a new serve mux and register the handlers
	r := mux.NewRouter()

	getRtr := r.Methods(http.MethodGet).Subrouter()
	getRtr.HandleFunc("/", ph.GetProducts)

	putRtr := r.Methods(http.MethodPut).Subrouter()
	putRtr.HandleFunc("/{id:[0-9]}", ph.UpdateProduct)
	putRtr.Use(ph.ProductValidationMiddleware)

	postRtr := r.Methods(http.MethodPost).Subrouter()
	postRtr.HandleFunc("/", ph.AddProduct)
	postRtr.Use(ph.ProductValidationMiddleware)

	s := &http.Server{
		Addr:         Port,
		Handler:      r,
		IdleTimeout:  20 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 8080")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 30 seconds.
	quit := make(chan os.Signal, 1) // For a channel used for notification of just one signal value, a buffer of size 1 is sufficient.
	signal.Notify(quit, os.Interrupt)
	<-quit
	l.Println("Gracefully Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		l.Println("Could not gracefully shutdown the server:", err)
		os.Exit(1)
	}
	l.Println("Server exiting")
}
