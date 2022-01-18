package dto

type DtoPromoCode struct {
	Code string `json:"code"`
}

type DtoPromoProduct struct {
	ProductId uint64 `json:"product_id"`
	PromoCode string `json:"promo_code"`
}

type DtoPromoPrice struct {
	TotalCost int `json:"total_cost"`
	BaseCost  int `json:"base_cost"`
}
