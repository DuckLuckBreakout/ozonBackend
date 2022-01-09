package user

import (
	"github.com/DuckLuckBreakout/ozonBackend/services/auth/pkg/models"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/services/auth/pkg/user Repository

type Repository interface {
	AddProfile(user *models.AuthUser) (uint64, error)
	SelectUserByEmail(email string) (*models.AuthUser, error)
}
