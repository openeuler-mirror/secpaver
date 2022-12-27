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

package project

import (
	pb "gitee.com/openeuler/secpaver/api/proto/project"
	"reflect"
	"testing"
)

var resourceValid = `{"resourceList":[{"path":"/bin/test"}]}`
var resourceEmpty = `{}`
var resourceWithInvalidFormat = `{"resourceList":{"path":"/bin/test"}`
var resourceWithSameMacroAndGroupName = `{
"resourceList":[{"path":"/bin/test"}],
"macroList":[{"name":"test","value":"/bin/test"}],
"groupList":[{"name":"test","resources":["/bin/test"]}]
}`

// test function
func Test_parseProjectResources(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		want    *pb.ResourceInfo
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{[]byte(resourceValid)},
			want: &pb.ResourceInfo{
				ResourceList: []*pb.Resource{{Path: "/bin/test", Extends: map[string][]byte{}}},
			},
		},
		{
			args:    args{[]byte(resourceEmpty)},
			wantErr: true,
		},
		{
			args:    args{[]byte(resourceWithInvalidFormat)},
			wantErr: true,
		},
		{
			args:    args{[]byte(resourceWithSameMacroAndGroupName)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseProjectResources(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseProjectResources() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseProjectResources() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var testMacroList = []*pb.Macro{
	{Name: "TEST", Value: "/test"},
	{Name: "D1", Value: "$(D2)"},
	{Name: "D2", Value: "$(D3)"},
	{Name: "D3", Value: "$(D4)"},
	{Name: "D4", Value: "$(D5)"},
	{Name: "D5", Value: "$(D6)"},
	{Name: "D6", Value: "$(D7)"},
	{Name: "D7", Value: "$(D8)"},
	{Name: "D8", Value: "$(D9)"},
	{Name: "D9", Value: "$(D10)"},
	{Name: "D10", Value: "$(D11)"},
	{Name: "D11", Value: "$(D12)"},
	{Name: "D12", Value: "/final"},
}

// test function
func Test_isMacro(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		args args
		name string
		want bool
	}{
		{
			args: args{str: "$(TEST)"},
			want: true,
		},
		{
			args: args{str: "$(TEST )"},
			want: false,
		},
		{
			args: args{str: "TEST"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMacro(tt.args.str); got != tt.want {
				t.Errorf("isMacro() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestMacroReplace(t *testing.T) {
	type args struct {
		macros []*pb.Macro
		path   string
	}
	tests := []struct {
		args    args
		name    string
		want    string
		wantErr bool
	}{
		{
			args: args{path: "$(TEST)/file", macros: testMacroList},
			want: "/test/file",
		},
		{
			args:    args{path: "$(TEST/file", macros: testMacroList},
			wantErr: true,
		},
		{
			args: args{path: "$(D3)/file", macros: testMacroList},
			want: "/final/file",
		},
		{
			args:    args{path: "$(D2)/file", macros: testMacroList},
			wantErr: true,
		},
		{
			args: args{
				path:   "$(TEST)$(D3)$(D4)$(D5)$(D6)/file",
				macros: testMacroList,
			},
			want: "/test/final/final/final/final/file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MacroReplace(tt.args.path, tt.args.macros)
			if (err != nil) != tt.wantErr {
				t.Errorf("macroReplace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("macroReplace() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_getMacroName(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		args args
		name string
		want string
	}{
		{
			args: args{str: "$(TEST)"},
			want: "TEST",
		},
		{
			args: args{str: "$(TEST"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMacroName(tt.args.str); got != tt.want {
				t.Errorf("getMacroName() = %v, want %v", got, tt.want)
			}
		})
	}
}

