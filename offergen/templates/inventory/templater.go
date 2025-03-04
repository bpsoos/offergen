package inventory

import (
	"context"
	"io"
	"offergen/endpoint/models"
	"offergen/templates/layouts"
)

type Templater struct {
	publicBaseURL string
}

func NewTemplater(publicBaseURL string) *Templater {
	return &Templater{
		publicBaseURL: publicBaseURL,
	}
}

func (t *Templater) Inventory(
	ctx context.Context,
	w io.Writer,
	userID string,
) error {
	return layouts.Base(
		InventoryBase(userID, InventoryMain(userID)),
	).Render(ctx, w)
}

func (t *Templater) ItemCreator(
	ctx context.Context,
	w io.Writer,
	userID string,
) error {
	return layouts.Base(
		InventoryBase(userID, ItemCreator()),
	).Render(ctx, w)
}

func (t *Templater) SettingsPage(
	ctx context.Context,
	w io.Writer,
	userID string,
	inv *models.Inventory,
) error {
	return layouts.Base(
		InventoryBase(userID, SettingsPage(userID, t.publicBaseURL, inv)),
	).Render(ctx, w)
}

func (t *Templater) Categories(
	ctx context.Context,
	w io.Writer,
	userID string,
	categories []models.CountedCategory,
) error {
	return layouts.Base(
		InventoryBase(userID, Categories(userID, categories)),
	).Render(ctx, w)
}

func (t *Templater) CreateCategoryForm(
	ctx context.Context,
	w io.Writer,
) error {
	return CreateCategoryForm().Render(ctx, w)
}

func (t *Templater) CreateCategoryInitLink(
	ctx context.Context,
	w io.Writer,
) error {
	return CreateCategoryInitLink().Render(ctx, w)
}

func (t *Templater) Items(
	ctx context.Context,
	w io.Writer,
	items []models.Item,
) error {
	return ItemsTable(items).Render(ctx, w)
}

func (t *Templater) Paginator(
	ctx context.Context,
	w io.Writer,
	current, last int,
) error {
	return Paginator(current, last).Render(ctx, w)
}

func (t *Templater) InventoryDetails(
	ctx context.Context,
	w io.Writer,
	inv *models.Inventory,
) error {
	return InventoryDetails(t.publicBaseURL, inv).Render(ctx, w)
}
