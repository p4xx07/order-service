package order

import (
	"gorm.io/gorm"
	"time"
)

type IStore interface {
	List(startDate, endDate time.Time, name, description string) ([]Order, error)
	Get(orderID uint) (*Order, error)
	Create(order *Order) error
	Update(orderID uint, order *Order) error
	Delete(orderID uint) error
	BeginTransaction() *gorm.DB
	Transaction(f func(tx *gorm.DB) error) error
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) IStore {
	return &store{db: db}
}

func (s *store) List(startDate, endDate time.Time, name, description string) ([]Order, error) {
	var orders []Order
	query := s.db.Model(&Order{})

	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if description != "" {
		query = query.Where("description ILIKE ?", "%"+description+"%")
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *store) Get(orderID uint) (*Order, error) {
	var order Order
	err := s.db.Preload("Items").Where("id = ?", orderID).First(&order).Error
	return &order, err
}

func (s *store) Create(order *Order) error {
	return s.db.Create(order).Error
}

func (s *store) Update(orderID uint, order *Order) error {
	return s.db.Model(&Order{}).Where("id = ?", orderID).Updates(order).Error
}

func (s *store) Delete(orderID uint) error {
	return s.db.Delete(&Order{}, orderID).Error
}

func (s *store) BeginTransaction() *gorm.DB {
	return s.db.Begin()
}

func (s *store) Transaction(f func(tx *gorm.DB) error) error {
	return s.db.Transaction(f)
}
