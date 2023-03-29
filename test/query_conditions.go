package test

import (
	. "github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Condition (#Where)", func() {
	var (
		conn          *Chain
		bob, ali, tom People
		user          People
		users         []People
	)

	get := func(out interface{}, condition interface{}) error {
		return conn.Where(condition).Find(out).Error
	}

	BeforeEach(func() {
		conn = Conn()

		bob = People{Age: 22, Name: "Bob"}
		ali = People{Age: 21, Name: "Alice"}
		tom = People{Age: 22, Name: "Tom"}
		user = People{}
		users = []People{}

		factory.Create(S{&bob, &ali, &tom}, Opt{SkipCallback: true})
	})

	AfterEach(func() {
		CleanData(Models...)
	})

	Describe("EQ", func() {
		It("works", func() {
			Expect(get(&user, EQ{"id": bob.ID})).To(Succeed())
			Expect(user).To(Equal(bob))
		})

		It("fails", func() {
			Expect(get(&user, EQ{"name": "wrong"})).NotTo(Succeed())
			Expect(get(&user, EQ{"wrong_col": ""})).NotTo(Succeed())
		})
	})

	Describe("PLAIN", func() {
		It("works", func() {
			Expect(get(&user, PLAIN{"id = ?", bob.ID})).To(Succeed())
			Expect(user).To(Equal(bob))
		})

		It("works with slice values", func() {
			Expect(get(&user, PLAIN{"id = ?", []interface{}{bob.ID}})).To(Succeed())
			Expect(user).To(Equal(bob))
		})
	})

	Describe("IN", func() {
		It("works", func() {
			_ = get(&users, IN{"id": S{bob.ID, ali.ID, -1}})
			Expect(users).To(Equal([]People{bob, ali}))
		})
	})

	Describe("LIKE", func() {
		It("works", func() {
			_ = get(&users, LIKE{"name": "o"})
			Expect(users).To(Equal([]People{bob, tom}))
		})
	})

	Describe("NOT", func() {
		It("works", func() {
			_ = get(&users, NOT{"id = ?", bob.ID})
			Expect(users).To(Equal([]People{ali, tom}))
		})
	})

	Describe("Combine", func() {
		It("works with diff types of conditioner", func() {
			Expect(get(&user, Combine{EQ{"id": bob.ID}, LIKE{"name": "a"}})).NotTo(Succeed())
			Expect(get(&user, Combine{EQ{"id": bob.ID}, LIKE{"name": "o"}})).To(Succeed())
			Expect(user).To(Equal(bob))
		})

		It("works with same type of conditioner", func() {
			Expect(get(&user, Combine{EQ{"id": bob.ID}, EQ{"name": ali.Name}})).NotTo(Succeed())
			Expect(get(&user, Combine{EQ{"id": bob.ID}, EQ{"name": bob.Name}, EQ{"age": bob.Age}})).To(Succeed())
			Expect(user).To(Equal(bob))
		})

		Describe("OR", func() {
			It("works", func() {
				_ = get(&users, Combine{EQ{"id": bob.ID}, OR{"id = ?", ali.ID}})
				Expect(users).To(Equal([]People{bob, ali}))
			})
		})
	})
})
