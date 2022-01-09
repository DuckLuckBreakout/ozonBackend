package session

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/session Handler

type Handler interface {
	CheckSession(w http.ResponseWriter, r *http.Request)
}
