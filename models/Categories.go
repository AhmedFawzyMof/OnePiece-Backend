package models

import (
	"database/sql"
	"fmt"
)

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Img  string `json:"image"`
}

func (c Category) GetAllCategories(db *sql.DB) ([]Category, error) {
	var Categories []Category

	categories, err := db.Query("SELECT * FROM Categories")

	if err != nil {
		return nil, fmt.Errorf("error while prossing categories")
	}

	defer categories.Close()

	for categories.Next() {
		var Category Category

		if err := categories.Scan(&Category.Id, &Category.Name, &Category.Img); err != nil {
			return nil, fmt.Errorf("error while prossing categories")
		}

		Category.Img = "http://localhost:5500/assets" + Category.Img

		Categories = append(Categories, Category)
	}

	return Categories, nil

}
