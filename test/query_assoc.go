package test

import (
	. "github.com/go-web-kits/dbx"
	"github.com/go-web-kits/dbx/dbx_model"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Assoc", func() {
	var (
		conn          *Chain
		bob, ali, tom People
		bobGroup      Group
		bobGroupAct   Activity
		users         []People
		cars          []Car
	)
	var bobCar, aliCar1, aliCar2 Car

	BeforeEach(func() {
		conn = Conn()

		bob = People{Age: 22, Name: "Bob"}
		ali = People{Age: 21, Name: "Alice"}
		tom = People{Age: 22, Name: "Tom"}
		users = []People{}
		factory.Create(S{&bob, &ali, &tom}, Opt{SkipCallback: true})

		bobCar = Car{Name: "bobCar", Number: "num1", PeopleId: bob.ID}
		aliCar1 = Car{Name: "aliCar1", Number: "num2", PeopleId: ali.ID}
		aliCar2 = Car{Name: "aliCar2", Number: "num3", PeopleId: ali.ID}
		cars = []Car{}
		factory.Create(S{&bobCar, &aliCar1, &aliCar2})
	})

	AfterEach(func() {
		CleanData(Models...)
	})

	Describe("Related", func() {
		It("works", func() {
			Expect(conn.Related(&cars, Opt{RelatedWith: &bob, Count: true})).To(HaveFound(1))
			Expect(cars).To(Equal([]Car{bobCar}))

			Expect(Model(&ali).Related(&cars, With{Count: true})).To(HaveFound(2))
			Expect(cars).To(Equal([]Car{aliCar1, aliCar2}))
		})
	})

	Describe("RelatedWith", func() {
		It("works", func() {
			Expect(Model(&cars).RelatedWith(&ali, With{Count: true})).To(HaveFound(2))
			Expect(cars).To(Equal([]Car{aliCar1, aliCar2}))
		})
	})

	{
		BeforeEach(func() {
			bobGroup = Group{Name: "bobGroup"}
			factory.Create(&bobGroup)
			bobGroupAct = Activity{Title: "bobGroupActivity", ItemID: bobGroup.ID, ItemType: "groups"}
			factory.Create(&bobGroupAct)
			factory.Association(&bob, "Groups").Append(&bobGroup)
		})
		AfterEach(func() {
			factory.Association(&bob, "Groups").Clear()
		})
	}

	Describe("Preload", func() {
		It("works with single preload requirement", func() {
			Expect(Model(&ali).Preload("Cars").FindOut()).To(HaveFound())
			Expect(ali.Cars).To(Equal([]Car{aliCar1, aliCar2}))
		})

		It("works with slice of preload requirements", func() {
			Expect(Model(&bob).Preload([]string{"Cars", "Groups"}).FindOut()).To(HaveFound())
			Expect(bob.Cars).To(Equal([]Car{bobCar}))
			Expect(bob.Groups).To(BeTheSameRecordsTo(bobGroup))
			Expect(bob.Groups[0].Activities).To(BeEmpty())
		})

		When("requesting nested preloading", func() {
			It("preload step by step", func() {
				Expect(Model(&bob).Preload([]string{"Cars", "Groups.Activities"}).FindOut()).To(HaveFound())
				Expect(bob.Cars).To(Equal([]Car{bobCar}))
				_bobGroup := bobGroup
				_bobGroup.Activities = []Activity{bobGroupAct}
				Expect(bob.Groups[0]).To(Equal(_bobGroup))
				Expect(bob.Groups[0].Activities).To(Equal([]Activity{bobGroupAct}))
			})
		})

		When("requesting preload with conditions", func() {
			BeforeEach(func() {
				dbx_model.DBxDefinitions = map[string]dbx_model.Definition{
					"Activity": {DefaultScope: dbx_model.Scope{Where: EQ{"title": "xyz"}}},
				}
			})
			AfterEach(func() {
				dbx_model.DBxDefinitions = map[string]dbx_model.Definition{}
			})

			It("uses the conditions and ignores the default preload definition of the preload object", func() {
				Expect(Model(&bob).Preload("Groups.Activities").FindOut()).To(HaveFound())
				_bobGroup := bobGroup
				_bobGroup.Activities = []Activity{bobGroupAct}
				Expect(bob.Groups).To(BeTheSameRecordsTo(bobGroup))
				Expect(bob.Groups[0].Activities).To(BeEmpty())

				Expect(Model(&bob).Preload(map[string][]interface{}{"Groups.Activities": {"title = ?", "bobGroupActivity"}}).
					FindOut()).To(HaveFound())
				Expect(bob.Groups).To(BeTheSameRecordsTo(bobGroup))
				Expect(bob.Groups[0].Activities).To(Equal([]Activity{bobGroupAct}))
			})
		})
	})

	Describe("Joins", func() {
		Context("Normal Joining", func() {
			It("works", func() {
				err := conn.Model(users).Joins("Cars").
					Where(PLAIN{"cars.name = 'bobCar'"}).Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(BeTheSameRecordsTo(bob))
			})
		})

		Context("Polymorphic Joining", func() {
			It("works", func() {
				activity := Activity{Title: "bob", ItemType: "peoples", ItemID: bob.ID}
				factory.Create(&activity)
				factory.Association(&bob, "Activities").Append(activity)

				err := conn.Model(users).Joins("Activities").
					Where(PLAIN{"activities.title = 'bob'"}).Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(BeTheSameRecordsTo(bob))
			})
		})

		Context("Join Table Joining", func() {
			It("works", func() {
				err := conn.Model(users).Joins("Groups").
					Where(PLAIN{"groups.name = 'bobGroup'"}).Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(BeTheSameRecordsTo(bob))
			})
		})

		Context("nested joining", func() {
			It("works", func() {
				err := conn.Model(users).Joins("Groups", "Activities").
					Where(PLAIN{"activities.title = 'bobGroupActivity'"}).Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(BeTheSameRecordsTo(bob))
			})
		})
	})
})
