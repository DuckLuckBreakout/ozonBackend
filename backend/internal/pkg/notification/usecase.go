package notification

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification UseCase

type UseCase interface {
	SubscribeUser(userId *usecase.UserId, credentials *usecase.NotificationCredentials) error
	UnsubscribeUser(userId *usecase.UserId, endpoint string) error
}
