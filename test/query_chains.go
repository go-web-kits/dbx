package test

import (
	. "github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query Chains", func() {
	var (
		conn, where   *Chain
		bob, ali, tom People
		user          People
		users         []People
	)

	BeforeEach(func() {
		conn = Conn()
		where = conn.Where(LIKE{"name": "o"})

		bob = People{Age: 22, Name: "Bob"}
		ali = People{Age: 22, Name: "Alice"}
		tom = People{Age: 20, Name: "Tom"}
		user = People{}
		users = []People{}
		factory.Create(S{&bob, &ali, &tom}, Opt{SkipCallback: true})
	})

	AfterEach(func() {
		CleanData(Models...)
	})

	Describe("Where", func() {
		It("works", func() {
			Expect(conn.Where(EQ{"id": -1}).Find(&user).Error).To(HaveOccurred())
			Expect(conn.Where(EQ{"id": bob.ID}).Find(&user).Error).NotTo(HaveOccurred())
		})
	})

	Describe("Order", func() {
		It("works", func() {
			Expect(where.Order("age desc").Find(&users).Error).NotTo(HaveOccurred())
			Expect(users).To(Equal([]People{bob, tom}))

			Expect(where.Order("age asc").Find(&users).Error).NotTo(HaveOccurred())
			Expect(users).To(Equal([]People{tom, bob}))
		})

		When("calling it multiple times", func() {
			It("orders by left to right", func() {
				Expect(bob.ID).To(BeNumerically("<", tom.ID))

				err := where.Order("age desc").Order("id desc").Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(Equal([]People{bob, tom}))

				err = where.Order("id desc").Order("age desc").Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(Equal([]People{tom, bob}))
			})

			It("should not cover the same field", func() {
				err := where.Order("age desc").Order("age asc").Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(Equal([]People{bob, tom}))

				err = where.Order("age asc").Order("age desc").Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(Equal([]People{tom, bob}))
			})
		})

		When("passing `reorder`", func() {
			It("covers the same field", func() {
				err := where.Order("age desc").Order("age asc", true).Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(Equal([]People{tom, bob}))

				err = where.Order("age asc").Order("age desc", true).Find(&users).Error
				Expect(err).NotTo(HaveOccurred())
				Expect(users).To(Equal([]People{bob, tom}))
			})
		})
	})

	Describe("Uniq", func() {
		It("works", func() {
			Expect(bob.ID).To(BeNumerically("<", ali.ID))

			err := conn.Where(EQ{"age": 22}).Uniq(Opt{UniqBy: "age", UniqOrder: "id asc"}).Find(&users).Error
			Expect(err).NotTo(HaveOccurred())
			Expect(users).To(Equal([]People{bob}))

			err = conn.Where(EQ{"age": 22}).Uniq(Opt{UniqBy: "age", UniqOrder: "id desc"}).Find(&users).Error
			Expect(err).NotTo(HaveOccurred())
			Expect(users).To(Equal([]People{ali}))
		})
	})

	Describe("Pagy", func() {
		It("works", func() {
			Expect(where.Pagy(Opt{Page: 1, Rows: 1}).Find(&users).Error).NotTo(HaveOccurred())
			Expect(users).To(Equal([]People{bob}))
			Expect(where.Pagy(Opt{Page: 2, Rows: 1}).Find(&users).Error).NotTo(HaveOccurred())
			Expect(users).To(Equal([]People{tom}))
		})

		It("works with default options", func() {
			page, rows := DefaultPage, DefaultRows
			DefaultPage, DefaultRows = 1, 1
			Expect(where.Pagy().Find(&users).Error).NotTo(HaveOccurred())
			Expect(users).To(Equal([]People{bob}))
			DefaultPage, DefaultRows = page, rows
		})
	})

	Describe("UnPagy", func() {
		It("works", func() {
			Expect(where.Pagy(Opt{Rows: 1}).Find(&users).Error).NotTo(HaveOccurred())
			Expect(users).To(Equal([]People{bob}))

			Expect(where.Pagy(Opt{Rows: 1}).Unpagy().Find(&users).Error).NotTo(HaveOccurred())
			Expect(users).To(Equal([]People{bob, tom}))
		})
	})
})
