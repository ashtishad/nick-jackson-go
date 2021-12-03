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
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to encode json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}
