package users

import (
	"offergen/common_deps"
	"offergen/endpoint/models"
	"offergen/logging"

	"github.com/gofiber/fiber/v2"
)

type (
	Handler struct {
		structValidator  common_deps.StructValidator
		tokenVerifier    TokenVerifier
		userManager      UserManager
		inventoryManager InventoryManager
		authCookieName   string
	}

	Config struct {
		AuthCookieName string
	}

	Deps struct {
		StructValidator  common_deps.StructValidator
		TokenVerifier    TokenVerifier
		UserManager      UserManager
		InventoryManager InventoryManager
	}

	UserManager interface {
		Save(id, address string) error
		GetEmail(id string) (string, error)
		Delete(id string, token []byte) error
	}

	InventoryManager interface {
		CreateInventory(inventory *models.Inventory) (*models.Inventory, error)
	}

	Authenticator interface {
		DeleteUser(authToken []byte) error
	}

	TokenVerifier interface {
		GetWebhookClaims(rawToken string) (data map[string]interface{}, evt string, err error)
		GetUserID(ctx *fiber.Ctx) string
		GetUserToken(ctx *fiber.Ctx) []byte
	}

	SaveUserParams struct {
		AuthID  string
		Address string
	}
)

func NewHandler(config *Config, deps *Deps) *Handler {
	return &Handler{
		authCookieName:   config.AuthCookieName,
		structValidator:  deps.StructValidator,
		tokenVerifier:    deps.TokenVerifier,
		inventoryManager: deps.InventoryManager,
		userManager:      deps.UserManager,
	}
}

var logger = logging.GetLogger()
