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
	_, err := ip.db.NamedExec(`
        INSERT INTO inventory (owner_id,title,is_published)
        VALUES (:owner_id,:title,:is_published)
        `,
		&Inventory{
			OwnerID:     inventory.OwnerID,
			Title:       inventory.Title,
			IsPublished: inventory.IsPublished,
		},
	)

	return nil, err
}

func (ip *InventoryPersister) Update(ownerID string, input *models.UpdateInventoryInput) (*models.Inventory, error) {
	return nil, nil
}

func (ip *InventoryPersister) Get(ownerID string) (*models.Inventory, error) {
	return nil, nil
}
