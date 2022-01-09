package repository

import (
	"database/sql"

	"github.com/DuckLuckBreakout/web/backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/web/backend/internal/pkg/promo_code"
	"github.com/DuckLuckBreakout/web/backend/internal/server/errors"
)

type PostgresqlRepository struct {
	db *sql.DB
}

func NewSessionPostgresqlRepository(db *sql.DB) promo_code.Repository {
	return &PostgresqlRepository{
		db: db,
	}
}

func (r *PostgresqlRepository) CheckPromo(promoCode string) error {
	row := r.db.QueryRow(
		"SELECT count(*) "+
			"FROM promo_codes "+
			"WHERE code = $1",
		promoCode,
	)

	var isExistPromo int
	err := row.Scan(
		&isExistPromo,
	)

	if err != nil || isExistPromo == 0 {
		return errors.ErrDBInternalError
	}

	return nil
}

func (r *PostgresqlRepository) GetDiscountPriceByPromo(productId uint64, promoCode string) (*models.PromoPrice, error) {
	row := r.db.QueryRow(
		"WITH pr AS ( "+
			"    SELECT id, sale "+
			"    FROM promo_codes "+
			"    WHERE code = $2 "+
			") "+
			"SELECT p.base_cost, p.total_cost, pr.sale "+
			"FROM products p "+
			"LEFT JOIN pr ON pr.id = ANY(p.sale_group) "+
			"WHERE p.id = $1",
		productId,
		promoCode,
	)

	promoSale := sql.NullInt64{}
	var baseCost, totalCost int
	err := row.Scan(
		&baseCost,
		&totalCost,
		&promoSale,
	)
	if err != nil {
		return nil, errors.ErrDBInternalError
	}
	price := &models.PromoPrice{
		BaseCost:  baseCost,
		TotalCost: totalCost,
	}

	sale := float32(promoSale.Int64)
	if sale == 0 {
		return price, errors.ErrProductNotInPromo
	}

	price.TotalCost = int(float32(price.TotalCost) * (1 - (sale / 100.0)))
	return price, nil
}
