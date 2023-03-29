package dbx_model_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/go-web-kits/dbx/dbx_model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDbxModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DbxModel Suite")
}

type User struct {
	ID        uint       `json:"id"`
	UserName  string     `json:"user_name"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Person struct{}

func (u User) Email() string {
	return fmt.Sprintf("%s@example.com", u.UserName)
}

func (u User) NickName() string {
	return fmt.Sprintf("%s@example.com", u.UserName)
}

var _ = BeforeSuite(func() {
	DBxDefinitions = map[string]Definition{
		"User": {
			Uniqueness: map[string][]string{
				"UserName": {"deleted_at"},
			},
			DefaultScope: Scope{
				Order: "updated_at DESC",
			},
			Serialization: Serialization{
				Rmv: []string{"deleted_at"},
				Add: map[string]string{
					"email": "Email",
				},
			},
		},
	}
})
