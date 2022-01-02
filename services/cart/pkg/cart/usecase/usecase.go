package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart"
	"github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/models"
)

type CartUseCase struct {
	CartRepo cart.Repository
}

func NewUseCase(cartRepo cart.Repository) cart.UseCase {
	return &CartUseCase{
		CartRepo: cartRepo,
	}
}

// Add product in user cart
func (u *CartUseCase) AddProduct(userId uint64, cartArticle *models.CartArticle) error {
	userCart, err := u.CartRepo.SelectCartById(userId)
	if err != nil {
		userCart = &models.Cart{}
		userCart.Products = make(map[uint64]*models.ProductPosition)
		userCart.Products[cartArticle.ProductId] = &cartArticle.ProductPosition
	} else {
		// If product position already exist then increment counter
		if _, ok := userCart.Products[cartArticle.ProductId]; ok {
			userCart.Products[cartArticle.ProductId].Count += cartArticle.Count
		} else {
			userCart.Products[cartArticle.ProductId] = &cartArticle.ProductPosition
		}
	}

	return u.CartRepo.AddCart(userId, userCart)
}

// Delete product from cart
func (u *CartUseCase) DeleteProduct(userId uint64, identifier *models.ProductIdentifier) error {
	userCart, err := u.CartRepo.SelectCartById(userId)
	if err != nil {
		return err
	}

	// Delete cart of current user
	if len(userCart.Products) == 1 {
		if err = u.CartRepo.DeleteCart(userId); err != nil {
			return err
		}
	}

	delete(userCart.Products, identifier.ProductId)
	return u.CartRepo.AddCart(userId, userCart)
}

// Change product in user cart
func (u *CartUseCase) ChangeProduct(userId uint64, cartArticle *models.CartArticle) error {
	userCart, err := u.CartRepo.SelectCartById(userId)
	if err != nil {
		return err
	}

	if _, ok := userCart.Products[cartArticle.ProductId]; !ok {
		return errors.ErrProductNotFoundInCart
	}
	userCart.Products[cartArticle.ProductId] = &cartArticle.ProductPosition

	return u.CartRepo.AddCart(userId, userCart)
}

// Get preview cart
func (u *CartUseCase) GetPreviewCart(userId uint64) (*models.Cart, error) {
	userCart, err := u.CartRepo.SelectCartById(userId)
	switch err {
	case errors.ErrCartNotFound:
		return &models.Cart{}, nil
	case nil:

	default:
		return nil, err
	}

	return userCart, nil
}

// Delete user cart
func (u *CartUseCase) DeleteCart(userId uint64) error {
	return u.CartRepo.DeleteCart(userId)
}
