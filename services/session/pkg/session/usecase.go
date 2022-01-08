package session

import "github.com/DuckLuckBreakout/ozonBackend/services/session/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/services/session/pkg/session UseCase

type UseCase interface {
	GetUserIdBySession(sessionCookieValue string) (uint64, error)
	CreateNewSession(userId *models.UserId) (*models.Session, error)
	DestroySession(sessionCookieValue string) error
}
