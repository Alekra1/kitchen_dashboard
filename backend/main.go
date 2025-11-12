package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Alekra1/kitchen_dashboard.git/db"
)

// getOrders responds with the list of all orders as JSON.
func getOrders(c *gin.Context) {
	orders, err := db.ListOrders(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch orders"})
		return
	}

	c.IndentedJSON(http.StatusOK, orders)
}

func postOrders(c *gin.Context) {
	var payload db.Order
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	stored, err := db.CreateOrder(c.Request.Context(), payload)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to create order"})
		return
	}

	c.IndentedJSON(http.StatusCreated, stored)
}

func getOrderByID(c *gin.Context) {
	order, err := db.GetOrder(c.Request.Context(), c.Param("id"))
	switch {
	case errors.Is(err, db.ErrOrderNotFound):
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "order not found"})
	case err != nil:
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch order"})
	default:
		c.IndentedJSON(http.StatusOK, order)
	}
}

func main() {
	if err := db.Connect(context.Background()); err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/orders", getOrders)
	router.POST("/orders", postOrders)
	router.GET("/orders/:id", getOrderByID)
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
