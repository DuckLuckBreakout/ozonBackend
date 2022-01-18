package api

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"

type ApiPreviewCart struct {
	Products []*usecase.PreviewCartArticle `json:"products" valid:"notnull"`
	Price    usecase.TotalPrice            `json:"price" valid:"notnull"`
}

type ApiProductPosition struct {
	Count uint64 `json:"count"`
}

type ApiProductIdentifier struct {
	ProductId uint64 `json:"product_id"`
}

type ApiCartArticle struct {
	ApiProductPosition
	ApiProductIdentifier
}
