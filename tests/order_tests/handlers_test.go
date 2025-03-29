package order_tests

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/p4xx07/order-service/app/domains/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Post(t *testing.T) {
	mockService := new(MockService)
	logger := zap.NewNop().Sugar()
	handler := order.NewHandler(mockService, logger)

	request := order.PostRequest{
		UserID: 1,
		Items: []order.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	}

	expectedResponse := &order.CreateOrderResponse{ID: 1}

	mockService.On("Create", mock.Anything, request).Return(expectedResponse, nil)

	app := fiber.New()
	app.Post("/orders", handler.Post)

	jRequest, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Error marshalling request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(jRequest))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response order.CreateOrderResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	assert.Equal(t, expectedResponse.ID, response.ID)

	mockService.AssertExpectations(t)
}

func TestHandler_Get(t *testing.T) {
	mockService := new(MockService)
	logger := zap.NewNop().Sugar()
	handler := order.NewHandler(mockService, logger)

	orderID := uint(1)
	expectedResponse := &order.OrderResponse{
		ID: orderID,
		Items: []order.OrderItemResponse{
			{Quantity: 2, Price: 10.0},
			{Quantity: 1, Price: 20.0},
		},
	}

	mockService.On("Get", mock.Anything, orderID).Return(expectedResponse, nil)

	app := fiber.New()
	app.Get("/orders/:id", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response order.Order
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, len(expectedResponse.Items), len(response.Items))

	mockService.AssertExpectations(t)
}

func TestHandler_Put(t *testing.T) {
	mockService := new(MockService)
	logger := zap.NewNop().Sugar()
	handler := order.NewHandler(mockService, logger)

	orderID := uint(1)
	request := order.PutRequest{
		ID: orderID,
		Items: []order.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	}

	mockService.On("Update", mock.Anything, request).Return(nil)

	app := fiber.New()
	app.Put("/orders/:id", handler.Put)

	jRequest, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Error marshalling request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPut, "/orders/1", bytes.NewReader(jRequest))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestHandler_Delete(t *testing.T) {
	mockService := new(MockService)
	logger := zap.NewNop().Sugar()
	handler := order.NewHandler(mockService, logger)

	orderID := uint(1)

	mockService.On("Delete", mock.Anything, orderID).Return(nil)

	app := fiber.New()
	app.Delete("/orders/:id", handler.Delete)

	req := httptest.NewRequest(http.MethodDelete, "/orders/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockService.AssertExpectations(t)
}
