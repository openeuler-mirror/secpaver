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
	"testing"
)

// test function
func TestCheckUnsafeArg(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			args:    args{"test-os_1.1"},
			wantErr: false,
		},
		{
			args:    args{"test-os_1+1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckUnsafeArg(tt.args.arg); (err != nil) != tt.wantErr {
				t.Errorf("CheckUnsafeArg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func TestCheckUnsafePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			args:    args{"/dir/path-123.txt"},
			wantErr: false,
		},
		{
			args:    args{"/dir../ect"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckUnsafePath(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CheckUnsafePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func TestCheckPort(t *testing.T) {
	type args struct {
		port string
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			args:    args{"65535"},
			wantErr: false,
		},
		{
			args:    args{"65536"},
			wantErr: true,
		},
		{
			args:    args{"-1"},
			wantErr: true,
		},
		{
			args:    args{"test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := CheckPort(tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("CheckPort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func TestCheckZipFileName(t *testing.T) {
	type args struct {
		zipName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{"test.zip"},
			wantErr: false,
		},
		{
			args:    args{"test@.zip"},
			wantErr: true,
		},
		{
			args:    args{"test"},
			wantErr: true,
		},
		{
			args:    args{"testdata/test.zip"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckZipFileName(tt.args.zipName); (err != nil) != tt.wantErr {
				t.Errorf("CheckZipFileName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func TestCheckValidSelinuxType(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{"unconfined_t"},
			wantErr: false,
		},
		{
			args:    args{"unconfined"},
			wantErr: false,
		},
		{
			args:    args{"test#_t"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckValidSelinuxType(tt.args.str); (err != nil) != tt.wantErr {
				t.Errorf("CheckValidSelinuxType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func TestCheckVersion(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{"1"},
			wantErr: false,
		},
		{
			args:    args{"1.1"},
			wantErr: false,
		},
		{
			args:    args{"1.12.1"},
			wantErr: false,
		},
		{
			args:    args{"1.12.1111"},
			wantErr: true,
		},
		{
			args:    args{"a"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckVersion(tt.args.str); (err != nil) != tt.wantErr {
				t.Errorf("CheckVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

