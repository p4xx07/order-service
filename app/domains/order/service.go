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
	Get(ctx context.Context, orderID uint) (*OrderResponse, error)
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
			return nil, ErrStockUpdateInProgress
		}
		lockedProducts = append(lockedProducts, redisKey)

		if inventories[item.ProductID].Stock < item.Quantity {
			s.logger.Errorw("stock update in progress for product %d", item.ProductID)
			return nil, ErrNoStockAvailable
		}
	}

	defer func() {
		for _, key := range lockedProducts {
			s.redisClient.Del(ctx, key)
		}
	}()

	updates := map[uint]int{}
	var orderItems []OrderItem

	for _, item := range request.Items {
		updates[item.ProductID] = item.Quantity

		orderItems = append(orderItems, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     inventories[item.ProductID].Product.Price,
		})
	}
	if err := s.inventoryService.DecreaseStockBulk(ctx, updates); err != nil {
		s.logger.Errorw("failed to decrease stock bulk", "error", err)
		return nil, err
	}

	order := NewOrder(request.UserID, orderItems)
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
			return ErrStockUpdateInProgress
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

	for _, item := range request.Items {
		if inventories[item.ProductID].Stock+existingUpdates[item.ProductID] < item.Quantity {
			s.logger.Errorw("not enough stock for product %d", item.ProductID)
			return ErrNoStockAvailable
		}
	}

	if err := s.inventoryService.IncreaseStockBulk(ctx, existingUpdates); err != nil {
		s.logger.Errorw("error increasing stock bulk", "error", err, "id", existingOrder.ID)
		return err
	}

	updates := map[uint]int{}
	var orderItems []OrderItem

	for _, item := range request.Items {
		updates[item.ProductID] = item.Quantity

		orderItems = append(orderItems, OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     inventories[item.ProductID].Product.Price,
		})
	}
	if err := s.inventoryService.DecreaseStockBulk(ctx, updates); err != nil {
		s.logger.Errorw("error decreasing stock", "error", err, "id", request.ID)
		return err
	}

	toDelete := make([]uint, len(existingOrder.Items))
	for i, item := range existingOrder.Items {
		toDelete[i] = item.ID
	}

	if err := s.store.DeleteOrderItems(ctx, toDelete); err != nil {
		s.logger.Errorw("error deleting items", "error", err, "id", request.ID)
		return err
	}

	existingOrder.Items = orderItems
	return s.store.Update(ctx, existingOrder)
}

func (s *service) Get(ctx context.Context, id uint) (*OrderResponse, error) {
	order, err := s.store.Get(ctx, id)
	if err != nil {
		s.logger.Errorw("error getting order", "error", err, "id", id)
		return nil, err
	}
	return order.ToResponse(), nil
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
			return ErrStockUpdateInProgress
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
