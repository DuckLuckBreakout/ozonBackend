package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/services/session/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/services/session/pkg/session"
	"github.com/DuckLuckBreakout/ozonBackend/services/session/server/errors"
)

type SessionUseCase struct {
	SessionRepo session.Repository
}

func NewUseCase(SessionRepo session.Repository) session.UseCase {
	return &SessionUseCase{
		SessionRepo: SessionRepo,
	}
}

// Get user id by session value
func (u *SessionUseCase) GetUserIdBySession(sessionCookieValue string) (uint64, error) {
	userId, err := u.SessionRepo.SelectUserIdBySession(sessionCookieValue)
	if err != nil {
		return 0, err
	}

	return userId.Id, nil
}

// Create new user session and save in repository
func (u *SessionUseCase) CreateNewSession(userId *models.UserId) (*models.Session, error) {
	sess := models.NewSession(userId.Id)
	err := u.SessionRepo.AddSession(&models.DtoSession{
		Value:  sess.Value,
		UserId: sess.UserId,
	})
	if err != nil {
		return nil, errors.ErrInternalError
	}

	return sess, nil
}

// Destroy session from repository by session value
func (u *SessionUseCase) DestroySession(sessionCookieValue string) error {
	return u.SessionRepo.DeleteSessionByValue(sessionCookieValue)
}
