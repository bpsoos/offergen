package users

import (
	"encoding/json"
	"net/http"
	"offergen/endpoint/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Create(ctx *fiber.Ctx) error {
	event := new(models.Event)

	if err := ctx.BodyParser(event); err != nil {
		logger.Error("could not parse event body for user created webhook")
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	if err := h.structValidator.Validate(event); err != nil {
		logger.Error("invalid event body for user created webhook")
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}

	logger.Info("parsed and validated body", "eventType", event.Event)

	data, evt, err := h.tokenVerifier.GetWebhookClaims(event.Token)
	if err != nil {
		logger.Error("error parsing webhook token", "errMsg", err.Error())
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	logger.Info("successfully parsed webhook token", "evt", evt)

	user := new(User)
	rawUser, err := json.Marshal(data)
	if err != nil {
		logger.Error("error marshaling user", "errMsg", err.Error())
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	err = json.Unmarshal(rawUser, user)
	if err != nil {
		logger.Error("error unmarshaling user", "errMsg", err.Error())
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	if err := h.structValidator.Validate(user); err != nil {
		logger.Error("invalid user")
		return ctx.SendStatus(http.StatusUnprocessableEntity)
	}
	logger.Info(
		"successfully parsed hanko user",
		"id", user.ID,
		"addressID", user.Emails[0].ID,
	)
	if !user.Emails[0].IsVerified {
		logger.Warn("user is not verified")
	}
	if !user.Emails[0].IsPrimary {
		logger.Warn("user address is not primary")
	}

	if err := h.userManager.Save(user.ID, user.Emails[0].Address); err != nil {
		logger.Error("could not save user", "id", user.ID)
	}
	logger.Info("user saved", "userID", user.ID)

	if err := h.inventoryManager.CreateInventory(&models.Inventory{
		OwnerID:     user.ID,
		Title:       "Offering",
		IsPublished: false,
	}); err != nil {
		logger.Error("could not create inventory for user", "id", user.ID)
	}
	logger.Info("inventory saved for user", "userID", user.ID)

	return ctx.SendStatus(http.StatusOK)
}

type User struct {
	ID     string `json:"id" validate:"required,uuid"`
	Emails []struct {
		Address    string    `json:"address" validate:"required,email"`
		ID         string    `json:"id" validate:"required,uuid"`
		IsPrimary  bool      `json:"is_primary"`
		IsVerified bool      `json:"is_verified"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	} `json:"emails" validate:"required,len=1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
