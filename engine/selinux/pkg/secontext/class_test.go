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
	"testing"
)

const invalidClassID FileClass = 100

// test function
func TestGetFileClassBySymbol(t *testing.T) {
	type args struct {
		symbol string
	}
	tests := []struct {
		args    args
		name    string
		want    FileClass
		wantErr bool
	}{
		{
			args:    args{symbol: "--"},
			want:    ComFile,
			wantErr: false,
		},
		{
			args:    args{symbol: "-d"},
			want:    DirFile,
			wantErr: false,
		},
		{
			args:    args{symbol: "---"},
			want:    UnknownFile,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileClassBySymbol(tt.args.symbol)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileClassBySymbol() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetFileClassBySymbol() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestGetFileSymbolbyID(t *testing.T) {
	type args struct {
		classID FileClass
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{SockFile},
			want: "-s",
		},
		{
			args: args{invalidClassID},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileSymbolbyID(tt.args.classID); got != tt.want {
				t.Errorf("GetFileSymbolbyID() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestGetFileClassByID(t *testing.T) {
	type args struct {
		classID FileClass
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{classID: ComFile},
			want: []string{"file"},
		},
		{
			args: args{classID: invalidClassID},
			want: []string{"file", "chr_file", "blk_file", "dir", "fifo_file", "lnk_file", "sock_file"},
		},
		{
			args: args{classID: UnknownFile},
			want: []string{"file", "chr_file", "blk_file", "dir", "fifo_file", "lnk_file", "sock_file"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileClassByID(tt.args.classID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileClassByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

