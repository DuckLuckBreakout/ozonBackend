package favorites

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites UseCase

type UseCase interface {
	AddProductToFavorites(productId, userId uint64) error
	DeleteProductFromFavorites(productId, userId uint64) error
	GetRangeFavorites(paginator *models.PaginatorFavorites,
		userId uint64) (*models.RangeFavorites, error)
	GetUserFavorites(userId uint64) (*models.UserFavorites, error)
}
