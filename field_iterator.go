package envloader

import (
	"errors"
	"reflect"
)

type fieldIterator struct {
	iterable   interface{}
	current    int
	total      int
	structElem reflect.Type
	valElem    reflect.Value
}

func iterableType(i interface{}) (f fieldIterator, err error) {
	f = fieldIterator{
		iterable: i,
		current:  -1,
	}

	ptrType := reflect.TypeOf(i)
	ptrVal := reflect.ValueOf(i)
	if ptrType.Kind() != reflect.Ptr {
		err = errors.New("ENVLOADER: Supplied variable is not a pointer")
		return
	}

	f.structElem = ptrType.Elem()
	f.valElem = ptrVal.Elem()

	f.total = f.valElem.NumField()

	return
}

// do not run without running Next() first
func (f *fieldIterator) Get() (reflect.Value, reflect.StructField) {
	return f.valElem.Field(f.current), f.structElem.Field(f.current)
}

// will return true until it runs out of elements
func (f *fieldIterator) Next() bool {
	f.current++

	return f.current < f.total
}
