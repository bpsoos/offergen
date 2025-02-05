package service

import (
	"offergen/endpoint/models"

	"github.com/google/uuid"
)

type (
	InventoryManager struct {
		itemPersister      ItemPersister
		inventoryPersister InventoryPersister
	}

	InventoryManagerDeps struct {
		ItemPersister      ItemPersister
		InventoryPersister InventoryPersister
	}

	InventoryPersister interface {
		Create(inventory *models.Inventory) error
		Update(input *UpdateInventoryInput) error
	}

	UpdateInventoryInput struct {
		Title       string
		IsPublished bool
	}

	ItemPersister interface {
		Create(item *models.Item, ownerID string) error
		BatchGet(from, amount uint, ownerID string) ([]models.Item, error)
		Delete(itemID, ownerID string) error
		ItemCount(ownerID string) (int, error)
	}
)

func NewInventoryManager(deps InventoryManagerDeps) *InventoryManager {
	return &InventoryManager{
		itemPersister:      deps.ItemPersister,
		inventoryPersister: deps.InventoryPersister,
	}
}

func (im *InventoryManager) CreateItem(input *models.AddItemInput, ownerID string) (*models.Item, error) {
	item := &models.Item{
		ID:    uuid.New(),
		Price: input.Price,
		Name:  input.Name,
	}

	if err := im.itemPersister.Create(
		item,
		ownerID,
	); err != nil {
		return nil, err
	}

	return item, nil
}

func (im *InventoryManager) BatchGetItem(from, amount uint, ownerID string) ([]models.Item, error) {
	return im.itemPersister.BatchGet(from, amount, ownerID)
}

func (im *InventoryManager) DeleteItem(itemID, ownerID string) error {
	return im.itemPersister.Delete(itemID, ownerID)
}

func (im *InventoryManager) ItemCount(ownerID string) (int, error) {
	return im.itemPersister.ItemCount(ownerID)
}

func (im *InventoryManager) CreateInventory(inventory *models.Inventory) error {
	return im.inventoryPersister.Create(inventory)
}

func (im *InventoryManager) UpdateInventory(input *models.UpdateInventoryInput) error {
	return nil
}
