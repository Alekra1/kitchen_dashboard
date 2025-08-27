package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type order struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var orders = []order{
	{ID: "1", Name: "California", Price: 56.99},
	{ID: "2", Name: "Filadelfia", Price: 69},
	{ID: "3", Name: "Oleksii", Price: 222},
}

// getOrders responds with the list of all orders as JSON.
func getOrders(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, orders)
}

func postOrders(c *gin.Context) {
	var newOrder order

	// Call BindJSON to bind the received JSON to
	// newOrder.
	if err := c.BindJSON(&newOrder); err != nil {
		return
	}

	// Add the new order to the slice.
	orders = append(orders, newOrder)
	c.IndentedJSON(http.StatusCreated, newOrder)
}

func getOrderByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of orders, looking for
	// an order whose ID value matches the parameter.
	for _, a := range orders {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
}

func main() {
	router := gin.Default()
	router.GET("/orders", getOrders)
	router.POST("/orders", postOrders)
	router.GET("/orders/:id", getOrderByID)
	router.Run("localhost:8080")
}
