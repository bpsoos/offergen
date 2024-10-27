package models

type Event struct {
	Token string `json:"token" validate:"required,jwt"`
	Event string `json:"event" validate:"required,oneof=user.create"`
}
