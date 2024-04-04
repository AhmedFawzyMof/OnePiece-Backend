package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
)

type FilterData struct {
	Min_price int    `json:"min_price"`
	Max_price int    `json:"max_price"`
	Category  string `json:"category"`
}

type Product struct {
	Id            int            `json:"id"`
	Tag           int            `json:"tag"`
	Category      int            `json:"category"`
	Name          string         `json:"name"`
	NameAr        string         `json:"nameAr"`
	Slug          string         `json:"slug"`
	Description   string         `json:"description"`
	DescriptionAr string         `json:"descriptionAr"`
	Price         float64        `json:"price"`
	Discount      float64        `json:"discount"`
	Image         string         `json:"image"`
	Color         sql.NullString `json:"color"`
	CategoryName  string         `json:"categoryName"`
	TagName       string         `json:"tagName"`
}

func (p Product) GetAllProduct(db *sql.DB, productChan chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	var Products []Product

	productsPre, err := db.Prepare("SELECT Products.id, Products.name, Products.nameAr, Products.slug, Products.description, Products.descriptionAr, Products.price, Products.discount, ProductImages.image  FROM Products INNER JOIN ProductImages ON Products.id = ProductImages.product GROUP BY Products.id")

	if err != nil {
		fmt.Println(err.Error())
	}

	products, err := productsPre.Query()

	if err != nil {
		fmt.Println(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var Product Product

		if err := products.Scan(&Product.Id, &Product.Name, &Product.NameAr, &Product.Slug, &Product.Description, &Product.DescriptionAr, &Product.Price, &Product.Discount, &Product.Image); err != nil {
			fmt.Println(err.Error())
		}

		Products = append(Products, Product)
	}

	ProductsBytes, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err.Error())
	}

	productChan <- ProductsBytes
}

func (p Product) FilteredProducts(db *sql.DB, filter FilterData, productChan chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	var Products []Product

	productsPre, err := db.Prepare("SELECT Products.id, Products.name, Products.nameAr, Products.slug, Products.description, Products.descriptionAr, Products.price, Products.discount, ProductImages.image FROM Products INNER JOIN ProductImages ON Products.id = ProductImages.product WHERE Products.category LIKE ? AND Products.price >= ? AND Products.price <= ? GROUP BY Products.id")

	if err != nil {
		fmt.Println(err.Error())
	}

	products, err := productsPre.Query(filter.Category, filter.Min_price, filter.Max_price)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var Product Product

		if err := products.Scan(&Product.Id, &Product.Name, &Product.NameAr, &Product.Slug, &Product.Description, &Product.DescriptionAr, &Product.Price, &Product.Discount, &Product.Image); err != nil {
			fmt.Println(err.Error())
		}

		Products = append(Products, Product)
	}

	ProductsBytes, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err.Error())
	}

	productChan <- ProductsBytes
}

func (p Product) ProductBySlug(db *sql.DB, productChan chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	var Products []Product

	productsPre, err := db.Prepare("SELECT Products.id, Products.name, Products.nameAr, Products.slug, Products.description, Products.descriptionAr, Products.price, Products.discount, ProductImages.image, ProductImages.color, Tags.name, Categories.name FROM Products LEFT JOIN ProductImages ON Products.id = ProductImages.product LEFT JOIN Tags ON Products.tag = Tags.id LEFT JOIN Categories ON Products.category = Categories.id WHERE Products.slug = ?")

	if err != nil {
		fmt.Println(err.Error())
	}

	products, err := productsPre.Query(p.Slug)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer products.Close()

	for products.Next() {

		var Product Product

		if err := products.Scan(&Product.Id, &Product.Name, &Product.NameAr, &Product.Slug, &Product.Description, &Product.DescriptionAr, &Product.Price, &Product.Discount, &Product.Image, &Product.Color, &Product.TagName, &Product.CategoryName); err != nil {
			fmt.Println(err.Error())
		}

		Products = append(Products, Product)
	}

	var product Product
	if len(Products) > 1 {
		for i := range Products {
			if i != len(Products)-1 {
				product.Image += Products[i].Image + ","
				if Products[i].Color.Valid {
					product.Color.String += Products[i].Color.String + ","
				}
			}
			if i == len(Products)-1 {
				product.Id = Products[i].Id
				product.Name = Products[i].Name
				product.Slug = Products[i].Slug
				product.Description = Products[i].Description
				product.Price = Products[i].Price
				product.TagName = Products[i].TagName
				product.CategoryName = Products[i].CategoryName
				if Products[i].Color.Valid {
					product.Color.String += Products[i].Color.String
				}
				product.Image += Products[i].Image
			}
		}
	} else {
		product = Products[0]
	}

	ProductsBytes, err := json.Marshal(product)

	if err != nil {
		fmt.Println(err.Error())
	}

	productChan <- ProductsBytes
}

func (p Product) ProductsByCategorys(db *sql.DB, productChan chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	var Products []Product

	productsPre, err := db.Prepare("SELECT Products.id, Products.name, Products.nameAr, Products.slug, Products.description, Products.descriptionAr, Products.price, Products.discount, ProductImages.image FROM Products INNER JOIN ProductImages ON Products.id = ProductImages.product WHERE Products.category = ? GROUP BY Products.id")

	if err != nil {
		fmt.Println(err.Error())
	}

	products, err := productsPre.Query(p.Category)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var Product Product

		if err := products.Scan(&Product.Id, &Product.Name, &Product.NameAr, &Product.Slug, &Product.Description, &Product.DescriptionAr, &Product.Price, &Product.Discount, &Product.Image); err != nil {
			fmt.Println(err.Error())
		}

		Products = append(Products, Product)
	}

	ProductsBytes, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err.Error())
	}

	productChan <- ProductsBytes
}

func (p Product) ProductsByTag(db *sql.DB, productChan chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	var Products []Product

	productsPre, err := db.Prepare("SELECT Products.id, Products.name, Products.nameAr, Products.slug, Products.description, Products.descriptionAr, Products.price, Products.discount, ProductImages.image FROM Products INNER JOIN ProductImages ON Products.id = ProductImages.product WHERE Products.tag = ? GROUP BY Products.id")

	if err != nil {
		fmt.Println(err.Error())
	}

	products, err := productsPre.Query(p.Tag)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var Product Product

		if err := products.Scan(&Product.Id, &Product.Name, &Product.NameAr, &Product.Slug, &Product.Description, &Product.DescriptionAr, &Product.Price, &Product.Discount, &Product.Image); err != nil {
			fmt.Println(err.Error())
		}

		Products = append(Products, Product)
	}

	ProductsBytes, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err.Error())
	}

	productChan <- ProductsBytes
}
func (p Product) ProductsBySearch(db *sql.DB, productChan chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	var Products []Product

	productsPre, err := db.Prepare("SELECT Products.id, Products.name, Products.nameAr, Products.slug, Products.description, Products.descriptionAr, Products.price, Products.discount, ProductImages.image FROM Products INNER JOIN ProductImages ON Products.id = ProductImages.product WHERE Products.name LIKE ? OR Products.description LIKE ? GROUP BY Products.id")

	if err != nil {
		fmt.Println(err.Error())
	}

	products, err := productsPre.Query(p.Name, p.Description)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer products.Close()

	for products.Next() {
		var Product Product

		if err := products.Scan(&Product.Id, &Product.Name, &Product.NameAr, &Product.Slug, &Product.Description, &Product.DescriptionAr, &Product.Price, &Product.Discount, &Product.Image); err != nil {
			fmt.Println(err.Error())
		}

		Products = append(Products, Product)
	}

	ProductsBytes, err := json.Marshal(Products)

	if err != nil {
		fmt.Println(err.Error())
	}

	productChan <- ProductsBytes
}
