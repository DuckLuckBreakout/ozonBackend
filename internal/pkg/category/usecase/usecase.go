package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category/repository"
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
	categories, err := u.CategoryRepo.GetCategoriesByLevel(&repository.DtoCategoryLevel{Level:uint64(1)})
	if err != nil {
		return nil, errors.ErrDBInternalError
	}

	for _, category := range categories.Catalog {
		nextLevel, err := u.CategoryRepo.GetNextLevelCategories(&repository.DtoCategoryId{Id: category.Id})
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		category.Next = nextLevel.Catalog

		for _, subCategory := range category.Next {
			nextLevel, err = u.CategoryRepo.GetNextLevelCategories(&repository.DtoCategoryId{Id: subCategory.Id})
			if err != nil {
				return nil, errors.ErrDBInternalError
			}
			subCategory.Next = nextLevel.Catalog
		}
	}

	return categories.Catalog, nil
}

// Get subcategories by category id
func (u *CategoryUseCase) GetSubCategoriesById(categoryId *models.CategoryId) ([]*models.CategoriesCatalog, error) {
	categories, err := u.CategoryRepo.GetNextLevelCategories(&repository.DtoCategoryId{Id: categoryId.Id})
	if err != nil {
		return nil, err
	}

	return categories.Catalog, nil
}
