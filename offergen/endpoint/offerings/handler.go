package offerings

import (
	"context"
	"io"
	"offergen/common_deps"
	"offergen/endpoint/models"
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
		BatchGetItem(
			ownerID string,
			input *models.GetItemsInput,
		) ([]models.Item, error)
		GetInventory(ownerID string) (*models.Inventory, error)
	}

	OfferingTemplater interface {
		Offering(
			ctx context.Context,
			w io.Writer,
			title string,
			items []models.Item,
		) error
	}

	OfferingManager interface {
		GetOffering(ownerID string) ([]models.Item, error)
	}
)

func NewHandler(deps *Dependencies) *Handler {
	return &Handler{
		inventoryManager:  deps.InventoryManager,
		offeringTemplater: deps.OfferingTemplater,
		renderer:          deps.Renderer,
	}
}
