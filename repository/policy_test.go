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

package repository

import (
	"gitee.com/openeuler/secpaver/domain"
	"gotest.tools/v3/fs"
	"reflect"
	"testing"
)

// test function
func Test_repo_FindAllPolicies(t *testing.T) {
	polRoot := fs.NewDir(t, "policies", fs.WithDir("selinux", fs.WithDir("pol")))
	defer polRoot.Remove()

	type fields struct {
		projectRoot string
		policyRoot  string
	}
	tests := []struct {
		want    []*domain.Policy
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{policyRoot: polRoot.Path()},
			want: []*domain.Policy{
				domain.NewPolicy(polRoot.Join("selinux", "pol")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repo{
				projectRoot: tt.fields.projectRoot,
				policyRoot:  tt.fields.policyRoot,
			}
			got, err := r.FindAllPolicies()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllPolicies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAllPolicies() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_repo_FindPolicyByName(t *testing.T) {
	polRoot := fs.NewDir(t, "policies", fs.WithDir("selinux", fs.WithDir("pol")))
	defer polRoot.Remove()

	type fields struct {
		projectRoot string
		policyRoot  string
	}
	type args struct {
		name   string
		engine string
	}
	tests := []struct {
		want    *domain.Policy
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields:  fields{policyRoot: polRoot.Path()},
			args:    args{name: "pol", engine: "selinux"},
			want:    domain.NewPolicy(polRoot.Join("selinux", "pol")),
			wantErr: false,
		},
		{
			fields:  fields{policyRoot: polRoot.Path()},
			args:    args{name: "pol", engine: "apparmor"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repo{
				projectRoot: tt.fields.projectRoot,
				policyRoot:  tt.fields.policyRoot,
			}
			got, err := r.FindPolicyByName(tt.args.name, tt.args.engine)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindPolicyByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPolicyByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

