package preview

import (
	"offergen/common_deps"
	"offergen/logging"
)

type Handler struct {
	formDecoder     common_deps.FormDecoder
	structValidator common_deps.StructValidator
}

type Deps struct {
	FormDecoder     common_deps.FormDecoder
	StructValidator common_deps.StructValidator
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		formDecoder:     deps.FormDecoder,
		structValidator: deps.StructValidator,
	}
}

var logger = logging.GetLogger()
