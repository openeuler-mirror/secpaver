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

package project

import (
	pb "gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/utils"
	"gotest.tools/v3/fs"
	"path/filepath"
	"testing"
)

func createProjectValid(path string) {
	_ = utils.CopyFile("testdata/pav.proj", filepath.Join(path, "pav.proj"))
	_ = utils.CopyFile("testdata/resources.json", filepath.Join(path, "resources.json"))
	_ = utils.CopyFile("testdata/selinux.json", filepath.Join(path, "selinux.json"))
	_ = utils.CopyFile("testdata/test.json", filepath.Join(path, "test.json"))
}

func createProjectWithoutMeta(path string) {
	_ = utils.CopyFile("testdata/resources.json", filepath.Join(path, "resources.json"))
	_ = utils.CopyFile("testdata/selinux.json", filepath.Join(path, "selinux.json"))
	_ = utils.CopyFile("testdata/test.json", filepath.Join(path, "test.json"))
}

func createProjectWithoutResource(path string) {
	_ = utils.CopyFile("testdata/pav.proj", filepath.Join(path, "pav.proj"))
	_ = utils.CopyFile("testdata/selinux.json", filepath.Join(path, "selinux.json"))
	_ = utils.CopyFile("testdata/test.json", filepath.Join(path, "test.json"))
}

func createProjectWithoutSpec(path string) {
	_ = utils.CopyFile("testdata/pav.proj", filepath.Join(path, "pav.proj"))
	_ = utils.CopyFile("testdata/resources.json", filepath.Join(path, "resources.json"))
	_ = utils.CopyFile("testdata/selinux.json", filepath.Join(path, "selinux.json"))
}

// test function
func TestParseProjectFromDir(t *testing.T) {
	prjValid := fs.NewDir(t, "test")
	defer prjValid.Remove()
	createProjectValid(prjValid.Path())

	prjWithoutMeta := fs.NewDir(t, "test")
	defer prjWithoutMeta.Remove()
	createProjectWithoutMeta(prjValid.Path())

	prjWithoutResource := fs.NewDir(t, "test")
	defer prjWithoutResource.Remove()
	createProjectWithoutResource(prjWithoutResource.Path())

	prjWithoutSpec := fs.NewDir(t, "test")
	defer prjWithoutSpec.Remove()
	createProjectWithoutSpec(prjValid.Path())

	type args struct {
		dir string
	}
	tests := []struct {
		want    *pb.ProjectInfo
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{prjValid.Path()},
		},
		{
			args:    args{prjWithoutMeta.Path()},
			wantErr: true,
		},
		{
			args:    args{prjWithoutResource.Path()},
			wantErr: true,
		},
		{
			args:    args{prjWithoutSpec.Path()},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseProjectFromDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProjectFromDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

