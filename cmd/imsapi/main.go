package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"mereb-crud/pkg/common/app"
	"mereb-crud/pkg/common/postgresql"
	"mereb-crud/pkg/controller"
	"mereb-crud/pkg/repository"
	"mereb-crud/pkg/service"
	"os"
	"path/filepath"
)

func main() {
	envPath := filepath.Join("..", "..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	configurationManager := app.NewConfigurationManager()
	dbPool := postgresql.GetConnectionPool(ctx, configurationManager.PostgresqlConfig)

	userRepository := repository.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	personRepository := repository.NewPersonRepository(dbPool)
	personService := service.NewPersonService(personRepository)
	personController := controller.NewPersonController(personService)

	e := echo.New()
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
