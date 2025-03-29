package http

import "github.com/gofiber/fiber/v2"

type ResponseWrapper struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

func JSON(c *fiber.Ctx, statusCode int, data any, err error) error {
	response := ResponseWrapper{
		Status: "success",
		Data:   data,
	}

	if err != nil {
		response.Status = "error"
		response.Error = err.Error()
	}

	return c.Status(statusCode).JSON(response)
}
