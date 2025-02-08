package persistence

import (
	"offergen/endpoint/models"

	"github.com/jmoiron/sqlx"
)

type (
	InventoryPersister struct {
		db *sqlx.DB
	}

	Inventory struct {
		OwnerID     string `db:"owner_id"`
		Title       string `db:"title"`
		IsPublished bool   `db:"is_published"`
	}
)

func NewInventoryPersister(db *sqlx.DB) *InventoryPersister {
	return &InventoryPersister{db: db}
}

func (ip *InventoryPersister) Create(inventory *models.Inventory) (*models.Inventory, error) {
	var createdInv Inventory
	err := ip.db.Get(
		&createdInv,
		`
            INSERT INTO inventories (owner_id,title,is_published)
            VALUES ($1,$2,$3)
            RETURNING owner_id,title,is_published
        `,
		inventory.OwnerID,
		inventory.Title,
		inventory.IsPublished,
	)
	if err != nil {
		return nil, err
	}

	return &models.Inventory{
		OwnerID:     createdInv.OwnerID,
		Title:       createdInv.Title,
		IsPublished: createdInv.IsPublished,
	}, nil
}

func (ip *InventoryPersister) Update(ownerID string, input *models.UpdateInventoryInput) (*models.Inventory, error) {
	var createdInv Inventory
	err := ip.db.Get(
		&createdInv,
		`
            UPDATE inventories
            SET title=$2, is_published=$3
            WHERE owner_id=$1
            RETURNING owner_ID, title, is_published
        `,
		ownerID,
		input.Title,
		input.IsPublished,
	)
	if err != nil {
		return nil, err
	}

	return &models.Inventory{
		OwnerID:     createdInv.OwnerID,
		Title:       createdInv.Title,
		IsPublished: createdInv.IsPublished,
	}, nil
}

func (ip *InventoryPersister) Get(ownerID string) (*models.Inventory, error) {
	return nil, nil
}
