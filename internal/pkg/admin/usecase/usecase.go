package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/admin"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/order"
	orderRepo "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/order/repository"
	notificationRepo "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/notification/repository"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/pkg/tools/server_push"
)

type AdminUseCase struct {
	NotificationRepo notification.Repository
	OrderRepo        order.Repository
}

func NewUseCase(notificationRepo notification.Repository, orderRepo order.Repository) admin.UseCase {
	return &AdminUseCase{
		NotificationRepo: notificationRepo,
		OrderRepo:        orderRepo,
	}
}

func (u *AdminUseCase) ChangeOrderStatus(updateOrder *models.UpdateOrder) error {
	changedOrder, err := u.OrderRepo.ChangeStatusOrder(&orderRepo.DtoUpdateOrder{
		OrderId: updateOrder.OrderId,
		Status:  updateOrder.Status,
	})
	if err != nil {
		return errors.ErrInternalError
	}

	subscribes, err := u.NotificationRepo.SelectCredentialsByUserId(&notificationRepo.DtoUserId{Id: changedOrder.UserId})
	if err == nil {
		for endpoint, keys := range subscribes.Credentials {
			err = server_push.Push(&server_push.Subscription{
				Endpoint: endpoint,
				Auth:     keys.Auth,
				P256dh:   keys.P256dh,
			}, models.OrderNotification{
				Number: models.OrderNumber{
					Number: changedOrder.Number,
				},
				Status: updateOrder.Status,
			})

			if err != nil {
				return errors.ErrInternalError
			}
		}
	}

	return nil
}
