package order

type postRequest struct {
	UserID uint        `json:"user_id,omitempty" validate:"min=0,nonnil" required:"true"`
	Items  []OrderItem `json:"items,omitempty" validate:"min=0,nonnil" required:"true"`
}

type putRequest struct {
	UpdatedOrder Order       `json:"updated_order" validate:"min=0,nonnil" required:"true"`
	UpdatedItems []OrderItem `json:"updated_items,omitempty" validate:"min=0,nonnil" required:"true"`
}
