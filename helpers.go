package dbx

import (
	"time"

	"github.com/go-web-kits/utils/logx"
)

func dealWithTx(c *Chain, result Result, opt Opt) Result {
	if opt.Tx == nil {
		return result
	}

	if result.Err != nil {
		c.Rollback()
		return result
	}

	if opt.TxCommit {
		cResult := c.Commit()
		result.Err = cResult.Err
		return result
	}

	result.Tx = c
	return result
}

func logging(s string) {
	if UnLog {
		return
	}
	DefaultLogger.Println("\n", logx.Yello("["+time.Now().Format("2006-01-02 15:04:05")+"]"), s)
}
