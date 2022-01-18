package promo_code

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/promo_code Repository

type Repository interface {
	CheckPromo(promoCode *dto.DtoPromoCode) error
	GetDiscountPriceByPromo(promoPrice *dto.DtoPromoProduct) (*dto.DtoPromoPrice, error)
}
