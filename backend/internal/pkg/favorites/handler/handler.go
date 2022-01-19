package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"

	"github.com/gorilla/mux"
)

type FavoritesHandler struct {
	FavoritesUCase favorites.UseCase
}

func NewHandler(favoritesUCase favorites.UseCase) favorites.Handler {
	return &FavoritesHandler{
		FavoritesUCase: favoritesUCase,
	}
}

// AddProductToFavorites godoc
// @Summary Добавление товара в избранное.
// @Description Добавление товара в список избранных товаров.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param id query int true "Id товара"
// @Success 200 {object} errors.Error "Товар успешно добавлен в избранное."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags favorites
// @Router /favorites/product/{id} [PATCH]
func (h *FavoritesHandler) AddProductToFavorites(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("favorites_handler", "AddProductToFavorites", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	if err = h.FavoritesUCase.AddProductToFavorites(&usecase.FavoriteProduct{
		ProductId: uint64(productId),
		UserId:    currentSession.UserData.Id,
	}); err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// DeleteProductFromFavorites godoc
// @Summary Удаление товара из избранного.
// @Description Удаление товара из списка избранных товаров.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param id query int true "Id товара"
// @Success 200 {object} errors.Error "Товар успешно удалён из избранного."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags favorites
// @Router /favorites/product/{id} [DELETE]
func (h *FavoritesHandler) DeleteProductFromFavorites(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("favorites_handler", "DeleteProductFromFavorites", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	productId, err := strconv.Atoi(vars["id"])
	if err != nil || productId < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	if err = h.FavoritesUCase.DeleteProductFromFavorites(&usecase.FavoriteProduct{
		ProductId: uint64(productId),
		UserId:    currentSession.UserData.Id,
	}); err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}

// GetListPreviewFavorites godoc
// @Summary Получение превью избранного.
// @Description Получение превью избранных товаров.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param PaginatorFavorites body usecase.PaginatorFavorites true "Получение превью списка избранных товаров."
// @Success 200 {object} usecase.RangeFavorites "Список товаров успешно получен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags favorites
// @Router /favorites [POST]
func (h *FavoritesHandler) GetListPreviewFavorites(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("favorites_handler", "GetListPreviewFavorites", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator usecase.PaginatorFavorites
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

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	listPreviewFavorites, err := h.FavoritesUCase.GetRangeFavorites(&usecase.PaginatorFavorite{
		Paginator: paginator,
		UserId:    currentSession.UserData.Id,
	})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, listPreviewFavorites, http.StatusOK)
}

// GetUserFavorites godoc
// @Summary Получение превью избранного.
// @Description Получение превью избранных товаров.
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "Authorization"
// @Param PaginatorFavorites body usecase.PaginatorFavorites true "Получение превью списка избранных товаров."
// @Success 200 {object} usecase.UserFavorites "Список товаров успешно получен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags favorites
// @Router /favorites [POST]
func (h *FavoritesHandler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("favorites_handler", "GetUserFavorites", requireId, err)
		}
	}()

	currentSession := http_utils.MustGetSessionFromContext(r.Context())

	listFavorites, err := h.FavoritesUCase.GetUserFavorites(&usecase.UserId{Id: currentSession.UserData.Id})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, listFavorites, http.StatusOK)
}
