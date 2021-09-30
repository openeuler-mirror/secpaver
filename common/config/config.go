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
Package config implements the parsing of config file.
*/
package config

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io/ioutil"
	"secpaver/common/global"
	"secpaver/common/utils"
)

// GlobalConfig is the global config information parsed by config file
type GlobalConfig struct {
	Connect    ConnectInfo    `json:"connect"`
	Repository RepositoryInfo `json:"repository"`
	Log        LogInfo        `json:"log"`
}

// ConnectInfo is the net connection config information parsed by config file
type ConnectInfo struct {
	Grpc Grpc `json:"grpc"`
}

// Grpc is grpc config
type Grpc struct {
	Socket string `json:"socket"`
}

// RepositoryInfo is the repository paths parsed by config file
type RepositoryInfo struct {
	ProjectRoot string `json:"projects"`
	PolicyRoot  string `json:"policies"`
}

// LogInfo is the log information parsed by config file
type LogInfo struct {
	FilePath    string `json:"path"`
	Level       string `json:"level"`
	MaxFileSize int    `json:"maxFileSize"`
	MaxFileNum  int    `json:"maxFileNum"`
	MaxFileAge  int    `json:"maxFileAge"`
}

// ParseConfig parse the Pavd global config file
func ParseConfig(file string) (*GlobalConfig, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("fail to read file")
	}

	config := &GlobalConfig{}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return config, errors.Wrap(err, "fail to unmarshal file")
	}

	// set the default parameters
	setIfStringParamIsEmpty(&config.Connect.Grpc.Socket, global.DefaultGrpcSocket)
	setIfStringParamIsEmpty(&config.Repository.ProjectRoot, global.DefaultProjectRoot)
	setIfStringParamIsEmpty(&config.Repository.PolicyRoot, global.DefaultPolicyRoot)
	setIfStringParamIsEmpty(&config.Log.FilePath, global.DefaultLogPath)
	setIfStringParamIsEmpty(&config.Log.Level, global.DefaultLogLevel)
	setIfIntParamIsEmpty(&config.Log.MaxFileAge, global.DefaultLogFileAge)
	setIfIntParamIsEmpty(&config.Log.MaxFileNum, global.DefaultLogFileNum)
	setIfIntParamIsEmpty(&config.Log.MaxFileSize, global.DefaultLogFileSize)

	// check all path safety
	allPath := []string{
		config.Repository.ProjectRoot,
		config.Repository.PolicyRoot,
	}

	if config.Log.FilePath != "" {
		allPath = append(allPath, config.Log.FilePath)
	}

	for _, path := range allPath {
		if err := utils.CheckUnsafePath(path); err != nil {
			return nil, errors.Wrap(err, "fail to check paths in config file")
		}
	}

	return config, nil
}

// MergeConfigInfo returns the global config info by merging cli and config file
func MergeConfigInfo(info *GlobalConfig, ctx *cli.Context) *GlobalConfig {
	if ctx == nil || info == nil {
		return info
	}

	setIfStringValueIsNotEmpty(&info.Log.Level, ctx.String("log-level"))
	setIfStringValueIsNotEmpty(&info.Connect.Grpc.Socket, ctx.String("socket"))

	return info
}

func setIfStringParamIsEmpty(param *string, value string) {
	if param != nil && *param == "" {
		*param = value
	}
}

func setIfIntParamIsEmpty(param *int, value int) {
	if param != nil && *param == 0 {
		*param = value
	}
}

func setIfStringValueIsNotEmpty(param *string, value string) {
	if param != nil && value != "" {
		*param = value
	}
}
