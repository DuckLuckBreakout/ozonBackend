package middleware

import (
	"net/http"
	"time"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/jwt_token"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
)

func CsrfCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			requireId := http_utils.MustGetRequireId(r.Context())
			if err != nil {
				logger.LogError("middleware", "CsrfCheck", requireId, err)
			}
		}()

		if r.Method == http.MethodPost || r.Method == http.MethodDelete ||
			r.Method == http.MethodPut || r.Method == http.MethodPatch {
			csrfTokenFromHeader := r.Header.Get(models.CsrfTokenHeaderName)
			if csrfTokenFromHeader == "" {
				http_utils.SetJSONResponse(w, errors.ErrNotFoundCsrfToken, http.StatusBadRequest)
				return
			}

			csrfTokenFromCookie := jwt_token.JwtToken{}
			jwtToken, err := jwt_token.ParseJwtToken(csrfTokenFromHeader, &csrfTokenFromCookie)
			if err != nil || !jwtToken.Valid {
				http_utils.SetJSONResponse(w, errors.ErrIncorrectJwtToken, http.StatusBadRequest)
				return
			}

			t := time.Now()
			if t.After(csrfTokenFromCookie.Expires) {
				http_utils.SetJSONResponse(w, errors.ErrIncorrectJwtToken, http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
