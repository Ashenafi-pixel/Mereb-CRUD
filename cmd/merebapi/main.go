package main

import (
	"context"
	"log"
	"mereb-crud/pkg/common/app"
	"mereb-crud/pkg/common/postgresql"
	"mereb-crud/pkg/controller"
	"mereb-crud/pkg/repository"
	"mereb-crud/pkg/service"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	envPath := filepath.Join("..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize context and configurations
	ctx := context.Background()
	configurationManager := app.NewConfigurationManager()
	dbPool := postgresql.GetConnectionPool(ctx, configurationManager.PostgresqlConfig)

	// Initialize repositories, services, and controllers
	userRepository := repository.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
	personRepository := repository.NewPersonRepository(dbPool)
	personService := service.NewPersonService(personRepository)
	personController := controller.NewPersonController(personService)
	e := echo.New()

	// Enable Access to from any host or client for mereb
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	userController.RegisterUserRoutes(e)
	personController.RegisterPersonRoutes(e)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server is running on port", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
