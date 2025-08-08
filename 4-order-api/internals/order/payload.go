package order

type OrderCreateRequest struct {
	UserID     uint        `json:"user_id"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderUpdateRequest struct {
	UserID     uint        `json:"user_id"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderGetAllResponse struct {
	ID         uint        `json:"id"`
	UserID     uint        `json:"user_id"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderGetOneResponse struct {
	ID         uint        `json:"id"`
	UserID     uint        `json:"user_id"`
	OrderItems []OrderItem `json:"order_items"`
}
