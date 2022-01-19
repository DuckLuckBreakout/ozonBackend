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

// AddProductInCart godoc
// @Summary Добавление товара в корзину.
// @Description Добавление товара в корзину юзера.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param CartArticle body api.ApiCartArticle true "Данные товара."
// @Success 200 {object} errors.Error "Товар успешно добавлен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags cart
// @Router /cart/product [POST]
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

// DeleteProductInCart godoc
// @Summary Удаление товара из корзины.
// @Description Удаление товара из корзины юзера.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param ProductIdentifier body api.ApiProductIdentifier true "Данные товара."
// @Success 200 {object} errors.Error "Товар успешно удалён."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags cart
// @Router /cart/product [DELETE]
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

// ChangeProductInCart godoc
// @Summary Изменение продуктов в корзине.
// @Description Изменение продуктов в пользовательской корзине.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param CartArticle body api.ApiCartArticle true "Данные товара."
// @Success 200 {object} errors.Error "Товар успешно удалён."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags cart
// @Router /cart/product [PUT]
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

// GetProductsFromCart godoc
// @Summary Получение корзины.
// @Description Получение пользовательской корзины.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Success 200 {object} api.ApiPreviewCart "Товары из пользовательской корзины успешно получены."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags cart
// @Router /cart [GET]
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

// DeleteProductsFromCart godoc
// @Summary Удаление корзины.
// @Description Удаление пользовательской корзины.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Success 200 {object} errors.Error "Товары из пользовательской корзины успешно удалены."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags cart
// @Router /cart [DELETE]
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
