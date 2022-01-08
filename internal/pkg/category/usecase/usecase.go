package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
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
func (u *CategoryUseCase) GetCatalogCategories() ([]*usecase.CategoriesCatalog, error) {
	categories, err := u.CategoryRepo.GetCategoriesByLevel(&dto.DtoCategoryLevel{Level: uint64(1)})
	if err != nil {
		return nil, errors.ErrDBInternalError
	}

	for _, category := range categories.Catalog {
		nextLevel, err := u.CategoryRepo.GetNextLevelCategories(&dto.DtoCategoryId{Id: category.Id})
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		category.Next = nextLevel.Catalog

		for _, subCategory := range category.Next {
			nextLevel, err = u.CategoryRepo.GetNextLevelCategories(&dto.DtoCategoryId{Id: subCategory.Id})
			if err != nil {
				return nil, errors.ErrDBInternalError
			}
			subCategory.Next = nextLevel.Catalog
		}
	}

	return categories.Catalog, nil
}

// Get subcategories by category id
func (u *CategoryUseCase) GetSubCategoriesById(categoryId *usecase.CategoryId) ([]*usecase.CategoriesCatalog, error) {
	categories, err := u.CategoryRepo.GetNextLevelCategories(&dto.DtoCategoryId{Id: categoryId.Id})
	if err != nil {
		return nil, err
	}

	return categories.Catalog, nil
}
