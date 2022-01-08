package dto

type DtoFavoriteProduct struct {
	ProductId uint64 `json:"product_id"`
	UserId    uint64 `json:"user_id"`
}

type DtoCountPages struct {
	Count  int    `json:"count"`
	UserId uint64 `json:"user_id"`
}

type DtoCounter struct {
	Count int `json:"count"`
}

type DtoUserId struct {
	Id uint64 `json:"id"`
}

type DtoUserFavorites struct {
	Products []uint64 `json:"products"`
}
