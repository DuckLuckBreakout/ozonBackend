package favorites

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites Repository

type Repository interface {
	AddProductToFavorites(favorite *dto.DtoFavoriteProduct) error
	DeleteProductFromFavorites(favorite *dto.DtoFavoriteProduct) error
	GetCountPages(countPages *dto.DtoCountPages) (*dto.DtoCounter, error)
	CreateSortString(sortKey, sortDirection string) (string, error)
	SelectRangeFavorites(paginator *usecase.PaginatorFavorites, sortString string, userId uint64) ([]*usecase.ViewFavorite, error)
	GetUserFavorites(userId *dto.DtoUserId) (*dto.DtoUserFavorites, error)
}
