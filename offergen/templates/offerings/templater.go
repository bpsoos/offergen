package offerings

import (
	"context"
	"io"
	"offergen/endpoint/models"
	"offergen/templates/layouts"
)

type Templater struct{}

func NewTemplater() *Templater {
	return &Templater{}
}

func (t *Templater) Offering(
	ctx context.Context,
	w io.Writer,
	title string,
	items []models.Item,
) error {
	return layouts.Base(
		Offering(title, items),
	).Render(ctx, w)
}
