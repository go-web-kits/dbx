package test

import (
	. "github.com/go-web-kits/dbx"
	"github.com/go-web-kits/dbx/dbx_model"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validation", func() {
	var user1 People

	BeforeEach(func() {
		user1 = People{
			Age:  22,
			Name: "Bob",
		}

		factory.Create(&user1)
	})

	AfterEach(func() {
		CleanData(&People{})
	})

	It("Validation Uniqueness by map", func() {
		p := IsExpectedToCall(dbx_model.DefinitionOf).AndReturn(
			dbx_model.Definition{
				Uniqueness: map[string][]string{"name": {"age"}},
			})
		defer p.Reset()

		result := IsDuplicateRecord(&People{ID: user1.ID + 1, Age: user1.Age}, H{"name": user1.Name})
		Expect(result).To(BeTrue())
	})

	It("Validation Uniqueness by string", func() {
		p := IsExpectedToCall(dbx_model.DefinitionOf).AndReturn(dbx_model.Definition{Uniqueness: "name"})
		defer p.Reset()

		result := IsDuplicateRecord(&People{ID: user1.ID + 1, Age: user1.Age}, H{"name": user1.Name})
		Expect(result).To(BeTrue())
	})

	It("Validation Uniqueness by slice of string", func() {
		p := IsExpectedToCall(dbx_model.DefinitionOf).AndReturn(dbx_model.Definition{Uniqueness: []string{"name", "age"}})
		defer p.Reset()

		result := IsDuplicateRecord(&People{ID: user1.ID + 1, Age: user1.Age}, H{"name": user1.Name})
		Expect(result).To(BeTrue())
	})
})
