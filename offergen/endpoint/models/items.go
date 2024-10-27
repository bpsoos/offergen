package models

import "github.com/google/uuid"

type Item struct {
	ID    uuid.UUID `form:"ID" validate:"required"`
	Price uint32    `form:"Price" validate:"required"`
	Name  string    `form:"Name" validate:"required,max=500"`
}

type ItemsForm struct {
	Items map[string]Item `form:"items" validate:"required"`
}

type AddItemInput struct {
	Name  string `form:"Name" validate:"required,max=500"`
	Price uint32 `form:"Price" validate:"required"`
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
