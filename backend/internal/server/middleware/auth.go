package middleware

import (
	"context"
	"github.com/DuckLuckBreakout/web/backend/internal/server/tools/jwt_token"
	"net/http"
	"time"

	"github.com/DuckLuckBreakout/web/backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/web/backend/internal/pkg/session"
	"github.com/DuckLuckBreakout/web/backend/internal/server/errors"
	"github.com/DuckLuckBreakout/web/backend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/web/backend/pkg/tools/logger"
)

func Auth(sm session.UseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var err error
			defer func() {
				requireId := http_utils.MustGetRequireId(r.Context())
				if err != nil {
					logger.LogError("middleware", "Auth", requireId, err)
				}
			}()

			sessionCookie, err := r.Cookie(models.SessionCookieName)
			if err != nil {
				http_utils.SetJSONResponse(w, errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			csrfTokenFromCookie := jwt_token.JwtToken{}
			jwtToken, err := jwt_token.ParseJwtToken(sessionCookie.Value, &csrfTokenFromCookie)
			if err != nil || !jwtToken.Valid {
				http_utils.SetJSONResponse(w, errors.ErrIncorrectJwtToken, http.StatusBadRequest)
				return
			}

			t := time.Now()
			if t.After(csrfTokenFromCookie.Expires) {
				http_utils.SetJSONResponse(w, errors.ErrIncorrectJwtToken, http.StatusBadRequest)
				return
			}

			userId, err := sm.GetUserIdBySession(string(csrfTokenFromCookie.Value))
			if err != nil {
				http_utils.SetJSONResponse(w, errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), models.SessionContextKey,
				&models.Session{
					Value: sessionCookie.Value,
					UserData: models.UserId{
						Id: userId,
					},
				})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
