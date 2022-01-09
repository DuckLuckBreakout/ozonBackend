package errors

import (
	"github.com/DuckLuckBreakout/ozonBackend/pkg/errors"
)

func CreateError(err error) error {
	return errors.CreateError(err)
}

var (
	ErrInternalError error = errors.Error{
		Message: "something went wrong",
	}
	ErrUserNotFound error = errors.Error{
		Message: "user not found",
	}
	ErrEmailAlreadyExist error = errors.Error{
		Message: "user email already exist",
	}
	ErrDBInternalError error = errors.Error{
		Message: "internal db error",
	}
	ErrHashFunctionFailed error = errors.Error{
		Message: "can't get hash of data",
	}
	ErrIncorrectAuthData error = errors.Error{
		Message: "incorrect auth user data",
	}
)
