package main

type ProductService interface {
	CreateProduct(product *Product) error
	GetProduct(id int) (Product, error)
	UpdateProduct(id int, product Product) error
	DeleteProduct(id int) error
	GetAllProducts() ([]Product, error)
}
