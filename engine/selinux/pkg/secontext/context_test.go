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

package secontext

import (
	"reflect"
	"testing"
)

// test function
func TestParseContextFromLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		want    *Context
		args    args
		name    string
		wantErr bool
	}{
		{
			args: args{"unconfined_u:unconfined_r:unconfined_t:s0:c0"},
			want: &Context{"unconfined_u", "unconfined_r", "unconfined_t", "s0", "c0"},
		},
		{
			args: args{"unconfined_u  :  unconfined_r  :  unconfined_t:s0:c0"},
			want: &Context{"unconfined_u", "unconfined_r", "unconfined_t", "s0", "c0"},
		},
		{
			args:    args{"unconfined_u:unconfined_r:unconfined_t:s0:c0:test"},
			want:    nil,
			wantErr: true,
		},
		{
			args: args{"unconfined_u:unconfined_r:unconfined_t:s0"},
			want: &Context{"unconfined_u", "unconfined_r", "unconfined_t", "s0", ""},
		},
		{
			args: args{"unconfined_u:unconfined_r:unconfined_t"},
			want: &Context{"unconfined_u", "unconfined_r", "unconfined_t", "", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseContextFromLine(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseContextFromLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseContextFromLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// test function
func TestContext_Text(t *testing.T) {
	type fields struct {
		User        string
		Role        string
		Type        string
		Sensitivity string
		Categories  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{User: "u", Role: "r", Type: "t"},
			want:   "u:r:t",
		},
		{
			fields: fields{User: "u", Role: "r", Type: "t", Sensitivity: "s"},
			want:   "u:r:t:s",
		},
		{
			fields: fields{User: "u", Role: "r", Type: "t", Sensitivity: "s"},
			want:   "u:r:t:s",
		},
		{
			fields: fields{User: "u", Role: "r", Type: "t", Sensitivity: "s", Categories: "c"},
			want:   "u:r:t:s:c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &Context{
				User:        tt.fields.User,
				Role:        tt.fields.Role,
				Type:        tt.fields.Type,
				Sensitivity: tt.fields.Sensitivity,
				Categories:  tt.fields.Categories,
			}
			if got := ctx.Text(); got != tt.want {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

