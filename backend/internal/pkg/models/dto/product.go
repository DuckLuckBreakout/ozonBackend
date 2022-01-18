package dto

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"

type DtoProductPrice struct {
	Discount  int `json:"discount"`
	BaseCost  int `json:"base_cost"`
	TotalCost int `json:"total_cost"`
}

type DtoProduct struct {
	Id           uint64                       `json:"id"`
	Title        string                       `json:"title" valid:"minstringlength(3)"`
	Price        DtoProductPrice              `json:"price" valid:"notnull, json"`
	Rating       float32                      `json:"rating" valid:"float, range(0, 10)"`
	CountReviews uint64                       `json:"count_reviews"`
	Properties   string                       `json:"properties"`
	Description  string                       `json:"description" valid:"utfletter"`
	Category     uint64                       `json:"category"`
	CategoryPath []*usecase.CategoriesCatalog `json:"category_path" valid:"notnull"`
	Images       []string                     `json:"images" valid:"notnull"`
}

type DtoProductId struct {
	ProductId uint64 `json:"product_id"`
}

type DtoRecommendations struct {
	ProductId uint64 `json:"product_id"`
	Count     int    `json:"count"`
}

type DtoRecommendationProduct struct {
	Id           uint64          `json:"id"`
	Title        string          `json:"title" valid:"minstringlength(3)"`
	Price        DtoProductPrice `json:"price" valid:"notnull, json"`
	PreviewImage string          `json:"preview_image" valid:"minstringlength(3)"`
}

type DtoSearchPages struct {
	Category     uint64 `json:"category"`
	Count        int    `json:"count"`
	SearchString string `json:"search_string"`
	FilterString string `json:"filter_string"`
}

type DtoProductFilter struct {
	MinPrice   uint64 `json:"min_price"`
	MaxPrice   uint64 `json:"max_price"`
	IsNew      bool   `json:"is_new"`
	IsRating   bool   `json:"is_rating"`
	IsDiscount bool   `json:"is_discount"`
}

type DtoPaginatorProducts struct {
	PageNum       int               `json:"page_num"`
	Count         int               `json:"count"`
	Category      uint64            `json:"category"`
	Filter        *DtoProductFilter `json:"filter"`
	SortKey       string            `json:"sort_key" valid:"in(cost|rating|date|discount)"`
	SortDirection string            `json:"sort_direction" valid:"in(ASC|DESC)"`
}

type DtoRageProducts struct {
	SortString   string `json:"sort_string"`
	FilterString string `json:"filter_string"`
}

type DtoViewProduct struct {
	Id           uint64          `json:"id"`
	Title        string          `json:"title" valid:"minstringlength(3)"`
	Price        DtoProductPrice `json:"price" valid:"notnull, json"`
	Rating       float32         `json:"rating" valid:"float, range(0, 10)"`
	CountReviews uint64          `json:"count_reviews"`
	PreviewImage string          `json:"preview_image" valid:"minstringlength(3)"`
}

type DtoSearchQuery struct {
	QueryString   string            `json:"query_string" valid:"minstringlength(2)"`
	PageNum       int               `json:"page_num"`
	Count         int               `json:"count"`
	Category      uint64            `json:"category"`
	Filter        *DtoProductFilter `json:"filter"`
	SortKey       string            `json:"sort_key" valid:"in(cost|rating|date|discount)"`
	SortDirection string            `json:"sort_direction" valid:"in(ASC|DESC)"`
}
