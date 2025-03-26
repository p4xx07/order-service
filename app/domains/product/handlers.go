package product

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type IHandler interface {
	List(ctx *fiber.Ctx) error
}

type handler struct {
	service IService
	logger  *zap.SugaredLogger
}

func NewHandler(service IService, logger *zap.SugaredLogger) IHandler {
	return &handler{service: service, logger: logger}
}

func (h *handler) List(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	name := c.Query("name")
	description := c.Query("description")
	limitString := c.Query("limit", "10")
	offsetString := c.Query("offset", "0")

	var start, end time.Time
	var err error
	if startDate != "" {
		start, err = time.Parse(time.RFC3339, startDate)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON("Invalid start_date format")
		}
	}
	if endDate != "" {
		end, err = time.Parse(time.RFC3339, endDate)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON("Invalid end_date format")
		}
	}

	limit, err := strconv.ParseInt(limitString, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Invalid limit format")
	}

	offset, err := strconv.ParseInt(offsetString, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON("Invalid offset format")
	}

	response, err := h.service.List(start, end, name, description, limit, offset)
	if err != nil {
		h.logger.Errorf("List error: %v", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(response)
}
