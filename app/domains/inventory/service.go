package inventory

import (
	"errors"
	"github.com/p4xx07/order-service/configuration"
	"go.uber.org/zap"
)

type IService interface {
	AdjustStock(productID uint, quantityChange int) error
}

type service struct {
	configuration *configuration.Configuration
	logger        *zap.SugaredLogger
	store         IStore
}

func NewService(store IStore, configuration *configuration.Configuration, logger *zap.SugaredLogger) IService {
	return &service{store: store, configuration: configuration, logger: logger}
}

func (s *service) AdjustStock(productID uint, quantityChange int) error {
	stock, err := s.store.GetStock(productID)
	if err != nil {
		return err
	}

	if stock+quantityChange < 0 {
		return errors.New("insufficient stock")
	}

	return s.store.UpdateStock(productID, quantityChange)
}
