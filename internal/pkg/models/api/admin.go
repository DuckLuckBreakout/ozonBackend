package api

import "github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/sanitizer"

type ApiUpdateOrder struct {
	OrderId uint64 `json:"order_id"`
	Status  string `json:"status" valid:"in(в пути|оформлен|получен)"`
}

func (u *ApiUpdateOrder) Sanitize() {
	sanitize := sanitizer.NewSanitizer()
	u.Status = sanitize.Sanitize(u.Status)
}
