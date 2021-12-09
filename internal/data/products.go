// Package data contains the products data
package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Product defines the structure of a product API
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// FromJSON populates the product struct from a JSON payload
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	err := e.Decode(p)
	return err
}

func (p *Product) generateNextID() int {
	idx := len(ProductList) - 1
	lp := ProductList[idx]
	return lp.ID + 1
}

func (p *Product) AddProduct() {
	p.ID = p.generateNextID()
	ProductList = append(ProductList, p)
}

var ErrProductNotFound = fmt.Errorf("product not found")

func (p *Product) UpdateProductByID(id int) error {
	_, idx, err := p.getProductByID(id)
	if err != nil {
		return err
	}
	p.ID = id
	ProductList[idx] = p
	return nil
}

// getProductByID returns a product by its ID and it's Index
func (p *Product) getProductByID(id int) (*Product, int, error) {
	for i, v := range ProductList {
		if v.ID == id {
			return v, i, nil
		}
	}
	return nil, -1, ErrProductNotFound
}

type Products []*Product

func (p *Products) Len() int {
	return len(ProductList)
}

func GetProducts() Products {
	return ProductList
}

// ToJSON populates the JSON payload from the product struct
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	err := e.Encode(p)
	return err
}

// ProductList Why pointer? - making ProductList mutable
var ProductList = Products{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
