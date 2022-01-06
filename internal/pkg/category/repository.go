package category

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category/repository"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category Repository

type Repository interface {
	GetNextLevelCategories(categoryId *repository.DtoCategoryId) (*repository.DtoCategoriesCatalog, error)
	GetCategoriesByLevel(categoryLevel *repository.DtoCategoryLevel) (*repository.DtoCategoriesCatalog, error)
	GetBordersOfBranch(categoryId *repository.DtoCategoryId) (*repository.DtoBranchBorders, error)
	GetPathToCategory(categoryId *repository.DtoCategoryId) (*repository.DtoCategoriesCatalog, error)
}
