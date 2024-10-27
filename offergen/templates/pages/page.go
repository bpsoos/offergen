package pages

import "offergen/templates/components/inventory"

type PageTemplater struct {
	cacheID            string
	inventoryTemplater *inventory.InventoryTemplater
}

type Deps struct {
	InventoryTemplater *inventory.InventoryTemplater
}

func NewPageTemplater(cacheID string, deps *Deps) *PageTemplater {
	return &PageTemplater{
		cacheID:            cacheID,
		inventoryTemplater: deps.InventoryTemplater,
	}
}
