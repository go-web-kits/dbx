package dbx

import (
	"reflect"
	"strings"

	"github.com/go-web-kits/utils/mapx"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

type Combine []interface{}

type EQ map[string]interface{}
type IN map[string]interface{}
type LIKE map[string]interface{}

type PLAIN []interface{}
type OR []interface{}
type NOT []interface{}

// =====================================
// The methods returns Result{Data, Err}
// =====================================

// @doc c.4.3 Query Interface #1
func Where(out interface{}, condition interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	conn := Conn(opt)
	c := conn.ScopingModel(out, opts...).Where(condition)

	if opt.RelatedWith != nil {
		result := c.Related(out, opt)
		return dealWithTx(conn, result, opt)
	}

	if err := c.pagy(opt).order(opt).Find(out).Error; err != nil {
		return dealWithTx(conn, Result{Err: err}, opt)
	}

	if opt.Count {
		count, err := c.Count()
		return dealWithTx(conn, Result{out, err, count, nil}, opt)
	}

	return dealWithTx(conn, Result{Data: out}, opt)
}

// @doc c.4.3 Query Interface #3
func First(out interface{}, n ...int) Result {
	limit := 1
	if len(n) == 1 {
		limit = n[0]
	}
	err := Conn().ScopingModel(out).Limit(limit).Find(out).Error
	return Result{Data: out, Err: err}
}

// @doc c.4.3 Query Interface #2
func Find(out interface{}, condition interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	conn := Conn(opt)
	c := conn.ScopingModel(out, opts...).Where(condition)

	if opt.RelatedWith != nil {
		// TODO
		result := c.Related(out, opt)
		return dealWithTx(conn, result, opt)
	}

	var err error
	// TODO default scope order
	if opt.Order == "" {
		err = c.First(out).Error
	} else {
		err = c.order(opt).Limit(1).Find(out).Error
	}
	return dealWithTx(conn, Result{Data: out, Err: err}, opt)
}

var FindBy = Find

// @doc c.4.3 Query Interface #4
func FindById(out interface{}, id interface{}, opts ...Opt) Result {
	return Find(out, EQ{"id": id}, opts...)
}

// @doc c.4.3 Query Interface #5
func FindOrInitBy(out, condition interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	conn := Conn(opt)
	err := conn.ScopingModel(out, opts...).FindOrInitBy(condition)
	return dealWithTx(conn, Result{Data: out, Err: err}, opt)
}

// @doc c.4.3 Query Interface #6
func IsExists(value interface{}, condition interface{}, opts ...Opt) bool {
	switch v := value.(type) {
	case string:
		c, table := 0, strcase.ToSnake(inflection.Plural(v))
		// TODO Scoping
		Conn(opts...).Where(condition).DB.Table(table).Count(&c)
		return c > 0
	default:
		newObj := reflect.New(reflect.Indirect(reflect.ValueOf(value)).Type()).Interface()
		return Find(newObj, condition, opts...).Err == nil
	}
}

// ===============================================
// The (chaining) methods returns *Chain{*gorm.DB}
// ===============================================

// @doc c.4.1 Query Chain #1
func (c *Chain) Where(condition interface{}) *Chain {
	switch cond := condition.(type) {
	case nil:
		return c
	case Conditioner:
		return cond.Build(c)
	case mapx.Map:
		return EQ(cond).Build(c)
	default:
		return &Chain{c.DB.Where(condition)}
	}
}

func (c *Chain) order(opt Opt) *Chain {
	if opt.UniqBy != "" {
		return c
	}
	return c.Order(opt.Order, true)
}

// @doc c.4.1 Query Chain #2
func (c *Chain) Order(value interface{}, reorder ...bool) *Chain {
	if value == "" {
		return c
	}
	if len(reorder) == 0 {
		reorder = append(reorder, false)
	}
	return &Chain{c.DB.Order(value, reorder[0])}
}

// @doc c.4.1 Query Chain #4
func (c *Chain) Uniq(opts ...Opt) *Chain {
	opt := OptsPackGet(opts)
	if opt.UniqBy == "" {
		return c
	}

	if opt.UniqOrder == "" {
		return &Chain{c.Select("DISTINCT ON (?) *", opt.UniqBy)}
	} else {
		orderField := strings.Split(opt.UniqOrder, " ")[0]
		return &Chain{c.
			Select("DISTINCT ON ("+opt.UniqBy+") "+orderField+", *").
			Order(opt.UniqBy+", "+opt.UniqOrder, true)}
	}
}

func (c *Chain) pagy(opt Opt) *Chain {
	if opt.Page == 0 && opt.Rows == 0 {
		return c
	}
	return c.Pagy(opt)
}

// @doc c.4.1 Query Chain #5
func (c *Chain) Pagy(opts ...Opt) *Chain {
	opt := OptsPackGet(opts)
	if opt.Rows == 0 {
		opt.Rows = DefaultRows
	}
	if opt.Page == 0 {
		opt.Page = DefaultPage
	}
	return &Chain{c.Limit(opt.Rows).Offset((opt.Page - 1) * opt.Rows)}
}

func (c *Chain) Unpagy() *Chain {
	return &Chain{c.Limit(-1).Offset(-1)}
}

// ====================
// Normal Query Methods
// ====================

// @doc c.4.2 Query Method #1
// dbx.Model(&user).Where(...).Count()
func (c *Chain) Count() (int64, error) {
	var total int64
	err := c.DB.Count(&total).Error
	return total, err
}

// @doc c.4.2 Query Method #2
func (c *Chain) Pluck(fields interface{}, out interface{}) error {
	obj := c.GetModel()
	switch f := fields.(type) {
	case string:
		return c.DB.Find(obj).Pluck(f, out).Error
	// case []string:
	// 	return nil // TODO
	default:
		panic("Pluck option error")
	}
}

// @doc c.4.2 Query Method #5
func (c *Chain) FindOrInitBy(condition interface{}) error {
	return c.Where(condition).FirstOrInit(c.GetModel()).Error
}

// @doc c.4.2 Query Method #4
func (c *Chain) FindOut(out ...interface{}) Result {
	if len(out) == 1 {
		err := c.Find(out[0]).Error
		return Result{Data: out[0], Err: err}
	} else {
		obj := c.GetModel()
		err := c.Find(obj).Error
		return Result{Data: obj, Err: err}
	}
}
