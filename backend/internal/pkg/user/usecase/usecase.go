package usecase

import (
	"context"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/tools/s3_utils"
	proto "github.com/DuckLuckBreakout/ozonBackend/services/auth/proto/user"
	"mime/multipart"
)

type UserUseCase struct {
	UserRepo   user.Repository
	AuthClient proto.AuthServiceClient
}

func NewUseCase(authClient proto.AuthServiceClient, repo user.Repository) user.UseCase {
	return &UserUseCase{
		AuthClient: authClient,
		UserRepo:   repo,
	}
}

// Auth user
func (u *UserUseCase) Authorize(authUser *usecase.LoginUser) (*usecase.UserId, error) {
	userId, err := u.AuthClient.Login(context.Background(), &proto.LoginUser{
		Email:    authUser.Email,
		Password: authUser.Password,
	})

	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return &usecase.UserId{Id: userId.Id}, nil
}

// Set new avatar
func (u *UserUseCase) SetAvatar(userId *usecase.UserId, file *multipart.File, header *multipart.FileHeader) (string, error) {
	// Upload new user avatar to S3
	fileName, err := s3_utils.UploadMultipartFile("avatar", file, header)
	if err != nil {
		return "", err
	}

	// Destroy old user avatar
	profileUser, err := u.UserRepo.SelectProfileById(&dto.DtoUserId{Id: userId.Id})
	if err == nil && profileUser.Avatar.Url != "" {
		if err = s3_utils.DeleteFile(profileUser.Avatar.Url); err != nil {
			return "", err
		}
	}

	err = u.UserRepo.UpdateAvatar(&dto.DtoUserId{Id: userId.Id}, fileName)
	if err != nil {
		return "", err
	}

	return s3_utils.PathToFile(fileName), nil
}

// Get user avatar
func (u *UserUseCase) GetAvatar(userId *usecase.UserId) (string, error) {
	profileUser, err := u.UserRepo.SelectProfileById(&dto.DtoUserId{Id: userId.Id})
	if err != nil {
		return "", errors.ErrUserNotFound
	}

	return s3_utils.PathToFile(profileUser.Avatar.Url), nil
}

// Update user profile in repo
func (u *UserUseCase) UpdateProfile(userId *usecase.UserId, updateUser *usecase.UpdateUser) error {
	return u.UserRepo.UpdateProfile(
		&dto.DtoUserId{Id: userId.Id},
		&dto.DtoUpdateUser{
			FirstName: updateUser.FirstName,
			LastName:  updateUser.LastName,
		},
	)
}

// Get user profile by id
func (u *UserUseCase) GetUserById(userId *usecase.UserId) (*usecase.ProfileUser, error) {
	userById, err := u.UserRepo.SelectProfileById(&dto.DtoUserId{Id: userId.Id})
	if err != nil {
		return nil, err
	}
	userById.Avatar.Url = s3_utils.PathToFile(userById.Avatar.Url)
	return &usecase.ProfileUser{
		Id:        userById.Id,
		FirstName: userById.FirstName,
		LastName:  userById.LastName,
		Avatar: usecase.Avatar{
			Url: userById.Avatar.Url,
		},
		AuthId: userById.AuthId,
		Email:  userById.Email,
	}, nil
}

// Create new user account
func (u *UserUseCase) AddUser(signupUser *usecase.SignupUser) (*usecase.UserId, error) {
	userId, err := u.AuthClient.Signup(context.Background(), &proto.SignupUser{
		Email:    signupUser.Email,
		Password: signupUser.Password,
	})
	if err != nil {
		return nil, errors.ErrInternalError
	}

	_, err = u.UserRepo.AddProfile(&dto.DtoProfileUser{
		AuthId: userId.Id,
		Email:  signupUser.Email,
	})
	if err != nil {
		return nil, errors.ErrInternalError
	}

	return &usecase.UserId{Id: userId.Id}, nil
}
