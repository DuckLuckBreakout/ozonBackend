package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/dto"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
)

type FavoritesUseCase struct {
	FavoritesRepo favorites.Repository
}

func NewUseCase(favoritesRepo favorites.Repository) favorites.UseCase {
	return &FavoritesUseCase{
		FavoritesRepo: favoritesRepo,
	}
}

func (u *FavoritesUseCase) AddProductToFavorites(favorite *usecase.FavoriteProduct) error {
	return u.FavoritesRepo.AddProductToFavorites(&dto.DtoFavoriteProduct{
		ProductId: favorite.ProductId,
		UserId:    favorite.UserId,
	})
}

func (u *FavoritesUseCase) DeleteProductFromFavorites(favorite *usecase.FavoriteProduct) error {
	return u.FavoritesRepo.DeleteProductFromFavorites(&dto.DtoFavoriteProduct{
		ProductId: favorite.ProductId,
		UserId:    favorite.UserId,
	})
}

func (u *FavoritesUseCase) GetRangeFavorites(pg *usecase.PaginatorFavorite) (*usecase.RangeFavorites, error) {
	if pg.Paginator.PageNum < 1 || pg.Paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	// Max count pages
	countPages, err := u.FavoritesRepo.GetCountPages(&dto.DtoCountPages{
		Count:  pg.Paginator.Count,
		UserId: pg.UserId,
	})
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Keys for sort items
	sortString, err := u.FavoritesRepo.CreateSortString(pg.Paginator.SortKey, pg.Paginator.SortDirection)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get range of favorites
	products, err := u.FavoritesRepo.SelectRangeFavorites(&pg.Paginator, sortString, pg.UserId)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	return &usecase.RangeFavorites{
		ListPreviewProducts: products,
		MaxCountPages:       countPages.Count,
	}, nil
}

func (u *FavoritesUseCase) GetUserFavorites(userId *usecase.UserId) (*usecase.UserFavorites, error) {
	favorite, err := u.FavoritesRepo.GetUserFavorites(&dto.DtoUserId{Id: userId.Id})
	if err != nil {
		return nil, err
	}

	return &usecase.UserFavorites{Products: favorite.Products}, nil
}
