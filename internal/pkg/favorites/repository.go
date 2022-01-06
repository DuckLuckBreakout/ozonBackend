package favorites

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites/repository"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites Repository

type Repository interface {
	AddProductToFavorites(favorite *repository.DtoFavoriteProduct) error
	DeleteProductFromFavorites(favorite *repository.DtoFavoriteProduct) error
	GetCountPages(countPages *repository.DtoCountPages) (*repository.DtoCounter, error)
	CreateSortString(sortKey, sortDirection string) (string, error)
	SelectRangeFavorites(paginator *models.PaginatorFavorites, sortString string, userId uint64) ([]*models.ViewFavorite, error)
	GetUserFavorites(userId *repository.DtoUserId) (*repository.DtoUserFavorites, error)
}
