package favorites

import "net/http"

//go:generate mockgen -destination=./mock/mock_handler.go -package=mock github.com/DuckLuckBreakout/web/backend/internal/pkg/favorites Handler

type Handler interface {
	AddProductToFavorites(w http.ResponseWriter, r *http.Request)
	DeleteProductFromFavorites(w http.ResponseWriter, r *http.Request)
	GetListPreviewFavorites(w http.ResponseWriter, r *http.Request)
	GetUserFavorites(w http.ResponseWriter, r *http.Request)
}
