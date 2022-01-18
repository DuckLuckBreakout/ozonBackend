package product

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/product Repository

type Repository interface {
	SelectProductById(productId *dto.DtoProductId) (*dto.DtoProduct, error)
	SelectRecommendationsByReviews(rec *dto.DtoRecommendations) ([]*dto.DtoRecommendationProduct, error)
	GetCountPages(cntPages *dto.DtoCountPages) (int, error)
	GetCountSearchPages(srcPages *dto.DtoSearchPages) (int, error)
	CreateSortString(sortStr *dto.DtoSortString) (string, error)
	CreateFilterString(filter *dto.DtoProductFilter) string
	SelectRangeProducts(
		paginator *dto.DtoPaginatorProducts,
		rageProducts *dto.DtoRageProducts,
	) ([]*dto.DtoViewProduct, error)
	SearchRangeProducts(
		searchQuery *dto.DtoSearchQuery,
		rageProducts *dto.DtoRageProducts,
	) ([]*dto.DtoViewProduct, error)
}
