// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package casbinhelper

import (
	"fmt"
	"reflect"
)

func Access(args ...interface{}) (interface{}, error) {
	vCurrent := reflect.ValueOf(args[0])
	for _, field := range args {

		if vCurrent.Kind() == reflect.Pointer {
			vCurrent = vCurrent.Elem()
		}

		//index?
		if indexFloat, ok := field.(float64); ok {
			index := int(indexFloat)
			if vCurrent.Kind() != reflect.Array && vCurrent.Kind() != reflect.Slice {
				return nil, fmt.Errorf("appling numeric accessor to non-array")
			}
			vCurrent = vCurrent.Index(index)
			continue
		}
		if attr, ok := field.(string); ok {
			if vCurrent.Kind() != reflect.Struct {
				return nil, fmt.Errorf("appling string accessor to non-struct")
			}
			vCurrent = vCurrent.FieldByName(attr)
			if vCurrent.IsZero() {
				return nil, fmt.Errorf("no attribut %s found", attr)
			}
		}
	}
	return vCurrent.Interface(), nil

}
