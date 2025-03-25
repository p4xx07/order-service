package order

import "time"

type OrderResponse struct {
	ID        uint        `json:"id,omitempty"`
	UserID    uint        `json:"user_id,omitempty"`
	Status    string      `json:"status,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Items     []OrderItem `json:"items,omitempty"`
}

func ToResponse(order Order) OrderResponse {
	return OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		CreatedAt: order.CreatedAt,
		Status:    order.Status,
		Items:     order.Items,
	}
}

func ToArrayResponse(orders []Order) []OrderResponse {
	responses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = ToResponse(order)
	}
	return responses
}
