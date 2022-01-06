package favorites

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites UseCase

type UseCase interface {
	AddProductToFavorites(favorite *models.FavoriteProduct) error
	DeleteProductFromFavorites(favorite *models.FavoriteProduct) error
	GetRangeFavorites(paginator *models.PaginatorFavorite) (*models.RangeFavorites, error)
	GetUserFavorites(userId *models.UserId) (*models.UserFavorites, error)
}
