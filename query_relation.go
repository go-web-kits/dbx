package dbx

import (
	"github.com/pkg/errors"

	"github.com/go-web-kits/utils"
	"github.com/go-web-kits/utils/structx"
	"github.com/jinzhu/inflection"
)

// @doc c.4.3 Query Interface #5
func Related(assoc interface{}, obj interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	c := Conn(opt)
	result := c.ScopingModel(assoc, opts...).RelatedWith(obj, opts...)
	return dealWithTx(c, result, opt)
}

// @doc c.4.2 Query Method #4
// db.Model(&user_groups).Related(&users) => will find out the users which associated with the user_groups
func (c *Chain) Related(assoc interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	obj := opt.RelatedWith
	if obj == nil {
		obj, _ = c.Get("dbx:model")
	}
	if obj == nil {
		return Result{Err: errors.New("related: model not set")}
	}

	return c.related(assoc, obj, opt)
}

// db.Model(&users).RelatedWith(&user_groups) => will find out the users which associated with the user_groups
func (c *Chain) RelatedWith(obj interface{}, opts ...Opt) Result {
	assoc, _ := c.Get("dbx:model")
	if assoc == nil {
		return Result{Err: errors.New("related: model not set")}
	}
	return c.related(assoc, obj, OptsPackGet(opts))
}

// ===========

func (c *Chain) related(assoc, obj interface{}, opt Opt) Result {
	assocName := AssocFieldName(assoc, obj, opt)
	var count int
	if opt.Count {
		assocs := c.Model(obj).Association(assocName)
		if assocs.Error != nil {
			return Result{Err: assocs.Error}
		}
		count = assocs.Count()
	}

	err := c.pagy(opt).order(opt).Model(obj).DB.Related(assoc, assocName).Error
	return Result{Data: assoc, Err: err, Total: count}
}

func AssocFieldName(assoc, obj interface{}, opt Opt) string {
	if opt.AssocField != "" {
		return opt.AssocField
	}

	typeName := utils.TypeNameOf(assoc)
	if structx.GetFieldValueOf(obj, typeName) == nil {
		pluraize := inflection.Plural(typeName)
		if structx.GetFieldValueOf(obj, pluraize) == nil {
			panic("db.Related.AssocFieldName: cannot get association")
		} else {
			return pluraize
		}
	} else {
		return typeName
	}
}
