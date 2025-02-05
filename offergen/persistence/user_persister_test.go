package persistence_test

import (
	"offergen/common_deps"
	"offergen/persistence"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type User struct {
	ID    string `db:"id"`
	Email string `db:"email"`
}

var _ = Describe("user persister", func() {
	var db = GetDB()

	BeforeEach(func() {
		CleanDB(db)
	})

	Describe("save user", func() {
		Context("save called with empty db", func() {
			var err error

			It("should insert a user", func() {
				userID := uuid.NewString()
				email := "dummy@email.com"

				err = persistence.NewUserPersister(db).Save(userID, email)
				Expect(err).ToNot(HaveOccurred())

				users := []User{}
				err = db.Select(&users, `SELECT * FROM users;`)
				Expect(err).ToNot(HaveOccurred())

				Expect(users).To(HaveLen(1))
				Expect(users[0].ID).To(Equal(userID))
				Expect(users[0].Email).To(Equal(email))
			})
		})
	})

	Describe("get email", func() {
		Context("get email called with empty db", func() {
			It("should return err not found", func() {
				userID := uuid.NewString()

				_, err := persistence.NewUserPersister(db).GetEmail(userID)
				Expect(err).To(MatchError(common_deps.ErrUserNotFound))
			})
		})

		Context("get email called with a single user in db", func() {
			It("should return err not found", func() {
				userID := uuid.NewString()
				expectedEmail := "dummy@email.com"

				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email)`,
					&User{
						ID:    userID,
						Email: expectedEmail,
					},
				)
				Expect(err).ToNot(HaveOccurred())

				email, err := persistence.NewUserPersister(db).GetEmail(userID)
				Expect(err).ToNot(HaveOccurred())
				Expect(email).To(Equal(expectedEmail))
			})
		})

		Context("get email called multiple users in db", func() {
			It("should return err not found", func() {
				userID := uuid.NewString()
				expectedEmail := "dummy@email.com"

				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email)`,
					[]User{
						{
							ID:    uuid.NewString(),
							Email: "other_dummy@email.com",
						},
						{
							ID:    userID,
							Email: expectedEmail,
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())

				email, err := persistence.NewUserPersister(db).GetEmail(userID)
				Expect(err).ToNot(HaveOccurred())
				Expect(email).To(Equal(expectedEmail))
			})
		})
	})

	Describe("delete user", func() {
		Context("called with empty db", func() {
			It("should return user not found error", func() {
				authID := uuid.NewString()

				err := persistence.NewUserPersister(db).Delete(authID)
				Expect(err).To(MatchError(common_deps.ErrUserNotFound))
			})
		})

		Context("called with single user in db", func() {
			var userID string
			BeforeEach(func() {
				userID = uuid.NewString()

				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email)`,
					&User{
						ID:    userID,
						Email: "dummy@email.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())
			})
			It("should delete the user", func() {
				err := persistence.NewUserPersister(db).Delete(userID)
				Expect(err).ToNot(HaveOccurred())

				var users []User
				err = db.Select(
					&users,
					`SELECT * FROM users;`,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(users).To(HaveLen(0))
			})

			Context("with user having a single inventory", func() {
				BeforeEach(func() {
					_, err := db.NamedExec(
						`INSERT INTO inventory (owner_id,title,is_published) VALUES (:owner_id,:title,:is_published)`,
						&Inventory{
							OwnerID:     userID,
							Title:       "dummy_title",
							IsPublished: true,
						},
					)
					Expect(err).ToNot(HaveOccurred())
				})
				It("should delete the user and the inventory", func() {
					err := persistence.NewUserPersister(db).Delete(userID)
					Expect(err).ToNot(HaveOccurred())

					var users []User
					err = db.Select(
						&users,
						`SELECT * FROM users;`,
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(users).To(HaveLen(0))
					var inventories []Inventory
					err = db.Select(
						&inventories,
						`SELECT * FROM inventory;`,
					)
					Expect(err).ToNot(HaveOccurred())
					Expect(users).To(HaveLen(0))
				})

				Context("with the inventory having a single item", func() {
					BeforeEach(func() {
						_, err := db.NamedExec(
							`INSERT INTO items (id,owner_id,price,name) VALUES (:id,:owner_id,:price,:name)`,
							&Item{
								ID:      uuid.NewString(),
								OwnerID: userID,
								Price:   1000,
								Name:    "dummy_name",
							},
						)
						Expect(err).ToNot(HaveOccurred())
					})
					It("should delete the user and the inventory", func() {
						err := persistence.NewUserPersister(db).Delete(userID)
						Expect(err).ToNot(HaveOccurred())

						var users []User
						err = db.Select(
							&users,
							`SELECT * FROM users;`,
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(users).To(HaveLen(0))
						var inventories []Inventory
						err = db.Select(
							&inventories,
							`SELECT * FROM inventory;`,
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(users).To(HaveLen(0))
						var items []Item
						err = db.Select(
							&items,
							`SELECT * FROM items;`,
						)
						Expect(err).ToNot(HaveOccurred())
						Expect(users).To(HaveLen(0))
					})
				})
			})
		})
		Context("called with multiple users in db", func() {
			It("delete the correct user", func() {
				userID := uuid.NewString()

				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email)`,
					[]User{
						{
							ID:    uuid.NewString(),
							Email: "other_dummy@email.com",
						},
						{
							ID:    userID,
							Email: "dummy@email.com",
						},
					},
				)
				Expect(err).ToNot(HaveOccurred())

				err = persistence.NewUserPersister(db).Delete(userID)
				Expect(err).ToNot(HaveOccurred())

				var users []User
				err = db.Select(
					&users,
					`SELECT * FROM users;`,
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(users).To(HaveLen(1))
				Expect(users[0].Email).To(Equal("other_dummy@email.com"))
			})
		})
	})
})
