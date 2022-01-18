package api

type ApiPromoCodeGroup struct {
	Products  []uint64 `json:"products" valid:"notnull"`
	PromoCode string   `json:"promo_code" valid:"stringlength(1|30)"`
}

type ApiDiscountedPrice struct {
	TotalDiscount int `json:"total_discount"`
	TotalCost     int `json:"total_cost"`
	TotalBaseCost int `json:"total_base_cost"`
}
