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

import (
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"reflect"
	"testing"
)

// test function
func TestCreateFileAllowRules(t *testing.T) {
	type args struct {
		sub       string
		obj       string
		FileClass secontext.FileClass
		acts      []string
	}
	tests := []struct {
		name    string
		args    args
		want    []Rule
		wantErr bool
	}{
		{
			args:    args{"subject", "object", secontext.ComFile, []string{"read", "write"}},
			want:    []Rule{&AvcRule{"allow", "subject", "object", "file", []string{"read", "write"}}},
			wantErr: false,
		},
		{
			args:    args{"subject", "object", secontext.ComFile, []string{"send"}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateFileAllowRules(tt.args.sub, tt.args.obj, tt.args.FileClass, tt.args.acts)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFileAllowRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFileAllowRules() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestCreateCapabilityAllowRule(t *testing.T) {
	type args struct {
		sub  string
		caps []string
	}
	tests := []struct {
		name    string
		args    args
		want    []Rule
		wantErr bool
	}{
		{
			args:    args{sub: "subject", caps: []string{"chown"}},
			want:    []Rule{&AvcRule{"allow", "subject", "subject", "capability", []string{"chown"}}},
			wantErr: false,
		},
		{
			args:    args{sub: "subject", caps: []string{"mac_override"}},
			want:    []Rule{&AvcRule{"allow", "subject", "subject", "capability2", []string{"mac_override"}}},
			wantErr: false,
		},
		{
			args:    args{sub: "subject", caps: []string{"invalid"}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateCapabilityAllowRule(tt.args.sub, tt.args.caps)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCapabilityAllowRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateCapabilityAllowRule() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestCreateDomainAutoTransRule(t *testing.T) {
	type args struct {
		sub string
		obj string
		tgt string
	}
	tests := []struct {
		name    string
		args    args
		want    []Rule
		wantErr bool
	}{
		{
			args:    args{"subject", "object", "target"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateDomainAutoTransRule(tt.args.sub, tt.args.obj, tt.args.tgt)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDomainAutoTransRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

