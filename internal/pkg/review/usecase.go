package review

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/review UseCase

type UseCase interface {
	GetStatisticsByProductId(productId *usecase.ProductId) (*usecase.ReviewStatistics, error)
	CheckReviewUserRights(userId *usecase.UserId, productId *usecase.ProductId) error
	AddNewReviewForProduct(userId *usecase.UserId, review *usecase.Review) error
	GetReviewsByProductId(productId *usecase.ProductId, paginator *usecase.PaginatorReviews) (*usecase.RangeReviews, error)
}
