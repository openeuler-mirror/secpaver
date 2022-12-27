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

package config

import (
	"gitee.com/openeuler/secpaver/common/global"
	"reflect"
	"testing"
)

var defConfigInfo = GlobalConfig{
	Connect: ConnectInfo{
		Grpc: Grpc{
			Socket: global.DefaultGrpcSocket,
		},
	},
	Repository: RepositoryInfo{
		ProjectRoot: global.DefaultProjectRoot,
		PolicyRoot:  global.DefaultPolicyRoot,
	},
	Log: LogInfo{
		FilePath:    global.DefaultLogPath,
		Level:       global.DefaultLogLevel,
		MaxFileAge:  global.DefaultLogFileAge,
		MaxFileSize: global.DefaultLogFileSize,
		MaxFileNum:  global.DefaultLogFileNum,
	},
}

// test function
func TestParseConfig(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		want    *GlobalConfig
		args    args
		name    string
		wantErr bool
	}{
		{
			args:    args{"testdata/file1.json"},
			want:    &defConfigInfo,
			wantErr: false,
		},
		{
			args:    args{"testdata/file2.json"},
			want:    &defConfigInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfig(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

