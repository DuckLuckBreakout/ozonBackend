package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/cart"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
)

type ApiPreviewCart struct {
	Products []*models.PreviewCartArticle `json:"products" valid:"notnull"`
	Price    models.TotalPrice            `json:"price" valid:"notnull"`
}

type CartHandler struct {
	CartUCase cart.UseCase
}

type ApiProductPosition struct {
	Count uint64 `json:"count"`
}

type ApiProductIdentifier struct {
	ProductId uint64 `json:"product_id"`
}

type ApiCartArticle struct {
	ApiProductPosition
	ApiProductIdentifier
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

	cartArticle := &ApiCartArticle{}
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
		&models.UserId{
			Id: currentSession.UserData.Id,
		},
		&models.CartArticle{
			ProductPosition:   models.ProductPosition{
				Count: cartArticle.Count,
			},
			ProductIdentifier: models.ProductIdentifier{
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

	identifier := &ApiProductIdentifier{}
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
		&models.UserId{
			Id: currentSession.UserData.Id,
		},
		&models.ProductIdentifier{
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

	cartArticle := &ApiCartArticle{}
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
		&models.UserId{
			Id: currentSession.UserData.Id,
		},
		&models.CartArticle{
			ProductPosition:   models.ProductPosition{
				Count: cartArticle.Count,
			},
			ProductIdentifier: models.ProductIdentifier{
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

	previewUserCart, err := h.CartUCase.GetPreviewCart(&models.UserId{Id: currentSession.UserData.Id})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, ApiPreviewCart{
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

	err = h.CartUCase.DeleteCart(&models.UserId{Id: currentSession.UserData.Id})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}
