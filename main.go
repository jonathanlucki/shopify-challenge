package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jonathanlucki/shopify-challenge/adapters"
	"github.com/jonathanlucki/shopify-challenge/handlers"
)

func setDatabase(db *adapters.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Database", db)
		c.Next()
	}
}

func setStorage(s *adapters.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Storage", s)
		c.Next()
	}
}

func initRoutes(router *gin.Engine) {
	// assign static routes
	router.Static("/public", "./public")

	// assign routes
	router.GET("/images", handlers.GetImages())
	router.POST("/images", handlers.UploadImage())
	router.GET("/images/:id", handlers.GetImage())
	router.DELETE("/images/:id", handlers.DeleteImage())
}

func main() {
	// get server port
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT is not set")
	}

	// initialize database
	db, err := adapters.InitDB()
	if err != nil {
		log.Printf("Database Error: %v", err)
		os.Exit(1)
	}

	// initialize storage
	storage, err := adapters.InitStorage()
	if err != nil {
		log.Printf("Storage Error: %v", err)
		os.Exit(1)
	}

	// set up gin router and set middleware
	router := gin.Default()
	router.Use(setDatabase(db))
	router.Use(setStorage(storage))

	// initialize routes
	initRoutes(router)

	// run http server
	router.Run(":" + port)
}