package usecase

import (
	"github.com/DuckLuckBreakout/web/backend/internal/pkg/admin"
	"github.com/DuckLuckBreakout/web/backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/web/backend/internal/pkg/notification"
	"github.com/DuckLuckBreakout/web/backend/internal/pkg/order"
	"github.com/DuckLuckBreakout/web/backend/internal/server/errors"
	"github.com/DuckLuckBreakout/web/backend/pkg/tools/server_push"
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
	orderNumber, userId, err := u.OrderRepo.ChangeStatusOrder(updateOrder.OrderId, updateOrder.Status)
	if err != nil {
		return errors.ErrInternalError
	}

	subscribes, err := u.NotificationRepo.SelectCredentialsByUserId(userId)
	if err == nil {
		for endpoint, keys := range subscribes.Credentials {
			err = server_push.Push(&server_push.Subscription{
				Endpoint: endpoint,
				Auth:     keys.Auth,
				P256dh:   keys.P256dh,
			}, models.OrderNotification{
				Number: *orderNumber,
				Status: updateOrder.Status,
			})

			if err != nil {
				return errors.ErrInternalError
			}
		}
	}

	return nil
}
