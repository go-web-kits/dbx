package test

import (
	. "github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UpdateExtend", func() {
	var user People

	Describe("Increment with Chain", func() {
		BeforeEach(func() {
			user = People{Name: "Bob", Age: 22, Cash: 10000}
			factory.Create(&user)
		})

		AfterEach(func() {
			CleanData(&People{})
		})

		It("Increment with string", func() {
			result := Conn().Increment(&user, "age")
			u := user
			factory.Find(&u)

			Expect(result.Err).To(BeNil())
			Expect(u.Age).To(Equal(int64(user.Age + 1)))
		})

		It("Increment with slice of string", func() {
			result := Conn().Increment(&user, []string{"age", "cash"})
			u := user
			factory.Find(&u)

			Expect(result.Err).To(BeNil())
			Expect(u.Age).To(Equal(int64(user.Age + 1)))
			Expect(u.Cash).To(Equal(int64(user.Cash + 1)))
		})

		It("Increment with map", func() {
			result := Conn().Increment(&user, map[string]int{"age": 1, "cash": 1000})
			u := user
			factory.Find(&u)

			Expect(result.Err).To(BeNil())
			Expect(u.Age).To(Equal(int64(user.Age + 1)))
			Expect(u.Cash).To(Equal(int64(user.Cash + 1000)))
		})
	})

	Describe("Increment", func() {
		BeforeEach(func() {
			user = People{Name: "Bob", Age: 22, Cash: 10000}
			factory.Create(&user)
		})

		AfterEach(func() {
			CleanData(&People{})
		})

		It("Increment", func() {
			result := Increment(&user, "age")
			u := user
			factory.Find(&u)

			Expect(result.Err).To(BeNil())
			Expect(u.Age).To(Equal(int64(user.Age + 1)))
		})
	})

	Describe("Decrement with Chain", func() {
		BeforeEach(func() {
			user = People{Name: "Bob", Age: 22, Cash: 10000}
			factory.Create(&user)
		})

		AfterEach(func() {
			CleanData(&People{})
		})
		It("Decrement with string", func() {
			result := Conn().Decrement(&user, "age")
			u := user
			factory.Find(&u)

			Expect(result.Err).To(BeNil())
			Expect(u.Age).To(Equal(int64(user.Age - 1)))
		})

		It("Decrement with slice of string", func() {
			result := Conn().Decrement(&user, []string{"age", "cash"})
			u := user
			factory.Find(&u)

			Expect(result.Err).To(BeNil())
			Expect(u.Age).To(Equal(int64(user.Age - 1)))
			Expect(u.Cash).To(Equal(int64(user.Cash - 1)))
		})

		It("Decrement with map", func() {
			result := Conn().Decrement(&user, map[string]int{"age": 1, "cash": 1000})
			u := user
			factory.Find(&u)

			Expect(result.Err).To(BeNil())
			Expect(u.Age).To(Equal(int64(user.Age - 1)))
			Expect(u.Cash).To(Equal(int64(user.Cash - 1000)))
		})
	})

	Describe("Decrement", func() {
		BeforeEach(func() {
			user = People{Name: "Bob", Age: 22, Cash: 10000}
			factory.Create(&user)
		})

		AfterEach(func() {
			CleanData(&People{})
		})

		It("Decrement", func() {
			result := Decrement(&user, "age")
			u := user
			factory.Find(&u)

			Expect(result.Err).To(BeNil())
			Expect(u.Age).To(Equal(int64(user.Age - 1)))
		})
	})
})
