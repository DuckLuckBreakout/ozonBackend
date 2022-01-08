package middleware

import (
	"context"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"net/http"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/session"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/http_utils"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/logger"
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

			sessionCookie, err := r.Cookie(usecase.SessionCookieName)
			if err != nil {
				http_utils.SetJSONResponse(w, errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			userId, err := sm.GetUserIdBySession(sessionCookie.Value)
			if err != nil {
				http_utils.SetJSONResponse(w, errors.ErrUserUnauthorized, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), usecase.SessionContextKey,
				&usecase.Session{
					Value: sessionCookie.Value,
					UserData: usecase.UserId{
						Id: userId,
					},
				})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
