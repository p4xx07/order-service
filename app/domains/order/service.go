package order

import (
	"errors"
	"github.com/p4xx07/order-service/app/domains/inventory"
	"github.com/p4xx07/order-service/configuration"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type IService interface {
	Create(request postRequest) (*createOrderResponse, error)
	Update(orderID uint, request putRequest) error
	List(from, to time.Time, name, description string, limit, offset int64) ([]OrderResponse, error)
	Get(orderID uint) (*OrderResponse, error)
	Delete(orderID uint) error
}

type service struct {
	configuration    *configuration.Configuration
	logger           *zap.SugaredLogger
	store            IStore
	inventoryService inventory.IService
}

func NewService(configuration *configuration.Configuration, logger *zap.SugaredLogger, store IStore, inventoryService inventory.IService) IService {
	return &service{configuration: configuration, logger: logger, store: store, inventoryService: inventoryService}
}

func (s *service) List(from, to time.Time, name, description string, limit, offset int64) ([]OrderResponse, error) {
	orders, err := s.store.List(from, to, name, description)
	if err != nil {
		return nil, err
	}

	response := ToArrayResponse(orders)

	return response, nil
}

func (s *service) Get(orderID uint) (*OrderResponse, error) {
	order, err := s.store.Get(orderID)
	if err != nil {
		s.logger.Errorw("error getting order", "error", err, "orderID", orderID)
		return nil, err
	}

	response := ToResponse(*order)

	return &response, nil
}

func (s *service) Delete(orderID uint) error {
	return s.store.
		Transaction(func(tx *gorm.DB) error {
			var order Order
			if err := tx.Preload("Items").First(&order, orderID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.logger.Warnw("order not found", "orderID", orderID)
					return errors.New("order not found")
				}
				s.logger.Errorw("error deleting order", "error", err, "orderID", orderID)
				return err
			}

			for _, item := range order.Items {
				if err := s.inventoryService.AdjustStock(item.ProductID, item.Quantity); err != nil {
					s.logger.Errorw("error adjusting stock", "error", err, "orderID", orderID)
					return err
				}
			}

			if err := tx.Where("order_id = ?", order.ID).Delete(&OrderItem{}).Error; err != nil {
				s.logger.Errorw("error deleting items", "error", err, "orderID", orderID)
				return err
			}

			if err := tx.Delete(&order).Error; err != nil {
				s.logger.Errorw("error deleting order", "error", err, "orderID", orderID)
				return err
			}

			return nil
		})
}

func (s *service) Create(request postRequest) (*createOrderResponse, error) {
	tx := s.store.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range request.Items {
		if err := s.inventoryService.AdjustStock(item.ProductID, -item.Quantity); err != nil {
			tx.Rollback()
			s.logger.Errorw("error adjusting stock", "error", err)
			return nil, err
		}
	}

	items := NewItems(request.Items)
	order := NewOrder(request.UserID, items)
	if err := s.store.Create(order); err != nil {
		tx.Rollback()
		s.logger.Errorw("error creating order", "error", err)
		return nil, err
	}

	commit := tx.Commit()
	if err := commit.Error; err != nil {
		s.logger.Errorw("error committing order", "error", err)
		return nil, err
	}

	return &createOrderResponse{OrderID: order.ID}, nil
}

func (s *service) Update(orderID uint, request putRequest) error {
	return s.store.
		Transaction(func(tx *gorm.DB) error {
			var existingOrder Order
			if err := tx.Preload("Items").First(&existingOrder, orderID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.logger.Warnw("order not found", "orderID", orderID)
					return errors.New("order not found")
				}
				s.logger.Errorw("error getting order", "error", err, "orderID", orderID)
				return err
			}

			if err := tx.Model(&existingOrder).Update("updated_at", time.Now().UTC()).Error; err != nil {
				s.logger.Errorw("error updating order", "error", err)
				return err
			}

			updatedItems := make([]OrderItem, len(request.UpdatedItems))
			for i, item := range request.UpdatedItems {
				updatedItems[i] = item.ToStore(existingOrder.ID)
			}
			if err := s.updateOrderItems(tx, &existingOrder, updatedItems); err != nil {
				s.logger.Errorw("error updating order", "error", err)
				return err
			}

			return nil
		})
}

func (s *service) updateOrderItems(tx *gorm.DB, existingOrder *Order, updatedItems []OrderItem) error {
	existingMap := make(map[uint]OrderItem)
	for _, item := range existingOrder.Items {
		existingMap[item.ProductID] = item
	}

	processed := make(map[uint]bool)

	for _, updatedItem := range updatedItems {
		existingItem, exists := existingMap[updatedItem.ProductID]
		quantityChange := updatedItem.Quantity

		if exists {
			quantityChange = updatedItem.Quantity - existingItem.Quantity
		}

		if err := s.inventoryService.AdjustStock(updatedItem.ProductID, -quantityChange); err != nil {
			s.logger.Errorw("error adjusting stock", err)
			return err
		}

		processed[updatedItem.ProductID] = true
	}

	for productID, oldItem := range existingMap {
		if !processed[productID] {
			if err := s.inventoryService.AdjustStock(productID, +oldItem.Quantity); err != nil {
				s.logger.Errorw("error adjusting stock", err)
				return err
			}
		}
	}

	if err := tx.Where("order_id = ?", existingOrder.ID).Delete(&OrderItem{}).Error; err != nil {
		s.logger.Errorw("error deleting order", err)
		return err
	}
	if len(updatedItems) > 0 {
		if err := tx.Create(&updatedItems).Error; err != nil {
			s.logger.Errorw("error deleting order", err)
			return err
		}
	}

	return nil
}
