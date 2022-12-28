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
	pb "gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/project"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/serule"
	"reflect"
	"testing"
)

var testSubject = &applicationItem{
	domain: "test_subject_t",
}

var testCapabilityRules = []serule.Rule{
	&serule.AvcRule{
		Prefix:  "allow",
		Subject: "test_subject_t",
		Object:  "test_subject_t",
		Class:   "capability",
		Actions: []string{
			"net_bind_service",
			"audit_control",
		},
	},
	&serule.AvcRule{
		Prefix:  "allow",
		Subject: "test_subject_t",
		Object:  "test_subject_t",
		Class:   "capability2",
		Actions: []string{
			"wake_alarm",
		},
	},
	&serule.AvcRule{
		Prefix:  "allow",
		Subject: "test_subject_t",
		Object:  "test_subject_t",
		Class:   "process",
		Actions: []string{
			"setcap",
			"getcap",
		},
	},
}

// test function
func Test_getCapabilityByActions(t *testing.T) {
	type args struct {
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
				[]string{
					project.CapWakeAlarm,
					project.CapNetBindService,
					project.CapAuditControl,
				},
			},
			want: []string{
				"wake_alarm",
				"net_bind_service",
				"audit_control",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCapabilityByActions(tt.args.actions)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCapPermsByActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCapPermsByActions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestBuilder_genCapabilityRules(t *testing.T) {
	type args struct {
		subject *applicationItem
		perm    *pb.Permission
	}
	tests := []struct {
		want    []serule.Rule
		args    args
		name    string
		wantErr bool
	}{
		{
			args: args{
				subject: testSubject,
				perm: &pb.Permission{
					Type: project.RuleCapability,
					Actions: []string{
						project.CapWakeAlarm,
						project.CapNetBindService,
						project.CapAuditControl,
					},
				},
			},
			want:    testCapabilityRules,
			wantErr: false,
		},
		{
			args: args{
				subject: testSubject,
				perm: &pb.Permission{
					Type:    project.RuleCapability,
					Actions: []string{"invalid"},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := (&Builder{}).getCapabilityRules(tt.args.subject, tt.args.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("genCapabilityRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genCapabilityRules() got = %v, want %v", got, tt.want)
			}
		})
	}
}

