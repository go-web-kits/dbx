package test

import (
	. "github.com/go-web-kits/dbx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Destroy", func() {
	var opt Opt
	var user People

	BeforeEach(func() {
		opt = Opt{}

		user = People{
			Age:  60,
			Name: "BobDestroy",
		}
	})

	Describe("Destroy: delete passing object", func() {
		When("destroying object", func() {
			It("destroy a created user", func() {
				Create(&user, opt)
				ret := Destroy(&user, opt)
				Expect(ret.Err).To(BeNil())
				Expect(ret.Data.(*People).ID).To(Equal(user.ID))
				// todo check user exists
			})
		})
	})

	Describe("Destroy: delete passing object in a Chain", func() {
		When("destroying object", func() {
			It("Destroy a created user", func() {
				Create(&user, opt)
				ret := Conn(opt).Model(&People{}).Where(EQ{"id": user.ID}).Destroy()
				Expect(ret.Err).To(BeNil())
				// todo check user exists
			})
		})
	})
})
