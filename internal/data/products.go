// Package data contains the products data
package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
	"regexp"
	"time"
)

// Product defines the structure of a product API
// uses Validator package
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// https://pkg.go.dev/github.com/go-playground/validator/v10
// sample: https://github.com/go-playground/validator/blob/master/_examples/simple/main.go
// RegisterValidation:  https://github.com/go-playground/validator/blob/v10.9.0/validator_instance.go#L198
var validate *validator.Validate

func (p *Product) Validate() error {
	validate = validator.New()
	err := validate.RegisterValidation("sku", validateSKU)
	if err != nil {
		fmt.Println(errors.New("error while registering custom validation function"))
	}
	return validate.Struct(p)
}

// validateSKU is our custom validate function
func validateSKU(fl validator.FieldLevel) bool {
	// regex sku format is abc-abs-dfs
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}
	return true
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

// getProductByID returns a product by its ID, and it's Index
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

// ToJSON populates the JSON payload from the product struct(serialization)
// NewEncoder provides better performance than json.Unmarshal
// as it doesn't have to buffer the output into an in memory slice of bytes
// thus it reduces allocations and overheads
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
		SKU:         "abc-abf-def",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd-ksf-jdt",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
