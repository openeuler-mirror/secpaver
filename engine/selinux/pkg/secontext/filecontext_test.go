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

package secontext

import (
	"reflect"
	"regexp"
	"testing"
)

var testFileContexts = []*FileContext{
	{
		Path:        "/.*",
		Context:     Context{"system_u", "object_r", "default_t", "s0", ""},
		Prefix:      "/",
		HasMetaChar: true,
		IsValid:     true,
		Class:       UnknownFile,
	},
	{
		Path:        "/[^/]+",
		Context:     Context{},
		Class:       ComFile,
		Prefix:      "/",
		HasMetaChar: true,
		IsValid:     false,
	},
	{
		Path:        "/usr/bin/bash",
		Context:     Context{"system_u", "object_r", "shell_exec_t", "s0", ""},
		Class:       ComFile,
		Prefix:      "",
		HasMetaChar: false,
		IsValid:     true,
	},
	{
		Path:        "/test\\.txt",
		Context:     Context{"system_u", "object_r", "default_t", "s0", ""},
		Class:       ComFile,
		Prefix:      "",
		HasMetaChar: false,
		IsValid:     true,
	},
}

// test function
func TestNewFileContextFromString(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		want    *FileContext
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{"/.* system_u:object_r:default_t:s0"},
			want:    testFileContexts[0],
			wantErr: false,
		},
		{
			args:    args{"/.* --- system_u:object_r:default_t:s0"},
			want:    nil,
			wantErr: true,
		},
		{
			args:    args{"/[^/]+   --   <<none>>"},
			want:    testFileContexts[1],
			wantErr: false,
		},
		{
			args:    args{"/usr/bin/bash -- system_u:object_r:shell_exec_t:s0"},
			want:    testFileContexts[2],
			wantErr: false,
		},
		{
			args:    args{`/test\.txt -- system_u:object_r:default_t:s0`},
			want:    testFileContexts[3],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFileContextFromString(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFileContextFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFileContextFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestFileContext_Match(t *testing.T) {
	type args struct {
		path  string
		class FileClass
	}
	tests := []struct {
		fields *FileContext
		name   string
		args   args
		want   bool
	}{
		{
			fields: testFileContexts[0],
			args:   args{"/test", ComFile},
			want:   true,
		},
		{
			fields: testFileContexts[0],
			args:   args{"test", ComFile},
			want:   false,
		},
		{
			fields: testFileContexts[1],
			args:   args{"/test", ChrFile},
			want:   false,
		},
		{
			fields: testFileContexts[1],
			args:   args{"/test/file", ComFile},
			want:   false,
		},
		{
			fields: testFileContexts[2],
			args:   args{"/usr/bin/bash", UnknownFile},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := tt.fields
			if got := fc.Match(tt.args.path, tt.args.class); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestFileContext_Text(t *testing.T) {
	type fields struct {
		Reg         *regexp.Regexp
		Path        string
		Context     Context
		Class       FileClass
		Prefix      string
		HasMetaChar bool
		IsValid     bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				Path:    "/path",
				Class:   UnknownFile,
				Context: *CreateDefaultObjectContext("test_t"),
			},
			want: "/path\t  \tsystem_u:object_r:test_t:s0\n",
		},
		{
			fields: fields{
				Path:    "/path",
				Class:   ComFile,
				Context: *CreateDefaultObjectContext("test_t"),
			},
			want: "/path\t--\tsystem_u:object_r:test_t:s0\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FileContext{
				Reg:         tt.fields.Reg,
				Path:        tt.fields.Path,
				Context:     tt.fields.Context,
				Class:       tt.fields.Class,
				Prefix:      tt.fields.Prefix,
				HasMetaChar: tt.fields.HasMetaChar,
				IsValid:     tt.fields.IsValid,
			}
			if got := fc.Text(); got != tt.want {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestFileContextCompare(t *testing.T) {
	type args struct {
		fc1 *FileContext
		fc2 *FileContext
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				fc1: nil,
				fc2: nil,
			},
			want: false,
		},
		{
			args: args{
				fc1: testFileContexts[0],
				fc2: testFileContexts[2],
			},
			want: false,
		},
		{
			args: args{
				fc1: testFileContexts[2],
				fc2: testFileContexts[3],
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileContextCompare(tt.args.fc1, tt.args.fc2); got != tt.want {
				t.Errorf("FileContextCompare() = %v, want %v", got, tt.want)
			}
		})
	}
}

