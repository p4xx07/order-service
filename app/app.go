package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/p4xx07/order-service/app/domains/order"
	"net/http"
)

type App struct {
	OrderHandler order.IHandler
}

func (a *App) Routes() *fiber.App {
	f := fiber.New()
	f.Use(logger.New())
	f.Use(recover.New())
	f.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Content-Type",
		AllowMethods: "GET, HEAD, OPTIONS, PUT, PATCH, POST, DELETE",
	}))
	f.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(http.StatusOK) })

	api := f.Group("/api/v1.0")

	order.SetRoutes(api, a.OrderHandler)

	return f
}
