package dto

import "github.com/DuckLuckBreakout/ozonBackend/internal/pkg/models/usecase"

type DtoCategoryId struct {
	Id uint64 `json:"id"`
}

type DtoCategoryLevel struct {
	Level uint64 `json:"level"`
}

type DtoCategoriesCatalog struct {
	Catalog []*usecase.CategoriesCatalog
}

type DtoBranchBorders struct {
	Left  uint64 `json:"left"`
	Right uint64 `json:"right"`
}
