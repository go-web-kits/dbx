// @doc c.9.1 Model
package dbx

func Model(obj interface{}, opts ...Opt) *Chain {
	return Conn(opts...).Model(obj)
}

func ScopingModel(obj interface{}, opts ...Opt) *Chain {
	return Conn(opts...).ScopingModel(obj, opts...)
}

func (c *Chain) Model(obj interface{}) *Chain {
	return &Chain{c.DB.Model(obj).InstantSet("dbx:model", obj)}
}

func (c *Chain) ScopingModel(obj interface{}, opts ...Opt) *Chain {
	opt := OptsPackGet(opts)
	return c.Model(obj).Scoping(obj, opts...).preload(opt).joins(opt).Uniq(opt)
}

// ====================

func (c *Chain) GetModel() interface{} {
	obj, modeling := c.Get("dbx:model")
	if !modeling {
		panic("model not set")
	}
	return obj
}
