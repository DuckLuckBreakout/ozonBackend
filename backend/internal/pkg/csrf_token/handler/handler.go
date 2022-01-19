package handler

import (
	"net/http"
	"time"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/csrf_token"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/jwt_token"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
)

type CsrfTokenHandler struct{}

func NewHandler() csrf_token.Handler {
	return &CsrfTokenHandler{}
}

// GetCsrfToken godoc
// @Summary Получение csrf токена.
// @Description Получение csrf токена.
// @Accept json
// @Produce json
// @Success 200 {object} usecase.CsrfToken "Токен csrf успешно сгенирован."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Tags csrf
// @Router /csrf [GET]
func (h *CsrfTokenHandler) GetCsrfToken(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("csrf_token_handler", "GetCsrfToken", requireId, err)
		}
	}()

	csrfToken := usecase.NewCsrfToken()
	jwtToken, err := jwt_token.CreateJwtToken([]byte(csrfToken.Value), time.Now().Add(usecase.ExpireCsrfToken*time.Second))

	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}
	csrfToken.Value = jwtToken

	http_utils.SetJSONResponse(w, csrfToken, http.StatusOK)
}
