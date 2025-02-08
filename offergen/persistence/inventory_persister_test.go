package persistence_test

import (
	"offergen/endpoint/models"
	"offergen/persistence"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Inventory struct {
	OwnerID     string `db:"owner_id"`
	Title       string `db:"title"`
	IsPublished bool   `db:"is_published"`
}

var _ = Describe("inventory persister", func() {
	var db = GetDB()

	Describe("create inventory", func() {
		BeforeEach(func() {
			CleanDB(db)
		})

		Context("called with empty db", func() {
			It("creates inventory", func() {
				userID := uuid.NewString()
				_, err := db.NamedExec(
					`INSERT INTO users (id,email) VALUES (:id,:email);`,
					&User{
						ID:    userID,
						Email: "dummy@email.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())

				_, err = persistence.NewInventoryPersister(db).Create(&models.Inventory{
					OwnerID:     userID,
					Title:       "offering",
					IsPublished: false,
				})
				Expect(err).To(Not(HaveOccurred()))

				inventories := []Inventory{}
				err = db.Select(&inventories, `SELECT * FROM inventory;`)
				Expect(err).To(Not(HaveOccurred()))

				Expect(inventories).To(HaveLen(1))
				Expect(inventories[0].OwnerID).To(Equal(userID))
				Expect(inventories[0].Title).To(Equal("offering"))
				Expect(inventories[0].IsPublished).To(BeFalse())
			})
		})
	})
})
