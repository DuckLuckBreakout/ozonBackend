package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification/repository"
)

type NotificationUseCase struct {
	NotificationRepo notification.Repository
}

func NewUseCase(notificationRepo notification.Repository) notification.UseCase {
	return &NotificationUseCase{
		NotificationRepo: notificationRepo,
	}
}

func (u *NotificationUseCase) SubscribeUser(userId *models.UserId, credentials *models.NotificationCredentials) error {
	var subscribes *repository.DtoSubscribes
	subscribes, err := u.NotificationRepo.SelectCredentialsByUserId(&repository.DtoUserId{Id: userId.Id})
	if err != nil || subscribes.Credentials == nil || subscribes == nil {
		subscribes = &repository.DtoSubscribes{}
		subscribes.Credentials = make(map[string]*repository.DtoNotificationKeys)
	}

	subscribes.Credentials[credentials.Endpoint] = &repository.DtoNotificationKeys{
		Auth:   credentials.Keys.Auth,
		P256dh: credentials.Keys.P256dh,
	}

	return u.NotificationRepo.AddSubscribeUser(&repository.DtoUserId{Id: userId.Id}, subscribes)
}

func (u *NotificationUseCase) UnsubscribeUser(userId *models.UserId, endpoint string) error {
	id := &repository.DtoUserId{Id: userId.Id}
	userSubscribes, err := u.NotificationRepo.SelectCredentialsByUserId(id)
	if err != nil {
		return err
	}

	if len(userSubscribes.Credentials) == 1 {
		return u.NotificationRepo.DeleteSubscribeUser(id)
	}

	delete(userSubscribes.Credentials, endpoint)
	return u.NotificationRepo.AddSubscribeUser(id, userSubscribes)
}
