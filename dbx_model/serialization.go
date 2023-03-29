package dbx_model

import (
	"reflect"
	"sync"

	"github.com/go-web-kits/utils/mapx"
	"github.com/go-web-kits/utils/structx"
)

func SerializeData(data interface{}, serializations ...Serialization) (interface{}, []error) {
	value := reflect.Indirect(reflect.ValueOf(data))
	switch value.Kind() {
	case reflect.Slice:
		result := make([]interface{}, value.Len())
		errs := []error{}

		waitGroup := sync.WaitGroup{}
		for i := 0; i < value.Len(); i++ {
			waitGroup.Add(1)
			go func(index int) {
				defer waitGroup.Done()
				r, e := Serialize(value.Index(index).Interface(), serializations...)
				result[index] = r
				errs = append(errs, e...)
			}(i)
		}
		waitGroup.Wait()
		return result, errs
	default:
		return Serialize(value.Interface(), serializations...)
	}
}

func Serialize(obj interface{}, serializations ...Serialization) (map[string]interface{}, []error) {
	serialized := structx.ToJsonizeMap(obj)
	serialization := DefinitionOf(obj).Serialization
	addition := mapx.Copy(serialization.Add).(map[string]string)
	prevention := serialization.Rmv
	if len(serializations) > 0 {
		prevention = append(prevention, serializations[0].Rmv...)
		for k, v := range serializations[0].Add {
			addition[k] = v
		}
	}

	for _, key := range prevention {
		delete(addition, key)
		delete(serialized, key)
	}

	var errs []error
	for key, methodName := range addition {
		results := reflect.ValueOf(obj).MethodByName(methodName).Call([]reflect.Value{})
		if len(results) == 2 {
			if err, ok := results[1].Interface().(error); ok && err != nil {
				errs = append(errs, err)
			}
		}
		serialized[key] = results[0].Interface()
	}

	return serialized, errs
}

func SerializeWithoutError(obj interface{}, serializations ...Serialization) map[string]interface{} {
	result, _ := Serialize(obj, serializations...)
	return result
}
