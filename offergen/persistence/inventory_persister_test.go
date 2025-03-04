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

type Inventory struct {
	OwnerID     string `db:"owner_id"`
	Title       string `db:"title"`
	IsPublished bool   `db:"is_published"`
}
type Category struct {
	OwnerID string `db:"owner_id"`
	Name    string `db:"name"`
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
					IsPublished: false,
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
							IsPublished: true,
						})
					Expect(err).To(Not(HaveOccurred()))

					inventories := []Inventory{}
					err = db.Select(&inventories, `SELECT * FROM inventories;`)

					Expect(err).To(Not(HaveOccurred()))
					Expect(inventories).To(ConsistOf(Inventory{
						OwnerID:     userID,
						Title:       "dummy_updated_title",
						IsPublished: true,
					}))
				})
				It("should return the inventory", func() {
					expectedInv := &models.Inventory{
						OwnerID:     userID,
						Title:       "dummy_updated_title",
						IsPublished: true,
					}
					updatedInv, err := persistence.NewInventoryPersister(db).Update(
						userID,
						&models.UpdateInventoryInput{
							Title:       "dummy_updated_title",
							IsPublished: true,
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
		Describe("add category", func() {
			It("should insert a new category", func() {
				err := persistence.NewInventoryPersister(db).CreateCategory(userID, "dummy_category")
				Expect(err).To(Not(HaveOccurred()))

				category := new(Category)
				err = db.Get(
					category,
					`SELECT * FROM categories`,
				)
				Expect(err).ToNot(HaveOccurred())
			})
		})
		Describe("batch get category", func() {
			Context("with empty categories table", func() {
				It("should return the category with proper count", func() {
					categories, err := persistence.NewInventoryPersister(db).
						BatchGetCategory(userID)

					Expect(err).To(Not(HaveOccurred()))
					Expect(categories).To(BeEmpty())
				})
			})
			Context("with a category existing", func() {
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
				Context("called with matching owner ID", func() {
					It("should return the category with proper count", func() {
						categories, err := persistence.NewInventoryPersister(db).
							BatchGetCategory(userID)

						Expect(err).To(Not(HaveOccurred()))
						Expect(categories).To(ConsistOf(category))
					})
				})
				Context("called with non matching owner ID", func() {
					It("should return nothing", func() {
						categories, err := persistence.NewInventoryPersister(db).
							BatchGetCategory(uuid.NewString())

						Expect(err).To(Not(HaveOccurred()))
						Expect(categories).To(BeEmpty())
					})
				})
				Context("with another category existing", func() {
					var otherCategory string
					BeforeEach(func() {
						otherCategory = "other_dummy_category"
						_, err := db.NamedExec(
							`INSERT INTO categories (owner_id,name) VALUES (:owner_id,:name)`,
							&Category{
								OwnerID: userID,
								Name:    otherCategory,
							},
						)
						Expect(err).ToNot(HaveOccurred())
					})
					It("should return the categories", func() {
						categories, err := persistence.NewInventoryPersister(db).
							BatchGetCategory(userID)

						Expect(err).To(Not(HaveOccurred()))
						Expect(categories).To(ConsistOf(category, otherCategory))
					})
				})
			})
		})
		Describe("batch get counted category", func() {
			Context("with a category existing", func() {
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
				It("should return the category with proper count", func() {
					categories, err := persistence.NewInventoryPersister(db).
						BatchGetCountedCategory(userID)

					Expect(err).To(Not(HaveOccurred()))
					Expect(categories).To(ConsistOf(models.CountedCategory{
						Name:  category,
						Count: 0,
					}))
				})
				Context("with a single item belonging to that category", func() {
					var itemID, itemName string
					BeforeEach(func() {
						itemID = uuid.NewString()
						itemName = "dummy item name"
						_, err := db.NamedExec(
							`INSERT INTO items (id,owner_id,name,price,category) VALUES (:id,:owner_id,:name,:price,:category);`,
							&Item{
								ID:       itemID,
								OwnerID:  userID,
								Name:     itemName,
								Price:    12000,
								Category: sql.NullString{String: category, Valid: true},
							},
						)
						Expect(err).ToNot(HaveOccurred())
					})
					It("should return the category with proper count", func() {
						categories, err := persistence.NewInventoryPersister(db).
							BatchGetCountedCategory(userID)

						Expect(err).To(Not(HaveOccurred()))
						Expect(categories).To(ConsistOf(models.CountedCategory{
							Name:  category,
							Count: 1,
						}))
					})
					Context("with multiple items belonging to that category", func() {
						var otherItemID, otherItemName string
						BeforeEach(func() {
							otherItemID = uuid.NewString()
							otherItemName = "other dummy item name"
							_, err := db.NamedExec(
								`INSERT INTO items (id,owner_id,name,price,category) VALUES (:id,:owner_id,:name,:price,:category);`,
								&Item{
									ID:       otherItemID,
									OwnerID:  userID,
									Name:     otherItemName,
									Price:    12000,
									Category: sql.NullString{String: category, Valid: true},
								},
							)
							Expect(err).ToNot(HaveOccurred())
						})
						It("should return the category with proper count", func() {
							categories, err := persistence.NewInventoryPersister(db).
								BatchGetCountedCategory(userID)

							Expect(err).To(Not(HaveOccurred()))
							Expect(categories).To(ConsistOf(
								models.CountedCategory{
									Name:  category,
									Count: 2,
								},
							))
						})
					})
					Context("with another category existing", func() {
						var otherCategory string
						BeforeEach(func() {
							otherCategory = "other_dummy_category"
							_, err := db.NamedExec(
								`INSERT INTO categories (owner_id,name) VALUES (:owner_id,:name)`,
								&Category{
									OwnerID: userID,
									Name:    otherCategory,
								},
							)
							Expect(err).ToNot(HaveOccurred())
						})
						It("should return the category with proper count", func() {
							categories, err := persistence.NewInventoryPersister(db).
								BatchGetCountedCategory(userID)

							Expect(err).To(Not(HaveOccurred()))
							Expect(categories).To(ConsistOf(
								models.CountedCategory{
									Name:  category,
									Count: 1,
								},
								models.CountedCategory{
									Name:  otherCategory,
									Count: 0,
								},
							))
						})
						Context("with a single item belonging to that other category", func() {
							var itemIDOtherCat, itemNameOtherCat string
							BeforeEach(func() {
								itemIDOtherCat = uuid.NewString()
								itemNameOtherCat = "dummy item name other category"
								_, err := db.NamedExec(
									`INSERT INTO items (id,owner_id,name,price,category) VALUES (:id,:owner_id,:name,:price,:category);`,
									&Item{
										ID:       itemIDOtherCat,
										OwnerID:  userID,
										Name:     itemNameOtherCat,
										Price:    12000,
										Category: sql.NullString{String: otherCategory, Valid: true},
									},
								)
								Expect(err).ToNot(HaveOccurred())
							})
							It("should return the category with proper count", func() {
								categories, err := persistence.NewInventoryPersister(db).
									BatchGetCountedCategory(userID)

								Expect(err).To(Not(HaveOccurred()))
								Expect(categories).To(ConsistOf(
									models.CountedCategory{
										Name:  category,
										Count: 1,
									},
									models.CountedCategory{
										Name:  otherCategory,
										Count: 1,
									},
								))
							})
							Context("with multiple items belonging to that other category", func() {
								var otherItemIDOtherCat, otherItemNameOtherCat string
								BeforeEach(func() {
									otherItemIDOtherCat = uuid.NewString()
									otherItemNameOtherCat = "other dummy item name other category"
									_, err := db.NamedExec(
										`INSERT INTO items (id,owner_id,name,price,category) VALUES (:id,:owner_id,:name,:price,:category);`,
										&Item{
											ID:       otherItemIDOtherCat,
											OwnerID:  userID,
											Name:     otherItemNameOtherCat,
											Price:    12000,
											Category: sql.NullString{String: otherCategory, Valid: true},
										},
									)
									Expect(err).ToNot(HaveOccurred())
								})
								It("should return the category with proper count", func() {
									categories, err := persistence.NewInventoryPersister(db).
										BatchGetCountedCategory(userID)

									Expect(err).To(Not(HaveOccurred()))
									Expect(categories).To(ConsistOf(
										models.CountedCategory{
											Name:  category,
											Count: 1,
										},
										models.CountedCategory{
											Name:  otherCategory,
											Count: 2,
										},
									))
								})
							})
						})
					})
				})
			})
		})
	})
})
