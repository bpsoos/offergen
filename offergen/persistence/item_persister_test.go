package persistence_test

import (
	"offergen/common_deps"
	"offergen/endpoint/models"
	"offergen/persistence"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Item struct {
	ID      string `db:"id"`
	OwnerID string `db:"owner_id"`
	Price   uint32 `db:"price"`
	Name    string `db:"name"`
}

var _ = Describe("items", func() {
	var db *sqlx.DB

	BeforeEach(func() {
		db = GetDB()
		CleanDB(db)
	})
	Describe("create", func() {
		Context("called with empty db", func() {
			It("should create item", func() {
				userID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					&User{
						ID:    userID,
						Email: "dummy@email.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())

				item := &models.Item{
					ID:    uuid.New(),
					Name:  "dummyItem",
					Price: 12000,
				}

				err = persistence.NewItemPersister(db).Create(item, userID)
				Expect(err).ToNot(HaveOccurred())

				items := []Item{}
				err = db.Select(&items, `SELECT * FROM items;`)
				Expect(err).ToNot(HaveOccurred())

				Expect(items).To(HaveLen(1))
				Expect(items[0].ID).To(Equal(item.ID.String()))
				Expect(items[0].Name).To(Equal(item.Name))
				Expect(items[0].Price).To(Equal(item.Price))
			})
		})
	})
	Describe("batch get", func() {
		Context("called with empty items table", func() {
			It("should return an empty list", func() {
				userID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					&User{
						ID:    userID,
						Email: "dummy@email.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())

				items, err := persistence.NewItemPersister(db).BatchGet(0, 10, userID)

				Expect(err).ToNot(HaveOccurred())
				Expect(items).To(HaveLen(0))
			})
		})

		Context("called with items table containing a single item", func() {
			It("should return the item", func() {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					&User{
						ID:    userID,
						Email: "dummy@email.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				_, err = db.NamedExec(
					`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
					&Item{
						ID:      itemID,
						OwnerID: userID,
						Name:    "dummy item name",
						Price:   12000,
					},
				)
				Expect(err).ToNot(HaveOccurred())

				items, err := persistence.NewItemPersister(db).BatchGet(0, 10, userID)

				Expect(err).ToNot(HaveOccurred())
				Expect(items).To(HaveLen(1))
				Expect(items[0].ID.String()).To(Equal(itemID))
				Expect(items[0].Name).To(Equal("dummy item name"))
				Expect(items[0].Price).To(Equal(uint32(12000)))
			})
		})

		Context("called with items table containing items for different users", func() {
			It("should return the correct item", func() {
				userID := uuid.NewString()
				otherUserID := uuid.NewString()
				itemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					[]User{
						{
							ID:    otherUserID,
							Email: "otherDummy@email.com",
						}, {
							ID:    userID,
							Email: "dummy@email.com",
						},
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
						}, {
							ID:      uuid.NewString(),
							OwnerID: otherUserID,
							Name:    "other dummy item name",
							Price:   24000,
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())

				items, err := persistence.NewItemPersister(db).BatchGet(0, 10, userID)

				Expect(err).ToNot(HaveOccurred())
				Expect(items).To(HaveLen(1))
				Expect(items[0].ID.String()).To(Equal(itemID))
				Expect(items[0].Name).To(Equal("dummy item name"))
				Expect(items[0].Price).To(Equal(uint32(12000)))
			})
		})
		Context("called with items table containing items that fall out of the specified range", func() {
			It("should return only the items in range", func() {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				otherItemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					&User{
						ID:    userID,
						Email: "dummy@email.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				_, err = db.NamedExec(
					`INSERT INTO items (id,owner_id,name,price) VALUES (:id,:owner_id,:name,:price);`,
					[]Item{
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
						{
							ID:      itemID,
							OwnerID: userID,
							Name:    "dummy item name 3",
							Price:   13000,
						},
						{
							ID:      otherItemID,
							OwnerID: userID,
							Name:    "dummy item name 4",
							Price:   14000,
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())

				fetchedItems, err := persistence.NewItemPersister(db).BatchGet(2, 2, userID)

				Expect(err).ToNot(HaveOccurred())
				Expect(fetchedItems).To(HaveLen(2))
				Expect(fetchedItems[0].ID.String()).To(Equal(itemID))
				Expect(fetchedItems[0].Name).To(Equal("dummy item name 3"))
				Expect(fetchedItems[0].Price).To(Equal(uint32(13000)))
				Expect(fetchedItems[1].ID.String()).To(Equal(otherItemID))
				Expect(fetchedItems[1].Name).To(Equal("dummy item name 4"))
				Expect(fetchedItems[1].Price).To(Equal(uint32(14000)))
			})
		})
	})

	Describe("delete", func() {
		Context("delete called with single item in table", func() {
			It("should delete the item", func() {
				userID := uuid.NewString()
				itemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					&User{
						ID:    userID,
						Email: "dummy@email.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())
				_, err = db.NamedExec(
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
				userID := uuid.NewString()
				otherUserID := uuid.NewString()
				itemID := uuid.NewString()
				otherItemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					[]User{
						{
							ID:    userID,
							Email: "dummy@email.com",
						}, {
							ID:    otherUserID,
							Email: "otherDummy@email.com",
						},
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
			userID := uuid.NewString()
			otherUserID := uuid.NewString()
			itemID := uuid.NewString()
			_, err := db.NamedExec(
				`INSERT INTO users (id,email) VALUES (:id,:email);`,
				&User{
					ID:    otherUserID,
					Email: "dummy@email.com",
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
			Expect(err).To(MatchError(common_deps.ErrItemNotFound))

			items := []Item{}
			err = db.Select(&items, `SELECT * FROM items;`)
			Expect(err).ToNot(HaveOccurred())
			Expect(items).To(HaveLen(1))
		})
	})
	Context("delete called with multiple items in the table that belong to a single user", func() {
		It("should delete the item", func() {
			userID := uuid.NewString()
			itemID := uuid.NewString()
			otherItemID := uuid.NewString()
			_, err := db.NamedExec(
				`INSERT INTO users (id,email) VALUES (:id,:email);`,
				&User{
					ID:    userID,
					Email: "dummy@email.com",
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
				userID := uuid.NewString()
				itemID := uuid.NewString()
				otherItemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					&User{
						ID:    userID,
						Email: "dummy@email.com",
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
				userID := uuid.NewString()
				otherUserID := uuid.NewString()
				itemID := uuid.NewString()
				otherItemID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					[]User{
						{
							ID:    otherUserID,
							Email: "otherDummy@email.com",
						},
						{
							ID:    userID,
							Email: "dummy@email.com",
						},
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
