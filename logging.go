package dbx

import (
	"github.com/jinzhu/gorm"
)

type logger interface {
	Print(v ...interface{})
	Println(v ...interface{})
}

type JSONLogger struct {
	gorm.LogWriter
}

func (logger JSONLogger) Print(values ...interface{}) {
	// TODO
	logger.Println(gorm.LogFormatter(values...)...)
}

func (c *Chain) Log(opt Opt) *Chain {
	db := c.DB
	if UnLog || opt.UnLog {
		return &Chain{db.LogMode(false).InstantSet("dbx:unlog", true)}
	} else {
		db = db.LogMode(true).InstantSet("dbx:unlog", false)
	}

	format, logger := DefaultLogFormat, logger(DefaultLogger)
	if opt.LogFormat != "" {
		format = opt.LogFormat
	}
	if opt.Logger != nil {
		logger = opt.Logger
	}

	switch format {
	case "json":
		logger = JSONLogger{logger}
	}

	db.SetLogger(logger)

	if opt.Debug {
		db = db.Debug()
	}
	return &Chain{db}
}
