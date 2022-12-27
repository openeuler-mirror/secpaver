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

package serule

import "testing"

// test function
func Test_checkClass(t *testing.T) {
	type args struct {
		class string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{class: "file"},
			wantErr: false,
		},
		{
			args:    args{class: "invalid_class"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkClass(tt.args.class); (err != nil) != tt.wantErr {
				t.Errorf("checkClass() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func Test_checkPermission(t *testing.T) {
	type args struct {
		class string
		perms []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				class: "file",
				perms: []string{"read", "write"},
			},
			wantErr: false,
		},
		{
			args: args{
				class: "file",
				perms: []string{"invalid_perm"},
			},
			wantErr: true,
		},
		{
			args: args{
				class: "invalid_class",
				perms: []string{"read"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkPermission(tt.args.class, tt.args.perms); (err != nil) != tt.wantErr {
				t.Errorf("checkPermission() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

