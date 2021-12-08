// Package handlers implements the handlers for Product API
package handlers

import (
	"log"
	"net/http"

	"github.com/ashtishad/go-microservice/internal/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// ServeHTTP : Products implicitly implements the http.Handler interface
func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		p.l.Println("Handle GET request")
		p.getProducts(w, r)
	case http.MethodPost:
		p.addProduct(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// getProducts returns all the products from datastore
func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p.l.Println("Handling GET request...")

	// fetch all products from the datastore
	lp := data.GetProducts()

	// serialize/encode the list of products to JSON
	if err := lp.ToJSON(w); err != nil {
		http.Error(w, "Unable to encode json", http.StatusInternalServerError)
		return
	}

	p.l.Printf("Total Products : %#v", lp.Len())
}

// addProduct adds a new product to the datastore
func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p.l.Println("Handling POST request...")

	//  create a new instance of product struct
	prod := &data.Product{}

	// deserialize the product struct from the request body
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
		p.l.Printf("Error while decoding json: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	p.l.Printf("Prod : %#v", prod)
}
