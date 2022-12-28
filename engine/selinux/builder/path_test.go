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

package builder

import "testing"

// test function
func Test_changeWildcardToRegExp(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		args args
		name string
		want string
	}{
		{
			args: args{"/dir"},
			want: "/dir",
		},
		{
			args: args{"/dir/*"},
			want: "/dir/[^/]*",
		},
		{
			args: args{"/dir/**"},
			want: "/dir/.*",
		},
		{
			args: args{"/dir{,/*}"},
			want: "/dir(/[^/]*)?",
		},
		{
			args: args{"/dir{,/**}"},
			want: "/dir(/.*)?",
		},
		{
			args: args{"/path/?/dir"},
			want: "/path/[^/]/dir",
		},
		{
			args: args{"/path\\*"},
			want: "/path\\*",
		},
		{
			args: args{"/path\\{"},
			want: "/path\\{",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := changeWildcardToRegExp(tt.args.path); got != tt.want {
				t.Errorf("changeWildcardToRegExp() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_getDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{"/dir/file"},
			want: "/dir",
		},
		{
			args: args{"/dir/file{,/**}"},
			want: "/dir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDir(tt.args.path); got != tt.want {
				t.Errorf("getDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_getBase(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{"/dir/file"},
			want: "file",
		},
		{
			args: args{"/dir/file{,/**}"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBase(tt.args.path); got != tt.want {
				t.Errorf("getBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_dealLinkPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{"/bin/test"},
			want: "/usr/bin/test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dealLinkPath(tt.args.path); got != tt.want {
				t.Errorf("dealLinkPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

