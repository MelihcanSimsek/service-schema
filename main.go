package main

import (
	"Service-schema/controller"
	"Service-schema/core/app"
	"Service-schema/core/postgresql"
	"Service-schema/persistence"
	"Service-schema/service"
	"context"
	"github.com/labstack/echo/v4"
)

func main() {
	ctx := context.Background()
	e := echo.New()

	configurationManager := app.NewConfigurationManager()

	dbPool := postgresql.GetConnectionPool(ctx, configurationManager.PostgresqlConfig)

	productRepository := persistence.NewProductRepository(dbPool)

	productService := service.NewProductService(productRepository)

	productController := controller.NewProductController(productService)

	productController.RegisterRoutes(e)

	e.Start("localhost:8080")
}
