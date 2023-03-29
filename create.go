package dbx

import (
	"github.com/go-web-kits/dbx/dbx_callback"
	. "github.com/go-web-kits/lab/business_error"
	"github.com/jinzhu/gorm"
)

// @doc
func Create(obj interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	if !opt.SkipUniqValidate && IsDuplicateRecord(obj) {
		return Result{Err: CommonErrors[NotUnique]}
	}

	c := Conn(opt)
	return doCreate(c, obj, opts...)
}

func FirstOrCreate(obj interface{}, conditionObj interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	c := Conn(opt)

	err := c.Where(conditionObj).Find(obj).Error
	if gorm.IsRecordNotFoundError(err) {
		return doCreate(c, obj, opts...)
	}
	return dealWithTx(c, Result{Data: obj, Err: err}, opt)
}

func doCreate(c *Chain, obj interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	if !opt.SkipCallback {
		c.AddCallback(obj, dbx_callback.Create)
	}
	err := c.Create(obj).Error

	return dealWithTx(c, Result{Data: obj, Err: err}, opt)
}
