package usecase

type UpdateOrder struct {
	OrderId uint64 `json:"order_id"`
	Status  string `json:"status"`
}
