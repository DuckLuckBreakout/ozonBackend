package review

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/review Repository

type Repository interface {
	SelectRangeReviews(
		rangeReviews *dto.DtoRangeReviews,
		paginator *dto.DtoPaginatorReviews,
	) ([]*dto.DtoViewReview, error)
	GetCountPages(countPages *dto.DtoCountPages) (*dto.DtoCounter, error)
	CreateSortString(sortString *dto.DtoSortString) (string, error)
	SelectStatisticsByProductId(productId *dto.DtoProductId) (*dto.DtoReviewStatistics, error)
	CheckReview(review *dto.DtoCheckReview) bool
	AddReview(userId *dto.DtoUserId, review *dto.DtoReview) (*dto.DtoReviewId, error)
}
