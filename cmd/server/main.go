package main

import (
	"log"
	"os"

	"github.com/anuzx/live-attendance-backend/internal/auth"
	"github.com/anuzx/live-attendance-backend/internal/database"
	"github.com/anuzx/live-attendance-backend/internal/http/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

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

	// Dependency wiring
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	router := gin.Default()

	authMiddleware := auth.AuthMiddleware()

	//routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", userHandler.Signup)
		authRoutes.POST("/login", userHandler.Login)
		authRoutes.GET("/me", authMiddleware, userHandler.Me)
	}

	//start server
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
