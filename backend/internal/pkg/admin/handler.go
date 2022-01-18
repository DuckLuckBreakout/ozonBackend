package admin

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/admin Handler

type Handler interface {
	ChangeOrderStatus(w http.ResponseWriter, r *http.Request)
}
