package product

import "github.com/DuckLuckBreakout/web/backend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/web/backend/internal/pkg/product UseCase

type UseCase interface {
	GetProductById(productId uint64) (*models.Product, error)
	GetRangeProducts(paginator *models.PaginatorProducts) (*models.RangeProducts, error)
	SearchRangeProducts(searchQuery *models.SearchQuery) (*models.RangeProducts, error)
	GetProductRecommendationsById(productId uint64,
		paginator *models.PaginatorRecommendations) ([]*models.RecommendationProduct, error)
}
