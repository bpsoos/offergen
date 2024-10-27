package validation

import (
	"offergen/common_deps"
	"offergen/logging"

	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validator validatorI
}

type validatorI interface {
	Struct(s interface{}) error
}

func NewStructValidator() *StructValidator {
	return &StructValidator{
		validator: validator.New(
			validator.WithRequiredStructEnabled(),
		),
	}
}

var logger = logging.GetLogger()

func (sv StructValidator) Validate(v interface{}) error {
	return sv.validator.Struct(v)
}

func (sv StructValidator) MustParseValidationErrors(err error) []common_deps.FieldError {
	if err == nil {
		logger.Error("err was nil")
		panic(err)
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		logger.Error("invalid validation error", "errMsg", err.Error())
		panic(err)
	}
	errValidation, ok := err.(validator.ValidationErrors)
	if !ok {
		logger.Info("unknown validation error", "errMsg", err.Error())
		panic(err)
	}
	if len(errValidation) == 0 {
		panic("len of err must not be 0")
	}

	parsedErrors := make([]common_deps.FieldError, len(errValidation))
	for i, fieldError := range errValidation {
		parsedErrors[i] = fieldError.(common_deps.FieldError)
	}

	return parsedErrors
}
