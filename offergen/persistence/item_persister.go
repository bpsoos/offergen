package persistence

import (
	"errors"
	"offergen/common_deps"
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
	ID      string `db:"id"`
	OwnerID string `db:"owner_id"`
	Price   uint32 `db:"price"`
	Name    string `db:"name"`
}

func (im *ItemPersister) Create(item *models.Item, ownerID string) error {
	_, err := im.db.NamedExec(
		`INSERT INTO items (id,owner_id,price,name) VALUES (:id,:owner_id,:price,:name)`,
		DBItem{
			ID:      item.ID.String(),
			OwnerID: ownerID,
			Price:   item.Price,
			Name:    item.Name,
		},
	)

	return err
}

func (im *ItemPersister) BatchGet(from, amount uint, ownerID string) ([]models.Item, error) {
	dbItems := make([]DBItem, 0, amount)
	err := im.db.Select(
		&dbItems,
		`SELECT id, name, price FROM items WHERE owner_id=$1 LIMIT $2 OFFSET $3`,
		ownerID,
		amount,
		from,
	)
	if err != nil {
		return nil, err
	}

	items := make([]models.Item, len(dbItems))
	for i, item := range dbItems {
		items[i] = models.Item{
			ID:    uuid.MustParse(item.ID),
			Price: item.Price,
			Name:  item.Name,
		}
	}

	return items, nil
}

func (im *ItemPersister) Delete(itemID, ownerID string) error {
	result, err := im.db.Exec(`DELETE FROM items WHERE id=$1 AND owner_id=$2`, itemID, ownerID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return common_deps.ErrItemNotFound
	}
	if count > 1 {
		return errors.New("unexpected rows affected")
	}

	return nil
}

func (im *ItemPersister) ItemCount(ownerID string) (int, error) {
	var count int
	row := im.db.QueryRow(`SELECT COUNT(id) FROM items WHERE owner_id=$1`, ownerID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
