package usecase

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/favorites/repository"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models"
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

func (u *FavoritesUseCase) AddProductToFavorites(favorite *models.FavoriteProduct) error {
	return u.FavoritesRepo.AddProductToFavorites(&repository.DtoFavoriteProduct{
		ProductId: favorite.ProductId,
		UserId:    favorite.UserId,
	})
}

func (u *FavoritesUseCase) DeleteProductFromFavorites(favorite *models.FavoriteProduct) error {
	return u.FavoritesRepo.DeleteProductFromFavorites(&repository.DtoFavoriteProduct{
		ProductId: favorite.ProductId,
		UserId:    favorite.UserId,
	})
}

func (u *FavoritesUseCase) GetRangeFavorites(pg *models.PaginatorFavorite) (*models.RangeFavorites, error) {
	if pg.Paginator.PageNum < 1 || pg.Paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	// Max count pages
	countPages, err := u.FavoritesRepo.GetCountPages(&repository.DtoCountPages{
		Count: pg.Paginator.Count,
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

	return &models.RangeFavorites{
		ListPreviewProducts: products,
		MaxCountPages:       countPages.Count,
	}, nil
}

func (u *FavoritesUseCase) GetUserFavorites(userId *models.UserId) (*models.UserFavorites, error) {
	favorite, err := u.FavoritesRepo.GetUserFavorites(&repository.DtoUserId{Id: userId.Id})
	if err != nil {
		return nil, err
	}

	return &models.UserFavorites{Products: favorite.Products}, nil
}
