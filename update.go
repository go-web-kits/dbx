package dbx

import (
	"reflect"

	"github.com/go-web-kits/dbx/dbx_callback"
	"github.com/go-web-kits/dbx/dbx_model"
	. "github.com/go-web-kits/lab/business_error"
	"github.com/go-web-kits/utils"
	"github.com/go-web-kits/utils/mapx"
	"github.com/go-web-kits/utils/structx"
)

// @doc c.5.2 Update Interface #1
func Update(obj interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	return update(ScopingModel(obj, opt.M(Opt{UnscopeDefaultOrder: true})), obj, opts...)
}

// @doc c.5.2 Update Interface #2
func UpdateBy(obj interface{}, values interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	return updateBy(ScopingModel(obj, opt.M(Opt{UnscopeDefaultOrder: true})), obj, getValues(values), opts...)
}

// @doc c.5.1 Update By Chain #1
func (c *Chain) Update(obj interface{}, opts ...Opt) Result {
	result := c.Find(utils.Clone(obj))
	if result.Error != nil {
		return Result{Err: result.Error}
	}
	return update(c, obj, opts...)
}

// @doc c.5.1 Update By Chain #2
func (c *Chain) UpdateBy(obj interface{}, values interface{}, opts ...Opt) Result {
	result := c.Find(obj)
	if result.Error != nil {
		return Result{Err: result.Error}
	}
	return updateBy(c, obj, getValues(values), opts...)
}

// ====================

func update(c *Chain, obj interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	if !opt.SkipUniqValidate && IsDuplicateRecord(obj) {
		return Result{Err: CommonErrors[NotUnique]}
	}
	return doUpdate(c, obj, obj, opt)
}

func updateBy(c *Chain, obj interface{}, values map[string]interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	if !opt.SkipUniqValidate && IsDuplicateRecord(obj, values) {
		return Result{Err: CommonErrors[NotUnique]}
	}
	return doUpdate(c, obj, values, opt)
}

func doUpdate(c *Chain, obj interface{}, values interface{}, opt Opt) Result {
	var err error
	reflectedObj := reflect.Indirect(reflect.ValueOf(obj))
	if opt.SkipCallback {
		if reflectedObj.Kind() == reflect.Slice {
			tableName := dbx_model.TableNameOf(reflectedObj.Index(0).Interface())
			err = c.Table(tableName).Where("id IN (?)", Result{Data: obj}.GetIds()).UpdateColumns(values).Error
		} else {
			err = c.Model(obj).UpdateColumns(values).Error
		}
	} else {
		// TODO: Batch update with callbacks
		c.AddCallback(obj, dbx_callback.Update)
		err = c.Model(obj).Updates(values).Error
	}

	return dealWithTx(c, Result{Data: obj, Err: err}, opt)
}

func getValues(values interface{}) map[string]interface{} {
	var _values map[string]interface{}
	if v, ok := values.(map[string]interface{}); ok {
		_values = v
	} else if v, ok := values.(mapx.Map); ok {
		_values = map[string]interface{}(v)
	} else {
		_values = structx.ToTagValue2FieldValueMap(values, "json")
		delete(_values, "id") // TODO
	}
	return _values
}
