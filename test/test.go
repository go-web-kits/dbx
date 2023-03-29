package test

type S = []interface{}

var Models = S{&People{}, &Car{}, &Group{}, &Activity{}}

type People struct {
	ID   uint   `gorm:"primary_key" db:"id"`
	Age  int64  `db:"age"`
	Cash int64  `db:"cash"`
	Name string `db:"name"`

	Cars       []Car
	Groups     []Group    `gorm:"many2many:people_and_groups"`
	Activities []Activity `gorm:"polymorphic:Item"`

	CallbackAction string `gorm:"-" db:"-"`
}

func (u *People) AfterCommit(action string) {
	u.CallbackAction = action
}

type Car struct {
	Id     uint `gorm:"primary_key" db:"id"`
	Name   string
	Number string

	People   People
	PeopleId uint
}

type Group struct {
	ID   uint   `gorm:"primary_key" db:"id"`
	Name string `db:"name"`

	Users      []People   `gorm:"many2many:people_and_groups"`
	Activities []Activity `gorm:"polymorphic:Item"`
}

type Activity struct {
	ID    uint   `gorm:"primary_key" db:"id"`
	Title string `db:"title"`

	ItemID   uint   `db:"item_id"`
	ItemType string `db:"item_type"`
}
