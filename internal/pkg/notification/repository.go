package notification

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification Repository

type Repository interface {
	AddSubscribeUser(userId *dto.DtoUserId, subscribes *dto.DtoSubscribes) error
	DeleteSubscribeUser(userId *dto.DtoUserId) error
	SelectCredentialsByUserId(userId *dto.DtoUserId) (*dto.DtoSubscribes, error)
}
