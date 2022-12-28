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

package sehandle

import (
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"os"
	"reflect"
	"regexp"
	"testing"
	"time"
)

var testHandle Handle

// test function
func Test_handleImpl_GetFileContext(t *testing.T) {
	h := testHandle

	type args struct {
		path string
	}
	tests := []struct {
		want *secontext.FileContext
		name string
		args args
	}{
		{
			args: args{path: "/bin/bash"},
			want: &secontext.FileContext{
				Path:        "/bin/bash",
				Context:     *secontext.CreateDefaultObjectContext("shell_exec_t"),
				Class:       secontext.ComFile,
				HasMetaChar: false,
				IsValid:     true,
			},
		},
		{
			args: args{path: ""},
			want: nil,
		},
	}
	for _, tt := range tests {
		if got := h.GetFileContext(tt.args.path); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("GetFileContext() = %v, want %v", got, tt.want)
		}
	}
}

// test function
func Test_handleImpl_LookupFileContext(t *testing.T) {
	h := testHandle

	type args struct {
		path  string
		class secontext.FileClass
	}
	tests := []struct {
		want *secontext.FileContext
		name string
		args args
	}{
		{
			args: args{path: "/bin/bash", class: secontext.ComFile},
			want: &secontext.FileContext{
				Reg:         regexp.MustCompile("^/bin/bash$"),
				Path:        "/bin/bash",
				Context:     *secontext.CreateDefaultObjectContext("shell_exec_t"),
				Class:       secontext.ComFile,
				HasMetaChar: false,
				IsValid:     true,
			},
		},
		{
			args: args{path: "/bin/test", class: secontext.ComFile},
			want: &secontext.FileContext{
				Reg:         regexp.MustCompile("^/bin/.*$"),
				Path:        "/bin/.*",
				Context:     *secontext.CreateDefaultObjectContext("bin_t"),
				Class:       secontext.UnknownFile,
				Prefix:      "/bin/",
				HasMetaChar: true,
				IsValid:     true,
			},
		},
		{
			args: args{},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.LookupFileContext(tt.args.path, tt.args.class); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LookupFileContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_handleImpl_TypeHasDefined(t *testing.T) {
	h := testHandle

	type args struct {
		tp string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{tp: "unconfined_t"},
			want: true,
		},
		{
			args: args{tp: "userdomain"},
			want: false,
		},
		{
			args: args{tp: ""},
			want: false,
		},
		{
			args: args{tp: "invalid"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.TypeHasDefined(tt.args.tp); got != tt.want {
				t.Errorf("TypeHasDefined() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_handleImpl_AttrHasDefined(t *testing.T) {
	h := testHandle

	type args struct {
		attr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{attr: "unconfined_t"},
			want: false,
		},
		{
			args: args{attr: "userdomain"},
			want: true,
		},
		{
			args: args{attr: ""},
			want: false,
		},
		{
			args: args{attr: "invalid"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := h.AttrHasDefined(tt.args.attr); got != tt.want {
				t.Errorf("AttrHasDefined() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_handleImpl_LookupPortContext(t *testing.T) {
	h := testHandle

	type args struct {
		port  uint32
		proto uint32
	}
	tests := []struct {
		want *secontext.PortContext
		name string
		args args
	}{
		{
			args: args{port: 8080, proto: secontext.ProtoTCP},
			want: &secontext.PortContext{
				LowPort:  8080,
				HighPort: 8080,
				Protocol: secontext.ProtoTCP,
				Context:  *secontext.CreateDefaultObjectContext("http_cache_port_t"),
			},
		},
		{
			args: args{port: 18120, proto: secontext.ProtoUDP},
			want: &secontext.PortContext{
				LowPort:  18120,
				HighPort: 18121,
				Protocol: secontext.ProtoUDP,
				Context:  *secontext.CreateDefaultObjectContext("radius_port_t"),
			},
		},
		{
			args: args{port: 999999, proto: secontext.ProtoTCP},
		},
		{
			args: args{port: 8080, proto: 9999},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := h.LookupPortContext(tt.args.port, tt.args.proto)
			if got != nil {
				got.Context.Sensitivity = "s0"
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LookupPortContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_handleImpl_GetRoles(t *testing.T) {
	h := testHandle

	roles := h.GetRoles()
	if !utils.IsExistItem("unconfined_r", roles) {
		t.Errorf("fail to test GetRoles")
	}
}

// main test function
func TestMain(m *testing.M) {
	// handle can't be opened too frequently, maybe other testcase just open the handle
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second * 5)
		testHandle, _ = HandleCreate()
		if testHandle != nil {
			break
		}
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}

