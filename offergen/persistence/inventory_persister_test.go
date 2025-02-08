package persistence_test

import (
	"offergen/endpoint/models"
	"offergen/persistence"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type Inventory struct {
	OwnerID     string `db:"owner_id"`
	Title       string `db:"title"`
	IsPublished bool   `db:"is_published"`
}

var _ = Describe("inventory persister", func() {
	var db *sqlx.DB
	BeforeEach(func() {
		db = GetDB()
		CleanDB(db)
	})
	Context("with a user existing", func() {
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
		Describe("create inventory", func() {
			Context("called with empty db", func() {
				It("should create the inventory", func() {
					createdInv, err := persistence.NewInventoryPersister(db).Create(&models.Inventory{
						OwnerID:     userID,
						Title:       "offering",
						IsPublished: false,
					})
					Expect(err).To(Not(HaveOccurred()))

					inventories := []Inventory{}
					err = db.Select(&inventories, `SELECT * FROM inventories;`)
					Expect(err).To(Not(HaveOccurred()))

					Expect(inventories).To(HaveLen(1))
					Expect(inventories[0].OwnerID).To(Equal(userID))
					Expect(inventories[0].Title).To(Equal("offering"))
					Expect(inventories[0].IsPublished).To(BeFalse())
					Expect(createdInv.OwnerID).To(Equal(userID))
					Expect(createdInv.Title).To(Equal("offering"))
					Expect(createdInv.IsPublished).To(BeFalse())
				})
			})
		})
		Describe("update inventory", func() {
			Context("called an inventory existing", func() {
				BeforeEach(func() {
					_, err := db.NamedExec(
						`
                            INSERT INTO inventories (owner_id,title,is_published)
                            VALUES (:owner_id,:title,:is_published)
                        `,
						&Inventory{
							OwnerID:     userID,
							Title:       "dummy_inventory",
							IsPublished: true,
						},
					)
					Expect(err).ToNot(HaveOccurred())
				})
				It("should update the inventory", func() {
					expectedInv := &models.Inventory{
						OwnerID:     userID,
						Title:       "dummy_updated_title",
						IsPublished: false,
					}
					updatedInv, err := persistence.NewInventoryPersister(db).Update(
						userID,
						&models.UpdateInventoryInput{
							Title:       "dummy_updated_title",
							IsPublished: models.Point(false),
						})
					Expect(err).To(Not(HaveOccurred()))

					inventories := []Inventory{}
					err = db.Select(&inventories, `SELECT * FROM inventories;`)
					Expect(err).To(Not(HaveOccurred()))

					Expect(inventories).To(ConsistOf(Inventory{
						OwnerID:     expectedInv.OwnerID,
						Title:       expectedInv.Title,
						IsPublished: expectedInv.IsPublished,
					}))
					Expect(updatedInv).To(Equal(expectedInv))
				})
			})
		})
	})

})
