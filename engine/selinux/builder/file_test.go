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

package builder

import (
	"reflect"
	pb "gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/project"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/serule"
	"sort"
	"testing"
)

var testLaunchersContext = []secontext.Context{
	{
		User:        "unconfined_u",
		Role:        "unconfined_r",
		Type:        "unconfined_t",
		Sensitivity: "s0",
	},
}

var testBuilder = &Builder{
	pcHandle: projectContextHandle{
		fileItems: testFileItems,
	},
}

var testApplicationItem = &applicationItem{
	file:   testFileItems[1],
	domain: "test_subject_t",
}

var testFileItems = []*fileItem{
	{
		path:           "/test/file",
		sePath:         "/test/file",
		class:          secontext.ComFile,
		context:        *secontext.CreateDefaultObjectContext("test_file_t"),
		contextInherit: *secontext.CreateDefaultObjectContext("test_dir_t"),
		execDomain:     "test_t",
	},
	{
		path:           "/test/exec",
		sePath:         "/test/exec",
		class:          secontext.ComFile,
		context:        *secontext.CreateDefaultObjectContext("test_exec_t"),
		contextInherit: *secontext.CreateDefaultObjectContext("test_dir_t"),
		execDomain:     "test_subject_t",
	},
}

var testFileExecRules = []serule.Rule{
	&serule.AvcRule{
		Prefix:  "allow",
		Subject: "test_subject_t",
		Object:  "test_file_t",
		Class:   "file",
		Actions: comFileExecPermsSet,
	},
}

var testFileCreateRules = []serule.Rule{
	&serule.AvcRule{
		Prefix:  "allow",
		Subject: "test_subject_t",
		Object:  "test_file_t",
		Class:   "file",
		Actions: baseFileCreatePermsSet,
	},
	&serule.AvcRule{
		Prefix:  "allow",
		Subject: "test_subject_t",
		Object:  "test_dir_t",
		Class:   "dir",
		Actions: dirFileWritePermsSet,
	},
	&serule.TypeRule{
		Prefix:     "type_transition",
		Subject:    "test_subject_t",
		Object:     "test_dir_t",
		Class:      "file",
		Target:     "test_file_t",
		ObjectName: "file",
	},
}

var testFileRemoveRules = []serule.Rule{
	&serule.AvcRule{
		Prefix:  "allow",
		Subject: "test_subject_t",
		Object:  "test_file_t",
		Class:   "file",
		Actions: baseFileRemovePermsSet,
	},
	&serule.AvcRule{
		Prefix:  "allow",
		Subject: "test_subject_t",
		Object:  "test_dir_t",
		Class:   "dir",
		Actions: dirFileWritePermsSet,
	},
}

var testFileExecPermission = &pb.Permission{
	Type: project.RuleFileSystem,
	Resources: []string{
		testFileItems[0].path,
	},
	Actions: []string{
		project.ActionFileExec,
	},
}

var testFileCreatePermission = &pb.Permission{
	Type: project.RuleFileSystem,
	Resources: []string{
		testFileItems[0].path,
	},
	Actions: []string{
		project.ActionFileCreate,
	},
}

var testFileRemovePermission = &pb.Permission{
	Type: project.RuleFileSystem,
	Resources: []string{
		testFileItems[0].path,
	},
	Actions: []string{
		project.ActionFileRemove,
	},
}

// test function
func Test_getFileClassByName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		args    args
		name    string
		want    fClass
		wantErr bool
	}{
		{
			args:    args{project.ComFile},
			want:    secontext.ComFile,
			wantErr: false,
		},
		{
			args:    args{"invalid"},
			want:    secontext.UnknownFile,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFileClassByName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileClassByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFileClassByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func Test_getFilePermsByActions(t *testing.T) {
	type args struct {
		classID fClass
		actions []string
	}
	tests := []struct {
		args    args
		name    string
		want    []string
		wantErr bool
	}{
		{
			args: args{
				classID: secontext.ComFile,
				actions: project.AllFileActions,
			},
			want: []string{
				"write", "create", "append", "setattr", "rename",
				"link", "getattr", "unlink", "lock", "map", "open",
				"read", "execute", "execute_no_trans", "ioctl", "mounton",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		sort.Strings(tt.want)
		t.Run(tt.name, func(t *testing.T) {
			got, err := getFilePermsByActions(tt.args.classID, tt.args.actions)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFilePermsByActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFilePermsByActions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestBuilder_getFileRules(t *testing.T) {
	rules, _ := serule.CreateDomainAutoTransRule("test_subject_t", "test_file_t", "test_t")
	testFileExecRules = append(testFileExecRules, rules...)

	type args struct {
		subject *applicationItem
		perm    *pb.Permission
	}
	tests := []struct {
		name    string
		args    args
		want    []serule.Rule
		wantErr bool
	}{
		{
			args: args{
				subject: testApplicationItem,
				perm:    testFileExecPermission,
			},
			want:    testFileExecRules,
			wantErr: false,
		},
		{
			args: args{
				subject: testApplicationItem,
				perm:    testFileCreatePermission,
			},
			want:    testFileCreateRules,
			wantErr: false,
		},
		{
			args: args{
				subject: testApplicationItem,
				perm:    testFileRemovePermission,
			},
			want:    testFileRemoveRules,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testBuilder.getFileRules(tt.args.subject, tt.args.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFileRules() got = %v, want %v", got, tt.want)
			}
		})
	}
}

