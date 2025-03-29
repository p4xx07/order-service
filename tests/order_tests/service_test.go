package order_tests

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/p4xx07/order-service/app/domains/inventory"
	"github.com/p4xx07/order-service/app/domains/order"
	"github.com/p4xx07/order-service/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
	mockStore := new(MockStore)
	mockInventoryService := new(MockInventoryService)
	mockMeilisearchService := new(MockMeilisearchService)

	logger := zap.NewNop().Sugar()

	orderID := uint(1)
	ord := &order.Order{
		ID: orderID,
		Items: []order.OrderItem{
			{ID: 1, ProductID: 1, Quantity: 2, Price: 10.0},
			{ID: 2, ProductID: 2, Quantity: 1, Price: 20.0},
		},
	}

	mockStore.On("Get", mock.Anything, orderID).Return(ord, nil)

	mockStore.On("DeleteOrderItems", mock.Anything, []uint{1, 2}).Return(nil)
	mockStore.On("Delete", mock.Anything, mock.Anything).Return(nil)

	mockInventoryService.On("IncreaseStockBulk", mock.Anything, mock.Anything).Return(nil)

	mockRedisClient, mockClient := redismock.NewClientMock()

	mockClient.ExpectSetNX("stock_lock_product_1", "locked", 5*time.Second).SetVal(true)
	mockClient.ExpectSetNX("stock_lock_product_2", "locked", 5*time.Second).SetVal(true)
	mockClient.ExpectDel(mock.Anything).RedisNil()

	service := order.NewService(mockMeilisearchService, mockRedisClient, &configuration.Configuration{}, logger, mockStore, mockInventoryService)

	err := service.Delete(context.Background(), orderID)

	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockInventoryService.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mockStore := new(MockStore)
	mockInventoryService := new(MockInventoryService)
	mockMeilisearchService := new(MockMeilisearchService)
	logger := zap.NewNop().Sugar()

	orderID := uint(1)
	ord := &order.Order{
		ID: orderID,
		Items: []order.OrderItem{
			{ID: 1, ProductID: 1, Quantity: 2, Price: 10.0},
			{ID: 2, ProductID: 2, Quantity: 1, Price: 20.0},
		},
	}

	mockStore.On("Get", mock.Anything, orderID).Return(ord, nil)

	mockStore.On("Update", mock.Anything, mock.Anything).Return(nil)

	mockInventoryService.On("IncreaseStockBulk", mock.Anything, mock.Anything).Return(nil)
	mockInventoryService.On("DecreaseStockBulk", mock.Anything, mock.Anything).Return(nil)
	mockInventoryService.On("GetMultiple", mock.Anything, mock.Anything).Return(map[uint]inventory.Inventory{
		1: {Stock: 5},
		2: {Stock: 5},
	}, nil)

	mockRedisClient, mockClient := redismock.NewClientMock()
	mockClient.ExpectSetNX("stock_lock_product_1", "locked", 5*time.Second).SetVal(true)
	mockClient.ExpectSetNX("stock_lock_product_2", "locked", 5*time.Second).SetVal(true)
	mockClient.ExpectDel(mock.Anything).RedisNil()

	service := order.NewService(mockMeilisearchService, mockRedisClient, &configuration.Configuration{}, logger, mockStore, mockInventoryService)

	err := service.Update(context.Background(), order.PutRequest{
		ID: orderID,
		Items: []order.OrderItemRequest{
			{ProductID: 1, Quantity: 3},
			{ProductID: 2, Quantity: 2},
		},
	})

	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockInventoryService.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	mockStore := new(MockStore)
	mockInventoryService := new(MockInventoryService)
	mockMeilisearchService := new(MockMeilisearchService)

	logger := zap.NewNop().Sugar()

	mockStore.On("Create", mock.Anything, mock.Anything).Return(nil)
	mockInventoryService.On("DecreaseStockBulk", mock.Anything, mock.Anything).Return(nil)

	mockInventoryService.On("GetMultiple", mock.Anything, mock.Anything).Return(map[uint]inventory.Inventory{
		1: {Stock: 10},
		2: {Stock: 10},
	}, nil)

	// Mock Redis client
	mockRedisClient, mockClient := redismock.NewClientMock()
	mockClient.ExpectSetNX("stock_lock_product_1", "locked", 5*time.Second).SetVal(true)
	mockClient.ExpectSetNX("stock_lock_product_2", "locked", 5*time.Second).SetVal(true)
	mockClient.ExpectDel(mock.Anything).RedisNil()

	service := order.NewService(mockMeilisearchService, mockRedisClient, &configuration.Configuration{}, logger, mockStore, mockInventoryService)

	_, err := service.Create(context.Background(), order.PostRequest{
		Items: []order.OrderItemRequest{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 1},
		},
	})

	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockInventoryService.AssertExpectations(t)
}
