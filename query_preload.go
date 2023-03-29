package dbx

import (
	"strings"

	"github.com/jinzhu/gorm"
)

func (c *Chain) preload(opt Opt) *Chain {
	if opt.Preload != nil {
		return c.Preload(opt.Preload, opt.PreloadWithoutDefault)
	}
	return c
}

// @doc c.4.1 Query Chain #3
func (c *Chain) Preload(value interface{}, withoutDefault ...bool) *Chain {
	if len(withoutDefault) == 0 {
		withoutDefault = append(withoutDefault, false)
	}

	db := c.DB
	for _, info := range getPreloadInfos(value, withoutDefault[0]) {
		db = db.Preload(info.Column, info.Conditions...)
	}

	return &Chain{db}
}

type preloadInfo struct {
	Column     string
	Conditions []interface{}
}

func getPreloadInfos(value interface{}, withoutDefault bool) []preloadInfo {
	infos := convertPreloadInfo(value)
	if withoutDefault {
		return infos
	}

	scopedInfos := []preloadInfo{}
	for _, info := range infos {
		columns := strings.Split(info.Column, ".")
		fullPath := ""
		for _, col := range columns {
			if fullPath != "" {
				fullPath += "."
			}
			fullPath += col

			if len(info.Conditions) > 0 {
				scopedInfos = append(scopedInfos, info)
			} else {
				_col := col // Go 闭包特性？如果不复制一下值，闭包只能取到循环最后一个值
				scopedInfos = append(scopedInfos, preloadInfo{fullPath, []interface{}{
					func(db *gorm.DB) *gorm.DB {
						return (&Chain{db}).Scoping(_col).DB
					}}})
			}
		}
	}

	return scopedInfos
}

func convertPreloadInfo(value interface{}) []preloadInfo {
	infos := []preloadInfo{}
	switch _value := value.(type) {
	case string:
		infos = []preloadInfo{{Column: _value}}
	case []string:
		for _, col := range _value {
			infos = append(infos, preloadInfo{Column: col})
		}
	case map[string][]interface{}:
		for col, conditions := range _value {
			infos = append(infos, preloadInfo{col, conditions})
		}
	default:
		panic("preloadInfo: invalid arguments")
	}
	return infos
}
