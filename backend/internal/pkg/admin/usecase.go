package admin

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/admin UseCase

type UseCase interface {
	ChangeOrderStatus(updateOrder *models.UpdateOrder) error
}
