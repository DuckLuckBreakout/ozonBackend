package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DuckLuckBreakout/web/backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/web/backend/internal/pkg/product/mock"
	"github.com/DuckLuckBreakout/web/backend/internal/server/errors"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestProductHandler_GetProduct(t *testing.T) {
	productId := uint64(4)
	product := models.Product{
		Id:    productId,
		Title: "test",
		Price: models.ProductPrice{
			Discount: 10,
			BaseCost: 20,
		},
		Rating:       4,
		Description:  "fdfdf",
		Category:     3,
		CategoryPath: nil,
		Images:       nil,
	}

	t.Run("GetProduct_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)
		productUCase.
			EXPECT().
			GetProductById(productId).
			Return(&product, nil)

		productHandler := NewHandler(productUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.GetProduct)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("GetProduct_without_args", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)

		productHandler := NewHandler(productUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.GetProduct)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("GetProduct_internal_error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)
		productUCase.
			EXPECT().
			GetProductById(productId).
			Return(nil, errors.ErrDBInternalError)

		productHandler := NewHandler(productUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/product/{id:[0-9]+}",
			bytes.NewBuffer(nil))

		vars := map[string]string{
			"id": fmt.Sprintf("%d", productId),
		}
		req = mux.SetURLVars(req, vars)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.GetProduct)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestProductHandler_GetListPreviewProducts(t *testing.T) {
	paginator := models.PaginatorProducts{
		PageNum: 2,
		Count:   10,
		SortOptions: models.SortOptions{
			SortKey:       "cost",
			SortDirection: "ASC",
		},
		Category: 1,
	}
	invalidPaginator := models.PaginatorProducts{
		PageNum: 2,
		Count:   10,
		SortOptions: models.SortOptions{
			SortKey:       "fdf",
			SortDirection: "df",
		},
		Category: 1,
	}
	rangeProduct := models.RangeProducts{
		ListPreviewProducts: []*models.ViewProduct{
			&models.ViewProduct{
				Id:    3,
				Title: "test",
				Price: models.ProductPrice{
					Discount: 10,
					BaseCost: 50,
				},
				Rating:       5,
				PreviewImage: "fdfdf",
			},
		},
		MaxCountPages: 3,
	}

	t.Run("GetListPreviewProducts_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)
		productUCase.
			EXPECT().
			GetRangeProducts(&paginator).
			Return(&rangeProduct, nil)

		productHandler := NewHandler(productUCase)

		bytesPaginator, _ := json.Marshal(paginator)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/product",
			bytes.NewBuffer(bytesPaginator))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.GetListPreviewProducts)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("GetListPreviewProducts_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)

		productHandler := NewHandler(productUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/product",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.GetListPreviewProducts)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("GetListPreviewProducts_incorrect_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)

		productHandler := NewHandler(productUCase)

		bytesPaginator, _ := json.Marshal(invalidPaginator)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/product",
			bytes.NewBuffer(bytesPaginator))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.GetListPreviewProducts)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("GetListPreviewProducts_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)
		productUCase.
			EXPECT().
			GetRangeProducts(&paginator).
			Return(nil, errors.ErrInternalError)

		productHandler := NewHandler(productUCase)

		bytesPaginator, _ := json.Marshal(paginator)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/product",
			bytes.NewBuffer(bytesPaginator))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.GetListPreviewProducts)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}

func TestProductHandler_SearchListPreviewProducts(t *testing.T) {
	searchQuery := models.SearchQuery{
		QueryString: "search",
		PageNum:     2,
		Count:       10,
		Category:    1,
		Filter:      nil,
		SortOptions: models.SortOptions{
			SortKey:       "cost",
			SortDirection: "ASC",
		},
	}
	invalidSearchQuery := models.SearchQuery{
		QueryString: "search",
		PageNum:     2,
		Count:       10,
		Category:    1,
		Filter:      nil,
		SortOptions: models.SortOptions{
			SortKey:       "fdf",
			SortDirection: "df",
		},
	}
	rangeProduct := models.RangeProducts{
		ListPreviewProducts: []*models.ViewProduct{
			&models.ViewProduct{
				Id:    3,
				Title: "test",
				Price: models.ProductPrice{
					Discount: 10,
					BaseCost: 50,
				},
				Rating:       5,
				PreviewImage: "fdfdf",
			},
		},
		MaxCountPages: 3,
	}

	t.Run("SearchListPreviewProducts_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)
		productUCase.
			EXPECT().
			SearchRangeProducts(&searchQuery).
			Return(&rangeProduct, nil)

		productHandler := NewHandler(productUCase)

		bytesPaginator, _ := json.Marshal(searchQuery)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/product/search",
			bytes.NewBuffer(bytesPaginator))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.SearchListPreviewProducts)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})

	t.Run("SearchListPreviewProducts_bad_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)

		productHandler := NewHandler(productUCase)

		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/product/search",
			bytes.NewBuffer(nil))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.SearchListPreviewProducts)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("SearchListPreviewProducts_incorrect_body", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)

		productHandler := NewHandler(productUCase)

		bytesPaginator, _ := json.Marshal(invalidSearchQuery)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/product/search",
			bytes.NewBuffer(bytesPaginator))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.SearchListPreviewProducts)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusBadRequest, "incorrect http code")
	})

	t.Run("SearchListPreviewProducts_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUCase := mock.NewMockUseCase(ctrl)
		productUCase.
			EXPECT().
			SearchRangeProducts(&searchQuery).
			Return(nil, errors.ErrInternalError)

		productHandler := NewHandler(productUCase)

		bytesPaginator, _ := json.Marshal(searchQuery)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "POST", "/api/v1/product/search",
			bytes.NewBuffer(bytesPaginator))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(productHandler.SearchListPreviewProducts)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusInternalServerError, "incorrect http code")
	})
}
