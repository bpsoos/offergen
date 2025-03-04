package inventory

import (
	"context"
	"io"
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
		inventoryTemplater InventoryTemplater
		errorTemplater     ErrorTemplater
		inventoryManager   InventoryManager
	}

	Deps struct {
		FormDecoder        common_deps.FormDecoder
		StructValidator    common_deps.StructValidator
		Renderer           common_deps.Renderer
		TokenVerifier      TokenVerifier
		InventoryTemplater InventoryTemplater
		ErrorTemplater     ErrorTemplater
		InventoryManager   InventoryManager
	}

	TokenVerifier interface {
		GetUserID(ctx *fiber.Ctx) string
	}

	InventoryTemplater interface {
		Items(
			ctx context.Context,
			w io.Writer,
			items []models.Item,
		) error
		ItemCreator(
			ctx context.Context,
			w io.Writer,
			userID string,
		) error
		Categories(
			ctx context.Context,
			w io.Writer,
			userID string,
			categories []models.CountedCategory,
		) error
		CreateCategoryForm(
			ctx context.Context,
			w io.Writer,
		) error
		CreateCategoryInitLink(
			ctx context.Context,
			w io.Writer,
		) error
		SettingsPage(
			ctx context.Context,
			w io.Writer,
			userID string,
			inv *models.Inventory,
		) error
		Paginator(
			ctx context.Context,
			w io.Writer,
			current, last int,
		) error
		InventoryDetails(
			ctx context.Context,
			w io.Writer,
			inv *models.Inventory,
		) error
		Inventory(
			ctx context.Context,
			w io.Writer,
			userID string,
		) error
	}

	ErrorTemplater interface {
		ErrorMessage(message string, styles ...string) templ.Component
	}

	InventoryManager interface {
		CreateItem(item *models.AddItemInput, ownerID string) (*models.Item, error)
		BatchGetItem(ownerID string, input *models.GetItemsInput) ([]models.Item, error)
		ItemCount(ownerID string) (int, error)
		DeleteItem(itemID, ownerID string) error
		CreateInventory(inventory *models.Inventory) (*models.Inventory, error)
		GetInventory(ownerID string) (*models.Inventory, error)
		UpdateInventory(ownerID string, input *models.UpdateInventoryInput) (*models.Inventory, error)
		CreateCategory(ownerID string, category string) error
		BatchGetCategory(ownerID string) ([]string, error)
		BatchGetCountedCategory(ownerID string) ([]models.CountedCategory, error)
	}
)

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		formDecoder:        deps.FormDecoder,
		structValidator:    deps.StructValidator,
		tokenVerifier:      deps.TokenVerifier,
		inventoryManager:   deps.InventoryManager,
		renderer:           deps.Renderer,
		inventoryTemplater: deps.InventoryTemplater,
		errorTemplater:     deps.ErrorTemplater,
	}
}

var logger = logging.GetLogger()
