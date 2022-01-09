package promo_code

import "github.com/DuckLuckBreakout/web/backend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/web/backend/internal/pkg/promo_code Repository

type Repository interface {
	GetDiscountPriceByPromo(productId uint64, promoCode string) (*models.PromoPrice, error)
	CheckPromo(promoCode string) error
}
