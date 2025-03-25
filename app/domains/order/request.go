package order

type postRequest struct {
	UserID uint               `json:"user_id,omitempty" validate:"min=1,nonnil" required:"true"`
	Items  []orderItemRequest `json:"items,omitempty" validate:"min=1,nonnil" required:"true"`
}

type putRequest struct {
	UpdatedItems []orderItemRequest `json:"updated_items,omitempty" validate:"min=1,nonnil" required:"true"`
}

type orderItemRequest struct {
	ProductID uint `json:"product_id,omitempty"`
	Quantity  int  `json:"quantity,omitempty"`
}

func (o orderItemRequest) ToStore(orderID uint) OrderItem {
	return OrderItem{
		OrderID:   orderID,
		ProductID: o.ProductID,
		Quantity:  o.Quantity,
	}
}
