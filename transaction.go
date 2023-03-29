package dbx

import (
	"time"

	"github.com/go-web-kits/utils/logx"
)

func Transaction(lambda func(c *Chain) error, opts ...Opt) Result {
	opt := OptsPackGet(opts)
	opt.Tx = "begin"
	c := Conn(opt)
	err := lambda(c)
	if err == nil {
		return c.Commit()
	}
	return Result{Err: err}
}

func (c *Chain) Begin() *Chain {
	logging(logx.Blod(logx.Magenta("↓↓↓ TX BEGIN ↓↓↓")))
	return &Chain{c.New().Begin().InstantSet("dbx:transaction", time.Now())}
}

func (c *Chain) Commit() Result {
	start, ok := c.Get("dbx:transaction")
	if !ok {
		panic("dbx.Commit: transaction not begin")
	}

	err := c.DB.Commit().Error
	if err == nil {
		c.RunCallbacks()
	}

	duration := float64(time.Since(start.(time.Time)).Nanoseconds()/1e4) / 100.0
	logging(logx.Blod(logx.Magenta(duration, "↑↑↑ TX COMMIT [total: %.2fms] ↑↑↑")))
	return Result{Err: err}
}

func (c *Chain) Rollback() Result {
	start, ok := c.Get("dbx:transaction")
	if !ok {
		panic("dbx.Commit: transaction not begin")
	}

	err := c.DB.Rollback().Error
	duration := float64(time.Since(start.(time.Time)).Nanoseconds()/1e4) / 100.0
	logging(logx.Blod(logx.Red(duration, "↑↑↑ TX ROLLBACK [total: %.2fms] ↑↑↑")))
	return Result{Err: err}
}
