package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/promo_code"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
)

type PromoCodeUseCase struct {
	PromoCodeRepo promo_code.Repository
}

func NewUseCase(promoCodeRepo promo_code.Repository) promo_code.UseCase {
	return &PromoCodeUseCase{
		PromoCodeRepo: promoCodeRepo,
	}
}

func (u *PromoCodeUseCase) ApplyPromoCodeToOrder(promoCodeGroup *usecase.PromoCodeGroup) (*usecase.DiscountedPrice, error) {
	err := u.PromoCodeRepo.CheckPromo(&dto.DtoPromoCode{
		Code: promoCodeGroup.PromoCode,
	})
	if err != nil {
		return nil, errors.ErrPromoCodeNotFound
	}

	productsInAction := 0
	discountedPrice := &usecase.DiscountedPrice{}
	for _, productId := range promoCodeGroup.Products {
		price, err := u.PromoCodeRepo.GetDiscountPriceByPromo(&dto.DtoPromoProduct{
			ProductId: productId,
			PromoCode: promoCodeGroup.PromoCode,
		})

		if err == nil {
			productsInAction += 1
		} else if err != errors.ErrProductNotInPromo {
			return nil, errors.ErrInternalError
		}

		discountedPrice.TotalBaseCost += price.BaseCost
		discountedPrice.TotalCost += price.TotalCost
	}
	if productsInAction == 0 {
		return nil, errors.ErrProductNotInPromo
	}
	discountedPrice.TotalDiscount = discountedPrice.TotalBaseCost - discountedPrice.TotalCost

	return discountedPrice, nil
}
