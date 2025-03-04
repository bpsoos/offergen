package common_deps

import (
	"errors"
	"net/url"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type FormDecoder interface {
	Decode(s interface{}, values url.Values) error
	MustParseDecodeErrors(err error) []FieldError
}

type StructValidator interface {
	Validate(v interface{}) error
	MustParseValidationErrors(err error) []FieldError
}

type FieldError interface {
	Error() string
	Field() string
}

type Renderer interface {
	Render(ctx *fiber.Ctx, component templ.Component) error
	RenderWithChildren(ctx *fiber.Ctx, children, component templ.Component) error
}

var ErrUserNotFound = errors.New("user not found")
var ErrItemNotFound = errors.New("item not found")
