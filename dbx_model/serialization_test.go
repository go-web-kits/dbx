package dbx_model_test

import (
	"errors"
	"reflect"
	"time"

	. "github.com/go-web-kits/dbx/dbx_model"
	. "github.com/go-web-kits/testx"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Serialization", func() {
	var p *MonkeyPatches

	AfterEach(func() {
		p.Check()
	})

	Describe("Serialization: get serialization of object based on definition", func() {
		When("passing object which has definition", func() {
			It("Can get serialization of User", func() {
				t := time.Now()
				s, err := Serialize(User{
					ID:        1,
					UserName:  "Bob",
					UpdatedAt: t,
					DeletedAt: &t,
				})
				Expect(len(s)).To(Equal(4))
				Expect(s["id"]).To(Equal(float64(1)))
				Expect(s["email"]).To(Equal("Bob@example.com"))
				Expect(len(err)).To(Equal(0))
			})
		})

		When("passing object with additional serialization config", func() {
			It("Can get serialization of User", func() {
				additionalS := Serialization{
					Rmv: []string{"user_name", "updated_at", "id"},
					Add: map[string]string{
						"nick_name": "NickName",
					},
				}
				s, err := Serialize(User{
					ID:        1,
					UserName:  "Bob",
					UpdatedAt: time.Now(),
				}, additionalS)
				Expect(len(s)).To(Equal(2))
				Expect(s["id"]).To(BeNil())
				Expect(len(err)).To(Equal(0))
			})
		})

		When("passing object with not existing function to add", func() {
			It("Panic", func() {
				additionalS := Serialization{
					Rmv: []string{"user_name", "drop_time"},
					Add: map[string]string{
						"nick_name":   "NickName",
						"family_name": "FamilyName",
					},
				}

				Expect(func() {
					Serialize(User{
						ID:        1,
						UserName:  "Bob",
						UpdatedAt: time.Now(),
					}, additionalS)
				}).To(Panic())

			})
		})

		When("passing object with not existing attribute to remove", func() {
			It("Neglect ", func() {
				additionalS := Serialization{
					Rmv: []string{"drop_time"},
				}
				s, err := Serialize(User{
					ID:        1,
					UserName:  "Bob",
					UpdatedAt: time.Now(),
				}, additionalS)
				Expect(len(s)).To(Equal(4))
				Expect(s["id"]).To(Equal(float64(1)))
				Expect(len(err)).To(Equal(0))
			})
		})
	})

	Describe("SerializeData: get serialization of objects in slice", func() {
		When("passing object in slice", func() {
			It("Can get serialization of object in slice", func() {
				t := time.Now()
				u := User{
					ID:        1,
					UserName:  "Bob",
					UpdatedAt: t,
					DeletedAt: &t,
				}
				s, err := SerializeData([]User{u, u})
				value := reflect.Indirect(reflect.ValueOf(s))
				Expect(value.Len()).To(Equal(2))
				Expect(len(err)).To(Equal(0))
			})
		})

		When("passing single object", func() {
			It("Can get serialization of object", func() {
				t := time.Now()
				u := User{
					ID:        1,
					UserName:  "Bob",
					UpdatedAt: t,
					DeletedAt: &t,
				}
				s, err := SerializeData(u)
				value, ok := s.(map[string]interface{})
				Expect(ok).To(Equal(true))
				Expect(len(value)).To(Equal(4))
				Expect(value["id"]).To(Equal(float64(1)))
				Expect(len(err)).To(Equal(0))
			})
		})

		When("passing object with serialization error", func() {
			BeforeEach(func() {
				p = IsExpectedToCall(Serialize).AndReturn(
					map[string]interface{}{"id": float64(1)}, []error{errors.New("errors")})
			})

			It("Can get serialization of object", func() {
				t := time.Now()
				u := User{
					ID:        1,
					UserName:  "Bob",
					UpdatedAt: t,
					DeletedAt: &t,
				}
				s, err := SerializeData(u)
				value, ok := s.(map[string]interface{})
				Expect(ok).To(Equal(true))
				Expect(len(value)).To(Equal(1))
				Expect(value["id"]).To(Equal(float64(1)))
				Expect(len(err)).To(Equal(1))
			})
		})
	})
})
