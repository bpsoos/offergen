package preview

import (
	"offergen/common_deps"
	"offergen/endpoint/models"
	"offergen/logging"

	"github.com/a-h/templ"
)

type (
	Handler struct {
		formDecoder       common_deps.FormDecoder
		structValidator   common_deps.StructValidator
		renderer          common_deps.Renderer
		offeringTemplater OfferingTemplater
	}

	Deps struct {
		FormDecoder       common_deps.FormDecoder
		StructValidator   common_deps.StructValidator
		Renderer          common_deps.Renderer
		OfferingTemplater OfferingTemplater
	}

	OfferingTemplater interface {
		Offering(items []models.Item) templ.Component
	}
)

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		formDecoder:       deps.FormDecoder,
		structValidator:   deps.StructValidator,
		renderer:          deps.Renderer,
		offeringTemplater: deps.OfferingTemplater,
	}
}

var logger = logging.GetLogger()
