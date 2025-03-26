package product

import (
	"github.com/p4xx07/order-service/configuration"
	"go.uber.org/zap"
	"time"
)

type IService interface {
	List(from, to time.Time, name, description string, limit, offset int64) ([]Product, error)
}

type service struct {
	configuration *configuration.Configuration
	logger        *zap.SugaredLogger
	store         IStore
}

func NewService(store IStore, configuration *configuration.Configuration, logger *zap.SugaredLogger) IService {
	return &service{store: store, configuration: configuration, logger: logger}
}

func (s *service) List(from, to time.Time, name, description string, limit, offset int64) ([]Product, error) {
	return s.store.List(from, to, name, description, limit, offset)
}
