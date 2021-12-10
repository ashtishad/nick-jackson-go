// Package handlers implements the handlers for Product API
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ashtishad/go-microservice/internal/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts returns all the products from datastore
func (p *Products) GetProducts(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p.l.Println("Handling GET request...")

	// fetch all products from the datastore
	lp := data.GetProducts()

	// serialize/encode the list of products to JSON
	if err := lp.ToJSON(w); err != nil {
		http.Error(w, "Unable to encode json", http.StatusBadRequest)
		p.l.Printf("Error while encoding json : %v", err)
		return
	}

	p.l.Printf("Total Products : %#v", lp.Len())
}

// AddProduct adds a new product to the datastore
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	p.l.Println("Handling POST request...")

	prod := r.Context().Value(KeyProd{}).(*data.Product)

	// add this new product to the datastore
	prod.AddProduct()

	//w.WriteHeader(http.StatusCreated)
	p.l.Printf("Prod : %#v", prod)
}

// UpdateProduct updates an existing product in the datastore
func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// gorilla generated id from request URI
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	w.Header().Set("Content-Type", "application/json")
	p.l.Println("Handling PUT request..., id = ", id)

	prod := r.Context().Value(KeyProd{}).(*data.Product)

	if err := prod.UpdateProductByID(id); err != nil {
		http.Error(w, "Unable to update by id", http.StatusBadRequest)
		p.l.Printf("Error while updating: %v", err)
		return
	}
}

type KeyProd struct {
}

// ProductValidationMiddleware validates products
func (p Products) ProductValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//  create a new instance of product struct
		prod := &data.Product{}

		// deserialize the product struct from the request body
		if err := prod.FromJSON(r.Body); err != nil {
			http.Error(w,
				fmt.Sprintf("Unable to read product %s", err),
				http.StatusBadRequest)
			p.l.Printf("Error while reading product: %v", err)
			return
		}

		// validate the product using the Validator package
		if err := prod.Validate(); err != nil {
			http.Error(w,
				fmt.Sprintf("Unable to validate product %s", err),
				http.StatusBadRequest)
			p.l.Printf("Error while validating product: %v", err)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProd{}, prod)
		nextReq := r.WithContext(ctx)
		// p.l.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, nextReq)
	})
}
