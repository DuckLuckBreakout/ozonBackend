package category

import "github.com/DuckLuckBreakout/web/backend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/web/backend/internal/pkg/category UseCase

type UseCase interface {
	GetCatalogCategories() ([]*models.CategoriesCatalog, error)
	GetSubCategoriesById(categoryId uint64) ([]*models.CategoriesCatalog, error)
}
