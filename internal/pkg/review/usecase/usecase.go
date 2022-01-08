package usecase

import (
	"fmt"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/review"
	reviewRepo "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/review/repository"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user"
	userRepo "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user/repository"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
)

type ReviewUseCase struct {
	ReviewRepo review.Repository
	UserRepo   user.Repository
}

func NewUseCase(reviewRepo review.Repository, userRepo user.Repository) review.UseCase {
	return &ReviewUseCase{
		ReviewRepo: reviewRepo,
		UserRepo:   userRepo,
	}
}

func (u *ReviewUseCase) GetStatisticsByProductId(productId *usecase.ProductId) (*usecase.ReviewStatistics, error) {
	reviews, err := u.ReviewRepo.SelectStatisticsByProductId(&reviewRepo.DtoProductId{Id: productId.Id})
	if err != nil {
		return nil, err
	}
	return &usecase.ReviewStatistics{
		Stars: reviews.Stars,
	}, nil
}

func (u *ReviewUseCase) CheckReviewUserRights(userId *usecase.UserId, productId *usecase.ProductId) error {
	rights := u.ReviewRepo.CheckReview(&reviewRepo.DtoCheckReview{
		UserId:    userId.Id,
		ProductId: productId.Id,
	})
	if !rights {
		return errors.ErrNoWriteRights
	}

	return nil
}

func (u *ReviewUseCase) AddNewReviewForProduct(userId *usecase.UserId, review *usecase.Review) error {
	rights := u.ReviewRepo.CheckReview(&reviewRepo.DtoCheckReview{
		UserId:    userId.Id,
		ProductId: uint64(review.ProductId),
	})
	if !rights {
		return errors.ErrNoWriteRights
	}

	_, err := u.ReviewRepo.AddReview(
		&reviewRepo.DtoUserId{Id: userId.Id},
		&reviewRepo.DtoReview{
			ProductId:     review.ProductId,
			Rating:        review.Rating,
			Advantages:    review.Advantages,
			Disadvantages: review.Disadvantages,
			Comment:       review.Comment,
			IsPublic:      review.IsPublic,
		},
	)
	if err != nil {
		return errors.ErrCanNotAddReview
	}
	return nil
}

func (u *ReviewUseCase) GetReviewsByProductId(
	productId *usecase.ProductId,
	paginator *usecase.PaginatorReviews,
) (*usecase.RangeReviews, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	// Max count pages in catalog
	countPages, err := u.ReviewRepo.GetCountPages(&reviewRepo.DtoCountPages{
		ProductId:         productId.Id,
		CountOrdersOnPage: paginator.Count,
	})
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Keys for sort reviews in catalog
	sortString, err := u.ReviewRepo.CreateSortString(&reviewRepo.DtoSortString{
		SortKey:       paginator.SortKey,
		SortDirection: paginator.SortDirection,
	})
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get range of reviews
	reviews, err := u.ReviewRepo.SelectRangeReviews(
		&reviewRepo.DtoRangeReviews{
			ProductId:  productId.Id,
			SortString: sortString,
		},
		&reviewRepo.DtoPaginatorReviews{
			PageNum:       paginator.PageNum,
			Count:         paginator.Count,
			SortKey:       paginator.SortKey,
			SortDirection: paginator.SortDirection,
		},
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get user data for review
	for _, userReview := range reviews {
		userInfo, err := u.UserRepo.SelectProfileById(&userRepo.DtoUserId{Id: uint64(userReview.UserId)})
		if err != nil {
			return nil, errors.ErrInternalError
		}

		if userReview.IsPublic {
			userReview.UserAvatar = userInfo.Avatar.Url
			userReview.UserName = fmt.Sprintf("%s %s", userInfo.FirstName, userInfo.LastName)
		}
	}

	var listReviews []*usecase.ViewReview
	for _, item := range reviews {
		listReviews = append(listReviews, &usecase.ViewReview{
			UserName:      item.UserName,
			UserAvatar:    item.UserAvatar,
			DateAdded:     item.DateAdded,
			Rating:        item.Rating,
			Advantages:    item.Advantages,
			Disadvantages: item.Disadvantages,
			Comment:       item.Comment,
			IsPublic:      item.IsPublic,
			UserId:        item.UserId,
		})
	}

	return &usecase.RangeReviews{
		ListPreviews:  listReviews,
		MaxCountPages: countPages.Count,
	}, nil
}
