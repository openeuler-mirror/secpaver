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

var defUserContext = &Context{
	User:        "unconfined_u",
	Role:        "unconfined_r",
	Type:        "unconfined_t",
	Sensitivity: "s0",
	Categories:  "",
}

// test function
func TestGetExecContext(t *testing.T) {
	type args struct {
		scon *Context
		tcon *Context
	}
	tests := []struct {
		want    *Context
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{
				scon: defUserContext,
				tcon: CreateDefaultObjectContext("shell_exec_t"),
			},
			want:    defUserContext,
			wantErr: false,
		},
		{
			args: args{
				scon: defUserContext,
				tcon: CreateDefaultObjectContext("invalid_t"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetExecContext(tt.args.scon, tt.args.tcon)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetExecContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetExecContext() got = %v, want %v", got, tt.want)
			}
		})
	}
}

