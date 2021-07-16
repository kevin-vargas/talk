package f

import (
	"errors"
	"reflect"
)

var (
	ErrSourceNotArray = errors.New("source value is not an array")
	ErrReducerNotFunc = errors.New("reducer argument must be a function")
)

func Reduce(reducer, initialValue, source interface{}) (interface{}, error) {
	sourceValues := reflect.ValueOf(source)
	if sourceValues.Kind() != reflect.Slice {
		return nil, ErrSourceNotArray
	}
	reducerValue := reflect.ValueOf(reducer)
	if reducerValue.Kind() != reflect.Func {
		return nil, ErrReducerNotFunc
	}

	accumulator := reflect.ValueOf(initialValue)
	for i := 0; i < sourceValues.Len(); i++ {
		entry := sourceValues.Index(i)
		result := reducerValue.Call([]reflect.Value{
			accumulator,
			entry,
			reflect.ValueOf(i),
		})
		accumulator = result[0]
	}
	return accumulator.Interface(), nil
}
