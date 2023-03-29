package dbx

import (
	"reflect"

	"github.com/go-web-kits/utils/slicex"
	"github.com/jinzhu/gorm"
)

type Result struct {
	Data  interface{}
	Err   error
	Total interface{}
	Tx    *Chain
}

// func (r Result) Error() string {
// 	if r.Err != nil {
// 		return r.Err.Error()
// 	}
// 	return ""
// }

func (r Result) Uniq() Result {
	if r.Err != nil {
		return r
	}

	ids := []uint{}
	data := []interface{}{}
	value := reflect.Indirect(reflect.ValueOf(r.Data))
	if value.Type().Kind() != reflect.Slice {
		return r
	}

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i).Interface()
		id := IdOf(item).(uint)
		if slicex.IncludeUint(ids, id) {
			continue
		}
		data = append(data, item)
		ids = append(ids, id)
	}

	r.Data = data
	return r
}

func (r Result) Ids() ([]uint, error) {
	if r.Err != nil {
		return []uint{}, r.Err
	}

	return r.GetIds(), nil
}

func (r Result) GetIds() []uint {
	ids := []uint{}
	value := reflect.Indirect(reflect.ValueOf(r.Data))
	if value.Type().Kind() != reflect.Slice {
		return ids
	}

	for i := 0; i < value.Len(); i++ {
		item := value.Index(i).Interface()
		id := IdOf(item).(uint)
		if id == 0 || slicex.IncludeUint(ids, id) {
			continue
		}
		ids = append(ids, id)
	}

	return ids
}

func (r Result) NotFound() bool {
	return gorm.IsRecordNotFoundError(r.Err)
}

func (r Result) IsNewRecord() bool {
	return IdOf(r.Data) == uint(0)
}
