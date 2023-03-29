package test

import (
	"math"

	. "github.com/go-web-kits/dbx"
	"github.com/go-web-kits/dbx/dbx_callback"
	. "github.com/go-web-kits/lab/business_error"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update", func() {
	var car1, car2, car3 Car
	var user1, user2, user3 People

	BeforeEach(func() {
		car1 = Car{
			Name:   "car1",
			Number: "num1",
		}
		car2 = Car{
			Name:   "car2",
			Number: "num2",
		}
		car3 = Car{
			Name:   "car3",
			Number: "num3",
		}
		user1 = People{
			Age:  22,
			Name: "Bob",
		}
		user2 = People{
			Age:  21,
			Name: "Alice",
		}
		user3 = People{
			Age:  22,
			Name: "Jhon",
		}

		factory.Create(&user1)
		factory.Create(&user2)
		factory.Create(&user3)

		car1.PeopleId = user1.ID
		car2.PeopleId = user2.ID
		car3.PeopleId = user2.ID

		factory.Create(&car1)
		factory.Create(&car2)
		factory.Create(&car3)
	})

	AfterEach(func() {
		CleanData(&Car{}, &People{})
	})

	Describe("Update by chain", func() {
		It("Update", func() {
			user1.Age = 23
			result := Conn().Update(&user1)

			Expect(result.Err).To(BeNil())
			Expect(result.Data.(*People).Age).To(Equal(int64(23)))
			Expect(result.Data.(*People).CallbackAction).To(Equal(dbx_callback.Update))
		})

		It("Update with SkipCallback", func() {
			user3.Age = 22
			result := Conn().Update(&user3, Opt{SkipCallback: true})

			Expect(result.Err).To(BeNil())
			Expect(result.Data.(*People).Age).To(Equal(int64(22)))
		})

		It("Update when not found", func() {
			user := People{ID: math.MaxInt32, Age: 11, Name: "11"}
			result := Conn().Update(&user, Opt{SkipCallback: true})

			Expect(result).NotTo(HaveFound())
		})

		It("Update when is duplicate", func() {
			p := IsExpectedToCall(IsDuplicateRecord).AndReturn(true)
			defer p.Reset()

			user := user3
			result := Conn().Update(&user, Opt{SkipCallback: true})

			Expect(result.Err).To(Equal(CommonErrors[NotUnique]))
		})

		It("UpdateBy", func() {
			result := Conn().UpdateBy(&user1, map[string]interface{}{"age": 23})

			Expect(result.Err).To(BeNil())
			Expect(result.Data.(*People).Age).To(Equal(int64(23)))
		})

		It("UpdateBy when not found", func() {
			user := People{ID: math.MaxInt32, Age: 11, Name: "11"}
			result := Conn().UpdateBy(&user, &user, Opt{SkipCallback: true})

			Expect(result).NotTo(HaveFound())
		})

		It("UpdateBy slice", func() {
			users := []People{user1, user3}
			result := Conn().UpdateBy(&users, map[string]interface{}{"age": "23"}, Opt{SkipCallback: true})

			Expect(result.Err).To(BeNil())
		})

		It("UpdateBy when is duplicate", func() {
			p := IsExpectedToCall(IsDuplicateRecord).AndReturn(true)
			defer p.Reset()

			user := user3
			result := Conn().UpdateBy(&user, Opt{SkipCallback: true})

			Expect(result.Err).To(Equal(CommonErrors[NotUnique]))
		})

		It("UpdateBy with log", func() {
			unlog := UnLog
			format := DefaultLogFormat
			UnLog = false
			DefaultLogFormat = "json"
			defer func() {
				UnLog = unlog
				DefaultLogFormat = format
			}()

			p := IsExpectedToCall(IsDuplicateRecord).AndReturn(true)
			defer p.Reset()

			user := user3
			result := Conn().UpdateBy(&user, Opt{SkipCallback: true})

			Expect(result.Err).To(Equal(CommonErrors[NotUnique]))
		})
	})

	Describe("Update by crud interface", func() {
		It("Update", func() {
			user1.Age = 24
			result := Update(&user1)

			Expect(result.Err).To(BeNil())
			Expect(result.Data.(*People).Age).To(Equal(int64(24)))
		})

		It("UpdateBy", func() {
			result := UpdateBy(&user1, map[string]interface{}{"age": 22})

			Expect(result.Err).To(BeNil())
			Expect(result.Data.(*People).Age).To(Equal(int64(22)))
		})
	})
})
