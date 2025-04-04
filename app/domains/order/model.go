package order

import (
	"github.com/p4xx07/order-service/app/domains/product"
	"time"
)

type Order struct {
	ID        uint        `gorm:"primaryKey;autoIncrement"`
	UserID    uint        `gorm:"index"`
	Status    string      `gorm:"type:varchar(20);default:'pending'"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	UpdatedAt time.Time   `gorm:"autoUpdateTime"`
	Items     []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

func NewOrder(userID uint, items []OrderItem) *Order {
	return &Order{
		UserID:    userID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Items:     items,
	}
}

type OrderItem struct {
	ID        uint            `gorm:"primaryKey;autoIncrement"`
	OrderID   uint            `gorm:"index"`
	ProductID uint            `gorm:"index"`
	Quantity  int             `gorm:"type:int;not null"`
	Price     float64         `gorm:"type:decimal(10,2);not null"`
	Product   product.Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type OrderMeilisearch struct {
	Order
	CreatedAtTimestamp int64 `gorm:"autoCreateTime"`
}

func (o *Order) toDocument() OrderMeilisearch {
	if o == nil {
		return OrderMeilisearch{}
	}

	return OrderMeilisearch{
		Order:              *o,
		CreatedAtTimestamp: o.CreatedAt.UnixMilli(),
	}
}
