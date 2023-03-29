package test

import (
	. "github.com/go-web-kits/dbx"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Result", func() {
	var user1, user2, user3 People
	var result Result

	BeforeEach(func() {
		user1 = People{ID: 1, Name: "bob"}
		user2 = People{ID: 2}
		user3 = People{ID: 3}

		result = Result{
			Data: []People{user1, {ID: 1, Name: "tom"}, user2, user3},
		}
	})

	Context("Uniq", func() {
		It("does de-duplication by ID", func() {
			Expect(result.Uniq().Data).To(Equal(S{user1, user2, user3}))
		})
	})

	Context("Ids", func() {
		It("returns the ids of unique data", func() {
			ids, err := result.Ids()

			Expect(err).NotTo(HaveOccurred())
			Expect(ids).To(Equal([]uint{1, 2, 3}))
		})
	})
})
