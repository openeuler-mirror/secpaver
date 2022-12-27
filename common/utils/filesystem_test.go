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
	"gotest.tools/v3/fs"
	"path/filepath"
	"reflect"
	"testing"
)

// test function
func TestGetBodyFileName(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		args args
		name string
		want string
	}{
		{
			args: args{"/dir/test.txt"},
			want: "test",
		},
		{
			args: args{"/dir.d/test.txt"},
			want: "test",
		},
		{
			args: args{"/dir.d/test.test.txt"},
			want: "test.test",
		},
		{
			args: args{"/dir.d/.test.test.txt"},
			want: ".test.test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBodyFileName(tt.args.fileName); got != tt.want {
				t.Errorf("GetBodyFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestCopyFile(t *testing.T) {
	testRoot := fs.NewDir(t, "test")
	defer testRoot.Remove()

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			args:    args{"testdata/testfile", filepath.Join(testRoot.Path(), "dst")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyFile(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func TestIsUnixFilePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		args args
		name string
		want bool
	}{
		{
			args: args{"/dir/file"},
			want: true,
		},
		{
			args: args{"/"},
			want: true,
		},
		{
			args: args{"C:\\file"},
			want: false,
		},
		{
			args: args{"file"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnixFilePath(tt.args.path); got != tt.want {
				t.Errorf("IsUnixFilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestPathExist(t *testing.T) {
	testDir := fs.NewDir(t, "test")
	defer testDir.Remove()

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			args:    args{testDir.Path()},
			want:    true,
			wantErr: false,
		},
		{
			args:    args{testDir.Join("invalid")},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PathExist(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("PathExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PathExist() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestDirExist(t *testing.T) {
	testDir := fs.NewDir(t, "test", fs.WithFile("testfile", ""))
	defer testDir.Remove()

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			args:    args{testDir.Path()},
			want:    true,
			wantErr: false,
		},
		{
			args:    args{testDir.Join("invalid")},
			want:    false,
			wantErr: false,
		},
		{
			args:    args{testDir.Join("testfile")},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DirExist(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("DirExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DirExist() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestZipDir(t *testing.T) {
	type args struct {
		dir     string
		zipFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{dir: "testdata/ziptest", zipFile: "testdata/ziptest.zip"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ZipDir(tt.args.dir, tt.args.zipFile); (err != nil) != tt.wantErr {
				t.Errorf("ZipDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// test function
func TestFindAllSubDir(t *testing.T) {
	testDir := fs.NewDir(t, "test", fs.WithDir("dir-1"), fs.WithDir("dir-2"))
	defer testDir.Remove()

	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			args:    args{testDir.Path()},
			want:    []string{testDir.Join("dir-1"), testDir.Join("dir-2")},
			wantErr: false,
		},
		{
			args:    args{testDir.Join("dir-1")},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindAllSubDir(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllSubDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAllSubDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestGetUIDOfFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr bool
	}{
		{
			args:    args{"/"},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUIDOfFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUIDOfFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUIDOfFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestGetModTime(t *testing.T) {
	GetModTime("testdata")
}

