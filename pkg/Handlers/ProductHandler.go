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

func (p *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	if err := lp.ToJSON(w); err != nil {
		http.Error(w, "Unable to encode json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p.l.Println("Handle POST request")

	prod := &data.Product{} // create a new instance of Product struct
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
		p.l.Printf("Error while decoding json: %v", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	p.l.Printf("Prod : %#v", prod)
}
