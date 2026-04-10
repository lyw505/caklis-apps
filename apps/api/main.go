package main

import (
	"log"
	"os"

	"github.com/aul-pkl/cakli/backend/config"
	"github.com/aul-pkl/cakli/backend/handlers"
	"github.com/aul-pkl/cakli/backend/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, trying root .env")
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("Warning: No .env file found")
		}
	}

	// Connect to database
	config.ConnectDatabase()

	// Connect to MinIO
	config.ConnectMinio()

	app := fiber.New(fiber.Config{
		AppName: "CAKLI v1 API",
	})

	// CORS middleware
	app.Use(middleware.CORSMiddleware())

	// API routes
	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "CAKLI v1 API is running",
			"data": fiber.Map{
				"version": "1.0.0",
				"status":  "healthy",
			},
		})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)
	auth.Post("/refresh", handlers.Refresh)
	auth.Post("/logout", handlers.Logout)
	auth.Get("/me", middleware.AuthMiddleware, handlers.GetCurrentAdmin)

	// User Auth routes (mobile/Flutter)
	userAuth := api.Group("/user/auth")
	userAuth.Post("/login", handlers.UserLogin)
	userAuth.Get("/me", middleware.AuthMiddleware, handlers.GetCurrentUser)

	// Upload routes (protected) - accessible by master_admin and operation_admin
	upload := api.Group("/upload", middleware.AuthMiddleware, middleware.RequireMasterOrOperation())
	upload.Post("/presigned-url", handlers.GenerateUploadURL)
	upload.Get("/presigned-view-url", handlers.GenerateViewURL)

	// Protected routes
	admin := api.Group("/admin", middleware.AuthMiddleware)
	
	// Dashboard stats - accessible by all admin roles
	admin.Get("/dashboard/stats", handlers.GetDashboardStats)

	// Banks reference data - accessible by all admin roles
	admin.Get("/banks", handlers.ListBanks)

	// Drivers CRUD - accessible by master_admin and operation_admin
	drivers := admin.Group("/drivers", middleware.RequireMasterOrOperation())
	drivers.Get("/", handlers.ListDrivers)
	drivers.Get("/:id", handlers.GetDriver)
	drivers.Post("/", handlers.CreateDriver)
	drivers.Put("/:id", handlers.UpdateDriver)
	drivers.Delete("/:id", handlers.DeleteDriver)
	drivers.Patch("/:id/verification-status", handlers.UpdateVerificationStatus)

	// Users CRUD - accessible by master_admin and operation_admin
	users := admin.Group("/users", middleware.RequireMasterOrOperation())
	users.Get("/", handlers.ListUsers)
	users.Get("/:id", handlers.GetUser)
	users.Post("/", handlers.CreateUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
