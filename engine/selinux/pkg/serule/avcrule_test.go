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

var testAvcRules = []*AvcRule{
	{
		Prefix:  "allow",
		Subject: "unconfined_t",
		Object:  "bin_t",
		Class:   "file",
		Actions: []string{"write", "read"},
	},
	{
		Prefix:  "allow",
		Subject: "unconfined_t",
		Object:  "bin_t",
		Class:   "file",
		Actions: []string{"read"},
	},
}

// test function
func TestParseAvcRule(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		want    *AvcRule
		name    string
		args    args
		wantErr bool
	}{
		{
			args:    args{"allow unconfined_t bin_t:file {write read};"},
			want:    testAvcRules[0],
			wantErr: false,
		},
		{
			args:    args{"allow unconfined_t bin_t:file read;"},
			want:    testAvcRules[1],
			wantErr: false,
		},
		{
			args:    args{"allow bin_t:file read;"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAvcRule(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAvcRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseAvcRule() got = %v, want %v", got, tt.want)
			}
		})
	}
}

