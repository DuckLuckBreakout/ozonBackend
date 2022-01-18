package session

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/session UseCase

type UseCase interface {
	GetUserIdBySession(sessionCookieValue string) (uint64, error)
	CreateNewSession(userId *usecase.UserId) (*usecase.Session, error)
	DestroySession(sessionCookieValue string) error
}
