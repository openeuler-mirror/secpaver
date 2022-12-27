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

package semodule

import (
	"reflect"
	"testing"
)

// test function
func TestGetModuleInfo(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		want    *Module
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{name: "abrt"},
			want: &Module{
				Name:     "abrt",
				Priority: 100,
				Status:   ModuleEnabled,
			},
			wantErr: false,
		},
		{
			args: args{name: "invalid"},
			want: &Module{
				Name:     "invalid",
				Priority: 0,
				Status:   ModuleNotExist,
			},
			wantErr: false,
		},
		{
			args:    args{name: ""},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetModuleInfo(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetModuleInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetModuleInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

