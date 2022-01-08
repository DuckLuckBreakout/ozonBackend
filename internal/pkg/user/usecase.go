package user

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"mime/multipart"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user UseCase

type UseCase interface {
	Authorize(authUser *usecase.LoginUser) (*usecase.UserId, error)
	SetAvatar(userId *usecase.UserId, file *multipart.File, header *multipart.FileHeader) (string, error)
	GetAvatar(userId *usecase.UserId) (string, error)
	UpdateProfile(userId *usecase.UserId, updateUser *usecase.UpdateUser) error
	GetUserById(userId *usecase.UserId) (*usecase.ProfileUser, error)
	AddUser(signupUser *usecase.SignupUser) (*usecase.UserId, error)
}
