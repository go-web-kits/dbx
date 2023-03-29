package test

import (
	. "github.com/go-web-kits/dbx"
	"github.com/go-web-kits/dbx/dbx_model"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Interfaces And Methods", func() {
	var users []People
	var cars []Car
	var user, bob, ali, tom People
	var bobCar, aliCar1, aliCar2 Car

	BeforeEach(func() {
		bobCar = Car{Name: "bobCar", Number: "num1"}
		aliCar1 = Car{Name: "aliCar1", Number: "num2"}
		aliCar2 = Car{Name: "aliCar2", Number: "num3"}
		bob = People{Age: 22, Name: "Bob"}
		ali = People{Age: 22, Name: "Alice"}
		tom = People{Age: 21, Name: "Tom"}
		user = People{}
		users = []People{}
		cars = []Car{}

		factory.Create(S{&bob, &ali, &tom}, Opt{SkipCallback: true})
		bobCar.PeopleId = bob.ID
		aliCar1.PeopleId = ali.ID
		aliCar2.PeopleId = ali.ID
		factory.Create(S{&bobCar, &aliCar1, &aliCar2})

		dbx_model.DBxDefinitions = map[string]dbx_model.Definition{
			"People": {DefaultScope: dbx_model.Scope{Order: "age DESC"}},
		}
	})

	AfterEach(func() {
		dbx_model.DBxDefinitions = map[string]dbx_model.Definition{}
		CleanData(Models...)
	})

	Describe("Where", func() {
		It("processes default scopes, pagy, count, tx and other options", func() {
			result := Where(&users, LIKE{"name": "o"}, With{Preload: "Cars", Page: 1, Rows: 1, Count: true})
			Expect(result.Err).NotTo(HaveOccurred())
			Expect(users).To(BeTheSameRecordsTo(bob))
			Expect(result.Data).To(Equal(&users))
			Expect(result.Total).To(BeEquivalentTo(2))
			Expect(users[0].Cars).To(Equal([]Car{bobCar}))

			Expect(Where(&users, LIKE{"name": "o"}, With{Page: 2, Rows: 1})).To(HaveFound())
			Expect(users).To(BeTheSameRecordsTo(tom))
		})

		It("processes Related", func() {
			result := Where(&cars, nil, Be{RelatedWith: &ali}, With{Count: true})
			Expect(result).To(HaveFound(2))
			Expect(cars).To(Equal([]Car{aliCar1, aliCar2}))
			Expect(result.Data).To(Equal(&cars))
		})
	})

	Describe("First", func() {
		It("processes default scopes, tx", func() {
			Expect(First(&[]People{}, 2).Data).To(BeTheSameRecordsTo(bob, ali))
			Expect(First(&People{}).Data).To(BeTheSameRecordTo(bob))
		})
	})

	Describe("Find", func() {
		It("processes default scopes, tx", func() {
			result := Find(&user, PLAIN{"cars.name = ?", "bobCar"}, With{Join: []string{"Cars"}})
			Expect(result.Err).NotTo(HaveOccurred())
			Expect(user).To(Equal(bob))
		})

		It("processes Related", func() {
			Expect(Find(&user, nil, Be{RelatedWith: bobCar})).To(HaveFound())
			Expect(user).To(Equal(bob))
		})
	})

	Describe("FindById", func() {
		It("works", func() {
			Expect(FindById(&user, bob.ID)).To(HaveFound())
			Expect(user).To(Equal(bob))
		})
	})

	Describe("Related", func() {
		It("processes default scopes (of association), tx", func() {
			Expect(Related(&cars, &ali)).To(HaveFound())
			Expect(cars).To(Equal([]Car{aliCar1, aliCar2}))
		})
	})

	Describe("FindOrInitBy", func() {
		It("processes default scopes, tx and returns the record existed", func() {
			Expect(FindOrInitBy(&user, EQ{"id": bob.ID})).To(HaveFound())
			Expect(user).To(Equal(bob))
		})

		It("returns the initialized record", func() {
			Expect(FindOrInitBy(&user, EQ{"name": "Abc"}).Err).NotTo(HaveOccurred())
			Expect(user.Name).To(Equal("Abc"))
		})
	})

	Describe("IsExists", func() {
		When("record not exists", func() {
			It("returns false", func() {
				Expect(IsExists(&People{}, EQ{"age": 11})).To(BeFalse())
				Expect(IsExists("People", EQ{"age": 11})).To(BeFalse())
			})
		})

		When("record exists", func() {
			It("returns true", func() {
				Expect(IsExists(&People{}, EQ{"id": bob.ID})).To(BeTrue())
				Expect(IsExists("People", EQ{"name": bob.Name})).To(BeTrue())
			})
		})
	})

	// Query methods

	Describe("Pluck", func() {
		It("works", func() {
			ages := []int{}
			Expect(Model(&[]People{}).Order("age DESC").Pluck("age", &ages)).To(Succeed())
			Expect(ages).To(Equal([]int{22, 22, 21}))
		})
	})
})
