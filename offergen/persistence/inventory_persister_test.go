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
		Context("when an inventory exists", func() {
			var inventory *Inventory
			BeforeEach(func() {
				inventory = &Inventory{
					OwnerID:     userID,
					Title:       "dummy_inventory",
					IsPublished: true,
				}
				_, err := db.NamedExec(
					`
                            INSERT INTO inventories (owner_id,title,is_published)
                            VALUES (:owner_id,:title,:is_published)
                        `,
					inventory,
				)
				Expect(err).ToNot(HaveOccurred())
			})
			Describe("update inventory", func() {
				It("should update the inventory", func() {
					_, err := persistence.NewInventoryPersister(db).Update(
						userID,
						&models.UpdateInventoryInput{
							Title:       "dummy_updated_title",
							IsPublished: models.Point(false),
						})
					Expect(err).To(Not(HaveOccurred()))

					inventories := []Inventory{}
					err = db.Select(&inventories, `SELECT * FROM inventories;`)

					Expect(err).To(Not(HaveOccurred()))
					Expect(inventories).To(ConsistOf(inventory))
				})
				It("should return the inventory", func() {
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
					Expect(updatedInv).To(Equal(expectedInv))
				})
			})
			Describe("get inventory", func() {
				It("should return the inventory", func() {
					expectedInv := &models.Inventory{
						OwnerID:     userID,
						Title:       inventory.Title,
						IsPublished: inventory.IsPublished,
					}
					returnedInv, err := persistence.NewInventoryPersister(db).Get(userID)
					Expect(err).To(Not(HaveOccurred()))

					Expect(returnedInv).To(Equal(expectedInv))
				})
			})
		})
	})
})
