package order

type PostRequest struct {
	UserID uint               `json:"user_id,omitempty" validate:"min=1,nonnil" required:"true"`
	Items  []OrderItemRequest `json:"items,omitempty" validate:"min=1,nonnil" required:"true"`
}

type PutRequest struct {
	ID    uint               `json:"id,omitempty" validate:"min=1,nonnil" required:"true"`
	Items []OrderItemRequest `json:"items,omitempty" validate:"min=1,nonnil" required:"true"`
}

type OrderItemRequest struct {
	ProductID uint `json:"product_id,omitempty" validate:"min=1,nonnil" required:"true"`
	Quantity  int  `json:"quantity,omitempty"`
}

func (o OrderItemRequest) ToStore(orderID uint) OrderItem {
	return OrderItem{
		OrderID:   orderID,
		ProductID: o.ProductID,
		Quantity:  o.Quantity,
	}
}
