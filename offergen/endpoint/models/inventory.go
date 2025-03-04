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

	CreateCategoryInput struct {
		Name string `json:"name" form:"name" validate:"userinput,max=50"`
	}

	CountedCategory struct {
		Name  string
		Count int
	}
)
