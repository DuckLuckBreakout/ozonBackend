package order

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/order UseCase

type UseCase interface {
	GetPreviewOrder(userId *usecase.UserId, previewCart *usecase.PreviewCart) (*usecase.PreviewOrder, error)
	AddCompletedOrder(order *usecase.Order, userId *usecase.UserId, previewCart *usecase.PreviewCart) (*usecase.OrderNumber, error)
	GetRangeOrders(userId *usecase.UserId, paginator *usecase.PaginatorOrders) (*usecase.RangeOrders, error)
}
