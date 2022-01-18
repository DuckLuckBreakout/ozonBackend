package favorites

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites UseCase

type UseCase interface {
	AddProductToFavorites(favorite *usecase.FavoriteProduct) error
	DeleteProductFromFavorites(favorite *usecase.FavoriteProduct) error
	GetRangeFavorites(paginator *usecase.PaginatorFavorite) (*usecase.RangeFavorites, error)
	GetUserFavorites(userId *usecase.UserId) (*usecase.UserFavorites, error)
}
