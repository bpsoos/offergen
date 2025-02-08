package components

import (
	"offergen/endpoint/models"

	"github.com/a-h/templ"
)

type OfferingTemplater struct{}

func NewOfferingTemplater() *OfferingTemplater {
	return &OfferingTemplater{}
}

func (ot *OfferingTemplater) Menu(items []models.Item) templ.Component {
	return Menu(items)
}
