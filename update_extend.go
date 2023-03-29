package dbx

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// @doc c.5.2 Update Interface #3
func Increment(obj interface{}, fields interface{}, opts ...Opt) Result {
	return chain().Increment(obj, fields, opts...)
}

// @doc c.5.2 Update Interface #3
func Decrement(obj interface{}, fields interface{}, opts ...Opt) Result {
	return chain().Decrement(obj, fields, opts...)
}

// @doc c.5.1 Update By Chain #3
func (c *Chain) Increment(obj interface{}, fields interface{}, opts ...Opt) Result {
	values := map[string]interface{}{}

	switch val := fields.(type) {
	case string:
		values[val] = gorm.Expr(val + " + 1")
	case []string:
		for _, k := range val {
			values[k] = gorm.Expr(k + " + 1")
		}
	case map[string]int:
		for k, v := range val {
			values[k] = gorm.Expr(k+" + ?", v)
		}
	default:
		return Result{Err: errors.New("Increment: do nothing")}
	}

	opt := OptsPackGet(opts)
	opt.SkipUniqValidate = true
	return c.UpdateBy(obj, values, opt)
}

// @doc c.5.1 Update By Chain #3
func (c *Chain) Decrement(obj interface{}, fields interface{}, opts ...Opt) Result {
	values := map[string]interface{}{}

	switch val := fields.(type) {
	case string:
		values[val] = gorm.Expr(val + " - 1")
	case []string:
		for _, k := range val {
			values[k] = gorm.Expr(k + " - 1")
		}
	case map[string]int:
		for k, v := range val {
			values[k] = gorm.Expr(k+" - ?", v)
		}
	default:
		return Result{Err: errors.New("Decrement: do nothing")}
	}

	opt := OptsPackGet(opts)
	opt.SkipUniqValidate = true
	return c.UpdateBy(obj, values, opt)
}
