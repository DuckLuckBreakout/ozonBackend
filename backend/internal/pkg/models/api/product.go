package api

type ApiPaginatorRecommendations struct {
	Count int `json:"count"`
}

type ApiProductFilter struct {
	MinPrice   uint64 `json:"min_price"`
	MaxPrice   uint64 `json:"max_price"`
	IsNew      bool   `json:"is_new"`
	IsRating   bool   `json:"is_rating"`
	IsDiscount bool   `json:"is_discount"`
}

type ApiPaginatorProducts struct {
	PageNum       int               `json:"page_num"`
	Count         int               `json:"count"`
	Category      uint64            `json:"category"`
	Filter        *ApiProductFilter `json:"filter"`
	SortKey       string            `json:"sort_key" valid:"in(cost|rating|date|discount)"`
	SortDirection string            `json:"sort_direction" valid:"in(ASC|DESC)"`
}

type ApiSearchQuery struct {
	QueryString   string            `json:"query_string" valid:"minstringlength(2)"`
	PageNum       int               `json:"page_num"`
	Count         int               `json:"count"`
	Category      uint64            `json:"category"`
	Filter        *ApiProductFilter `json:"filter"`
	SortKey       string            `json:"sort_key" valid:"in(cost|rating|date|discount)"`
	SortDirection string            `json:"sort_direction" valid:"in(ASC|DESC)"`
}
