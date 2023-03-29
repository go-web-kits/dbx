package dbx_model_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/go-web-kits/dbx/dbx_model"
)

var _ = AfterSuite(func() {
	DBxDefinitions = map[string]Definition{}
})

var _ = Describe("Definition", func() {
	Describe("DefinitionOf: get Definition by passing parameter", func() {
		When("passing object which has definition", func() {
			It("takes the definition", func() {
				Expect(DefinitionOf(User{})).NotTo(BeZero())
				Expect(DefinitionOf(User{}).Uniqueness).NotTo(BeZero())
				Expect(DefinitionOf(User{}).Serialization.Rmv).To(Equal([]string{"deleted_at"}))
			})
		})

		When("passing object which has not definition", func() {
			It("cannot take definition", func() {
				Expect(DefinitionOf(Person{})).To(BeZero())
				Expect(DefinitionOf(Person{}).Uniqueness).To(BeZero())
			})
		})

		When("passing string which has definition", func() {
			It("takes the definition", func() {
				Expect(DefinitionOf("user")).NotTo(BeZero())
				Expect(DefinitionOf("User")).NotTo(BeZero())
				Expect(DefinitionOf("User").Uniqueness).NotTo(BeZero())
			})
		})

		When("passing object which has not definition", func() {
			It("cannot take definition", func() {
				Expect(DefinitionOf("Person")).To(BeZero())
				Expect(DefinitionOf("person").Uniqueness).To(BeZero())
			})
		})
	})
})
