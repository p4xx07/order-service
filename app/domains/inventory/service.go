package inventory

import (
	"context"
	"errors"
	"fmt"
	"github.com/p4xx07/order-service/configuration"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type IService interface {
	AdjustStock(productID uint, quantityChange int) error
}

type service struct {
	configuration *configuration.Configuration
	logger        *zap.SugaredLogger
	store         IStore
	redisClient   *redis.Client
}

func NewService(redisClient *redis.Client, store IStore, configuration *configuration.Configuration, logger *zap.SugaredLogger) IService {
	return &service{redisClient: redisClient, store: store, configuration: configuration, logger: logger}
}

func (s *service) AdjustStock(productID uint, quantityChange int) error {
	ctx := context.Background()
	lockKey := fmt.Sprintf("inventory_lock:%d", productID)
	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())
	lockTTL := 5 * time.Second

	ok, err := s.redisClient.SetNX(ctx, lockKey, lockValue, lockTTL).Result()
	if err != nil {
		s.logger.Errorw("error acquiring Redis lock", "error", err)
		return err
	}
	if !ok {
		return errors.New("inventory is locked by another process")
	}

	defer s.redisClient.Del(ctx, lockKey)

	stock, err := s.store.GetStock(productID)
	if err != nil {
		return err
	}

	if stock+quantityChange < 0 {
		return errors.New("insufficient stock")
	}

	return s.store.UpdateStock(productID, quantityChange)
}
