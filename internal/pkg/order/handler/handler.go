package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/cart"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
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

func (h *OrderHandler) GetOrderFromCart(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("order_handler", "GetOrderFromCart", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	previewCart, err := h.CartUCase.GetPreviewCart(currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	previewOrder, err := h.OrderUCase.GetPreviewOrder(currentSession.UserData.Id, previewCart)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, previewOrder, http.StatusOK)
}

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

	var userOrder models.Order
	err = json.Unmarshal(body, &userOrder)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}
	userOrder.Sanitize()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	previewCart, err := h.CartUCase.GetPreviewCart(currentSession.UserData.Id)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	orderNumber, err := h.OrderUCase.AddCompletedOrder(&userOrder, currentSession.UserData.Id, previewCart)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, orderNumber, http.StatusOK)
}

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

	var paginator models.PaginatorOrders
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

	orders, err := h.OrderUCase.GetRangeOrders(currentSession.UserData.Id, &paginator)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, orders, http.StatusOK)
}
