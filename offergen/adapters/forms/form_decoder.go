package forms

import (
	"net/url"

	"offergen/common_deps"
	"offergen/logging"

	"github.com/go-playground/form"
)

type FormDecoder struct {
	decoder decoderI
}

type decoderI interface {
	Decode(v interface{}, values url.Values) error
}

type FieldError struct {
	name  string
	error error
}

func (fe FieldError) Error() string {
	return fe.error.Error()
}

func (fe FieldError) Field() string {
	return fe.name
}

func parseDecodeErrors(err *form.DecodeErrors) []common_deps.FieldError {
	if err == nil {
		return nil
	}
	if len(*err) == 0 {
		panic("len of err must not be 0")
	}

	parsedErrors := make([]common_deps.FieldError, len(*err))
	count := 0
	for field, err := range *err {
		parsedErrors[count] = FieldError{
			name:  field,
			error: err,
		}
		count++
	}

	return parsedErrors
}

func NewDecoder() *FormDecoder {
	return &FormDecoder{
		decoder: form.NewDecoder(),
	}
}

var logger = logging.GetLogger()

func (decoder FormDecoder) Decode(v interface{}, values url.Values) error {
	return decoder.decoder.Decode(v, values)
}

func (decoder FormDecoder) MustParseDecodeErrors(err error) []common_deps.FieldError {
	if err == nil {
		logger.Error("err was nil")
		panic("err was nil")
	}
	if err, ok := err.(*form.InvalidDecoderError); ok {
		logger.Error("Invalid decode error", "errMsg", err.Error())
		panic(err)
	}

	errDecode, ok := err.(form.DecodeErrors)
	if !ok {
		logger.Error("unknown error", "errMsg", err.Error())
		panic(err)
	}

	return parseDecodeErrors(&errDecode)
}
