package persistence

import (
	"offergen/endpoint/models"
	"offergen/service"

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

func (ip *InventoryPersister) Create(inventory *models.Inventory) error {
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

	return err
}

func (ip *InventoryPersister) Update(input *service.UpdateInventoryInput) error {
	return nil
}
