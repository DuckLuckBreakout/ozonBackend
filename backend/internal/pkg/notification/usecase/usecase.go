package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification"
)

type NotificationUseCase struct {
	NotificationRepo notification.Repository
}

func NewUseCase(notificationRepo notification.Repository) notification.UseCase {
	return &NotificationUseCase{
		NotificationRepo: notificationRepo,
	}
}

func (u *NotificationUseCase) SubscribeUser(userId *usecase.UserId, credentials *usecase.NotificationCredentials) error {
	var subscribes *dto.DtoSubscribes
	subscribes, err := u.NotificationRepo.SelectCredentialsByUserId(&dto.DtoUserId{Id: userId.Id})
	if err != nil || subscribes.Credentials == nil || subscribes == nil {
		subscribes = &dto.DtoSubscribes{}
		subscribes.Credentials = make(map[string]*dto.DtoNotificationKeys)
	}

	subscribes.Credentials[credentials.Endpoint] = &dto.DtoNotificationKeys{
		Auth:   credentials.Keys.Auth,
		P256dh: credentials.Keys.P256dh,
	}

	return u.NotificationRepo.AddSubscribeUser(&dto.DtoUserId{Id: userId.Id}, subscribes)
}

func (u *NotificationUseCase) UnsubscribeUser(userId *usecase.UserId, endpoint string) error {
	id := &dto.DtoUserId{Id: userId.Id}
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
