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
	"fmt"
	"reflect"
	"strings"
)

// IsExistItem determines if an item exists in the array
func IsExistItem(value interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}

	default:
		return false
	}

	return false
}

// RemoveRepeatedElement delete repeate element in string slice
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0, len(arr))
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}

		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// PrintTabulate print tabulate format parameters
func PrintTabulate(table [][]string, headers []string) {
	var widths = []int{}

	if len(headers) == 0 {
		return
	}

	for _, header := range headers {
		widths = append(widths, len(header))
	}

	for _, arr := range table {
		for i := range widths {
			if i >= len(arr) {
				arr = append(arr, "")
			}

			if widths[i] < len(arr[i]) {
				widths[i] = len(arr[i])
			}
		}
	}

	fmt.Println("")

	for i, header := range headers {
		if i < len(widths) {
			fmt.Printf("%s%*s", header, widths[i]+4-len(header), "")
		}
	}

	fmt.Println("")

	for _, arr := range table {
		for j, data := range arr {
			if j < len(widths) {
				fmt.Printf("%s%*s", data, widths[j]+4-len(data), "")
			}
		}

		fmt.Println("")
	}
}

// ShowStringsWithSpace change a string slice to a space-delimited string
func ShowStringsWithSpace(strs []string) string {
	var text string
	for i, str := range strs {
		text = text + str
		if i != len(strs)-1 {
			text = text + " "
		}
	}

	return text
}

// TrimSpaceAndTab trim all the space " " and tab "\t" in the string
func TrimSpaceAndTab(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	return strings.Replace(str, " ", "", -1)
}

// PrintAck print ack with level and message information
func PrintAck(level, msg string) {
	fmt.Printf("[%s]: %s\n", level, msg)
}

// SplitByLastUnderscore split a string by last underscore
func SplitByLastUnderscore(str string) (string, string) {
	if !strings.Contains(str, "_") {
		return str, ""
	}

	idx := strings.LastIndex(str, "_")
	if idx == len(str)-1 {
		return str[:idx], ""
	}

	return str[:idx], str[idx+1:]
}
