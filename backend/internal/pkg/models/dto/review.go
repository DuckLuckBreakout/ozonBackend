package dto

import "time"

type DtoRangeReviews struct {
	ProductId  uint64 `json:"product_id"`
	SortString string `json:"sort_string"`
}

type DtoPaginatorReviews struct {
	PageNum       int    `json:"page_num"`
	Count         int    `json:"count"`
	SortKey       string `json:"sort_key" valid:"in(rating|date)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DESC)"`
}

type DtoViewReview struct {
	UserName      string    `json:"user_name" valid:"minstringlength(1)"`
	UserAvatar    string    `json:"user_avatar" valid:"minstringlength(1)"`
	DateAdded     time.Time `json:"date_added" valid:"notnull"`
	Rating        int       `json:"rating" valid:"int"`
	Advantages    string    `json:"advantages"`
	Disadvantages string    `json:"disadvantages"`
	Comment       string    `json:"comment"`
	IsPublic      bool      `json:"-"`
	UserId        int       `json:"-"`
}

type DtoReviewStatistics struct {
	Stars []int `json:"stars"`
}

type DtoCheckReview struct {
	UserId    uint64 `json:"user_id"`
	ProductId uint64 `json:"product_id"`
}

type DtoReviewId struct {
	Id uint64 `json:"id"`
}

type DtoReview struct {
	ProductId     int    `json:"product_id"`
	Rating        int    `json:"rating" valid:"int"`
	Advantages    string `json:"advantages"`
	Disadvantages string `json:"disadvantages"`
	Comment       string `json:"comment"`
	IsPublic      bool   `json:"is_public"`
}
