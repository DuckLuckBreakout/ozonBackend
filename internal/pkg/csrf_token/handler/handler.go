package handler

import (
	"net/http"
	"time"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/csrf_token"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/jwt_token"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
)

type CsrfTokenHandler struct{}

func NewHandler() csrf_token.Handler {
	return &CsrfTokenHandler{}
}

// Get new csrf token for client
func (h *CsrfTokenHandler) GetCsrfToken(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		requireId := http_utils.MustGetRequireId(r.Context())
		if err != nil {
			logger.LogError("csrf_token_handler", "GetCsrfToken", requireId, err)
		}
	}()

	csrfToken := models.NewCsrfToken()
	jwtToken, err := jwt_token.CreateJwtToken([]byte(csrfToken.Value), time.Now().Add(models.ExpireCsrfToken*time.Second))

	if err != nil {
		http_utils.SetJSONResponse(w, errors.CreateError(err), http.StatusInternalServerError)
		return
	}
	csrfToken.Value = jwtToken

	http_utils.SetJSONResponse(w, csrfToken, http.StatusOK)
}
