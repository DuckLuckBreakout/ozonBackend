package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/services/auth/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/services/auth/pkg/user"
	"github.com/DuckLuckBreakout/ozonBackend/services/auth/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/services/auth/server/tools/password_hasher"
)

type UserUseCase struct {
	UserRepo user.Repository
}

func NewUseCase(repo user.Repository) user.UseCase {
	return &UserUseCase{
		UserRepo: repo,
	}
}

// Login user
func (u *UserUseCase) Login(loginUser *models.LoginUser) (uint64, error) {
	authUser, err := u.UserRepo.SelectUserByEmail(loginUser.Email)
	if err != nil {
		return 0, errors.ErrIncorrectAuthData
	}

	if ok := password_hasher.CompareHashAndPassword(authUser.Password, loginUser.Password); !ok {
		return 0, errors.ErrIncorrectAuthData
	}

	return authUser.Id, nil
}

// Signup user
func (u *UserUseCase) Signup(signupUser *models.SignupUser) (uint64, error) {
	hashOfPassword, err := password_hasher.GenerateHashFromPassword(signupUser.Password)
	if err != nil {
		return 0, errors.ErrHashFunctionFailed
	}

	userId, err := u.UserRepo.AddProfile(&models.ApiAuthUser{
		Email:    signupUser.Email,
		Password: hashOfPassword,
	})
	if err != nil {
		return 0, errors.ErrDBInternalError
	}

	return userId, nil
}
