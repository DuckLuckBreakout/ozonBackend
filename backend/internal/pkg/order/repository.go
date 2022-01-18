package order

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/order Repository

type Repository interface {
	AddOrder(
		order *dto.DtoOrder,
		userId *dto.DtoUserId,
		products []*dto.DtoPreviewCartArticle,
		price *dto.DtoTotalPrice,
	) (*dto.DtoOrderNumber, error)
	SelectRangeOrders(
		rangeOrders *dto.DtoRangeOrders,
		paginator *dto.DtoPaginatorOrders,
	) ([]*dto.DtoPlacedOrder, error)
	CreateSortString(sortStr *dto.DtoSortString) (string, error)
	GetCountPages(cnt *dto.DtoOrderCountPages) (int, error)
	GetProductsInOrder(orderId *dto.DtoOrderId) ([]*dto.DtoPreviewOrderedProducts, error)
	ChangeStatusOrder(order *dto.DtoUpdateOrder) (*dto.DtoChangedOrder, error)
}
