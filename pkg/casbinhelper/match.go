package casbinhelper

import (
	"fmt"
	"strings"
)

func HasPrefix(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("HasPrefix requires 2 parameters, currently %d", len(args))
	}
	str, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("HasPrefix requires 1st parameter to be string")
	}
	prefix, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("HasPrefix requires 2nd parameter to be string")
	}
	return strings.HasPrefix(str, prefix), nil
}
