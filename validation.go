package dbx

import (
	"reflect"

	"github.com/go-web-kits/dbx/dbx_model"
	"github.com/go-web-kits/utils/mapx"
	"github.com/go-web-kits/utils/structx"
)

// db.IsDuplicateRecord(User{ID: 1, Name: "abc"})
// db.IsDuplicateRecord(User{ID: 1}, map[string]interface{}{"name": "abc"})
// @notice Model must be `model.Unique`
func IsDuplicateRecord(obj interface{}, values ...map[string]interface{}) bool {
	uniqueness := dbx_model.DefinitionOf(obj).Uniqueness
	if uniqueness == nil {
		return false
	}
	// TODO: zero value judge
	// TODO: db? json?
	objMap := structx.ToTagValue2FieldValueMap(obj, "db")
	if len(values) > 0 {
		objMap = mapx.Merge(objMap, values[0])
	}

	combine := Combine{}
	if !IsNewRecord(obj) {
		combine = append(combine, NOT{"id", IdOf(obj)})
	}

	// TODO judge nil
	// FIXME 改成兼容 sec_storage 的 OR
	switch config := uniqueness.(type) {
	case string:
		cond := EQ{}
		cond[config] = objMap[config]
		return IsExists(obj, append(combine, cond))
	case []string:
		cond := EQ{}
		for _, k := range config {
			cond[k] = objMap[k]
		}
		return IsExists(obj, append(combine, cond))
	case map[string][]string:
		for primaryFiled, fileds := range config {
			cond := EQ{}
			cond[primaryFiled] = objMap[primaryFiled]
			for _, f := range fileds {
				// FIXME: gorm不会进行拆箱判空，需要提前拆箱
				val := reflect.ValueOf(objMap[f])
				if val.Kind() == reflect.Ptr && val.IsNil() {
					cond[f] = nil
				} else {
					cond[f] = objMap[f]
				}
			}
			combine = append(combine, cond)
		}
		return IsExists(obj, combine, Opt{Unscoped: true})
	}

	panic("dbx.IsDuplicateRecord: definition error")
}
