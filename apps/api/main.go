package main

import (
	"log"
	"os"

	"github.com/aul-pkl/cakli/backend/config"
	"github.com/aul-pkl/cakli/backend/handlers"
	"github.com/aul-pkl/cakli/backend/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("🛠️ Starting CAKLI Backend...")

	// Load .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("ℹ️ No .env file found in current directory, trying root...")
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("⚠️ Warning: No .env file found anywhere. Using system environment variables.")
		}
	}

	// Print some important env vars for debugging (BE CAREFUL NOT TO PRINT SECRETS IN PROD)
	log.Printf("📍 DB_HOST: %s", os.Getenv("DB_HOST"))
	log.Printf("📍 DB_PORT: %s", os.Getenv("DB_PORT"))
	log.Printf("📍 API_PORT: %s", os.Getenv("API_PORT"))

	// Connect to database
	log.Printf("🔗 Connecting to database at %s:%s...", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	config.ConnectDatabase()

	// Connect to MinIO
	log.Printf("☁️ Connecting to MinIO at %s...", os.Getenv("MINIO_ENDPOINT"))
	config.ConnectMinio()

	app := fiber.New(fiber.Config{
		AppName: "CAKLI v1 API",
	})

	log.Println("🚀 Initializing middleware and routes...")

	// 1. Recovery middleware (MUST BE FIRST)
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	// 2. Logger middleware
	app.Use(logger.New())

	// 3. CORS middleware
	app.Use(middleware.CORSMiddleware())

	// 4. API routes
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

	// User auth routes (public)
	userAuth := api.Group("/user/auth")
	userAuth.Post("/login", handlers.LoginUser)
	userAuth.Get("/me", middleware.AuthMiddleware, handlers.GetCurrentUser)

	// Driver auth routes (public)
	driverAuth := api.Group("/driver/auth")
	driverAuth.Post("/login", handlers.LoginDriver)
	driverAuth.Get("/me", middleware.AuthMiddleware, handlers.GetCurrentDriver)

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

	log.Printf("🚀 CAKLI v1 API is ready and starting on 0.0.0.0:%s", port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
