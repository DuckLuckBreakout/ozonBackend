package promo_code

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/promo_code UseCase

type UseCase interface {
	ApplyPromoCodeToOrder(promoCodeGroup *usecase.PromoCodeGroup) (*usecase.DiscountedPrice, error)
}
