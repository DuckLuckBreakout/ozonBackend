package product

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/product UseCase

type UseCase interface {
	GetProductById(productId *usecase.ProductId) (*usecase.Product, error)
	GetRangeProducts(paginator *usecase.PaginatorProducts) (*usecase.RangeProducts, error)
	SearchRangeProducts(searchQuery *usecase.SearchQuery) (*usecase.RangeProducts, error)
	GetProductRecommendationsById(productId *usecase.ProductId, paginator *usecase.PaginatorRecommendations) ([]*usecase.RecommendationProduct, error)
}
