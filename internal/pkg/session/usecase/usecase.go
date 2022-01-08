package usecase

import (
	"context"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/session"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	proto "github.com/DuckLuckBreakout/ozonBackend/services/session/proto/session"
)

type SessionUseCase struct {
	SessionClient proto.SessionServiceClient
}

func NewUseCase(sessionConn proto.SessionServiceClient) session.UseCase {
	return &SessionUseCase{
		SessionClient: sessionConn,
	}
}

// Get user id by session value
func (u *SessionUseCase) GetUserIdBySession(sessionCookieValue string) (uint64, error) {
	userId, err := u.SessionClient.GetUserIdBySession(context.Background(), &proto.SessionValue{
		Value: sessionCookieValue,
	})
	if err != nil {
		return 0, errors.ErrSessionNotFound
	}

	return userId.Id, nil
}

// Create new user session and save in repository
func (u *SessionUseCase) CreateNewSession(userId *usecase.UserId) (*usecase.Session, error) {
	userSession, err := u.SessionClient.CreateNewSession(context.Background(), &proto.UserId{
		Id: userId.Id,
	})
	if err != nil {
		return nil, errors.ErrInternalError
	}

	return &usecase.Session{
		Value: userSession.Value.Value,
		UserData: usecase.UserId{
			Id: userSession.Id.Id,
		},
	}, nil
}

// Destroy session from repository by session value
func (u *SessionUseCase) DestroySession(sessionCookieValue string) error {
	_, err := u.SessionClient.DestroySession(context.Background(), &proto.SessionValue{
		Value: sessionCookieValue,
	})

	if err != nil {
		return errors.ErrSessionNotFound
	}

	return nil
}
