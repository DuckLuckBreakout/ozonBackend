package handler

import (
	"net/http"
	"strconv"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"

	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	CategoryUCase category.UseCase
}

func NewHandler(UCase category.UseCase) category.Handler {
	return &CategoryHandler{
		CategoryUCase: UCase,
	}
}

// GetCatalogCategories godoc
// @Summary Получение каталога категорий.
// @Description Получение каталога категорий.
// @Accept json
// @Produce json
// @Success 200 {array} errors.Error "Каталог категорий успешно получен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags category
// @Router /category [GET]
func (h *CategoryHandler) GetCatalogCategories(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("category_handler", "GetCatalogCategories", requireId, err)
		}
	}()

	categories, err := h.CategoryUCase.GetCatalogCategories()
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, categories, http.StatusOK)
}

// GetSubCategories godoc
// @Summary Получение подкатегории.
// @Description Получение подкатегории.
// @Accept json
// @Produce json
// @Param id query int true "Id подкатегории"
// @Success 200 {array}  errors.Error "Каталог категорий успешно получен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags category
// @Router /category/{id} [GET]
func (h *CategoryHandler) GetSubCategories(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("category_handler", "GetSubCategories", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	categories, err := h.CategoryUCase.GetSubCategoriesById(&usecase.CategoryId{Id: uint64(id)})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, categories, http.StatusOK)
}
