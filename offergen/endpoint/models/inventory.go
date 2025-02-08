package models

type (
	UpdateInventoryInput struct {
		Title       string `form:"Title" validate:"required"`
		IsPublished bool   `form:"Published"`
	}

	Inventory struct {
		OwnerID     string
		Title       string
		IsPublished bool
	}
)
