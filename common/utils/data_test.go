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
	"reflect"
	"sort"
	"testing"
)

// test function
func TestRemoveRepeatedElement(t *testing.T) {
	type args struct {
		arr []string
	}
	tests := []struct {
		name       string
		args       args
		wantNewArr []string
		want       []string
	}{
		{
			args:       args{[]string{"aa", "bb", "aa", "aa"}},
			wantNewArr: nil,
			want:       []string{"aa", "bb"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := RemoveRepeatedElement(tt.args.arr)
			sort.Strings(res)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(res, tt.want) {
				t.Errorf("RemoveRepeatedElement() = %v, want %v", res, tt.want)
			}
		})
	}
}

// test function
func TestIsExistItem(t *testing.T) {
	type args struct {
		value interface{}
		array interface{}
	}
	tests := []struct {
		args args
		name string
		want bool
	}{
		{
			args: args{"a", []string{"a", "b", "c", "d"}},
			want: true,
		},
		{
			args: args{"1", []string{"a", "b", "c", "d"}},
			want: false,
		},
		{
			args: args{1, []string{"a", "b", "c", "d"}},
			want: false,
		},
		{
			args: args{nil, []string{"a", "b", "c", "d"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExistItem(tt.args.value, tt.args.array); got != tt.want {
				t.Errorf("IsExistItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestShowStringsWithSpace(t *testing.T) {
	type args struct {
		strs []string
	}
	tests := []struct {
		args args
		name string
		want string
	}{
		{
			args: args{[]string{"a", "b", "c"}},
			want: "a b c",
		},
		{
			args: args{[]string{"a"}},
			want: "a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShowStringsWithSpace(tt.args.strs); got != tt.want {
				t.Errorf("ShowStringsWithSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestTrimSpaceAndTab(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		args args
		name string
		want string
	}{
		{
			args: args{" a b c    d	 "},
			want: "abcd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimSpaceAndTab(tt.args.str); got != tt.want {
				t.Errorf("TrimSpaceAndTab() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestPrintTabulate(t *testing.T) {
	type args struct {
		table   [][]string
		headers []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				headers: []string{"name", "version", "mod time"},
				table: [][]string{
					{"1", "123", "1234"},
					{"123456", "123456", "123456"},
				},
			},
		},
		{
			args: args{
				headers: []string{"name", "version", "mod time"},
				table:   [][]string{},
			},
		},
		{
			args: args{
				headers: []string{"name", "version", "mod time"},
				table:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintTabulate(tt.args.table, tt.args.headers)
		})
	}
}

// test function
func TestSplitByLastUnderscore(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			args:  args{str: "test_selinux"},
			want:  "test",
			want1: "selinux",
		},
		{
			args:  args{str: "test_test_selinux"},
			want:  "test_test",
			want1: "selinux",
		},
		{
			args:  args{str: "test"},
			want:  "test",
			want1: "",
		},
		{
			args:  args{str: "test_"},
			want:  "test",
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SplitByLastUnderscore(tt.args.str)
			if got != tt.want {
				t.Errorf("SplitByLastUnderscore() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitByLastUnderscore() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

