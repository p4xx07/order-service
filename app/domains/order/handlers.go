package order

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	http2 "github.com/p4xx07/order-service/internal/http"
	"go.uber.org/zap"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type IHandler interface {
	List(ctx *fiber.Ctx) error
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
		if errors.Is(err, ErrNoStockAvailable) {
			return http2.JSON(c, http.StatusInternalServerError, nil, err)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http2.JSON(c, http.StatusNotFound, nil, err)
		}

		h.logger.Error(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return http2.JSON(c, http.StatusOK, response, nil)
}

func (h *handler) List(c *fiber.Ctx) error {
	input := c.Query("input")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limit := c.Query("limit")
	offset := c.Query("offset")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			h.logger.Error("Invalid start_date format: ", err)
			return http2.JSON(c, http.StatusBadRequest, nil, err)
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			h.logger.Error("Invalid end_date format: ", err)
			return http2.JSON(c, http.StatusBadRequest, nil, err)
		}
	}

	limitInt, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		limitInt = 0
	}

	offsetInt, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		offsetInt = 0
	}

	request := ListRequest{
		Input:     input,
		StartDate: &startDate,
		EndDate:   &endDate,
		Limit:     limitInt,
		Offset:    offsetInt,
	}

	response, err := h.service.List(c.Context(), request)
	if err != nil {
		h.logger.Error(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return http2.JSON(c, http.StatusOK, response, nil)
}

func (h *handler) Get(c *fiber.Ctx) error {
	orderIDString := c.Params("id")
	orderID, err := strconv.ParseUint(orderIDString, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	response, err := h.service.Get(c.Context(), uint(orderID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http2.JSON(c, http.StatusNotFound, nil, err)
		}
		h.logger.Error(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return http2.JSON(c, http.StatusOK, response, err)
}

func (h *handler) Put(c *fiber.Ctx) error {
	orderIDString := c.Params("id")
	orderID, err := strconv.ParseUint(orderIDString, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}

	var request PutRequest
	request.ID = uint(orderID)
	if err := c.BodyParser(&request); err != nil {
		h.logger.Errorf("bodyRequest error: %v", err.Error())
		return c.Status(http.StatusBadRequest).JSON("Invalid request body")
	}

	if errs := validator.Validate(request); errs != nil {
		return c.Status(http.StatusBadRequest).JSON(errs)
	}

	err = h.service.Update(c.Context(), request)
	if err != nil {
		if errors.Is(err, ErrNoStockAvailable) {
			return http2.JSON(c, http.StatusInternalServerError, nil, err)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http2.JSON(c, http.StatusNotFound, nil, err)
		}

		h.logger.Error(err)
		return c.SendStatus(http.StatusInternalServerError)
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
		if errors.Is(err, ErrNoStockAvailable) {
			return http2.JSON(c, http.StatusInternalServerError, nil, err)
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return http2.JSON(c, http.StatusNotFound, nil, err)
		}

		h.logger.Error(err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}
