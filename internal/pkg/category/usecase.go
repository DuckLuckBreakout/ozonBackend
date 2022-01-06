package category

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category UseCase

type UseCase interface {
	GetCatalogCategories() ([]*models.CategoriesCatalog, error)
	GetSubCategoriesById(categoryId *models.CategoryId) ([]*models.CategoriesCatalog, error)
}
