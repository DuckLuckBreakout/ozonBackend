package notification

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification UseCase

type UseCase interface {
	SubscribeUser(userId *models.UserId, credentials *models.NotificationCredentials) error
	UnsubscribeUser(userId *models.UserId, endpoint string) error
}
