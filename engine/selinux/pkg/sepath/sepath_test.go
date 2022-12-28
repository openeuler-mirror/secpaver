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

package sepath

import "testing"

// test function
func TestGetFixedPrefix(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{path: "/test"},
			want: "/test",
		},
		{
			args: args{path: "/test.*"},
			want: "/test",
		},
		{
			args: args{path: "/test\\*.*"},
			want: "/test\\*",
		},
		{
			args: args{path: "/test/*"},
			want: "/test/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFixedPrefix(tt.args.path); got != tt.want {
				t.Errorf("GetFixedPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

