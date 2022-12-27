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

package log

import (
	"gitee.com/openeuler/secpaver/common/config"
	"gitee.com/openeuler/secpaver/common/global"
	"gotest.tools/v3/fs"
	"runtime"
	"testing"
)

// test function
func TestInitLogger(t *testing.T) {
	logDir := fs.NewDir(t, "test")
	defer logDir.Remove()

	type args struct {
		info *config.LogInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			args: args{info: &config.LogInfo{
				FilePath:    "",
				Level:       "debug",
				MaxFileSize: global.DefaultLogFileSize,
				MaxFileNum:  global.DefaultLogFileNum,
				MaxFileAge:  global.DefaultLogFileAge,
			}},
			wantErr: false,
		},
		{
			args: args{info: &config.LogInfo{
				FilePath:    logDir.Join("test.log"),
				Level:       "debug",
				MaxFileSize: global.DefaultLogFileSize,
				MaxFileNum:  global.DefaultLogFileNum,
				MaxFileAge:  global.DefaultLogFileAge,
			}},
			wantErr: false,
		},
		{
			args: args{info: &config.LogInfo{
				FilePath:    logDir.Join("test.log"),
				Level:       "invalid",
				MaxFileSize: global.DefaultLogFileSize,
				MaxFileNum:  global.DefaultLogFileNum,
				MaxFileAge:  global.DefaultLogFileAge,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitLogger(tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("InitLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

		Close()
	}
}

// test function
func TestLogPrint(t *testing.T) {
	logDir := fs.NewDir(t, "test")
	defer logDir.Remove()

	if err := InitLogger(&config.LogInfo{
		FilePath:    logDir.Join("test.log"),
		Level:       "debug",
		MaxFileSize: global.DefaultLogFileSize,
		MaxFileNum:  global.DefaultLogFileNum,
		MaxFileAge:  global.DefaultLogFileAge,
	}); err != nil {

		t.Error(err)
	}

	Debug(funcName())
	Debugf(funcName())
	Debugln(funcName())
	Info(funcName())
	Infof(funcName())
	Infoln(funcName())
	Warn(funcName())
	Warnf(funcName())
	Warnln(funcName())
	Error(funcName())
	Errorf(funcName())
	Errorln(funcName())
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name()
}

