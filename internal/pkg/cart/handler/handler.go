package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/cart"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/api"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
)

type CartHandler struct {
	CartUCase cart.UseCase
}

func NewHandler(cartUCase cart.UseCase) cart.Handler {
	return &CartHandler{
		CartUCase: cartUCase,
	}
}

// Add product in user cart
func (h *CartHandler) AddProductInCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("cart_handler", "AddProductInCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cartArticle := &api.ApiCartArticle{}
	err = json.Unmarshal(body, cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.AddProduct(
		&usecase.UserId{
			Id: currentSession.UserData.Id,
		},
		&usecase.CartArticle{
			ProductPosition: usecase.ProductPosition{
				Count: cartArticle.Count,
			},
			ProductIdentifier: usecase.ProductIdentifier{
				ProductId: cartArticle.ProductId,
			},
		},
	)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Delete product from user cart
func (h *CartHandler) DeleteProductInCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("cart_handler", "DeleteProductInCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	identifier := &api.ApiProductIdentifier{}
	err = json.Unmarshal(body, identifier)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(identifier)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.DeleteProduct(
		&usecase.UserId{
			Id: currentSession.UserData.Id,
		},
		&usecase.ProductIdentifier{
			ProductId: identifier.ProductId,
		},
	)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Change product characteristics in user cart
func (h *CartHandler) ChangeProductInCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("cart_handler", "ChangeProductInCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cartArticle := &api.ApiCartArticle{}
	err = json.Unmarshal(body, cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(cartArticle)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	err = h.CartUCase.ChangeProduct(
		&usecase.UserId{
			Id: currentSession.UserData.Id,
		},
		&usecase.CartArticle{
			ProductPosition: usecase.ProductPosition{
				Count: cartArticle.Count,
			},
			ProductIdentifier: usecase.ProductIdentifier{
				ProductId: cartArticle.ProductId,
			},
		},
	)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// Get all preview products from user cart
func (h *CartHandler) GetProductsFromCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("cart_handler", "GetProductsFromCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	previewUserCart, err := h.CartUCase.GetPreviewCart(&usecase.UserId{Id: currentSession.UserData.Id})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, api.ApiPreviewCart{
		Products: previewUserCart.Products,
		Price:    previewUserCart.Price,
	}, http.StatusOK)
}

// Delete user cart
func (h *CartHandler) DeleteProductsFromCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("cart_handler", "GetProductsFromCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	err = h.CartUCase.DeleteCart(&usecase.UserId{Id: currentSession.UserData.Id})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}
