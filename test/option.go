package test

import (
	. "github.com/go-web-kits/dbx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Option", func() {
	var (
		opts                   []Opt
		opt1, opt2, opt3, opt4 Opt
	)

	BeforeEach(func() {
		opts = []Opt{
			{
				Page: 1,
				Rows: 2,
			},
			{
				Count: true,
				Order: "id",
			},
			{
				ReOrder: true,
				Rows:    4,
			},
		}
		opt1 = Opt{
			Page:    1,
			Rows:    2,
			Count:   true,
			Order:   "id",
			ReOrder: true,
		}
		opt2 = Opt{
			Debug: true,
		}
		opt3 = Opt{
			UnLog: true,
		}
		opt4 = Opt{
			Rows:  5,
			Debug: false,
		}
	})

	Describe("OptsPack: packs multiple opts", func() {
		When("passing opt slice", func() {
			It("Should get one opt with merged config", func() {
				optOut, ok := OptsPack(opts)
				Expect(ok).To(Equal(true))
				Expect(optOut).To(Equal(opt1))
				Expect(optOut.Rows).To(Equal(2))
			})
		})
	})

	Describe("OptsPackGet", func() {
		When("merge opt to opt", func() {
			It("Should get a merged opt", func() {
				o := OptsPackGet(opts)
				Expect(o).To(Equal(opt1))

				oM := o.M(opt2)
				Expect(oM.Debug).To(Equal(true))

				oMerged := oM.Merge([]Opt{opt3, opt4})
				Expect(oMerged.UnLog).To(Equal(true))
				Expect(oMerged.Debug).NotTo(Equal(false)) // set debug to false neglect
				Expect(oMerged.Rows).NotTo(Equal(5))      // neglect while rows already set
			})
		})
	})
})
