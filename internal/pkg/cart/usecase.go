package cart

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
)

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/internal/pkg/cart UseCase

type UseCase interface {
	AddProduct(userId *usecase.UserId, cartArticle *usecase.CartArticle) error
	DeleteProduct(userId *usecase.UserId, identifier *usecase.ProductIdentifier) error
	ChangeProduct(userId *usecase.UserId, cartArticle *usecase.CartArticle) error
	GetPreviewCart(userId *usecase.UserId) (*usecase.PreviewCart, error)
	DeleteCart(userId *usecase.UserId) error
}
