package main

import (
	"log"
	"os"

	"github.com/anuzx/live-attendance-backend/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	//connect to db
	db, err := database.NewPostgresPool(databaseURL)

	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	defer db.Close()

	router := gin.Default()

	//routes
	auth := router.Group("/auth")
	{
		auth.POST("/signup")
	}

	//start server
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
