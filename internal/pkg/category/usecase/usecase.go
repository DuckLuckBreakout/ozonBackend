package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
)

type CategoryUseCase struct {
	CategoryRepo category.Repository
}

func NewUseCase(repo category.Repository) category.UseCase {
	return &CategoryUseCase{
		CategoryRepo: repo,
	}
}

// Get first three levels of categories tree
func (u *CategoryUseCase) GetCatalogCategories() ([]*models.CategoriesCatalog, error) {
	categories, err := u.CategoryRepo.GetCategoriesByLevel(1)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}

	for _, category := range categories {
		nextLevel, err := u.CategoryRepo.GetNextLevelCategories(category.Id)
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		category.Next = nextLevel

		for _, subCategory := range category.Next {
			nextLevel, err = u.CategoryRepo.GetNextLevelCategories(subCategory.Id)
			if err != nil {
				return nil, errors.ErrDBInternalError
			}
			subCategory.Next = nextLevel
		}
	}

	return categories, nil
}

// Get subcategories by category id
func (u *CategoryUseCase) GetSubCategoriesById(categoryId uint64) ([]*models.CategoriesCatalog, error) {
	return u.CategoryRepo.GetNextLevelCategories(categoryId)
}
