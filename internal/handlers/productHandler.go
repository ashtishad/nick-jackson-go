// Package handlers implements the handlers for Product API
package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
		p.getProducts(w, r)
	case http.MethodPost:
		p.addProduct(w, r)
	case http.MethodPut:
		// expects a product id in the URL path
		// localhost:8080/products/1
		rgx := regexp.MustCompile(`/([0-9]+)`)
		grp := rgx.FindAllStringSubmatch(r.URL.Path, -1)
		// p.l.Printf("Captured Group : %#v", grp)
		// grp = [][]string{[]string{"/3", "3"}}
		if len(grp) != 1 || len(grp[0]) != 2 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idStr := grp[0][1] // "3"
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id, can't convert to number", http.StatusBadRequest)
			return
		}
		p.updateProduct(w, r, id)
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

	// add this new product to the datastore
	prod.AddProduct()

	w.WriteHeader(http.StatusCreated)
	p.l.Printf("Prod : %#v", prod)
}

// updateProduct updates an existing product in the datastore
func (p *Products) updateProduct(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")
	p.l.Println("Handling PUT request...")

	//  create a new instance of product struct
	prod := &data.Product{}

	// deserialize the product struct from the request body
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
		p.l.Printf("Error while decoding json: %v", err)
		return
	}

	err := prod.UpdateProductByID(id)
	if err != nil {
		http.Error(w, "Unable to update by id", http.StatusBadRequest)
		p.l.Printf("Error while updating: %v", err)
		return
	}
}
