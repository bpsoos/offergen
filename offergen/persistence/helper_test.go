package persistence_test

import (
	"offergen/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func GetDB() *sqlx.DB {
	config := config.NewConfig()
	db, err := sqlx.Connect("postgres", config.PostgresURL)
	if err != nil {
		panic("ERROR: could not connect to test db")
	}

	return db
}

func CleanDB(db *sqlx.DB) {
	_, err := db.Exec(`
		DELETE FROM inventory;
		DELETE FROM items;
		DELETE FROM users;
    `)
	if err != nil {
		panic(err)
	}
}
