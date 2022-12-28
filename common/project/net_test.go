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
	"reflect"
	"testing"
)

// test function
func TestParseNetwork(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		want    *NetInfo
		args    args
		name    string
		wantErr bool
	}{
		{
			args: args{line: "domain:inet,type:stream,protocol:tcp,port:123"},
			want: &NetInfo{Domain: "inet", Type: "stream", Protocol: "tcp", Port: 123},
		},
		{
			args: args{line: " domain :		inet,type: stream, protocol :tcp , port:  123"},
			want: &NetInfo{Domain: "inet", Type: "stream", Protocol: "tcp", Port: 123},
		},
		{
			args: args{line: "domain  :inet,type:stream, "},
			want: &NetInfo{Domain: "inet", Type: "stream"},
		},
		{
			args: args{line: "domain  :inet,  type  :  "},
			want: &NetInfo{Domain: "inet"},
		},
		{
			args:    args{line: ""},
			wantErr: true,
		},
		{
			args:    args{line: "domain:inet,type:stream,protocol:tcp,port:aaa"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNetwork(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNetwork() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseNetwork() got = %v, want %v", got, tt.want)
			}
		})
	}
}

