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
	"reflect"
	"testing"
)

var testMetaEmpty = []byte("")

var testMetaValid = []byte(`{
	"name": "testPrj",
	"version": "1.0",
	"resources": "resources.json",
	"specs": ["specs/test.json"],
	"selinux": {
		"config": "selinux.json"
	}
}`)

var testMetaWithoutVersion = []byte(`{
	"name": "testPrj",
	"resources": "resources.json",
	"specs": ["specs/test.json"],
	"selinux": {
		"config": "selinux.json"
	}
}`)

var testMetaWithoutResourceFile = []byte(`{
	"name": "testPrj",
	"version": "1.0",
	"specs": ["specs/test.json"],
	"selinux": {
		"config": "selinux.json"
	}
}`)

var testMetaWithoutSpecFile = []byte(`{
	"name": "testPrj",
	"version": "1.0",
	"resources": "resources.json",
	"selinux": {
		"config": "selinux.json"
	}
}`)

var testMetaUnsafe = []byte(`{
	"name": "testPrj",
	"version": "1.0",
	"resources": "../resources.json",
	"specs": ["specs/test.json"],
	"selinux": {
		"config": "selinux.json"
	}
}`)

var testMetaValidVal = &pb.MetaInfo{
	Name:      "testPrj",
	Version:   "1.0",
	Resources: "resources.json",
	Specs:     []string{"specs/test.json"},
	Selinux: &pb.SelinuxFile{
		Config: "selinux.json",
	},
}

// test function
func Test_parseProjectMeta(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		want    *pb.MetaInfo
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{testMetaEmpty},
			wantErr: true,
		},
		{
			args: args{testMetaWithoutVersion},
			want: testMetaValidVal,
		},
		{
			args: args{testMetaValid},
			want: testMetaValidVal,
		},
		{
			args:    args{testMetaWithoutResourceFile},
			wantErr: true,
		},
		{
			args:    args{testMetaWithoutSpecFile},
			wantErr: true,
		},
		{
			args:    args{testMetaUnsafe},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseProjectMeta(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseProjectMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseProjectMeta() got = %v, want %v", got, tt.want)
			}
		})
	}
}

