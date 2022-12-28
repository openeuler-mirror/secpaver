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
	"gitee.com/openeuler/secpaver/common/project"
	"reflect"
	"testing"
)

// test function
func Test_getSocketClasses(t *testing.T) {
	type args struct {
		info *project.NetInfo
	}
	tests := []struct {
		args args
		name string
		want []string
	}{
		{args: args{info: &project.NetInfo{Domain: "inet", Type: "stream", Protocol: "tcp"}},
			want: []string{"tcp_socket"}},
		{args: args{info: &project.NetInfo{Domain: "inet", Type: "stream"}},
			want: []string{"tcp_socket", "sctp_socket"}},
		{args: args{info: &project.NetInfo{Domain: "inet", Type: "dgram"}},
			want: []string{"udp_socket"}},
		{args: args{info: &project.NetInfo{Domain: "inet"}},
			want: inetSocketSet},
		{args: args{info: &project.NetInfo{Domain: "netlink", Protocol: "route"}},
			want: []string{"netlink_route_socket"}},
		{args: args{info: &project.NetInfo{Domain: "netlink"}},
			want: netlinkSocketSet},
		{args: args{info: &project.NetInfo{Domain: "unix", Type: "stream"}},
			want: []string{"unix_stream_socket"}},
		{args: args{info: &project.NetInfo{Domain: "unix"}},
			want: unixSocketSet},
		{args: args{info: &project.NetInfo{Type: "stream"}},
			want: []string{"tcp_socket", "sctp_socket", "unix_stream_socket"}},
		{args: args{info: &project.NetInfo{Domain: "ax25"}},
			want: []string{"ax25_socket"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSocketClasses(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSocketClasses() = %v, want %v", got, tt.want)
			}
		})
	}
}

