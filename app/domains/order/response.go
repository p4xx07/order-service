package order

import (
	"github.com/p4xx07/order-service/app/domains/product"
	"time"
)

type CreateOrderResponse struct {
	ID uint `json:"id"`
}

type OrderResponse struct {
	ID        uint                `json:"id,omitempty"`
	UserID    uint                `json:"user_id,omitempty"`
	Status    string              `json:"status,omitempty"`
	CreatedAt time.Time           `json:"created_at"`
	Items     []OrderItemResponse `json:"items,omitempty"`
}

func (o *Order) ToResponse() *OrderResponse {
	items := make([]OrderItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = item.ToResponse()
	}
	return &OrderResponse{
		ID:        o.ID,
		UserID:    o.UserID,
		Status:    o.Status,
		CreatedAt: o.CreatedAt,
		Items:     items,
	}
}

type OrderItemResponse struct {
	Quantity int                     `json:"quantity,omitempty"`
	Price    float64                 `json:"price,omitempty"`
	Product  product.ProductResponse `json:"product"`
}

func (o *OrderItem) ToResponse() OrderItemResponse {
	return OrderItemResponse{
		Quantity: o.Quantity,
		Price:    o.Price,
		Product:  o.Product.ToResponse(),
	}
}
