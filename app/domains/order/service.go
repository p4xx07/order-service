package order

import (
	"context"
	"fmt"
	"github.com/p4xx07/order-service/app/domains/inventory"
	"github.com/p4xx07/order-service/configuration"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type IService interface {
	Get(ctx context.Context, orderID uint) (*Order, error)
	Create(ctx context.Context, request PostRequest) (*CreateOrderResponse, error)
	Update(ctx context.Context, request PutRequest) error
	Delete(ctx context.Context, id uint) error
}

type service struct {
	configuration    *configuration.Configuration
	logger           *zap.SugaredLogger
	store            IStore
	inventoryService inventory.IService
	redisClient      *redis.Client
}

func NewService(redisClient *redis.Client, configuration *configuration.Configuration, logger *zap.SugaredLogger, store IStore, inventoryService inventory.IService) IService {
	return &service{redisClient: redisClient, configuration: configuration, logger: logger, store: store, inventoryService: inventoryService}
}

func (s *service) Create(ctx context.Context, request PostRequest) (*CreateOrderResponse, error) {
	productIDs := make([]uint, len(request.Items))
	for i, item := range request.Items {
		productIDs[i] = item.ProductID
	}

	inventories, err := s.inventoryService.GetMultiple(ctx, productIDs)
	if err != nil {
		s.logger.Errorw("error getting inventory", "error", err)
		return nil, err
	}

	var lockedProducts []string
	for _, item := range request.Items {
		redisKey := s.getLockProductKey(item.ProductID)
		lock := s.redisClient.SetNX(ctx, redisKey, "locked", 5*time.Second)
		if !lock.Val() {
			s.logger.Errorw("stock update in progress", "productID", item.ProductID)
			return nil, fmt.Errorf("stock update in progress for product %d", item.ProductID)
		}
		lockedProducts = append(lockedProducts, redisKey)

		if inventories[item.ProductID].Stock < item.Quantity {
			s.logger.Errorw("stock update in progress for product %d", item.ProductID)
			return nil, fmt.Errorf("not enough stock for product %d", item.ProductID)
		}
	}

	defer func() {
		for _, key := range lockedProducts {
			s.redisClient.Del(ctx, key)
		}
	}()

	updates := map[uint]int{}
	for _, item := range request.Items {
		updates[item.ProductID] = item.Quantity
	}

	if err := s.inventoryService.DecreaseStockBulk(ctx, updates); err != nil {
		s.logger.Errorw("failed to decrease stock bulk", "error", err)
		return nil, err
	}

	items := NewItems(request.Items)
	order := NewOrder(request.UserID, items)
	err = s.store.Create(ctx, order)
	if err != nil {
		s.logger.Errorw("failed to store order", "error", err)
		return nil, err
	}

	return &CreateOrderResponse{ID: order.ID}, nil
}

func (s *service) Update(ctx context.Context, request PutRequest) error {
	existingOrder, err := s.store.Get(ctx, request.ID)
	if err != nil {
		s.logger.Errorw("error getting existing order", "error", err, "id", request.ID)
		return fmt.Errorf("order not found: %w", err)
	}

	productIDs := make(map[uint]struct{})
	for _, item := range existingOrder.Items {
		productIDs[item.ProductID] = struct{}{}
	}
	for _, item := range request.Items {
		productIDs[item.ProductID] = struct{}{}
	}

	var ids []uint
	for id := range productIDs {
		ids = append(ids, id)
	}

	inventories, err := s.inventoryService.GetMultiple(ctx, ids)
	if err != nil {
		s.logger.Errorw("error getting inventory", "error", err, "id", ids)
		return err
	}

	var lockedProducts []string
	for _, productID := range ids {
		redisKey := s.getLockProductKey(productID)
		lock := s.redisClient.SetNX(ctx, redisKey, "locked", 5*time.Second)
		if !lock.Val() {
			s.logger.Errorw("stock update in progress", "productID", productID)
			return fmt.Errorf("stock update in progress for product %d", productID)
		}
		lockedProducts = append(lockedProducts, redisKey)
	}

	defer func() {
		for _, key := range lockedProducts {
			s.redisClient.Del(ctx, key)
		}
	}()

	existingUpdates := map[uint]int{}
	for _, item := range existingOrder.Items {
		existingUpdates[item.ProductID] = item.Quantity
	}
	if err := s.inventoryService.IncreaseStockBulk(ctx, existingUpdates); err != nil {
		s.logger.Errorw("error increasing stock bulk", "error", err, "id", existingOrder.ID)
		return err
	}

	for _, item := range request.Items {
		if inventories[item.ProductID].Stock < item.Quantity {
			s.logger.Errorw("not enough stock for product %d", item.ProductID)
			return fmt.Errorf("not enough stock for product %d", item.ProductID)
		}
	}

	updatedUpdates := map[uint]int{}
	for _, item := range existingOrder.Items {
		updatedUpdates[item.ProductID] = item.Quantity
	}
	if err := s.inventoryService.DecreaseStockBulk(ctx, updatedUpdates); err != nil {
		s.logger.Errorw("error decreasing stock", "error", err, "id", request.ID)
		return err
	}

	var removedItemIDs []uint
	for _, item := range existingOrder.Items {
		found := false
		for _, updatedItem := range request.Items {
			if item.ProductID == updatedItem.ProductID {
				found = true
				break
			}
		}
		if !found {
			removedItemIDs = append(removedItemIDs, item.ID)
		}
	}

	if len(removedItemIDs) > 0 {
		if err := s.store.DeleteOrderItems(ctx, removedItemIDs); err != nil {
			s.logger.Errorw("error bulk deleting removed order items", "error", err, "orderItemIDs", removedItemIDs)
			return err
		}
	}

	existingOrder.Items = NewItems(request.Items)
	return s.store.Update(ctx, existingOrder)
}

func (s *service) Get(ctx context.Context, id uint) (*Order, error) {
	return s.store.Get(ctx, id)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	order, err := s.store.Get(ctx, id)
	if err != nil {
		s.logger.Errorw("failed to get order", "error", err, "id", id)
		return fmt.Errorf("order not found: %w", err)
	}

	productIDs := make([]uint, len(order.Items))
	for i, item := range order.Items {
		productIDs[i] = item.ProductID
	}

	var lockedProducts []string
	for _, productID := range productIDs {
		redisKey := s.getLockProductKey(productID)
		lock := s.redisClient.SetNX(ctx, redisKey, "locked", 5*time.Second)
		if !lock.Val() {
			s.logger.Errorw("stock update in progress", "productID", productID)
			return fmt.Errorf("stock update in progress for product %d", productID)
		}
		lockedProducts = append(lockedProducts, redisKey)
	}

	defer func() {
		for _, key := range lockedProducts {
			s.redisClient.Del(ctx, key)
		}
	}()

	updates := map[uint]int{}
	for _, item := range order.Items {
		updates[item.ProductID] = item.Quantity
	}
	if err := s.inventoryService.IncreaseStockBulk(ctx, updates); err != nil {
		s.logger.Errorw("error increasing stock", "error", err, "id", id)
		return err
	}

	var orderItemIDs []uint
	for _, item := range order.Items {
		orderItemIDs = append(orderItemIDs, item.ID)
	}

	if len(orderItemIDs) > 0 {
		if err := s.store.DeleteOrderItems(ctx, orderItemIDs); err != nil {
			s.logger.Errorw("error bulk deleting order items", "error", err, "orderItemIDs", orderItemIDs)
			return err
		}
	}

	return s.store.Delete(ctx, id)
}

func (s *service) getLockProductKey(productID uint) string {
	return fmt.Sprintf("stock_lock_product_%d", productID)
}
