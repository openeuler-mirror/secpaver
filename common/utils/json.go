/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2021. All rights reserved.
 * secPaver is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// UnmarshalJSONWithExtendParams unmarshal a json data, and store the undefined data
// to a map, keeping the raw json data
func UnmarshalJSONWithExtendParams(data []byte, s interface{}) (map[string][]byte, error) {
	if !isValidStructPtr(s) {
		return nil, fmt.Errorf("parameter s should be a valid struct pointer")
	}

	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}

	allData := map[string]interface{}{}
	if err := json.Unmarshal(data, &allData); err != nil {
		return nil, err
	}

	vElems := reflect.ValueOf(s).Elem()
	tElems := reflect.TypeOf(s).Elem()

	if vElems.NumField() != tElems.NumField() {
		return nil, fmt.Errorf("invalid struct data, fields number dosen't match")
	}

	extData := map[string][]byte{}
	for k, v := range allData {
		isExtendData := true

		for i := 0; i < tElems.NumField(); i++ {
			vField := vElems.Field(i)
			tField := tElems.Field(i)

			if (!vField.CanSet()) || (!vField.IsValid()) {
				continue
			}

			if getJSONFieldNameInTag(tField.Tag.Get("json")) == k {
				isExtendData = false
				break
			}
		}

		if isExtendData {
			data, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}

			extData[k] = data
		}
	}

	return extData, nil
}

// MarshalJSONWithExtendParams marshal a struct, and expand the undefined data
// to the top
func MarshalJSONWithExtendParams(s interface{}, extData map[string][]byte, extKey string) ([]byte, error) {
	if !isValidStructPtr(s) {
		return nil, fmt.Errorf("parameter s should be a valid struct pointer")
	}

	allData, err := structToMap(s)
	if err != nil {
		return nil, err
	}

	delete(allData, extKey)

	for k, v := range extData {
		allData[k] = json.RawMessage(v)
	}

	return json.Marshal(allData)
}

func getJSONFieldNameInTag(tag string) string {
	tags := strings.Split(tag, ",")
	if len(tags) >= 1 {
		return tags[0]
	}

	return ""
}

func structToMap(s interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return m, nil
}

func isValidStructPtr(s interface{}) bool {
	if s == nil {
		return false
	}

	if (reflect.TypeOf(s).Kind() == reflect.Ptr) &&
		(reflect.TypeOf(s).Elem().Kind() == reflect.Struct) {

		return true
	}

	return false
}
