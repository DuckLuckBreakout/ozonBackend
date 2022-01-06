package cart

import "github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/models"

//go:generate mockgen -destination=./mock/mock_usecase.go -package=mock github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart UseCase

type UseCase interface {
	AddProduct(userId *models.UserId, cartArticle *models.CartArticle) error
	DeleteProduct(userId *models.UserId, identifier *models.ProductIdentifier) error
	ChangeProduct(userId *models.UserId, cartArticle *models.CartArticle) error
	GetPreviewCart(userId *models.UserId) (*models.Cart, error)
	DeleteCart(userId *models.UserId) error
}
