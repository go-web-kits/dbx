package dbx_callback

const (
	Create = Action("create")
	Update = Action("update")
	Delete = Action("delete")
	Multi  = Action("multi")
)

type Action = string

type Info struct {
	Func   func(string)
	Action Action
}

// TODO error
type AfterCommitI interface {
	AfterCommit(Action)
}
