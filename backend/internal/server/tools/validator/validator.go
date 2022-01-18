package validator

import (
	"github.com/DuckLuckBreakout/ozonBackend/internal/server/errors"

	"github.com/asaskevich/govalidator"
	Errors "github.com/pkg/errors"
)

func ValidateStruct(data interface{}) error {
	if _, err := govalidator.ValidateStruct(data); err != nil {
		return Errors.Wrap(errors.ErrInvalidData, err.Error())
	}

	return nil
}
