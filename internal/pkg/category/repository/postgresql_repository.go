package repository

import (
	"database/sql"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

type DtoCategoryId struct {
	Id   uint64               `json:"id"`
}

type DtoCategoryLevel struct {
	Level   uint64               `json:"level"`
}

type DtoCategoriesCatalog struct {
	Catalog []*models.CategoriesCatalog
}

type DtoBranchBorders struct {
	Left uint64 `json:"left"`
	Right uint64 `json:"right"`
}

func NewSessionPostgresqlRepository(db *sql.DB) category.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

// Get lower level in categories tree
func (r *PostgresqlRepository) GetNextLevelCategories(categoryId *DtoCategoryId) (*DtoCategoriesCatalog, error) {
	rows, err := r.db.Query(
		"WITH current_node AS ( "+
			"SELECT c.left_node, c.right_node, c.level + 1  as level "+
			"FROM categories c "+
			"WHERE c.id = $1 "+
			") "+
			"SELECT c.id, c.name "+
			"FROM categories c, current_node "+
			"WHERE (c.left_node > current_node.left_node "+
			"AND c.right_node < current_node.right_node "+
			"AND c.level = current_node.level)",
		categoryId.Id,
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}
	defer rows.Close()

	categories := make([]*models.CategoriesCatalog, 0)
	for rows.Next() {
		nextLevelCategory := &models.CategoriesCatalog{}
		err = rows.Scan(
			&nextLevelCategory.Id,
			&nextLevelCategory.Name,
		)
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		categories = append(categories, nextLevelCategory)
	}

	return &DtoCategoriesCatalog{categories}, nil
}

// Get categories in select level
func (r *PostgresqlRepository) GetCategoriesByLevel(categoryLevel *DtoCategoryLevel) (*DtoCategoriesCatalog, error) {
	rows, err := r.db.Query(
		"SELECT c.id, c.name "+
			"FROM categories c "+
			"WHERE c.level = $1",
		categoryLevel.Level,
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	categories := make([]*models.CategoriesCatalog, 0)
	for rows.Next() {
		nextLevelCategory := &models.CategoriesCatalog{}
		err = rows.Scan(
			&nextLevelCategory.Id,
			&nextLevelCategory.Name,
		)
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		categories = append(categories, nextLevelCategory)
	}

	return &DtoCategoriesCatalog{categories}, nil
}

// Get left and right border of branch
func (r *PostgresqlRepository) GetBordersOfBranch(categoryId *DtoCategoryId) (*DtoBranchBorders, error) {
	row := r.db.QueryRow(
		"SELECT c.left_node, c.right_node "+
			"FROM categories c "+
			"WHERE c.id = $1",
		categoryId.Id,
	)

	var left, right uint64
	err := row.Scan(
		&left,
		&right,
	)

	if err != nil {
		return nil, errors.ErrDBInternalError
	}

	return &DtoBranchBorders{
		Left:  left,
		Right: right,
	}, nil
}

// Get path from root to category
func (r *PostgresqlRepository) GetPathToCategory(categoryId *DtoCategoryId) (*DtoCategoriesCatalog, error) {
	rows, err := r.db.Query(
		"WITH current_node AS ( "+
			"SELECT c.left_node, c.right_node, c.level + 1  as level "+
			"FROM categories c "+
			"WHERE c.id = $1 "+
			") "+
			"SELECT c.id, c.name "+
			"FROM categories c, current_node "+
			"WHERE (c.left_node <= current_node.left_node "+
			"AND c.right_node >= current_node.right_node "+
			"AND c.level BETWEEN 1 AND current_node.level)",
		categoryId.Id,
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	defer rows.Close()

	categories := make([]*models.CategoriesCatalog, 0)
	for rows.Next() {
		nextLevelCategory := &models.CategoriesCatalog{}
		err = rows.Scan(
			&nextLevelCategory.Id,
			&nextLevelCategory.Name,
		)
		if err != nil {
			return nil, errors.ErrDBInternalError
		}
		categories = append(categories, nextLevelCategory)
	}

	return &DtoCategoriesCatalog{categories}, nil
}
