package order

import (
	"github.com/gofiber/fiber/v2"
)

func SetRoutes(router fiber.Router, handler IHandler) {
	liveclip := router.Group("order")
	liveclip.Get("/", handler.List)
	liveclip.Post("/", handler.Post)
	liveclip.Get("/:id", handler.Get)
	liveclip.Put("/:id", handler.Put)
	liveclip.Delete("/:id", handler.Delete)
}
