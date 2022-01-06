package usecase

import (
	"context"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/cart"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/product"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	proto "github.com/DuckLuckBreakout/ozonBackend/services/cart/proto/cart"
)

type CartUseCase struct {
	CartClient  proto.CartServiceClient
	ProductRepo product.Repository
}

func NewUseCase(cartClient proto.CartServiceClient, productRepo product.Repository) cart.UseCase {
	return &CartUseCase{
		CartClient:  cartClient,
		ProductRepo: productRepo,
	}
}

// Add product in user cart
func (u *CartUseCase) AddProduct(userId *models.UserId, cartArticle *models.CartArticle) error {
	_, err := u.CartClient.AddProduct(
		context.Background(),
		&proto.ReqCartArticle{
			Position:  &proto.ProductPosition{Count: cartArticle.Count},
			ProductId: cartArticle.ProductId,
			UserId:    userId.Id,
		},
	)

	if err != nil {
		return errors.ErrInternalError
	}

	return nil
}

// Delete product from cart
func (u *CartUseCase) DeleteProduct(userId *models.UserId, identifier *models.ProductIdentifier) error {
	_, err := u.CartClient.DeleteProduct(
		context.Background(),
		&proto.ReqProductIdentifier{
			ProductId: identifier.ProductId,
			UserId:    userId.Id,
		},
	)

	if err != nil {
		return errors.ErrInternalError
	}

	return nil
}

// Change product in user cart
func (u *CartUseCase) ChangeProduct(userId *models.UserId, cartArticle *models.CartArticle) error {
	_, err := u.CartClient.ChangeProduct(
		context.Background(),
		&proto.ReqCartArticle{
			Position:  &proto.ProductPosition{Count: cartArticle.Count},
			ProductId: cartArticle.ProductId,
			UserId:    userId.Id,
		},
	)

	if err != nil {
		return errors.ErrInternalError
	}

	return nil
}

// Get preview cart
func (u *CartUseCase) GetPreviewCart(userId *models.UserId) (*models.PreviewCart, error) {
	userCart, err := u.CartClient.GetPreviewCart(
		context.Background(),
		&proto.ReqUserId{UserId: userId.Id},
	)

	if err != nil {
		return nil, errors.ErrInternalError
	}

	previewUserCart := &models.PreviewCart{}
	for id, productPosition := range userCart.Products {
		productById, err := u.ProductRepo.SelectProductById(id)
		if err != nil {
			return nil, err
		}

		previewUserCart.Products = append(previewUserCart.Products,
			&models.PreviewCartArticle{
				Id:    productById.Id,
				Title: productById.Title,
				Price: models.CartProductPrice{
					Discount:  productById.Price.Discount,
					BaseCost:  productById.Price.BaseCost,
					TotalCost: productById.Price.TotalCost,
				},
				PreviewImage: productById.Images[0],
				Count:        productPosition.Count,
			})

		previewUserCart.Price.TotalBaseCost += productById.Price.BaseCost * int(productPosition.Count)
		previewUserCart.Price.TotalCost += productById.Price.TotalCost * int(productPosition.Count)
	}
	previewUserCart.Price.TotalDiscount = previewUserCart.Price.TotalBaseCost - previewUserCart.Price.TotalCost

	return previewUserCart, nil
}

// Delete user cart
func (u *CartUseCase) DeleteCart(userId *models.UserId) error {
	_, err := u.CartClient.DeleteCart(
		context.Background(),
		&proto.ReqUserId{UserId: userId.Id},
	)

	if err != nil {
		return errors.ErrInternalError
	}

	return nil
}
