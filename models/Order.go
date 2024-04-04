package models

import (
	"Onepiece/middleware"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

type OrderItem struct {
	Order        string  `json:"order"`
	Product      int     `json:"product"`
	Quantity     int     `json:"quantity"`
	Color        string  `json:"color"`
	ProductName  string  `json:"name"`
	ProductPrice float64 `json:"price"`
	ProductSlug  string  `json:"slug"`
}

type Order struct {
	Id          string      `json:"id"`
	User        string      `json:"email"`
	Name        string      `json:"name"`
	Method      string      `json:"method"`
	Phone       string      `json:"phone"`
	SparePhone  string      `json:"spare_phone"`
	Address     string      `json:"address"`
	IsPaid      bool        `json:"ispaid"`
	IsDelivered bool        `json:"isdelivered"`
	OrderStatus string      `json:"order_status"`
	OrderItems  []OrderItem `json:"items"`
	CreatedAt   time.Time   `json:"created_at"`
}

type OrdersS []Order

func (o Order) Create(db *sql.DB, resChan chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	Response := make(map[string]interface{})

	_, err := db.Exec("INSERT INTO Orders(id, user, name, method, phone, spare_phone, address, isPaid, isDelivered, order_status, created_at) VALUES(?,?,?,?,?,?,?,?,?,?,?)", o.Id, o.User, o.Name, o.Method, o.Phone, o.SparePhone, o.Address, o.IsPaid, o.IsDelivered, "pending", o.CreatedAt)

	if err != nil {
		Response["Error"] = true
		Response["Message"] = "error occurred while processing the order please try again or contact us"

		middleware.SendResponse(resChan, Response)
		return
	}

	items, err := db.Prepare("INSERT INTO OrdersItem(`order`, product, quantity, color) VALUES(?, ?, ?, ?)")

	if err != nil {
		Response["Error"] = true
		Response["Message"] = "error occurred while processing the order please try again or contact us"

		middleware.SendResponse(resChan, Response)
		return
	}

	for _, item := range o.OrderItems {
		_, err := items.Exec(o.Id, item.Product, item.Quantity, item.Color)

		if err != nil {
			Response["Error"] = true
			Response["Message"] = "error occurred while processing the order please try again or contact us"

			middleware.SendResponse(resChan, Response)
			return
		}
	}

	Response["Done"] = true

	middleware.SendResponse(resChan, Response)
}

func (o Order) GetHistory(db *sql.DB) []Order {

	var Orders []Order

	orders, err := db.Query("SELECT * FROM Orders WHERE Orders.user = ?", o.User)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer orders.Close()

	for orders.Next() {
		var Order Order

		if err := orders.Scan(&Order.Id, &Order.User, &Order.Name, &Order.Method, &Order.Phone, &Order.SparePhone, &Order.Address, &Order.IsPaid, &Order.IsDelivered, &Order.OrderStatus, &Order.CreatedAt); err != nil {
			fmt.Println(err.Error())
		}

		Orders = append(Orders, Order)
	}

	return Orders
}

func (o *OrdersS) GetOrderProducts(db *sql.DB) {

	productsPre, err := db.Prepare("SELECT OrdersItem.`order`, OrdersItem.quantity, OrdersItem.color, Products.name, Products.slug, Products.price FROM OrdersItem INNER JOIN Products ON OrdersItem.product = Products.id WHERE OrdersItem.`order` = ?")

	if err != nil {
		fmt.Println(err.Error())
	}
	for i := range *o {

		products, err := productsPre.Query((*o)[i].Id)

		if err != nil {
			fmt.Println(err.Error())
		}

		defer products.Close()

		for products.Next() {
			var OrderItem OrderItem

			if err := products.Scan(&OrderItem.Order, &OrderItem.Quantity, &OrderItem.Color, &OrderItem.ProductName, &OrderItem.ProductSlug, &OrderItem.ProductPrice); err != nil {
				fmt.Println(err.Error())
			}
			(*o)[i].OrderItems = append((*o)[i].OrderItems, OrderItem)
		}
	}
}
