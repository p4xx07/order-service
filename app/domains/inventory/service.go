package inventory

import (
	"context"
	"github.com/p4xx07/order-service/configuration"
	"go.uber.org/zap"
)

type IService interface {
	Get(ctx context.Context, productID uint) (*Inventory, error)
	GetMultiple(ctx context.Context, productIDs []uint) (map[uint]Inventory, error) // New function
	DecreaseStockBulk(ctx context.Context, updates map[uint]int) error              // New function
	IncreaseStockBulk(ctx context.Context, updates map[uint]int) error              // New function
}

type service struct {
	configuration *configuration.Configuration
	logger        *zap.SugaredLogger
	store         IStore
}

func NewService(store IStore, configuration *configuration.Configuration, logger *zap.SugaredLogger) IService {
	return &service{store: store, configuration: configuration, logger: logger}
}

func (s *service) Get(ctx context.Context, productID uint) (*Inventory, error) {
	return s.store.Get(ctx, productID)
}

func (s *service) GetMultiple(ctx context.Context, productIDs []uint) (map[uint]Inventory, error) {
	return s.store.GetMultiple(ctx, productIDs)
}

func (s *service) DecreaseStockBulk(ctx context.Context, updates map[uint]int) error {
	return s.store.DecreaseStockBulk(ctx, updates)
}

func (s *service) IncreaseStockBulk(ctx context.Context, updates map[uint]int) error {
	return s.store.IncreaseStockBulk(ctx, updates)
}
