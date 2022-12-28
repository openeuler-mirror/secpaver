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
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"testing"
)

// test function
func Test_autoGenFileTypeByPathAndClass(t *testing.T) {
	type args struct {
		path   string
		class  fClass
		isExec bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{"/test", secontext.ComFile, false},
			want: "auto_test_b306_file_t",
		},
		{
			args: args{"/test/a*", secontext.ComFile, false},
			want: "auto_test_a_6f3e_file_t",
		},
		{
			args: args{"/test/a*/b*", secontext.ComFile, false},
			want: "auto_test_a_b_f8d4_file_t",
		},
		{
			args: args{"/test/*", secontext.ComFile, false},
			want: "auto_test_38f0_file_t",
		},
		{
			args: args{"/test/**", secontext.ComFile, false},
			want: "auto_test_e935_file_t",
		},
		{
			args: args{"/test{,/*}", secontext.ComFile, false},
			want: "auto_test_3e38_file_t",
		},
		{
			args: args{"/a*", secontext.ComFile, false},
			want: "auto_root_a_8d61_file_t",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := autoGenFileTypeByPathAndClass(tt.args.path, tt.args.class, tt.args.isExec); got != tt.want {
				t.Errorf("autoGenFileTypeByPathAndClass() = %v, want %v", got, tt.want)
			}
		})
	}
}

