package dto

import "time"

type DtoUpdateOrder struct {
	OrderId uint64 `json:"order_id"`
	Status  string `json:"status"`
}

type DtoChangedOrder struct {
	Number string `json:"number"`
	UserId uint64 `json:"user_id"`
}

type DtoRangeOrders struct {
	UserId     uint64 `json:"user_id"`
	SortString string `json:"sort_string"`
}

type DtoSortOrdersOptions struct {
	SortKey       string `json:"sort_key" valid:"in(date)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DESC)"`
}

type DtoPaginatorOrders struct {
	PageNum int `json:"page_num"`
	Count   int `json:"count"`
	DtoSortOrdersOptions
}

type DtoOrderAddress struct {
	Address string `json:"address" valid:"utfletter, stringlength(1|30)"`
}

type DtoPreviewOrderedProducts struct {
	Id           uint64 `json:"id"`
	PreviewImage string `json:"preview_image" valid:"minstringlength(1)"`
}

type DtoOrderNumber struct {
	Number string `json:"number"`
}

type DtoSortString struct {
	SortKey       string `json:"sort_key"`
	SortDirection string `json:"sort_direction"`
}

type DtoOrderId struct {
	Id uint64 `json:"id"`
}

type DtoCartProductPrice struct {
	Discount  int `json:"discount"`
	BaseCost  int `json:"base_cost"`
	TotalCost int `json:"total_cost"`
}

type DtoPreviewCartArticle struct {
	Id           uint64              `json:"id"`
	Title        string              `json:"title" valid:"minstringlength(1)"`
	Price        DtoCartProductPrice `json:"price" valid:"notnull"`
	PreviewImage string              `json:"preview_image" valid:"minstringlength(1)"`
	Count        uint64              `json:"count"`
}

type DtoPlacedOrder struct {
	Id           uint64                       `json:"id"`
	Address      DtoOrderAddress              `json:"address" valid:"notnull"`
	TotalCost    int                          `json:"total_cost"`
	Products     []*DtoPreviewOrderedProducts `json:"product_images" valid:"notnull"`
	DateAdded    time.Time                    `json:"date_added"`
	DateDelivery time.Time                    `json:"date_delivery"`
	OrderNumber  DtoOrderNumber               `json:"order_number"`
	Status       string                       `json:"status"`
}

type DtoOrderCountPages struct {
	UserId            uint64 `json:"user_id"`
	CountOrdersOnPage int    `json:"count_orders_on_page"`
}

type DtoOrderRecipient struct {
	FirstName string `json:"first_name" valid:"utfletter, stringlength(1|30)"`
	LastName  string `json:"last_name" valid:"utfletter, stringlength(1|30)"`
	Email     string `json:"email" valid:"email"`
}

type DtoOrder struct {
	Recipient DtoOrderRecipient `json:"recipient" valid:"notnull"`
	Address   DtoOrderAddress   `json:"address" valid:"notnull"`
	PromoCode string            `json:"promo_code"`
}

type DtoTotalPrice struct {
	TotalDiscount int `json:"total_discount"`
	TotalCost     int `json:"total_cost"`
	TotalBaseCost int `json:"total_base_cost"`
}
