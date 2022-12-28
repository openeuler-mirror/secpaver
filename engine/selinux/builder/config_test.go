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
	"testing"
)

var confValid = `{"policy":{"monolithic":true},
    "extraRules":["allow test_t test_t file: write;"]}`
var confWithInvalidFormat = `{"policy":{"monolithic":true},
    "extraRules":["allow test_t test_t file: write;"}`
var confEmpty = `{}`

// test function
func Test_parseProjectSelinuxConfigFile(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		want    *config
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{[]byte(confValid)},
			want: &config{
				Policy:     &policyOpt{Monolithic: true},
				ExtraRules: []string{"allow test_t test_t file: write;"},
			},
		},
		{
			args:    args{[]byte(confWithInvalidFormat)},
			wantErr: true,
		},
		{
			args: args{[]byte(confEmpty)},
			want: &config{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseProjectSelinuxConfigFile(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseProjectSelinuxConfigFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseProjectSelinuxConfigFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

