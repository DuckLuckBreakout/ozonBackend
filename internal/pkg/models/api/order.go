package api

import "github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/sanitizer"

type ApiPaginatorOrders struct {
	PageNum       int    `json:"page_num"`
	Count         int    `json:"count"`
	SortKey       string `json:"sort_key" valid:"in(date)"`
	SortDirection string `json:"sort_direction" valid:"in(ASC|DESC)"`
}

type ApiOrderRecipient struct {
	FirstName string `json:"first_name" valid:"utfletter, stringlength(1|30)"`
	LastName  string `json:"last_name" valid:"utfletter, stringlength(1|30)"`
	Email     string `json:"email" valid:"email"`
}

type ApiOrderAddress struct {
	Address string `json:"address" valid:"utfletter, stringlength(1|30)"`
}

type ApiOrder struct {
	Recipient ApiOrderRecipient `json:"recipient" valid:"notnull"`
	Address   ApiOrderAddress   `json:"address" valid:"notnull"`
	PromoCode string            `json:"promo_code"`
}

func (o *ApiOrder) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	o.Recipient.FirstName = sanitizer.Sanitize(o.Recipient.FirstName)
	o.Recipient.LastName = sanitizer.Sanitize(o.Recipient.LastName)
	o.Recipient.Email = sanitizer.Sanitize(o.Recipient.Email)
	o.Address.Address = sanitizer.Sanitize(o.Address.Address)
}
