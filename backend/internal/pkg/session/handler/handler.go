package handler

import (
	"net/http"

	"github.com/DuckLuckBreakout/web/backend/internal/pkg/session"
	"github.com/DuckLuckBreakout/web/backend/internal/server/tools/http_utils"
)

type SessionHandler struct {
	SessionUCase session.UseCase
}

func NewHandler(sessionUCase session.UseCase) session.Handler {
	return &SessionHandler{
		SessionUCase: sessionUCase,
	}
}

func (h *SessionHandler) CheckSession(w http.ResponseWriter, r *http.Request) {
	http_utils.SetJSONResponseSuccess(w, http.StatusOK)
}
