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
	"gotest.tools/v3/fs"
	"reflect"
	"gitee.com/openeuler/secpaver/domain"
	"testing"
)

// test function
func Test_repo_FindAllProjects(t *testing.T) {
	prjRoot := fs.NewDir(t, "projects", fs.WithDir("prj-1"),
		fs.WithDir("prj-2"), fs.WithDir("prj-3"))
	defer prjRoot.Remove()

	type fields struct {
		projectRoot string
		policyRoot  string
	}
	tests := []struct {
		want    []*domain.Project
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{projectRoot: prjRoot.Path()},
			want: []*domain.Project{
				domain.NewProject(prjRoot.Join("prj-1")),
				domain.NewProject(prjRoot.Join("prj-2")),
				domain.NewProject(prjRoot.Join("prj-3")),
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
			got, err := r.FindAllProjects()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAllProjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_repo_FindProjectByName(t *testing.T) {
	prjRoot := fs.NewDir(t, "projects", fs.WithDir("prj-1"),
		fs.WithDir("prj-2"), fs.WithDir("prj-3"))
	defer prjRoot.Remove()

	type fields struct {
		projectRoot string
		policyRoot  string
	}
	type args struct {
		prjName string
	}
	tests := []struct {
		want    *domain.Project
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields:  fields{projectRoot: prjRoot.Path()},
			args:    args{prjName: "prj-1"},
			want:    domain.NewProject(prjRoot.Join("prj-1")),
			wantErr: false,
		},
		{
			fields:  fields{projectRoot: prjRoot.Path()},
			args:    args{prjName: "prj-invalid"},
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
			got, err := r.FindProjectByName(tt.args.prjName)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindProjectByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindProjectByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

