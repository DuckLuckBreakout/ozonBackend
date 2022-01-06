package cart

import (
	"github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart/repository"
)

//go:generate mockgen -destination=./mock/mock_repository.go -package=mock github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart Repository

type Repository interface {
	SelectCartById(userId *repository.DtoUserId) (*repository.DtoCart, error)
	AddCart(userId *repository.DtoUserId, userCart *repository.DtoCart) error
	DeleteCart(userId *repository.DtoUserId) error
}
