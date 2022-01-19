package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/cart"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/api"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/order"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
)

type OrderHandler struct {
	OrderUCase order.UseCase
	CartUCase  cart.UseCase
}

func NewHandler(orderUCase order.UseCase, cartUCase cart.UseCase) order.Handler {
	return &OrderHandler{
		OrderUCase: orderUCase,
		CartUCase:  cartUCase,
	}
}

// GetOrderFromCart godoc
// @Summary Получение заказа на основании корзины.
// @Description Получение заказа на основании корзины.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Success 200 {object} usecase.PreviewOrder "Заказ успешно получен на основании корзины."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags order
// @Router /order [GET]
func (h *OrderHandler) GetOrderFromCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("order_handler", "GetOrderFromCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	id := &usecase.UserId{Id: currentSession.UserData.Id}
	previewCart, err := h.CartUCase.GetPreviewCart(id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	previewOrder, err := h.OrderUCase.GetPreviewOrder(id, previewCart)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, previewOrder, http.StatusOK)
}

// AddCompletedOrder godoc
// @Summary Оформление заказа.
// @Description Оформление пользовательского заказа.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param Order body api.ApiOrder true "Данные заказа."
// @Success 200 {object} usecase.OrderNumber "Заказ успешно оформлен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags order
// @Router /order [POST]
func (h *OrderHandler) AddCompletedOrder(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("order_handler", "AddCompletedOrder", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var userOrder api.ApiOrder
	err = json.Unmarshal(body, &userOrder)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	userOrder.Sanitize()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	id := &usecase.UserId{Id: currentSession.UserData.Id}
	previewCart, err := h.CartUCase.GetPreviewCart(id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	orderNumber, err := h.OrderUCase.AddCompletedOrder(&usecase.Order{
		Recipient: usecase.OrderRecipient{
			FirstName: userOrder.Recipient.FirstName,
			LastName:  userOrder.Recipient.LastName,
			Email:     userOrder.Recipient.Email,
		},
		Address: usecase.OrderAddress{
			Address: userOrder.Address.Address,
		},
		PromoCode: userOrder.PromoCode,
	}, id, previewCart)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, orderNumber, http.StatusOK)
}

// GetUserOrders godoc
// @Summary Получение заказов.
// @Description Получение заказов пользователя.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param PaginatorOrders body api.ApiPaginatorOrders true "Параметры пагинатора."
// @Success 200 {object} usecase.RangeOrders "Список заказов данного пользователя успешно получен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags order
// @Router /user/order [POST]
func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("order_handler", "GetUserOrders", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator api.ApiPaginatorOrders
	err = json.Unmarshal(body, &paginator)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(paginator)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	orders, err := h.OrderUCase.GetRangeOrders(
		&usecase.UserId{Id: currentSession.UserData.Id},
		&usecase.PaginatorOrders{
			PageNum: paginator.PageNum,
			Count:   paginator.Count,
			SortOrdersOptions: usecase.SortOrdersOptions{
				SortKey:       paginator.SortKey,
				SortDirection: paginator.SortDirection,
			},
		},
	)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, orders, http.StatusOK)
}
