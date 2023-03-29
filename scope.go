package dbx

import (
	"github.com/go-web-kits/dbx/dbx_model"
)

// @doc c.4.1 Query Chain #6
func (c *Chain) Unscoped() *Chain {
	return &Chain{c.DB.Unscoped()}
}

// @doc c.4.1 Query Chain #7
func (c *Chain) Scoping(obj interface{}, opts ...Opt) *Chain {
	opt := OptsPackGet(opts)
	if opt.Unscoped {
		return c.Unscoped()
	}

	_c := c
	if opt.WithDeleted {
		_c = _c.Unscoped()
	}

	if !opt.UnscopeDefault {
		_c = _c.defaultScope(dbx_model.DefinitionOf(obj).DefaultScope, opts...)
	}

	return _c
}

func (c *Chain) ScopingJoined(joinedModel string) *Chain {
	scope := dbx_model.DefinitionOf(joinedModel).DefaultScope
	cond := scope.WhereBeJoined
	if len(cond) > 0 {
		c = &Chain{c.DB.Where(cond[0], cond[1:]...)}
	}

	if scope.OrderBeJoined != "" {
		c = c.Order(scope.OrderBeJoined)
	}

	return c
}

// process all the scopes defined by a `model.Scope`
func (c *Chain) defaultScope(scope dbx_model.Scope, opts ...Opt) *Chain {
	_c, opt := c, OptsPackGet(opts)

	if !opt.UnscopeDefaultOrder && scope.Order != "" {
		_c = _c.Order(scope.Order, true)
	}

	if scope.Where != nil {
		_c = _c.Where(scope.Where)
	}

	if scope.Preload != nil {
		_c = _c.Preload(scope.Preload, opt.PreloadWithoutDefault)
	}

	if len(scope.Join) > 0 {
		_c = _c.joins(Opt{Join: scope.Join})
	}

	return _c
}
