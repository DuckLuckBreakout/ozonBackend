package promo_code

import "github.com/DuckLuckBreakout/web/backend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/web/backend/internal/pkg/promo_code UseCase

type UseCase interface {
	ApplyPromoCodeToOrder(promoCodeGroup *models.PromoCodeGroup) (*models.DiscountedPrice, error)
}
