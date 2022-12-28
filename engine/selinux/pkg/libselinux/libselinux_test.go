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

package libselinux

import (
	"testing"
)

// test function
func Test_FileContextPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			want: "/etc/selinux/targeted/contexts/files/file_contexts",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileContextPath(); got != tt.want {
				t.Errorf("selinuxFileContextPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_FileContextHomedirPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			want: "/etc/selinux/targeted/contexts/files/file_contexts.homedirs",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileContextHomedirPath(); got != tt.want {
				t.Errorf("selinuxFileContextHomedirPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_FileContextLocalPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			want: "/etc/selinux/targeted/contexts/files/file_contexts.local",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileContextLocalPath(); got != tt.want {
				t.Errorf("selinuxFileContextLocalPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestGetPolicyType(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{
			want:    "targeted",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPolicyType()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPolicyType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPolicyType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

