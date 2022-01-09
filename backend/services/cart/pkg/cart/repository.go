package cart

import "github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/models"

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart Repository

type Repository interface {
	SelectCartById(userId uint64) (*models.Cart, error)
	AddCart(userId uint64, userCart *models.Cart) error
	DeleteCart(userId uint64) error
}
