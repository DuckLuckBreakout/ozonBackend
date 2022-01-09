package errors

import (
	"github.com/DuckLuckBreakout/web/backend/pkg/errors"
)

func CreateError(err error) error {
	return errors.CreateError(err)
}

var (
	ErrInternalError error = errors.Error{
		Message: "something went wrong",
	}
	ErrDBInternalError error = errors.Error{
		Message: "internal db error",
	}
)
