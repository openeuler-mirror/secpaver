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

/*
Package log implements log function based on logrus
*/
package log

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"path/filepath"
	"gitee.com/openeuler/secpaver/common/config"
	"gitee.com/openeuler/secpaver/common/global"
)

var rotOpt *lumberjack.Logger

// InitLogger initialize the logger with the work mode settings
func InitLogger(info *config.LogInfo) error {
	if info == nil {
		return fmt.Errorf("nil log info parameter")
	}

	lev, err := logrus.ParseLevel(info.Level)
	if err != nil {
		return fmt.Errorf("invalid log level %s", info.Level)
	}

	logrus.SetLevel(lev)

	if info.FilePath == "" {
		logrus.SetOutput(os.Stdout)
		return nil
	}

	rotOpt = &lumberjack.Logger{
		Filename:   info.FilePath,
		MaxSize:    info.MaxFileSize,
		MaxBackups: info.MaxFileNum,
		MaxAge:     info.MaxFileAge,
	}

	if info.MaxFileSize > global.MaxLogFileSize {
		info.MaxFileSize = global.MaxLogFileSize
	}

	if info.MaxFileNum > global.MaxLogFileNum {
		info.MaxFileNum = global.MaxLogFileNum
	}

	if info.MaxFileAge > global.MaxLogFileAge {
		info.MaxFileAge = global.MaxLogFileAge
	}

	if err := os.MkdirAll(path.Dir(info.FilePath), global.DefaultDirPerm); err != nil {
		return errors.Wrapf(err, "fail to create directory for log file")
	}

	// write for test
	if _, err := rotOpt.Write([]byte("")); err != nil {
		return errors.Wrapf(err, "fail to write log file %s", filepath.Base(info.FilePath))
	}

	if err := os.Chmod(info.FilePath, global.LogFilePerm); err != nil {
		return errors.Wrapf(err, "fail to chmod log file %s", filepath.Base(info.FilePath))
	}

	writers := []io.Writer{
		os.Stdout,
		rotOpt,
	}

	logrus.SetOutput(io.MultiWriter(writers...))

	return nil
}

// Close closes the opened log file
func Close() {
	if rotOpt != nil {
		_ = rotOpt.Close()
		rotOpt = nil
	}
}

// Debug log debug level message
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Debugln log debug level message with newline
func Debugln(args ...interface{}) {
	logrus.Debugln(args...)
}

// Debugf log debug level message with format string
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Info log info level message
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Infoln log info level message with newline
func Infoln(args ...interface{}) {
	logrus.Infoln(args...)
}

// Infof log info level message with format string
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Warn log warn level message
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warnln log warn level message with newline
func Warnln(args ...interface{}) {
	logrus.Warnln(args...)
}

// Warnf log warn level message with format string
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Error log error level message
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorln log error level message with newline
func Errorln(args ...interface{}) {
	logrus.Errorln(args...)
}

// Errorf log error level message with format string
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}
