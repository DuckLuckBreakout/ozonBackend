package category

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category UseCase

type UseCase interface {
	GetCatalogCategories() ([]*usecase.CategoriesCatalog, error)
	GetSubCategoriesById(categoryId *usecase.CategoryId) ([]*usecase.CategoriesCatalog, error)
}
