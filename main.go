package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	prefix = "item:"
	port   = "8080"
)

type Item struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

var DB map[string]*Item = map[string]*Item{}

func main() {
	// Set Gin mode from environment variable
	gin.SetMode(os.Getenv(gin.ReleaseMode))

	// Initialize Gin router
	router := gin.Default()

	// Define routes
	router.GET("/", homeHandler)

	api := router.Group("/crud")
	api.POST("/", createHandler)
	api.GET("/", readHandler)
	api.PUT("/", updateHandler)
	api.DELETE("/", deleteHandler)

	// Start the server using port from environment variable
	p := os.Getenv("PORT")
	if p != "" {
		port = p
	}
	router.Run(":" + port)
}

// Home Route Handler
func homeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the CRUD API using Gin",
		"port":    port,
	})
}

// Create Handler
func createHandler(c *gin.Context) {
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if item already exists in DB
	if _, exists := DB[prefix+item.ID]; exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item: ID already present"})
		return
	}

	// Store item in DB
	id := prefix + item.ID
	DB[id] = &item
	c.JSON(http.StatusCreated, gin.H{"message": "Item created"})
}

// Read Handler
func readHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID parameter"})
		return
	}

	item := DB[prefix+id]
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Send the item back
	c.JSON(http.StatusOK, item)
}

// Update Handler
func updateHandler(c *gin.Context) {
	var item Item

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if item exists
	oldItem := DB[prefix+item.ID]
	if oldItem == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	// Update item in DB
	DB[prefix+item.ID].Value = item.Value
	c.JSON(http.StatusOK, DB[prefix+item.ID])
}

// Delete Handler
func deleteHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID parameter"})
		return
	}

	// Check if item exists
	if _, exists := DB[prefix+id]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found"})
		return
	}

	// Delete item from DB
	delete(DB, prefix+id)
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
