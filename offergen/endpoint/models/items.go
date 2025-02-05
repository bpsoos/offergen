package models

import "github.com/google/uuid"

type Item struct {
	ID    uuid.UUID `json:"id" form:"ID" validate:"required"`
	Price uint32    `json:"price" form:"Price" validate:"required"`
	Name  string    `json:"name" form:"Name" validate:"required,max=500"`
}

type ItemsForm struct {
	Items map[string]Item `form:"items" validate:"required"`
}

type AddItemInput struct {
	Name  string `form:"Name" validate:"required,max=500" json:"name"`
	Price uint32 `form:"Price" validate:"required" json:"price"`
}

type DeleteItemInput struct {
	ItemID string `validate:"required,uuid"`
}

type GetItemsInput struct {
	From   uint `query:"from"`
	Amount uint `query:"amount" validate:"oneof=10"`
}

type GetItemPagesInput struct {
	Current uint `query:"current" validate:"gt=0"`
}

const AllowedNamePattern = `^[a-zA-Z0-9 ]+$`
