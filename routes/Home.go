package routes

import (
	"Onepiece/database"
	"Onepiece/middleware"
	"Onepiece/models"
	"encoding/json"
	"net/http"
)

func Home(res http.ResponseWriter, req *http.Request, params map[string]string) {
	res.WriteHeader(http.StatusOK)

	db := database.Connect()
	defer db.Close()

	category := models.Category{}
	Categories, err := category.GetAllCategories(db)

	if err != nil {
		middleware.SendError(err, res)
		return
	}

	carousel := models.Carousel{}
	Carousels, err := carousel.GetHomeCarousel(db)

	if err != nil {
		middleware.SendError(err, res)
		return
	}

	Response := make(map[string]interface{})
	Response["Categories"] = Categories
	Response["Carousels"] = Carousels

	if err := json.NewEncoder(res).Encode(Response); err != nil {
		middleware.SendError(err, res)
		return
	}
}
