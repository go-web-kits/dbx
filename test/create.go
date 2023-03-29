package test

import (
	. "github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/lab/business_error"
	. "github.com/go-web-kits/testx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {
	var user1, user2 People

	BeforeEach(func() {
		user1 = People{
			Age:  60,
			Name: "Bob",
		}

		user2 = People{
			Age:  6,
			Name: "Alice",
		}

	})

	AfterEach(func() {
		CleanData(&People{})
	})

	Describe("Create", func() {
		When("creating an object which is duplicated with existing record", func() {
			var p *MonkeyPatches
			BeforeEach(func() {
				p = IsExpectedToCall(IsDuplicateRecord).AndReturn(true)
			})
			AfterEach(func() {
				p.Reset()
			})

			It("Should get NotUnique error", func() {
				ret := Create(&user1)
				Expect(ret.Err).To(Equal(CommonErrors[NotUnique]))
			})

			It("Should create succeed by skipping uniq check", func() {
				ret := Create(&user1, Opt{SkipUniqValidate: true})
				Expect(ret.Err).To(BeNil())
			})
		})

		When("creating an object which is duplicated with existing record", func() {
			var p *MonkeyPatches
			BeforeEach(func() {
				p = IsExpectedToCall(IsDuplicateRecord).AndReturn(false)
			})
			AfterEach(func() {
				p.Reset()
			})

			It("Should create succeed", func() {
				ret := Create(&user1)
				Expect(ret.Err).To(BeNil())
			})
		})
	})

	Describe("FirstOrCreate", func() {
		When("passing object which already exists", func() {
			It("Should not create, get existed one", func() {
				ret := FirstOrCreate(&user1, EQ{"name": user1.Name})
				Expect(ret.Err).To(BeNil())
				Expect(ret.Data.(*People).ID).To(Equal(user1.ID))
			})
		})

		When("passing object which does not exist", func() {
			It("Should create new one, since get no exists one", func() {
				ret := FirstOrCreate(&user2, EQ{"name": user2.Name})
				Expect(ret.Err).To(BeNil())
				Expect(ret.Data.(*People).ID).NotTo(Equal(0))
			})
		})
	})
})
