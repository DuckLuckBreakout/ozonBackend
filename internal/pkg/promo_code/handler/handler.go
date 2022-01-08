package handler

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/api"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"io/ioutil"
	"net/http"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/promo_code"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
)

type PromoCodeHandler struct {
	PromoCodeUCase promo_code.UseCase
}

func NewHandler(UCase promo_code.UseCase) promo_code.Handler {
	return &PromoCodeHandler{
		PromoCodeUCase: UCase,
	}
}

func (h *PromoCodeHandler) ApplyPromoCodeToOrder(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("promo_code_handler", "ApplyPromoCodeToOrder", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	promoCodeGroup := &api.ApiPromoCodeGroup{}
	err = json.Unmarshal(body, promoCodeGroup)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(promoCodeGroup)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	discountedPrice, err := h.PromoCodeUCase.ApplyPromoCodeToOrder(&usecase.PromoCodeGroup{
		Products:  promoCodeGroup.Products,
		PromoCode: promoCodeGroup.PromoCode,
	})
	if err == errors.ErrProductNotInPromo {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusConflict)
		return
	} else if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, api.ApiDiscountedPrice{
		TotalDiscount: discountedPrice.TotalDiscount,
		TotalCost:     discountedPrice.TotalCost,
		TotalBaseCost: discountedPrice.TotalBaseCost,
	}, http.StatusOK)
}
