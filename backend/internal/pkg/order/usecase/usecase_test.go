package usecase

import (
	"testing"

	"github.com/DuckLuckBreakout/web/backend/internal/pkg/models"
	order_repo "github.com/DuckLuckBreakout/web/backend/internal/pkg/order/mock"
	product_repo "github.com/DuckLuckBreakout/web/backend/internal/pkg/product/mock"
	promo_code_repo "github.com/DuckLuckBreakout/web/backend/internal/pkg/promo_code/mock"
	user_repo "github.com/DuckLuckBreakout/web/backend/internal/pkg/user/mock"
	"github.com/DuckLuckBreakout/web/backend/internal/server/errors"
	cart_service "github.com/DuckLuckBreakout/web/backend/services/cart/proto/cart"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOrderUseCase_GetSubCategoriesById(t *testing.T) {
	userId := uint64(12)
	previewCart := models.PreviewCart{
		Products: []*models.PreviewCartArticle{
			{},
		},
		Price: models.TotalPrice{},
	}
	userProfile := models.ProfileUser{}

	t.Run("GetPreviewOrder_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := order_repo.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)
		promoCodeRepo := promo_code_repo.NewMockRepository(ctrl)

		userRepo := user_repo.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(userId).
			Return(&userProfile, nil)

		orderUCase := NewUseCase(orderRepo, cartClient, productRepo, userRepo, promoCodeRepo)

		_, err := orderUCase.GetPreviewOrder(userId, &previewCart)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetPreviewOrder_not_found_user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := order_repo.NewMockRepository(ctrl)
		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)
		promoCodeRepo := promo_code_repo.NewMockRepository(ctrl)

		userRepo := user_repo.NewMockRepository(ctrl)
		userRepo.
			EXPECT().
			SelectProfileById(userId).
			Return(&userProfile, errors.ErrInternalError)

		orderUCase := NewUseCase(orderRepo, cartClient, productRepo, userRepo, promoCodeRepo)

		_, err := orderUCase.GetPreviewOrder(userId, &previewCart)
		assert.Error(t, err, "expected error")
	})
}

func TestOrderUseCase_GetRangeOrders(t *testing.T) {
	userId := uint64(12)
	paginator := models.PaginatorOrders{
		PageNum: 1,
		Count:   43,
		SortOrdersOptions: models.SortOrdersOptions{
			SortKey:       "date",
			SortDirection: "ASC",
		},
	}
	orders := []*models.PlacedOrder{
		{
			Id: 1,
		},
	}
	products := []*models.PreviewOrderedProducts{{}}

	t.Run("GetRangeOrders_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := order_repo.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			GetCountPages(userId, paginator.Count).
			Return(10, nil)

		orderRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", nil)

		orderRepo.
			EXPECT().
			SelectRangeOrders(userId, "", &paginator).
			Return(orders, nil)

		orderRepo.
			EXPECT().
			GetProductsInOrder(uint64(1)).
			Return(products, nil)

		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)
		userRepo := user_repo.NewMockRepository(ctrl)
		promoCodeRepo := promo_code_repo.NewMockRepository(ctrl)

		orderUCase := NewUseCase(orderRepo, cartClient, productRepo, userRepo, promoCodeRepo)

		_, err := orderUCase.GetRangeOrders(userId, &paginator)
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetRangeOrders_incorrect_count", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := order_repo.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			GetCountPages(userId, paginator.Count).
			Return(10, errors.ErrInternalError)

		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)
		userRepo := user_repo.NewMockRepository(ctrl)
		promoCodeRepo := promo_code_repo.NewMockRepository(ctrl)

		orderUCase := NewUseCase(orderRepo, cartClient, productRepo, userRepo, promoCodeRepo)

		_, err := orderUCase.GetRangeOrders(userId, &paginator)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeOrders_incorrect_sort_string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := order_repo.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			GetCountPages(userId, paginator.Count).
			Return(10, nil)

		orderRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", errors.ErrInternalError)

		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)
		userRepo := user_repo.NewMockRepository(ctrl)
		promoCodeRepo := promo_code_repo.NewMockRepository(ctrl)

		orderUCase := NewUseCase(orderRepo, cartClient, productRepo, userRepo, promoCodeRepo)

		_, err := orderUCase.GetRangeOrders(userId, &paginator)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeOrders_not_select_range", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := order_repo.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			GetCountPages(userId, paginator.Count).
			Return(10, nil)

		orderRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", nil)

		orderRepo.
			EXPECT().
			SelectRangeOrders(userId, "", &paginator).
			Return(orders, errors.ErrIncorrectPaginator)

		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)
		userRepo := user_repo.NewMockRepository(ctrl)
		promoCodeRepo := promo_code_repo.NewMockRepository(ctrl)

		orderUCase := NewUseCase(orderRepo, cartClient, productRepo, userRepo, promoCodeRepo)

		_, err := orderUCase.GetRangeOrders(userId, &paginator)
		assert.Error(t, err, "expected error")
	})

	t.Run("GetRangeOrders_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		orderRepo := order_repo.NewMockRepository(ctrl)
		orderRepo.
			EXPECT().
			GetCountPages(userId, paginator.Count).
			Return(10, nil)

		orderRepo.
			EXPECT().
			CreateSortString(paginator.SortKey, paginator.SortDirection).
			Return("", nil)

		orderRepo.
			EXPECT().
			SelectRangeOrders(userId, "", &paginator).
			Return(orders, nil)

		orderRepo.
			EXPECT().
			GetProductsInOrder(uint64(1)).
			Return(products, errors.ErrInternalError)

		cartClient := cart_service.NewMockCartServiceClient(ctrl)
		productRepo := product_repo.NewMockRepository(ctrl)
		userRepo := user_repo.NewMockRepository(ctrl)
		promoCodeRepo := promo_code_repo.NewMockRepository(ctrl)

		orderUCase := NewUseCase(orderRepo, cartClient, productRepo, userRepo, promoCodeRepo)

		_, err := orderUCase.GetRangeOrders(userId, &paginator)
		assert.Error(t, err, "expected error")
	})
}
