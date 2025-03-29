package inventory

import "github.com/p4xx07/order-service/app/domains/product"

type Inventory struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	ProductID uint `gorm:"index;unique"`
	Stock     int  `gorm:"type:int;not null"`

	Product product.Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
