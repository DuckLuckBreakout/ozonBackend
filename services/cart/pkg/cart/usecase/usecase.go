package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	"github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart"
	"github.com/DuckLuckBreakout/ozonBackend/services/cart/pkg/cart/repository"
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
func (u *CartUseCase) AddProduct(userId *models.UserId, cartArticle *models.CartArticle) error {
	userIdentifier := &repository.DtoUserId{Id: userId.Id}
	userCart, err := u.CartRepo.SelectCartById(userIdentifier)
	if err != nil {
		userCart = &repository.DtoCart{}
		userCart.Products = make(map[uint64]*repository.DtoProductPosition)
		userCart.Products[cartArticle.ProductId] = &repository.DtoProductPosition{
			Count: cartArticle.ProductPosition.Count,
		}
	} else {
		// If product position already exist then increment counter
		if _, ok := userCart.Products[cartArticle.ProductId]; ok {
			userCart.Products[cartArticle.ProductId].Count += cartArticle.Count
		} else {
			userCart.Products[cartArticle.ProductId] = &repository.DtoProductPosition{
				Count: cartArticle.ProductPosition.Count,
			}
		}
	}

	return u.CartRepo.AddCart(userIdentifier, userCart)
}

// Delete product from cart
func (u *CartUseCase) DeleteProduct(userId *models.UserId, identifier *models.ProductIdentifier) error {
	userIdentifier := &repository.DtoUserId{Id: userId.Id}
	userCart, err := u.CartRepo.SelectCartById(userIdentifier)
	if err != nil {
		return err
	}

	// Delete cart of current user
	if len(userCart.Products) == 1 {
		if err = u.CartRepo.DeleteCart(userIdentifier); err != nil {
			return err
		}
	}

	delete(userCart.Products, identifier.ProductId)
	return u.CartRepo.AddCart(userIdentifier, userCart)
}

// Change product in user cart
func (u *CartUseCase) ChangeProduct(userId *models.UserId, cartArticle *models.CartArticle) error {
	userIdentifier := &repository.DtoUserId{Id: userId.Id}
	userCart, err := u.CartRepo.SelectCartById(userIdentifier)
	if err != nil {
		return err
	}

	if _, ok := userCart.Products[cartArticle.ProductId]; !ok {
		return errors.ErrProductNotFoundInCart
	}
	userCart.Products[cartArticle.ProductId] = &repository.DtoProductPosition{
		Count: cartArticle.ProductPosition.Count,
	}

	return u.CartRepo.AddCart(userIdentifier, userCart)
}

// Get preview cart
func (u *CartUseCase) GetPreviewCart(userId *models.UserId) (*models.Cart, error) {
	userCart, err := u.CartRepo.SelectCartById(&repository.DtoUserId{Id: userId.Id})
	switch err {
	case errors.ErrCartNotFound:
		return &models.Cart{}, nil
	case nil:

	default:
		return nil, err
	}

	previewCart := &models.Cart{}
	for key, item := range userCart.Products {
		previewCart.Products[key] = &models.ProductPosition{Count: item.Count}
	}

	return previewCart, nil
}

// Delete user cart
func (u *CartUseCase) DeleteCart(userId *models.UserId) error {
	return u.CartRepo.DeleteCart(&repository.DtoUserId{Id: userId.Id})
}
