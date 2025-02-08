package offerings

import (
	"offergen/common_deps"
	"offergen/endpoint/models"

	"github.com/a-h/templ"
)

type (
	Handler struct {
		inventoryManager  InventoryManager
		offeringTemplater OfferingTemplater
		renderer          common_deps.Renderer
	}

	Dependencies struct {
		InventoryManager  InventoryManager
		OfferingTemplater OfferingTemplater
		Renderer          common_deps.Renderer
	}

	InventoryManager interface {
		BatchGetItem(from, amount uint, ownerID string) ([]models.Item, error)
		GetInventory(ownerID string) (*models.Inventory, error)
	}

	OfferingTemplater interface {
		Menu(items []models.Item) templ.Component
	}

	OfferingManager interface {
		GetOffering(ownerID string) ([]models.Item, error)
	}
)

func NewHandler(deps *Dependencies) *Handler {
	return &Handler{
		inventoryManager: deps.InventoryManager,
        offeringTemplater: deps.OfferingTemplater,
        renderer: deps.Renderer,
	}
}
