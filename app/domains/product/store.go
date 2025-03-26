package product

import (
	"gorm.io/gorm"
	"time"
)

type IStore interface {
	List(from time.Time, to time.Time, name string, description string, limit int64, offset int64) ([]Product, error)
}

type store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) IStore {
	return &store{db: db}
}

func (s *store) List(from time.Time, to time.Time, name string, description string, limit int, offset int) ([]Product, error) {
	var products []Product

	query := s.db.Where("created_at BETWEEN ? AND ?", from, to)

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if description != "" {
		query = query.Where("description ILIKE ?", "%"+description+"%")
	}

	err := query.Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}
