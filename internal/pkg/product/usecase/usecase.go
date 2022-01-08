package product

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category"
	categoryRepo "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/category/repository"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"
	"github.com/DuckLuckBreakout/ozonBackend/internal/pkg/product"
	productRepo "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/product/repository"
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"
)

type ProductUseCase struct {
	ProductRepo  product.Repository
	CategoryRepo category.Repository
}

func NewUseCase(productRepo product.Repository, categoryRepo category.Repository) product.UseCase {
	return &ProductUseCase{
		ProductRepo:  productRepo,
		CategoryRepo: categoryRepo,
	}
}

// Get product by id from repo
func (u *ProductUseCase) GetProductById(productId *usecase.ProductId) (*usecase.Product, error) {
	productById, err := u.ProductRepo.SelectProductById(&productRepo.DtoProductId{ProductId: productId.Id})
	if err != nil {
		return nil, errors.ErrProductNotFound
	}

	categories, err := u.CategoryRepo.GetPathToCategory(&categoryRepo.DtoCategoryId{Id: productById.Category})
	if err != nil {
		return nil, errors.ErrCategoryNotFound
	}

	productById.CategoryPath = categories.Catalog

	return &usecase.Product{
		Id:    productById.Id,
		Title: productById.Title,
		Price: usecase.ProductPrice{
			Discount:  productById.Price.Discount,
			BaseCost:  productById.Price.BaseCost,
			TotalCost: productById.Price.TotalCost,
		},
		Rating:       productById.Rating,
		CountReviews: productById.CountReviews,
		Properties:   productById.Properties,
		Description:  productById.Description,
		Category:     productById.Category,
		CategoryPath: productById.CategoryPath,
		Images:       productById.Images,
	}, nil
}

// Get recommendations for product
func (u *ProductUseCase) GetProductRecommendationsById(
	productId *usecase.ProductId,
	paginator *usecase.PaginatorRecommendations,
) ([]*usecase.RecommendationProduct, error) {
	recommendationsByReviews, err := u.ProductRepo.SelectRecommendationsByReviews(&productRepo.DtoRecommendations{
		ProductId: productId.Id,
		Count:     paginator.Count,
	})
	if err != nil {
		return nil, errors.ErrProductNotFound
	}

	var recommendationsProducts []*usecase.RecommendationProduct
	for _, item := range recommendationsByReviews {
		recommendationsProducts = append(recommendationsProducts, &usecase.RecommendationProduct{
			Id:    item.Id,
			Title: item.Title,
			Price: usecase.ProductPrice{
				Discount:  item.Price.Discount,
				BaseCost:  item.Price.BaseCost,
				TotalCost: item.Price.TotalCost,
			},
			PreviewImage: item.PreviewImage,
		})
	}

	return recommendationsProducts, nil
}

// Get range products by paginator settings from repo
func (u *ProductUseCase) GetRangeProducts(paginator *usecase.PaginatorProducts) (*usecase.RangeProducts, error) {
	if paginator.PageNum < 1 || paginator.Count < 1 {
		return nil, errors.ErrIncorrectPaginator
	}

	var filterString string
	if paginator.Filter != nil {
		filterString = u.ProductRepo.CreateFilterString(&productRepo.DtoProductFilter{
			MinPrice:   paginator.Filter.MinPrice,
			MaxPrice:   paginator.Filter.MaxPrice,
			IsNew:      paginator.Filter.IsNew,
			IsRating:   paginator.Filter.IsRating,
			IsDiscount: paginator.Filter.IsDiscount,
		})
	}

	// Max count pages in catalog
	countPages, err := u.ProductRepo.GetCountPages(&productRepo.DtoCountPages{
		Category:     paginator.Category,
		Count:        paginator.Count,
		FilterString: filterString,
	})
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Keys for sort items in catalog
	sortString, err := u.ProductRepo.CreateSortString(&productRepo.DtoSortString{
		SortKey:       paginator.SortKey,
		SortDirection: paginator.SortDirection,
	})
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	// Get range of products
	products, err := u.ProductRepo.SelectRangeProducts(
		&productRepo.DtoPaginatorProducts{
			PageNum:  paginator.PageNum,
			Count:    paginator.Count,
			Category: paginator.Category,
			Filter: &productRepo.DtoProductFilter{
				MinPrice:   paginator.Filter.MinPrice,
				MaxPrice:   paginator.Filter.MaxPrice,
				IsNew:      paginator.Filter.IsNew,
				IsRating:   paginator.Filter.IsRating,
				IsDiscount: paginator.Filter.IsDiscount,
			},
			SortKey:       paginator.SortKey,
			SortDirection: paginator.SortDirection,
		},
		&productRepo.DtoRageProducts{
			SortString:   sortString,
			FilterString: filterString,
		})
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	var rangeProducts []*usecase.ViewProduct
	for _, item := range products {
		rangeProducts = append(rangeProducts, &usecase.ViewProduct{
			Id:    item.Id,
			Title: item.Title,
			Price: usecase.ProductPrice{
				Discount:  item.Price.Discount,
				BaseCost:  item.Price.BaseCost,
				TotalCost: item.Price.TotalCost,
			},
			Rating:       item.Rating,
			CountReviews: item.CountReviews,
			PreviewImage: item.PreviewImage,
		})
	}
	return &usecase.RangeProducts{
		ListPreviewProducts: rangeProducts,
		MaxCountPages:       countPages,
	}, nil
}

// Find range products by search settings from repo
func (u *ProductUseCase) SearchRangeProducts(searchQuery *usecase.SearchQuery) (*usecase.RangeProducts, error) {
	if searchQuery.PageNum < 1 || searchQuery.Count < 1 {
		return nil, errors.ErrIncorrectSearchQuery
	}

	var filterString string
	if searchQuery.Filter != nil {
		filterString = u.ProductRepo.CreateFilterString(&productRepo.DtoProductFilter{
			MinPrice:   searchQuery.Filter.MinPrice,
			MaxPrice:   searchQuery.Filter.MaxPrice,
			IsNew:      searchQuery.Filter.IsNew,
			IsRating:   searchQuery.Filter.IsRating,
			IsDiscount: searchQuery.Filter.IsDiscount,
		})
	}

	// Max count pages for this search
	countPages, err := u.ProductRepo.GetCountSearchPages(
		&productRepo.DtoSearchPages{
			Category:     searchQuery.Category,
			Count:        searchQuery.Count,
			SearchString: searchQuery.QueryString,
			FilterString: filterString,
		},
	)
	if err != nil {
		return nil, errors.ErrIncorrectSearchQuery
	}

	// Keys for sort items in result of search
	sortString, err := u.ProductRepo.CreateSortString(&productRepo.DtoSortString{
		SortKey:       searchQuery.SortKey,
		SortDirection: searchQuery.SortDirection,
	})
	if err != nil {
		return nil, errors.ErrIncorrectSearchQuery
	}

	// Get range of products
	products, err := u.ProductRepo.SearchRangeProducts(
		&productRepo.DtoSearchQuery{
			QueryString: searchQuery.QueryString,
			PageNum:     searchQuery.PageNum,
			Count:       searchQuery.Count,
			Category:    searchQuery.Category,
			Filter: &productRepo.DtoProductFilter{
				MinPrice:   searchQuery.Filter.MinPrice,
				MaxPrice:   searchQuery.Filter.MaxPrice,
				IsNew:      searchQuery.Filter.IsNew,
				IsRating:   searchQuery.Filter.IsRating,
				IsDiscount: searchQuery.Filter.IsDiscount,
			},
			SortKey:       searchQuery.SortKey,
			SortDirection: searchQuery.SortDirection,
		},
		&productRepo.DtoRageProducts{
			SortString:   sortString,
			FilterString: filterString,
		},
	)
	if err != nil {
		return nil, errors.ErrIncorrectPaginator
	}

	var rangeProducts []*usecase.ViewProduct
	for _, item := range products {
		rangeProducts = append(rangeProducts, &usecase.ViewProduct{
			Id:    item.Id,
			Title: item.Title,
			Price: usecase.ProductPrice{
				Discount:  item.Price.Discount,
				BaseCost:  item.Price.BaseCost,
				TotalCost: item.Price.TotalCost,
			},
			Rating:       item.Rating,
			CountReviews: item.CountReviews,
			PreviewImage: item.PreviewImage,
		})
	}
	return &usecase.RangeProducts{
		ListPreviewProducts: rangeProducts,
		MaxCountPages:       countPages,
	}, nil
}
