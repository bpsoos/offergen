package persistence

import (
	"database/sql"
	"errors"
	"offergen/common_deps"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserPersister struct {
	db *sqlx.DB
}

func NewUserPersister(db *sqlx.DB) *UserPersister {
	return &UserPersister{db: db}
}

func (up *UserPersister) Save(id, address string) error {
	if err := uuid.Validate(id); err != nil {
		return err
	}

	if err := uuid.Validate(id); err != nil {
		return err
	}

	_, err := up.db.NamedExec(
		`INSERT INTO users (id,email) VALUES (:id,:email)`,
		map[string]interface{}{
			"id":    id,
			"email": address,
		},
	)

	return err
}

func (up *UserPersister) GetEmail(id string) (string, error) {
	row := up.db.QueryRow(`SELECT email FROM users WHERE id=$1`, id)

	var email string
	err := row.Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", common_deps.ErrUserNotFound
		}

		return "", err
	}

	return email, nil
}

func (up *UserPersister) Delete(id string) error {
	result, err := up.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return common_deps.ErrUserNotFound
	}
	if count > 1 {
		return errors.New("unexpected rows affected")
	}

	return nil
}
