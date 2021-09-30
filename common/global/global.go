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
Package global defines global variables of the application
*/
package global

// common config
const (
	DefaultGrpcSocket = "/var/run/secpaver/pavd.sock"

	DefaultDirPerm  = 0700
	DefaultFilePerm = 0600
	SocketFilePerm  = 0600
	LogFilePerm     = 0600
	PidFilePerm     = 0600
)

// global config of pav application
const (
	PavVersion     = "1.0.2"
	DefaultTimeout = 120
)

// global config of pavd application
const (
	PavdVerison = "1.0.2"

	// default file path
	DefaultPidFile     = "/var/run/pavd.pid"
	DefaultConfigFile  = "/etc/secpaver/pavd/config.json"
	DefaultProjectRoot = "/var/local/secpaver/projects"
	DefaultPolicyRoot  = "/var/local/secpaver/policies"
	DefaultPluginRoot  = "/usr/lib64/secpaver/"
	DefaultScriptRoot  = "/usr/share/secpaver/scripts"

	// default log config
	DefaultLogPath     = "/var/log/secpaver/pavd.log"
	DefaultLogLevel    = "info"
	DefaultLogFileSize = 10 // maximum size in megabytes of a log file
	DefaultLogFileAge  = 30 // maximum number of days to retain old log files
	DefaultLogFileNum  = 20 // maximum number of old log files to retain

	MaxLogFileSize = 100
	MaxLogFileAge  = 3650
	MaxLogFileNum  = 100
)
