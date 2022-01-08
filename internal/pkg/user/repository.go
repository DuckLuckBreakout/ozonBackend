package user

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user Repository

type Repository interface {
	AddProfile(user *dto.DtoProfileUser) (*dto.DtoUserId, error)
	UpdateProfile(userId *dto.DtoUserId, user *dto.DtoUpdateUser) error
	SelectProfileById(userId *dto.DtoUserId) (*dto.DtoProfileUser, error)
	UpdateAvatar(userId *dto.DtoUserId, avatarUrl string) error
}
