package usecase

import (
	"testing"

	category_mock "github.com/DuckLuckBreakout/web/backend/internal/pkg/category/mock"
	"github.com/DuckLuckBreakout/web/backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/web/backend/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCategoryUseCase_GetSubCategoriesById(t *testing.T) {
	categoryId := uint64(1)
	categories := []*models.CategoriesCatalog{
		{
			Id:   categoryId,
			Name: "test",
			Next: nil,
		},
	}

	t.Run("GetCatalogCategories_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_mock.NewMockRepository(ctrl)

		categoryRepo.
			EXPECT().
			GetCategoriesByLevel(uint64(1)).
			Return(categories, nil)

		categoryRepo.
			EXPECT().
			GetNextLevelCategories(categoryId).
			Return(categories, nil)

		categoryRepo.
			EXPECT().
			GetNextLevelCategories(categoryId).
			Return(categories, nil)

		userUCase := NewUseCase(categoryRepo)

		_, err := userUCase.GetCatalogCategories()
		assert.NoError(t, err, "unexpected error")
	})

	t.Run("GetCatalogCategories_incorrect_first_lvl", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_mock.NewMockRepository(ctrl)

		categoryRepo.
			EXPECT().
			GetCategoriesByLevel(uint64(1)).
			Return(categories, errors.ErrInternalError)

		userUCase := NewUseCase(categoryRepo)

		_, err := userUCase.GetCatalogCategories()
		assert.Error(t, err, "expected error")
	})

	t.Run("GetCatalogCategories_incorrect_second_lvl", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_mock.NewMockRepository(ctrl)

		categoryRepo.
			EXPECT().
			GetCategoriesByLevel(uint64(1)).
			Return(categories, nil)

		categoryRepo.
			EXPECT().
			GetNextLevelCategories(categoryId).
			Return(categories, errors.ErrInternalError)

		userUCase := NewUseCase(categoryRepo)

		_, err := userUCase.GetCatalogCategories()
		assert.Error(t, err, "expected error")
	})

	t.Run("GetCatalogCategories_incorrect_last_lvl", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		categoryRepo := category_mock.NewMockRepository(ctrl)

		categoryRepo.
			EXPECT().
			GetCategoriesByLevel(uint64(1)).
			Return(categories, nil)

		categoryRepo.
			EXPECT().
			GetNextLevelCategories(categoryId).
			Return(categories, nil)

		categoryRepo.
			EXPECT().
			GetNextLevelCategories(categoryId).
			Return(categories, errors.ErrInternalError)

		userUCase := NewUseCase(categoryRepo)

		_, err := userUCase.GetCatalogCategories()
		assert.Error(t, err, "expected error")
	})
}
