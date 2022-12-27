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

package main

import (
	"gotest.tools/v3/fs"
	"net"
	"testing"
)

// test function
func Test_listenSocket(t *testing.T) {
	dir := fs.NewDir(t, "test")
	defer dir.Remove()

	type args struct {
		sock string
	}
	tests := []struct {
		name    string
		args    args
		want    net.Listener
		wantErr bool
	}{
		{
			args:    args{dir.Join("test.sock")},
			wantErr: false,
		},
		{
			args:    args{"invalid"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := listenSocket(tt.args.sock)
			if (err != nil) != tt.wantErr {
				t.Errorf("listenSocket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				_ = got.Close()
			}
		})
	}
}

