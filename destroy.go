// @doc c.6.2 Destroy
package dbx

import (
	"github.com/go-web-kits/dbx/dbx_callback"
)

func Destroy(obj interface{}, opts ...Opt) Result {
	if IsNewRecord(obj) {
		panic("Do not destroy all by using `Destroy` method, check if the object you passed has an ID")
	}
	return DestroyAll(obj, opts...)
}

func DestroyAll(model interface{}, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	c := Conn(opt).ScopingModel(model, opt.M(Opt{UnscopeDefaultOrder: true}))

	if !opt.SkipCallback {
		c.AddCallback(model, dbx_callback.Delete)
	}
	err := c.Delete(model).Error

	return dealWithTx(c, Result{Data: model, Err: err}, opt)
}

func (c *Chain) Destroy(opts ...Opt) Result {
	opt, obj := OptsPackGet(opts), c.GetModel()
	result := c.Find(obj)
	if result.Error != nil {
		return Result{Err: result.Error}
	}

	if !opt.SkipCallback {
		c.AddCallback(obj, dbx_callback.Delete)
	}
	err := c.Delete(obj).Error

	return dealWithTx(c, Result{Data: obj, Err: err}, opt)
}
