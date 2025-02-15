package inventory

import (
	"offergen/common_deps"
	"offergen/endpoint/models"
	"offergen/logging"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type (
	Handler struct {
		formDecoder        common_deps.FormDecoder
		structValidator    common_deps.StructValidator
		renderer           common_deps.Renderer
		tokenVerifier      TokenVerifier
		pageTemplater      PageTemplater
		inventoryTemplater InventoryTemplater
		errorTemplater     ErrorTemplater
		inventoryManager   InventoryManager
	}

	Deps struct {
		FormDecoder        common_deps.FormDecoder
		StructValidator    common_deps.StructValidator
		Renderer           common_deps.Renderer
		TokenVerifier      TokenVerifier
		PageTemplater      PageTemplater
		InventoryTemplater InventoryTemplater
		ErrorTemplater     ErrorTemplater
		InventoryManager   InventoryManager
	}

	TokenVerifier interface {
		GetUserID(ctx *fiber.Ctx) string
	}

	PageTemplater interface {
		Inventory(userID string) templ.Component
	}

	InventoryTemplater interface {
		Items(items []models.Item) templ.Component
		ItemCreator() templ.Component
		Paginator(current, last int) templ.Component
		SettingsPage(inv *models.Inventory) templ.Component
		InventoryDetails(inv *models.Inventory) templ.Component
		ItemsPage(userID string) templ.Component
	}

	ErrorTemplater interface {
		ErrorMessage(message string, styles ...string) templ.Component
	}

	InventoryManager interface {
		CreateItem(item *models.AddItemInput, ownerID string) (*models.Item, error)
		BatchGetItem(from, amount uint, ownerID string) ([]models.Item, error)
		ItemCount(ownerID string) (int, error)
		DeleteItem(itemID, ownerID string) error
		CreateInventory(inventory *models.Inventory) (*models.Inventory, error)
		GetInventory(ownerID string) (*models.Inventory, error)
		UpdateInventory(ownerID string, input *models.UpdateInventoryInput) (*models.Inventory, error)
	}
)

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		formDecoder:        deps.FormDecoder,
		structValidator:    deps.StructValidator,
		tokenVerifier:      deps.TokenVerifier,
		inventoryManager:   deps.InventoryManager,
		renderer:           deps.Renderer,
		pageTemplater:      deps.PageTemplater,
		inventoryTemplater: deps.InventoryTemplater,
		errorTemplater:     deps.ErrorTemplater,
	}
}

var logger = logging.GetLogger()
