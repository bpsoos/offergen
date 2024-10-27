package models

type (
	UpdateInventoryInput struct {
		Title       string `form:"Title" validate:"required"`
		IsPublished *bool  `form:"Published" validate:"required"`
	}

	Inventory struct {
		OwnerID     string
		Title       string
		IsPublished bool
	}
)
