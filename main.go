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

func main() {
	router := gin.Default()
	router.GET("/orders", getOrders)

	router.Run("localhost:8080")
}
