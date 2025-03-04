package models

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Item struct {
	ID       uuid.UUID `json:"id" form:"id" validate:"required"`
	Price    uint32    `json:"price" form:"price" validate:"required"`
	Name     string    `json:"name" form:"name" validate:"required,userinput,max=150"`
	Desc     string    `json:"desc" form:"desc" validate:"userinput,max=500"`
	Category string    `json:"category" form:"category" validate:"userinput,max=50"`
}

type ItemsForm struct {
	Items map[string]Item `form:"items" validate:"required,dive"`
}

type AddItemInput struct {
	Name  string `json:"name" form:"name" validate:"required,max=150"`
	Price uint32 `json:"price" form:"price" validate:"required"`
	Desc  string `json:"desc" form:"desc" validate:"userinput,max=500"`
}

type DeleteItemInput struct {
	ItemID string `validate:"required,uuid"`
}

type GetItemsInput struct {
	From     uint   `query:"from"`
	Amount   uint   `query:"amount" validate:"oneof=10"`
	Category string `query:"category" validate:"userinput,max=50"`
}

type GetItemPagesInput struct {
	Current uint `query:"current"`
}

type ValidatorI interface {
	RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) error
}

var AllowedUserInputPattern = `[A-Za-z0-9 .,!?;'":()&\-]*`

func RegisterUserInputCheck(validator ValidatorI) {
	err := validator.RegisterValidation("userinput", userInputCheck)
	if err != nil {
		panic(err)
	}
}

func userInputCheck(fl validator.FieldLevel) bool {
	r, err := regexp.Compile(AllowedUserInputPattern)
	if err != nil {
		panic(err)
	}

	return r.MatchString(fl.Field().String())
}
