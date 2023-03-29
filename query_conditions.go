package dbx

import (
	"fmt"
	"strings"

	"github.com/go-web-kits/secure_storage"
)

type Conditioner interface {
	Build(c *Chain) *Chain
}

func (cond Combine) Build(c *Chain) *Chain {
	_c := c
	for _, cond := range cond {
		_c = _c.Where(cond)
	}
	return _c
}

func (cond EQ) Build(c *Chain) *Chain {
	return &Chain{c.DB.Where(secure_storage.FilterMap(c.DB, cond))}
}

func (cond IN) Build(c *Chain) *Chain {
	db := c.DB
	for k, v := range secure_storage.FilterMap(c.DB, cond) {
		db = db.Where(k+" IN (?)", v)
	}
	return &Chain{db}
}

func (cond LIKE) Build(c *Chain) *Chain {
	keys, vals := []string{}, []interface{}{}
	for k, v := range cond {
		if v != "" {
			// TODO: For MySQL
			keys = append(keys, "LOWER("+k+") LIKE LOWER(?)")
			vals = append(vals, fmt.Sprint("%", v, "%"))
		}
	}

	if len(keys) == 0 {
		return c
	} else {
		return &Chain{c.DB.Where("("+strings.Join(keys, " OR ")+")", vals...)}
	}
}

func (cond PLAIN) Build(c *Chain) *Chain {
	db := c.DB

	if len(cond) > 1 {
		if vals, ok := cond[1].([]interface{}); ok {
			// PLAIN{"name = ? AND age = ?", []interface{"abc", 18}}
			db = db.Where(cond[0], vals...)
		} else {
			// PLAIN{"name = ? AND age = ?", "abc", 18}
			db = db.Where(cond[0], cond[1:]...)
		}
	} else {
		// PLAIN{"name = 'will'"}
		db = db.Where(cond[0])
	}

	return &Chain{db}
}

func (cond OR) Build(c *Chain) *Chain {
	return &Chain{c.DB.Or(cond[0], cond[1:]...)}
}

func (cond NOT) Build(c *Chain) *Chain {
	return &Chain{c.DB.Not(cond[0], cond[1:]...)}
}
