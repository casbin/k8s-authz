package casbinhelper

import (
	"fmt"
	"reflect"
)

func Index(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("Index requires 2 parameters, currently %d", len(args))
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("1st parameter should be array, currently %s", v.Kind().String())
	}

	indexString := args[1].(float64)
	index := int(indexString)

	if index >= v.Len() {
		return nil, fmt.Errorf("index out of range")
	}

	intf := v.Index(index).Interface()
	return intf, nil
}

func Len(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Index requires 2 parameters, currently %d", len(args))
	}
	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Array && v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("1st parameter should be array, currently %s", v.Kind().String())
	}
	return v.Len(), nil

}
