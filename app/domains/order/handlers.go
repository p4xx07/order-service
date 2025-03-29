package order

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gopkg.in/validator.v2"
	"net/http"
	"strconv"
)

type IHandler interface {
	Post(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Put(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type handler struct {
	service IService
	logger  *zap.SugaredLogger
}

func NewHandler(service IService, logger *zap.SugaredLogger) IHandler {
	return &handler{service: service, logger: logger}
}

func (h *handler) Post(c *fiber.Ctx) error {
	var request PostRequest
	if err := c.BodyParser(&request); err != nil {
		h.logger.Errorf("bodyRequest error %v | %v", request, err.Error())
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	if errs := validator.Validate(request); errs != nil {
		return c.Status(http.StatusBadRequest).JSON(errs)
	}

	response, err := h.service.Create(c.Context(), request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h *handler) Get(c *fiber.Ctx) error {
	orderIDString := c.Params("id")
	orderID, err := strconv.ParseUint(orderIDString, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	response, err := h.service.Get(c.Context(), uint(orderID))
	if err != nil {
		h.logger.Errorf("Get order error: %v", err.Error())
		return c.Status(http.StatusNotFound).JSON("Order not found")
	}

	return c.Status(http.StatusOK).JSON(response)
}

func (h *handler) Put(c *fiber.Ctx) error {
	orderIDString := c.Params("id")
	orderID, err := strconv.ParseUint(orderIDString, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	var request PutRequest
	if err := c.BodyParser(&request); err != nil {
		h.logger.Errorf("bodyRequest error: %v", err.Error())
		return c.Status(http.StatusBadRequest).JSON("Invalid request body")
	}

	if errs := validator.Validate(request); errs != nil {
		return c.Status(http.StatusBadRequest).JSON(errs)
	}

	request.ID = uint(orderID)
	err = h.service.Update(c.Context(), request)
	if err != nil {
		h.logger.Errorf("Update error: %v", err.Error())
		return c.Status(http.StatusInternalServerError).JSON("Failed to update order")
	}

	return c.SendStatus(http.StatusOK)
}

func (h *handler) Delete(c *fiber.Ctx) error {
	orderIDString := c.Params("id")
	orderID, err := strconv.ParseUint(orderIDString, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	err = h.service.Delete(c.Context(), uint(orderID))
	if err != nil {
		h.logger.Errorf("Delete error: %v", err.Error())
		return c.Status(http.StatusInternalServerError).JSON("Failed to delete order")
	}

	return c.SendStatus(http.StatusOK)
}
