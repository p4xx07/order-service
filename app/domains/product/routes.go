package product

import (
	"github.com/gofiber/fiber/v2"
)

func SetRoutes(router fiber.Router, handler IHandler) {
	liveclip := router.Group("product")
	liveclip.Get("/", handler.List)
}
