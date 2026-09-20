package main

import (
	"log"
	"os"

	"github.com/anuzx/live-attendance-backend/internal/auth"
	"github.com/anuzx/live-attendance-backend/internal/database"
	"github.com/anuzx/live-attendance-backend/internal/http/attendance"
	"github.com/anuzx/live-attendance-backend/internal/http/class"
	"github.com/anuzx/live-attendance-backend/internal/http/students"
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

	// Dependency wiring for auth endpoints
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	//class endpoints
	classRepo := class.NewRepository(db)
	classService := class.NewService(classRepo)
	classHandler := class.NewHandler(classService)

	//students endpoint
	studentRepo := students.NewRepository(db)
	studentService := students.NewService(studentRepo)
	studentHandler := students.NewHandler(studentService)

	//attendance endpoint
	store := attendance.NewSessionStore()
	attService := attendance.NewService(classService, store)
	attHandler := attendance.NewHandler(attService)

	router := gin.Default()

	authMiddleware := auth.AuthMiddleware()

	//routes
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", userHandler.Signup)
		authRoutes.POST("/login", userHandler.Login)
		authRoutes.GET("/me", authMiddleware, userHandler.Me)
	}

	classRoutes := router.Group("/class")
	{
		classRoutes.POST("/", authMiddleware, auth.OnlyTeacher(), classHandler.CreateClass)
		classRoutes.POST("/:id/add-student", authMiddleware, auth.OnlyTeacher(), classHandler.AddStudent)
		classRoutes.GET("/:id", authMiddleware, classHandler.GetClassDetails)
		classRoutes.GET("/:id/my-attendance", authMiddleware, auth.OnlyStudent(), classHandler.MyAttendance)

	}

	router.GET("/students", authMiddleware, auth.OnlyTeacher(), studentHandler.GetStudents)

	router.POST("/attendance/start", authMiddleware, auth.OnlyTeacher(), attHandler.StartSession)

	//start server
	if err := router.Run(":3000"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}

/**
 * type -> repository -> service -> handler -> route
 */
