package order_tests

import (
	"context"
	"github.com/p4xx07/order-service/app/domains/inventory"
	"github.com/p4xx07/order-service/app/domains/order"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Get(ctx context.Context, id uint) (*order.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockStore) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStore) DeleteOrderItems(ctx context.Context, orderItemIDs []uint) error {
	args := m.Called(ctx, orderItemIDs)
	return args.Error(0)
}

func (m *MockStore) Create(ctx context.Context, order *order.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockStore) Update(ctx context.Context, order *order.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

type MockInventoryService struct {
	mock.Mock
}

func (m *MockInventoryService) Get(ctx context.Context, productID uint) (*inventory.Inventory, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).(*inventory.Inventory), args.Error(0)
}

func (m *MockInventoryService) DecreaseStockBulk(ctx context.Context, updates map[uint]int) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *MockInventoryService) IncreaseStockBulk(ctx context.Context, updates map[uint]int) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func (m *MockInventoryService) GetMultiple(ctx context.Context, ids []uint) (map[uint]inventory.Inventory, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).(map[uint]inventory.Inventory), args.Error(1)
}

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, request order.PostRequest) (*order.CreateOrderResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*order.CreateOrderResponse), args.Error(1)
}

func (m *MockService) Get(ctx context.Context, id uint) (*order.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, request order.PutRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
