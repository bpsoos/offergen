package persistence

import (
	"database/sql"
	"fmt"
	"offergen/endpoint/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ItemPersister struct {
	db *sqlx.DB
}

func NewItemPersister(db *sqlx.DB) *ItemPersister {
	return &ItemPersister{db: db}
}

type DBItem struct {
	ID       string         `db:"id"`
	OwnerID  string         `db:"owner_id"`
	Desc     sql.NullString `db:"description"`
	Category sql.NullString `db:"category"`
	Price    uint32         `db:"price"`
	Name     string         `db:"name"`
}

const (
	createItemQuery = `
        INSERT INTO items (id,owner_id,price,name,description)
        VALUES (:id,:owner_id,:price,:name,:description)
    `
	createItemQueryWithCategory = `
        INSERT INTO items (id,owner_id,price,name,description,category)
        SELECT :id,:owner_id,:price,:name,:description,:category
        WHERE EXISTS (
            SELECT 1 FROM categories WHERE owner_id=:owner_id and name=:category
        )
    `
)

func (im *ItemPersister) Create(item *models.Item, ownerID string) error {
	query := createItemQuery
	params := DBItem{
		ID:      item.ID.String(),
		OwnerID: ownerID,
		Price:   item.Price,
		Name:    item.Name,
	}
	if item.Desc != "" {
		params.Desc = sql.NullString{
			String: item.Desc,
			Valid:  true,
		}
	}
	if item.Category != "" {
		query = createItemQueryWithCategory
		params.Category = sql.NullString{
			String: item.Category,
			Valid:  true,
		}
	}

	result, err := im.db.NamedExec(
		query,
		params,
	)
	if err != nil {
		return fmt.Errorf("named exec: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %v", err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("unexpected affected row count: %d", rowsAffected)
	}

	return nil
}

const (
	batchGetQuery = `
        SELECT id, name, price, description
        FROM items
        WHERE
            owner_id=:owner_id
        LIMIT :limit OFFSET :offset
    `
	batchGetQueryByCategory = `
        SELECT id, name, price, description, category
        FROM items
        WHERE
            owner_id=:owner_id
            and category=:category
        LIMIT :limit OFFSET :offset
    `
)

func (im *ItemPersister) BatchGet(
	ownerID string,
	input *models.GetItemsInput,
) ([]models.Item, error) {
	dbItems := make([]DBItem, 0, input.Amount)
	query := batchGetQuery
	params := map[string]interface{}{
		"owner_id": ownerID,
		"limit":    input.Amount,
		"offset":   input.From,
	}
	if input.Category != "" {
		query = batchGetQueryByCategory
		params["category"] = input.Category
	}
	stmt, err := im.db.PrepareNamed(query)
	if err != nil {
		return nil, fmt.Errorf("prepare named: %v", err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			panic(err)
		}
	}()

	err = stmt.Select(
		&dbItems,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("select: %v", err)
	}

	items := make([]models.Item, len(dbItems))
	for i := range dbItems {
		items[i] = models.Item{
			ID:    uuid.MustParse(dbItems[i].ID),
			Price: dbItems[i].Price,
			Name:  dbItems[i].Name,
		}
		if dbItems[i].Category.Valid {
			items[i].Category = dbItems[i].Category.String
		}
		if dbItems[i].Desc.Valid {
			items[i].Desc = dbItems[i].Desc.String
		}
	}

	return items, nil
}

func (im *ItemPersister) Delete(itemID, ownerID string) error {
	result, err := im.db.Exec(`DELETE FROM items WHERE id=$1 AND owner_id=$2`, itemID, ownerID)
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %v", err)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("unexpected affected row count: %d", rowsAffected)
	}

	return nil
}

func (im *ItemPersister) ItemCount(ownerID string) (int, error) {
	var count int
	row := im.db.QueryRow(`SELECT COUNT(id) FROM items WHERE owner_id=$1`, ownerID)
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("scan: %v", err)
	}

	return count, nil
}
