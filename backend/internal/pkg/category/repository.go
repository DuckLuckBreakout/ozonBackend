package category

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category Repository

type Repository interface {
	GetNextLevelCategories(categoryId uint64) ([]*models.CategoriesCatalog, error)
	GetCategoriesByLevel(level uint64) ([]*models.CategoriesCatalog, error)
	GetBordersOfBranch(categoryId uint64) (uint64, uint64, error)
	GetPathToCategory(categoryId uint64) ([]*models.CategoriesCatalog, error)
}
