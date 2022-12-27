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
	"reflect"
	"testing"
)

var testTypeRule = []*TypeRule{
	{
		Prefix:     "type_transition",
		Subject:    "unconfined_t",
		Object:     "object_t",
		Class:      "file",
		Target:     "target_t",
		ObjectName: "",
	},
	{
		Prefix:     "type_transition",
		Subject:    "unconfined_t",
		Object:     "object_t",
		Class:      "file",
		Target:     "target_t",
		ObjectName: "name",
	},
}

// test function
func TestParseTypeRule(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		want    *TypeRule
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{"type_transition unconfined_t object_t:file target_t;"},
			want:    testTypeRule[0],
			wantErr: false,
		},
		{
			args:    args{"type_transition unconfined_t object_t:file target_t \"name\";"},
			want:    testTypeRule[1],
			wantErr: false,
		},
		{
			args:    args{"type_transition unconfined_t:file target_t;"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTypeRule(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTypeRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseTypeRule() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestTypeRule_Text(t *testing.T) {
	tests := []struct {
		fields *TypeRule
		name   string
		want   string
	}{
		{
			fields: testTypeRule[0],
			want:   "type_transition unconfined_t object_t : file target_t;\n",
		},
		{
			fields: testTypeRule[1],
			want:   "type_transition unconfined_t object_t : file target_t \"name\";\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.fields
			if got := r.Text(); got != tt.want {
				t.Errorf("Text() = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

// test function
func TestTypeRule_CheckConflict(t *testing.T) {
	type args struct {
		rule *TypeRule
	}
	tests := []struct {
		fields  *TypeRule
		name    string
		args    args
		wantErr bool
	}{
		{
			fields: &TypeRule{
				Prefix:     "type_transition",
				Subject:    "subject",
				Object:     "object",
				Class:      "file",
				Target:     "target1",
				ObjectName: "",
			},
			args: args{rule: &TypeRule{
				Prefix:     "type_transition",
				Subject:    "subject",
				Object:     "object",
				Class:      "file",
				Target:     "target2",
				ObjectName: "",
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TypeRule{
				Prefix:     tt.fields.Prefix,
				Subject:    tt.fields.Subject,
				Object:     tt.fields.Object,
				Class:      tt.fields.Class,
				Target:     tt.fields.Target,
				ObjectName: tt.fields.ObjectName,
			}
			if err := r.CheckConflict(tt.args.rule); (err != nil) != tt.wantErr {
				t.Errorf("CheckConflict() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

