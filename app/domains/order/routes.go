package order

import (
	"github.com/gofiber/fiber/v2"
)

func SetRoutes(router fiber.Router, handler IHandler) {
	g := router.Group("order")
	g.Get("/", handler.List)
	g.Post("/", handler.Post)
	g.Get("/:id", handler.Get)
	g.Put("/:id", handler.Put)
	g.Delete("/:id", handler.Delete)
}
