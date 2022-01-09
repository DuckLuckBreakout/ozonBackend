package models

import "github.com/DuckLuckBreakout/web/backend/internal/server/tools/sanitizer"

type UpdateOrder struct {
	OrderId uint64 `json:"order_id"`
	Status  string `json:"status" valid:"in(в пути|оформлен|получен)"`
}

func (u *UpdateOrder) Sanitize() {
	sanitizer := sanitizer.NewSanitizer()
	u.Status = sanitizer.Sanitize(u.Status)
}
