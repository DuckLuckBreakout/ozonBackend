package usecase

import (
	"context"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"

	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/order"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/product"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/promo_code"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user"
	userRepo "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/user/repository"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
	proto "github.com/DuckLuckBreakout/ozonBackend/services/cart/proto/cart"
)

type OrderUseCase struct {
	OrderRepo     order.Repository
	CartClient    proto.CartServiceClient
	ProductRepo   product.Repository
	UserRepo      user.Repository
	PromoCodeRepo promo_code.Repository
}

func NewUseCase(orderRepo order.Repository, cartClient proto.CartServiceClient,
	productRepo product.Repository, userRepo user.Repository, promoCodeRepo promo_code.Repository) order.UseCase {
	return &OrderUseCase{
		OrderRepo:     orderRepo,
		CartClient:    cartClient,
		ProductRepo:   productRepo,
		UserRepo:      userRepo,
		PromoCodeRepo: promoCodeRepo,
	}
}

func (u *OrderUseCase) GetPreviewOrder(userId *usecase.UserId, previewCart *usecase.PreviewCart) (*usecase.PreviewOrder, error) {
	// Get all info about product in cart
	previewOrder := &usecase.PreviewOrder{}
	for _, item := range previewCart.Products {
		previewOrder.Products = append(previewOrder.Products,
			&usecase.PreviewOrderedProducts{
				Id:           item.Id,
				PreviewImage: item.PreviewImage,
			})
	}
	previewOrder.Price = previewCart.Price

	// Get info about user account for order
	userProfile, err := u.UserRepo.SelectProfileById(&userRepo.DtoUserId{Id: userId.Id})
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	previewOrder.Recipient = usecase.OrderRecipient{
		FirstName: userProfile.FirstName,
		LastName:  userProfile.LastName,
		Email:     userProfile.Email,
	}

	return previewOrder, nil
}

func (u *OrderUseCase) AddCompletedOrder(
	order *usecase.Order,
	userId *usecase.UserId,
	previewCart *usecase.PreviewCart,
) (*usecase.OrderNumber, error) {
	price := usecase.TotalPrice{}

	if order.PromoCode == "" {
		price = previewCart.Price
	} else {
		err := u.PromoCodeRepo.CheckPromo(&dto.DtoPromoCode{
			Code: order.PromoCode,
		})
		if err != nil {
			return nil, errors.ErrPromoCodeNotFound
		}

		for _, product := range previewCart.Products {
			promoPrice, err := u.PromoCodeRepo.GetDiscountPriceByPromo(&dto.DtoPromoProduct{
				ProductId: product.Id,
				PromoCode: order.PromoCode,
			})
			if err != nil && err != errors.ErrProductNotInPromo {
				return nil, errors.ErrProductNotFound
			}
			price.TotalBaseCost += promoPrice.BaseCost
			price.TotalCost += promoPrice.TotalCost
		}
		price.TotalDiscount = price.TotalBaseCost - price.TotalCost
	}
	products := previewCart.Products

	var orderProducts []*dto.DtoPreviewCartArticle
	for _, item := range products {
		orderProducts = append(orderProducts, &dto.DtoPreviewCartArticle{
			Id:    item.Id,
			Title: item.Title,
			Price: dto.DtoCartProductPrice{
				Discount:  item.Price.Discount,
				BaseCost:  item.Price.BaseCost,
				TotalCost: item.Price.TotalCost,
			},
			PreviewImage: item.PreviewImage,
			Count:        item.Count,
		})
	}
	orderNumber, err := u.OrderRepo.AddOrder(
		&dto.DtoOrder{
			Recipient: dto.DtoOrderRecipient{
				FirstName: order.Recipient.FirstName,
				LastName:  order.Recipient.LastName,
				Email:     order.Recipient.Email,
			},
			Address: dto.DtoOrderAddress{
				Address: order.Address.Address,
			},
			PromoCode: order.PromoCode,
		},
		&dto.DtoUserId{Id: userId.Id},
		orderProducts,
		&dto.DtoTotalPrice{
			TotalDiscount: price.TotalDiscount,
			TotalCost:     price.TotalCost,
			TotalBaseCost: price.TotalBaseCost,
		},
	)
	if err != nil {
		return nil, errors.ErrInternalError
	}

	if _, err = u.CartClient.DeleteCart(context.Background(), &proto.ReqUserId{UserId: userId.Id}); err != nil {
		return nil, errors.ErrCartNotFound
	}

	return &usecase.OrderNumber{Number: orderNumber.Number}, nil
}

func (u *OrderUseCase) GetRangeOrders(userId *usecase.UserId, paginator *usecase.PaginatorOrders) (*usecase.RangeOrders, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	// Max count pages in catalog
	countPages, err := u.OrderRepo.GetCountPages(&dto.DtoCountPages{
		ProductId:         userId.Id,
		CountOrdersOnPage: paginator.Count,
	})
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Keys for sort items in catalog
	sortString, err := u.OrderRepo.CreateSortString(&dto.DtoSortString{
		SortKey:       paginator.SortKey,
		SortDirection: paginator.SortDirection,
	})
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get range of products
	orders, err := u.OrderRepo.SelectRangeOrders(
		&dto.DtoRangeOrders{
			UserId:     userId.Id,
			SortString: sortString,
		},
		&dto.DtoPaginatorOrders{
			PageNum: paginator.PageNum,
			Count:   paginator.Count,
			DtoSortOrdersOptions: dto.DtoSortOrdersOptions{
				SortKey:       paginator.SortOrdersOptions.SortKey,
				SortDirection: paginator.SortOrdersOptions.SortDirection,
			},
		},
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get product for this order
	for _, item := range orders {
		products, err := u.OrderRepo.GetProductsInOrder(&dto.DtoOrderId{Id: item.Id})
		if err != nil {
			return nil, errors.ErrInternalError
		}

		item.Products = products
	}

	var placedOrder []*usecase.PlacedOrder
	for _, item := range orders {
		var placedProducts []*usecase.PreviewOrderedProducts
		for _, orderedProduct := range item.Products {
			placedProducts = append(placedProducts, &usecase.PreviewOrderedProducts{
				Id:           orderedProduct.Id,
				PreviewImage: orderedProduct.PreviewImage,
			})
		}

		placedOrder = append(placedOrder, &usecase.PlacedOrder{
			Id: item.Id,
			Address: usecase.OrderAddress{
				Address: item.Address.Address,
			},
			TotalCost:    item.TotalCost,
			Products:     placedProducts,
			DateAdded:    item.DateAdded,
			DateDelivery: item.DateDelivery,
			OrderNumber: usecase.OrderNumber{
				Number: item.OrderNumber.Number,
			},
			Status: item.Status,
		})
	}
	return &usecase.RangeOrders{
		ListPreviewOrders: placedOrder,
		MaxCountPages:     countPages,
	}, nil
}
