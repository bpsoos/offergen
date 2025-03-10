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
		Create(inventory *models.Inventory) (*models.Inventory, error)
		Get(ownerID string) (*models.Inventory, error)
		Update(ownerId string, input *models.UpdateInventoryInput) (*models.Inventory, error)
		CreateCategory(ownerID string, category string) error
		BatchGetCategory(ownerID string) ([]string, error)
		BatchGetCountedCategory(ownerID string) ([]models.CountedCategory, error)
	}

	ItemPersister interface {
		Create(item *models.Item, ownerID string) error
		BatchGet(ownerID string, input *models.GetItemsInput) ([]models.Item, error)
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
		Desc:  input.Desc,
	}

	if err := im.itemPersister.Create(
		item,
		ownerID,
	); err != nil {
		return nil, err
	}

	return item, nil
}

func (im *InventoryManager) BatchGetItem(ownerID string, input *models.GetItemsInput) ([]models.Item, error) {
	return im.itemPersister.BatchGet(ownerID, input)
}

func (im *InventoryManager) DeleteItem(itemID, ownerID string) error {
	return im.itemPersister.Delete(itemID, ownerID)
}

func (im *InventoryManager) ItemCount(ownerID string) (int, error) {
	return im.itemPersister.ItemCount(ownerID)
}

func (im *InventoryManager) CreateInventory(inventory *models.Inventory) (*models.Inventory, error) {
	return im.inventoryPersister.Create(inventory)
}

func (im *InventoryManager) GetInventory(ownerID string) (*models.Inventory, error) {
	return im.inventoryPersister.Get(ownerID)
}

func (im *InventoryManager) UpdateInventory(ownerID string, input *models.UpdateInventoryInput) (*models.Inventory, error) {
	return im.inventoryPersister.Update(ownerID, input)
}

func (im *InventoryManager) CreateCategory(ownerID string, category string) error {
	return im.inventoryPersister.CreateCategory(ownerID, category)
}

func (im *InventoryManager) BatchGetCategory(ownerID string) ([]string, error) {
	return im.inventoryPersister.BatchGetCategory(ownerID)
}

func (im *InventoryManager) BatchGetCountedCategory(ownerID string) ([]models.CountedCategory, error) {
	return im.inventoryPersister.BatchGetCountedCategory(ownerID)
}
