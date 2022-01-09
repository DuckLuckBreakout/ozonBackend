package order

import "github.com/DuckLuckBreakout/web/backend/internal/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/web/backend/internal/pkg/order UseCase

type UseCase interface {
	GetPreviewOrder(userId uint64, previewCart *models.PreviewCart) (*models.PreviewOrder, error)
	AddCompletedOrder(order *models.Order, userId uint64, previewCart *models.PreviewCart) (*models.OrderNumber, error)
	GetRangeOrders(userId uint64, paginator *models.PaginatorOrders) (*models.RangeOrders, error)
}
