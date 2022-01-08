package category

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category Repository

type Repository interface {
	GetNextLevelCategories(categoryId *dto.DtoCategoryId) (*dto.DtoCategoriesCatalog, error)
	GetCategoriesByLevel(categoryLevel *dto.DtoCategoryLevel) (*dto.DtoCategoriesCatalog, error)
	GetBordersOfBranch(categoryId *dto.DtoCategoryId) (*dto.DtoBranchBorders, error)
	GetPathToCategory(categoryId *dto.DtoCategoryId) (*dto.DtoCategoriesCatalog, error)
}
