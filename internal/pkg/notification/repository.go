package notification

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification/repository"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification Repository

type Repository interface {
	AddSubscribeUser(userId *repository.DtoUserId, subscribes *repository.DtoSubscribes) error
	DeleteSubscribeUser(userId *repository.DtoUserId) error
	SelectCredentialsByUserId(userId *repository.DtoUserId) (*repository.DtoSubscribes, error)
}
