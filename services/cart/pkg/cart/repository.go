package cart

import (
	"github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/models"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart Repository

type Repository interface {
	SelectCartById(userId *models.DtoUserId) (*models.DtoCart, error)
	AddCart(userId *models.DtoUserId, userCart *models.DtoCart) error
	DeleteCart(userId *models.DtoUserId) error
}
