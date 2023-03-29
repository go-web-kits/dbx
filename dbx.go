// CRUD methods optimization and encapsulation based on Gorm
package dbx

import (
	"log"
	"os"

	"github.com/go-web-kits/dbx/dbx_callback"
	"github.com/go-web-kits/utils/logx"
	"github.com/jinzhu/gorm"
)

// Configuration
var zeroCli = &gorm.DB{}
var Client = zeroCli
var DefaultPage = 1
var DefaultRows = 10
var DefaultLogger = gorm.Logger{log.New(os.Stdout, "\r\n", 0)}
var DefaultLogFormat = "" // "" (default) / "json"
var UnLog = false

type H = map[string]interface{}

type Chain struct {
	*gorm.DB
}

var c Chain

// @doc c.2
func chain(opts ...Opt) *Chain {
	if c == (Chain{}) {
		c = initChain()
	}

	db := c.New()
	opt := OptsPackGet(opts)
	switch tx := opt.Tx.(type) {
	case *gorm.DB:
		db = tx
	case string, bool:
		db = c.Begin().DB
	case Chain:
		db = tx.DB
	case *Chain:
		if tx != (*Chain)(nil) {
			db = tx.DB
		}
	}

	if !opt.SaveAssoc {
		db.InstantSet("gorm:save_associations", false)
	}

	for k, v := range opt.Set {
		db.InstantSet(k, v)
	}

	if opt.Model != nil {
		db = db.Model(opt.Model).InstantSet("dbx:model", opt.Model)
	}

	return (&Chain{db}).Log(opt)
}

var Conn = chain

// @doc c.9.2 Other Interface #1
func IdOf(obj interface{}) interface{} {
	return Conn().NewScope(obj).PrimaryKeyValue()
}

// @doc c.9.2 Other Interface #2
func IsNewRecord(obj interface{}) bool {
	return Conn().NewRecord(obj)
}

func (c *Chain) AddCallback(obj interface{}, action dbx_callback.Action) {
	o, ok := obj.(dbx_callback.AfterCommitI)
	if !ok {
		return
	}

	callbacks := []*dbx_callback.Info{}
	if v, ok := c.Get("dbx:callbacks"); ok {
		callbacks = v.([]*dbx_callback.Info)
	}
	c.InstantSet("dbx:callbacks", append(callbacks, &dbx_callback.Info{Func: o.AfterCommit, Action: action}))
}

func (c *Chain) RunCallbacks() {
	callbacks, ok := c.Get("dbx:callbacks")
	if !ok {
		return
	}
	for _, callback := range callbacks.([]*dbx_callback.Info) {
		callback.Func(callback.Action)
	}
}

func initChain() Chain {
	if Client == zeroCli {
		panic("Please Config DBx.Client")
	}

	c := Chain{Client}
	begin := func(scope *gorm.Scope) {
		// if _, ok := defaultScope.Get("dbx:transaction"); !ok {
		unlog, ok := scope.Get("dbx:unlog")
		if !(ok && unlog.(bool)) {
			logging(logx.Blod(logx.Magenta("BEGIN " + scope.TableName())))
		}
	}
	ending := func(scope *gorm.Scope) {
		unlog, ok := scope.Get("dbx:unlog")
		if scope.HasError() {
			if !(ok && unlog.(bool)) {
				logging(logx.Blod(logx.Red("ROLLBACK (OR COMMIT FAILED) " + scope.TableName())))
			}
			return
		}

		if !(ok && unlog.(bool)) {
			logging(logx.Blod(logx.Magenta("COMMIT " + scope.TableName())))
		}
		if _, ok := scope.Get("dbx:transaction"); !ok {
			(&Chain{scope.DB()}).RunCallbacks()
		}
	}

	c.Callback().Create().Before("gorm:begin_transaction").Register("dbx:gorm-begin_transaction", begin)
	c.Callback().Update().Before("gorm:begin_transaction").Register("dbx:gorm-begin_transaction", begin)
	c.Callback().Delete().Before("gorm:begin_transaction").Register("dbx:gorm-begin_transaction", begin)

	c.Callback().Create().After("gorm:commit_or_rollback_transaction").Register("dbx:gorm-cort", ending)
	c.Callback().Update().After("gorm:commit_or_rollback_transaction").Register("dbx:gorm-cort", ending)
	c.Callback().Delete().After("gorm:commit_or_rollback_transaction").Register("dbx:gorm-cort", ending)

	return c
}
