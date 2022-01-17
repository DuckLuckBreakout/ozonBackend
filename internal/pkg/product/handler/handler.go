package handler

import (
	"encoding/json"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/api"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/product"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/validator"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	ProductUCase product.UseCase
}

func NewHandler(UCase product.UseCase) product.Handler {
	return &ProductHandler{
		ProductUCase: UCase,
	}
}

// Get product info by id
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("product_handler", "GetProduct", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	productById, err := h.ProductUCase.GetProductById(&usecase.ProductId{Id: uint64(id)})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, productById, http.StatusOK)
}

// Get product recommendations by id
func (h *ProductHandler) GetProductRecommendations(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("product_handler", "GetProductRecommendations", requireId, err)
		}
	}()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator api.ApiPaginatorRecommendations
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

	listProducts, err := h.ProductUCase.GetProductRecommendationsById(
		&usecase.ProductId{
			Id: uint64(id),
		},
		&usecase.PaginatorRecommendations{
			Count: paginator.Count,
		},
	)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrProductNotFound, http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, listProducts, http.StatusOK)
}

// Get range of preview products
func (h *ProductHandler) GetListPreviewProducts(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("product_handler", "GetListPreviewProducts", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var paginator api.ApiPaginatorProducts
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

	listPreviewProducts, err := h.ProductUCase.GetRangeProducts(&usecase.PaginatorProducts{
		PageNum:  paginator.PageNum,
		Count:    paginator.Count,
		Category: paginator.Category,
		Filter: &usecase.ProductFilter{
			MinPrice:   paginator.Filter.MinPrice,
			MaxPrice:   paginator.Filter.MaxPrice,
			IsNew:      paginator.Filter.IsNew,
			IsRating:   paginator.Filter.IsRating,
			IsDiscount: paginator.Filter.IsDiscount,
		},
		SortOptions: usecase.SortOptions{
			SortKey:       paginator.SortKey,
			SortDirection: paginator.SortDirection,
		},
	})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, listPreviewProducts, http.StatusOK)
}

// Search range of preview products
func (h *ProductHandler) SearchListPreviewProducts(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("product_handler", "SearchListPreviewProducts", requireId, err)
		}
	}()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrBadRequest, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var searchQuery api.ApiSearchQuery
	err = json.Unmarshal(body, &searchQuery)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.ErrCanNotUnmarshal, http.StatusBadRequest)
		return
	}

	err = validator.ValidateStruct(searchQuery)
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusBadRequest)
		return
	}

	listPreviewProducts, err := h.ProductUCase.SearchRangeProducts(&usecase.SearchQuery{
		QueryString: searchQuery.QueryString,
		PageNum:     searchQuery.PageNum,
		Count:       searchQuery.Count,
		Category:    searchQuery.Category,
		Filter: &usecase.ProductFilter{
			MinPrice:   searchQuery.Filter.MaxPrice,
			MaxPrice:   searchQuery.Filter.MinPrice,
			IsNew:      searchQuery.Filter.IsNew,
			IsRating:   searchQuery.Filter.IsRating,
			IsDiscount: searchQuery.Filter.IsDiscount,
		},
		SortOptions: usecase.SortOptions{
			SortKey:       searchQuery.SortKey,
			SortDirection: searchQuery.SortDirection,
		},
	})
	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}

	http_utils.SetJSONResponse(w, listPreviewProducts, http.StatusOK)
}