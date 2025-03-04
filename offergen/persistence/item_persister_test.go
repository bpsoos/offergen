package persistence_test

import (
	"database/sql"
	"offergen/endpoint/models"
	"offergen/persistence"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Item struct {
	ID       string         `db:"id"`
	OwnerID  string         `db:"owner_id"`
	Desc     sql.NullString `db:"description"`
	Category sql.NullString `db:"category"`
	Price    uint32         `db:"price"`
	Name     string         `db:"name"`
}

var _ = Describe("items", func() {
	var db *sqlx.DB

	BeforeEach(func() {
		db = GetDB()
		CleanDB(db)
	})
	var userID string
	BeforeEach(func() {
		userID = uuid.NewString()
		_, err := db.NamedExec(
			`INSERT INTO users (id,email) VALUES (:id,:email);`,
			&User{
				ID:    userID,
				Email: "dummy@email.com",
			},
		)
		Expect(err).ToNot(HaveOccurred())
	})
	Describe("create", func() {
		Context("called with empty db", func() {
			It("should create item", func() {
				item := &models.Item{
					ID:    uuid.New(),
					Name:  "dummyItem",
					Price: 12000,
				}

				err := persistence.NewItemPersister(db).Create(item, userID)
				Expect(err).ToNot(HaveOccurred())

				items := []Item{}
				err = db.Select(&items, `SELECT * FROM items;`)
				Expect(err).ToNot(HaveOccurred())

				Expect(items).To(HaveLen(1))
				Expect(items[0].ID).To(Equal(item.ID.String()))
				Expect(items[0].Name).To(Equal(item.Name))
				Expect(items[0].Price).To(Equal(item.Price))
			})

			It("should create item with descripiton", func() {
				item := &models.Item{
					ID:    uuid.New(),
					Name:  "dummyItem",
					Desc:  "dummy_description",
					Price: 12000,
				}

				err := persistence.NewItemPersister(db).Create(item, userID)
				Expect(err).ToNot(HaveOccurred())

				items := []Item{}
				err = db.Select(&items, `SELECT * FROM items;`)
				Expect(err).ToNot(HaveOccurred())

				Expect(items).To(HaveLen(1))
				Expect(items[0].ID).To(Equal(item.ID.String()))
				Expect(items[0].Name).To(Equal(item.Name))
				Expect(items[0].Price).To(Equal(item.Price))
				Expect(items[0].Desc.String).To(Equal(item.Desc))
			})

			It("should return error when a category is specified", func() {
				item := &models.Item{
					ID:       uuid.New(),
					Name:     "dummyItem",
					Price:    12000,
					Category: "nonexistent_category",
				}

				err := persistence.NewItemPersister(db).Create(item, userID)
				Expect(err).To(HaveOccurred())

				items := []Item{}
				err = db.Select(&items, `SELECT * FROM items;`)
				Expect(err).ToNot(HaveOccurred())

				Expect(items).To(BeEmpty())
			})

			Context("with a category existing for user", func() {
				var category string
				BeforeEach(func() {
					category = "dummy_category"
					_, err := db.NamedExec(
						`INSERT INTO categories (owner_id,name) VALUES (:owner_id,:name)`,
						&Category{
							OwnerID: userID,
							Name:    category,
						},
					)
					Expect(err).ToNot(HaveOccurred())
				})

				It("should create item with category", func() {
					category := "dummy_category"
					item := &models.Item{
						ID:       uuid.New(),
						Name:     "dummyItem",
						Price:    12000,
						Category: category,
					}

					err := persistence.NewItemPersister(db).Create(item, userID)
					Expect(err).ToNot(HaveOccurred())

					items := []Item{}
					err = db.Select(&items, `SELECT * FROM items;`)
					Expect(err).ToNot(HaveOccurred())

					Expect(items).To(HaveLen(1))
					Expect(items[0].ID).To(Equal(item.ID.String()))
					Expect(items[0].Name).To(Equal(item.Name))
					Expect(items[0].Price).To(Equal(item.Price))
					Expect(items[0].Category.String).To(Equal(item.Category))
				})
			})

		})
	})
	Describe("batch get with a user existing", func() {
		Context("called with empty items table", func() {
			It("should return an empty list", func() {
				items, err := persistence.NewItemPersister(db).
					BatchGet(
						userID,
						&models.GetItemsInput{
							From:     0,
							Amount:   10,
							Category: "",
						},
					)

				Expect(err).ToNot(HaveOccurred())
				Expect(items).To(HaveLen(0))
			})
		})

		Context("called with items table containing an item without description", func() {
			var itemID string
			BeforeEach(func() {
				itemID = uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
					&Item{
						ID:      itemID,
						OwnerID: userID,
						Name:    "dummy item name",
						Price:   12000,
					},
				)
				Expect(err).ToNot(HaveOccurred())
			})
			It("should return the item", func() {
				items, err := persistence.NewItemPersister(db).
					BatchGet(
						userID,
						&models.GetItemsInput{
							From:     0,
							Amount:   10,
							Category: "",
						},
					)

				Expect(err).ToNot(HaveOccurred())
				Expect(items).To(HaveLen(1))
				Expect(items[0].ID.String()).To(Equal(itemID))
				Expect(items[0].Name).To(Equal("dummy item name"))
				Expect(items[0].Price).To(Equal(uint32(12000)))
				Expect(items[0].Desc).To(Equal(""))
			})
		})
		Context("called with items table containing an item with description", func() {
			var itemID string
			BeforeEach(func() {
				itemID = uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO
                        items (id,owner_id,name,price,description)
                    VALUES (:id,:owner_id,:name,:price,:description);`,
					&Item{
						ID:      itemID,
						OwnerID: userID,
						Name:    "dummy item name",
						Price:   12000,
						Desc: sql.NullString{
							String: "dummy description",
							Valid:  true,
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())
			})
			It("should return the item", func() {
				items, err := persistence.NewItemPersister(db).
					BatchGet(
						userID,
						&models.GetItemsInput{
							From:     0,
							Amount:   10,
							Category: "",
						},
					)

				Expect(err).ToNot(HaveOccurred())
				Expect(items).To(HaveLen(1))
				Expect(items[0].ID.String()).To(Equal(itemID))
				Expect(items[0].Name).To(Equal("dummy item name"))
				Expect(items[0].Price).To(Equal(uint32(12000)))
				Expect(items[0].Desc).To(Equal("dummy description"))
			})

			Context("with another user existing having an item", func() {
				BeforeEach(func() {
					otherUserID := uuid.NewString()
					_, err := db.NamedExec(
						`INSERT INTO users (id,email) VALUES (:id,:email);`,
						&User{
							ID:    otherUserID,
							Email: "other.dummy@email.com",
						},
					)
					Expect(err).ToNot(HaveOccurred())

					_, err = db.NamedExec(
						`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
						[]Item{
							{
								ID:      uuid.NewString(),
								OwnerID: otherUserID,
								Name:    "other dummy item name",
								Price:   24000,
							},
						},
					)
					Expect(err).ToNot(HaveOccurred())
				})
				It("should return the correct item", func() {
					items, err := persistence.NewItemPersister(db).
						BatchGet(
							userID,
							&models.GetItemsInput{
								From:     0,
								Amount:   10,
								Category: "",
							},
						)

					Expect(err).ToNot(HaveOccurred())
					Expect(items).To(HaveLen(1))
					Expect(items[0].ID.String()).To(Equal(itemID))
					Expect(items[0].Name).To(Equal("dummy item name"))
					Expect(items[0].Price).To(Equal(uint32(12000)))
				})
			})

			Context("with a category specified and an item belonging to that category", func() {
				var (
					category, otherItemID string
				)
				BeforeEach(func() {
					category = "dummy_category"
					_, err := db.NamedExec(
						`INSERT INTO categories (owner_id,name) VALUES (:owner_id,:name)`,
						&Category{
							OwnerID: userID,
							Name:    category,
						},
					)
					Expect(err).ToNot(HaveOccurred())

					otherItemID = uuid.NewString()
					_, err = db.NamedExec(
						`INSERT INTO
                            items (id,owner_id,name,price,category)
                        VALUES (:id,:owner_id,:name,:price,:category)`,
						[]Item{
							{
								ID:      otherItemID,
								OwnerID: userID,
								Name:    "other dummy item name",
								Category: sql.NullString{
									String: category,
									Valid:  true,
								},
								Price: 24000,
							},
						},
					)
					Expect(err).ToNot(HaveOccurred())
				})
				It("should return the correct item", func() {
					items, err := persistence.NewItemPersister(db).
						BatchGet(
							userID,
							&models.GetItemsInput{
								From:     0,
								Amount:   10,
								Category: category,
							},
						)

					Expect(err).ToNot(HaveOccurred())
					Expect(items).To(HaveLen(1))
					Expect(items[0].ID.String()).To(Equal(otherItemID))
					Expect(items[0].Name).To(Equal("other dummy item name"))
					Expect(items[0].Price).To(Equal(uint32(24000)))
					Expect(items[0].Category).To(Equal(category))
				})
			})

			Context("called with items table containing items that fall out of the specified range", func() {
				It("should return only the items in range", func() {
					otherItemID := uuid.NewString()
					_, err := db.NamedExec(
						`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
						[]Item{
							{
								ID:      otherItemID,
								OwnerID: userID,
								Name:    "dummy item name 4",
								Price:   14000,
							},
							{
								ID:      uuid.NewString(),
								OwnerID: userID,
								Name:    "dummy item name 1",
								Price:   11000,
							},
							{
								ID:      uuid.NewString(),
								OwnerID: userID,
								Name:    "dummy item name 2",
								Price:   12000,
							},
						},
					)
					Expect(err).ToNot(HaveOccurred())

					fetchedItems, err := persistence.NewItemPersister(db).
						BatchGet(
							userID,
							&models.GetItemsInput{
								From:     0,
								Amount:   2,
								Category: "",
							},
						)

					Expect(err).ToNot(HaveOccurred())
					Expect(fetchedItems).To(HaveLen(2))
					Expect(fetchedItems[0].ID.String()).To(Equal(itemID))
					Expect(fetchedItems[0].Name).To(Equal("dummy item name"))
					Expect(fetchedItems[0].Price).To(Equal(uint32(12000)))
					Expect(fetchedItems[1].ID.String()).To(Equal(otherItemID))
					Expect(fetchedItems[1].Name).To(Equal("dummy item name 4"))
					Expect(fetchedItems[1].Price).To(Equal(uint32(14000)))
				})
			})
		})
	})

	Describe("delete", func() {
		Context("delete called with single item in table", func() {
			It("should delete the item", func() {
				itemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
					&Item{

						ID:      itemID,
						OwnerID: userID,
						Name:    "dummy item name",
						Price:   12000,
					},
				)
				Expect(err).ToNot(HaveOccurred())

				err = persistence.NewItemPersister(db).Delete(itemID, userID)
				Expect(err).ToNot(HaveOccurred())

				items := []Item{}
				err = db.Select(&items, `SELECT * FROM items;`)
				Expect(err).ToNot(HaveOccurred())
				Expect(items).To(HaveLen(0))
			})
		})
		Context("delete called with items from different users in table", func() {
			It("should delete only the item that belongs to the given user", func() {
				otherUserID := uuid.NewString()
				itemID := uuid.NewString()
				otherItemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					&User{
						ID:    otherUserID,
						Email: "otherDummy@email.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				_, err = db.NamedExec(
					`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
					[]Item{
						{
							ID:      itemID,
							OwnerID: userID,
							Name:    "dummy item name",
							Price:   12000,
						},
						{
							ID:      otherItemID,
							OwnerID: otherUserID,
							Name:    "other dummy item name",
							Price:   24000,
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())

				err = persistence.NewItemPersister(db).Delete(itemID, userID)
				Expect(err).ToNot(HaveOccurred())

				items := []Item{}
				err = db.Select(&items, `SELECT * FROM items;`)
				Expect(err).ToNot(HaveOccurred())

				Expect(items).To(HaveLen(1))
				Expect(items[0].ID).To(Equal(otherItemID))
				Expect(items[0].OwnerID).To(Equal(otherUserID))
				Expect(items[0].Price).To(Equal(uint32(24000)))
				Expect(items[0].Name).To(Equal("other dummy item name"))
			})
		})
	})
	Context("delete called with single item in table belonging to a different user", func() {
		It("should delete the item", func() {
			otherUserID := uuid.NewString()
			itemID := uuid.NewString()
			_, err := db.NamedExec(
				`INSERT INTO users (id,email) VALUES (:id,:email);`,
				&User{
					ID:    otherUserID,
					Email: "other.dummy@email.com",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			_, err = db.NamedExec(
				`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
				&Item{
					ID:      itemID,
					OwnerID: otherUserID,
					Name:    "dummy item name",
					Price:   12000,
				},
			)
			Expect(err).ToNot(HaveOccurred())

			err = persistence.NewItemPersister(db).Delete(itemID, userID)
			Expect(err).To(HaveOccurred())

			items := []Item{}
			err = db.Select(&items, `SELECT * FROM items;`)
			Expect(err).ToNot(HaveOccurred())
			Expect(items).To(HaveLen(1))
		})
	})
	Context("delete called with multiple items in the table that belong to a single user", func() {
		It("should delete the item", func() {
			itemID := uuid.NewString()
			otherItemID := uuid.NewString()
			_, err := db.NamedExec(
				`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
				[]Item{
					{
						ID:      otherItemID,
						OwnerID: userID,
						Name:    "other dummy item name",
						Price:   12000,
					},
					{
						ID:      itemID,
						OwnerID: userID,
						Name:    "dummy item name",
						Price:   44000,
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())

			err = persistence.NewItemPersister(db).Delete(itemID, userID)
			Expect(err).ToNot(HaveOccurred())

			items := []Item{}
			err = db.Select(&items, `SELECT * FROM items;`)
			Expect(err).ToNot(HaveOccurred())
			Expect(items).To(HaveLen(1))
			Expect(items[0].ID).To(Equal(otherItemID))
			Expect(items[0].OwnerID).To(Equal(userID))
			Expect(items[0].Name).To(Equal("other dummy item name"))
			Expect(items[0].Price).To(Equal(uint32(12000)))
		})
	})

	Describe("item count", func() {
		Context("called with items table empty", func() {
			It("return 0", func() {
				count, err := persistence.NewItemPersister(db).ItemCount(uuid.NewString())
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(0))
			})
		})
		Context("called with items table containing multiple items", func() {
			It("should return the item count", func() {
				itemID := uuid.NewString()
				otherItemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
					[]Item{
						{
							ID:      otherItemID,
							OwnerID: userID,
							Name:    "other dummy item name",
							Price:   12000,
						},
						{
							ID:      itemID,
							OwnerID: userID,
							Name:    "dummy item name",
							Price:   44000,
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())

				count, err := persistence.NewItemPersister(db).ItemCount(userID)
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(2))
			})
		})
		Context("called with items table containing multiple items from different users", func() {
			It("should return the item count", func() {
				otherUserID := uuid.NewString()
				itemID := uuid.NewString()
				otherItemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					User{
						ID:    otherUserID,
						Email: "otherDummy@email.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				_, err = db.NamedExec(
					`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
					[]Item{
						{
							ID:      otherItemID,
							OwnerID: userID,
							Name:    "other dummy item name",
							Price:   12000,
						},
						{
							ID:      uuid.NewString(),
							OwnerID: otherUserID,
							Name:    "other users dummy item name",
							Price:   600,
						},
						{
							ID:      itemID,
							OwnerID: userID,
							Name:    "dummy item name",
							Price:   44000,
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())

				count, err := persistence.NewItemPersister(db).ItemCount(userID)
				Expect(err).ToNot(HaveOccurred())
				Expect(count).To(Equal(2))
			})
		})
	})
})
