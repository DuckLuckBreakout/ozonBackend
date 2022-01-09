package promo_code

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/promo_code Repository

type Repository interface {
	GetDiscountPriceByPromo(productId uint64, promoCode string) (*models.PromoPrice, error)
	CheckPromo(promoCode string) error
}
