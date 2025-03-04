package persistence

import (
	"fmt"
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

	Category struct {
		OwnerID string `db:"owner_id"`
		Name    string `db:"name"`
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
	var updatedInv Inventory
	err := ip.db.Get(
		&updatedInv,
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
		OwnerID:     updatedInv.OwnerID,
		Title:       updatedInv.Title,
		IsPublished: updatedInv.IsPublished,
	}, nil
}

func (ip *InventoryPersister) Get(ownerID string) (*models.Inventory, error) {
	var inv Inventory
	err := ip.db.Get(
		&inv,
		`
            SELECT owner_id,title,is_published FROM inventories
            WHERE owner_id=$1
        `,
		ownerID,
	)
	if err != nil {
		return nil, err
	}

	return &models.Inventory{
		OwnerID:     inv.OwnerID,
		Title:       inv.Title,
		IsPublished: inv.IsPublished,
	}, nil
}

func (ip *InventoryPersister) CreateCategory(ownerID string, category string) error {
	_, err := ip.db.NamedExec(
		`
            INSERT INTO categories (owner_id, name) VALUES (:owner_id, :name)
        `,
		&Category{
			OwnerID: ownerID,
			Name:    category,
		},
	)
	if err != nil {
		return fmt.Errorf("insert add category error: %v", err)
	}

	return nil
}

type CountedCategory struct {
	Name  string `db:"name"`
	Count int    `db:"item_count"`
}

func (ip *InventoryPersister) BatchGetCategory(ownerID string) ([]string, error) {
	categories := make([]string, 0)
	err := ip.db.Select(
		&categories,
		`
            SELECT name FROM categories WHERE categories.owner_id=$1;
        `,
		ownerID,
	)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, nil
	}

	return categories, nil

}

func (ip *InventoryPersister) BatchGetCountedCategory(ownerID string) ([]models.CountedCategory, error) {
	categories := make([]CountedCategory, 0)
	err := ip.db.Select(
		&categories,
		`
            SELECT
                c.name,
                count(items.owner_id) as item_count
            FROM (SELECT name, owner_id FROM categories WHERE categories.owner_id=$1) as c
            LEFT JOIN items ON items.owner_id=c.owner_id and items.category=c.name
            GROUP BY c.name, items.owner_id;
        `,
		ownerID,
	)
	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, nil
	}

	parsedCategories := make([]models.CountedCategory, len(categories))
	for i := range categories {
		parsedCategories[i].Count = categories[i].Count
		parsedCategories[i].Name = categories[i].Name
	}

	return parsedCategories, nil

}
