package dbx

import (
	"github.com/imdario/mergo"
)

// 请查阅文档中的使用说明或示例
type Opt struct {
	Page    int
	Rows    int
	Count   bool
	Order   string
	ReOrder bool

	Preload               interface{} // doc c.4.1 #3
	PreloadWithoutDefault bool
	Join                  []string
	RelatedWith           interface{}
	AssocField            string

	UniqBy    string
	UniqOrder string

	SkipCallback        bool
	SkipUniqValidate    bool
	SaveAssoc           bool
	UnscopeDefault      bool
	WithDeleted         bool
	Unscoped            bool // without deleted_at filter and default scopes
	UnscopeDefaultOrder bool

	Model    interface{}
	Tx       interface{} // "begin" or true | Chain{tx}
	TxCommit bool
	Set      map[string]interface{}

	Debug     bool
	UnLog     bool
	Logger    logger
	LogFormat string // "normal" / "json"
}

type With = Opt
type Be = Opt

func OptsPack(opts []Opt) (Opt, bool) {
	if len(opts) == 0 {
		return Opt{}, false
	}

	dest := opts[0]
	for _, opt := range opts[1:] {
		// TODO performance
		if err := mergo.Merge(&dest, opt); err != nil {
			return Opt{}, false
		}
	}

	return dest, true
}

func OptsPackGet(opts []Opt) Opt {
	opt, _ := OptsPack(opts)
	return opt
}

func (o Opt) Merge(opts []Opt) Opt {
	for _, opt := range opts {
		_ = mergo.Merge(&o, opt)
	}
	return o
}

func (o Opt) M(opt Opt) Opt {
	_ = mergo.Merge(&o, opt)
	return o
}
