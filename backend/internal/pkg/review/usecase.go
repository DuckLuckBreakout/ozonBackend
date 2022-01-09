package review

import "github.com/DuckLuckBreakout/web/backend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/web/backend/internal/pkg/review UseCase

type UseCase interface {
	GetStatisticsByProductId(productId uint64) (*models.ReviewStatistics, error)
	CheckReviewUserRights(userId uint64, productId uint64) error
	AddNewReviewForProduct(userId uint64, review *models.Review) error
	GetReviewsByProductId(productId uint64, paginator *models.PaginatorReviews) (*models.RangeReviews, error)
}
