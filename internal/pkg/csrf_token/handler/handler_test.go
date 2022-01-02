package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"

	"github.com/golang/mock/gomock"
	"github.com/lithammer/shortuuid"
	"github.com/stretchr/testify/assert"
)

func TestCsrfTokenHandler_GetCsrfToken(t *testing.T) {
	cartArticle := models.CartArticle{
		ProductPosition: models.ProductPosition{
			Count: 3,
		},
		ProductIdentifier: models.ProductIdentifier{
			ProductId: 2,
		},
	}

	t.Run("GetCsrfToken_success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		csrfHandler := NewHandler()

		bytesArticle, _ := json.Marshal(cartArticle)
		ctx := context.WithValue(context.Background(), models.RequireIdKey, shortuuid.New())
		req, _ := http.NewRequestWithContext(ctx, "GET", "/api/v1/csrf",
			bytes.NewBuffer(bytesArticle))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(csrfHandler.GetCsrfToken)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, http.StatusOK, "incorrect http code")
	})
}
