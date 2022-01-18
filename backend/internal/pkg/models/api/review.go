package api

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/sanitizer"
)

type ApiReviewStatistics struct {
	Stars []int `json:"stars"`
}

type ApiRangeReviews struct {
	ListPreviews  []*usecase.ViewReview `json:"list_reviews" valid:"notnull"`
	MaxCountPages int                   `json:"max_count_pages"`
}

type ApiReview struct {
	ProductId     int    `json:"product_id"`
	Rating        int    `json:"rating" valid:"int"`
	Advantages    string `json:"advantages"`
	Disadvantages string `json:"disadvantages"`
	Comment       string `json:"comment"`
	IsPublic      bool   `json:"is_public"`
}

func (r *ApiReview) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	r.Advantages = sanitizer.Sanitize(r.Advantages)
	r.Disadvantages = sanitizer.Sanitize(r.Disadvantages)
	r.Comment = sanitizer.Sanitize(r.Comment)
}

type ApiPaginatorReviews struct {
	PageNum       int    `json:"page_num"`
	Count         int    `json:"count"`
	SortKey       string `json:"sort_key" valid:"in(rating|date)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DESC)"`
}
