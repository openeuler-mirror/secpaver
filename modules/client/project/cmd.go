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
Package project implements the 'project' subcommand of Pav
*/
package project

import (
	"github.com/urfave/cli"
)

// NewProjectCmd creates a command for `project` subcommand
func NewProjectCmd() *cli.Command {
	return &cli.Command{
		Name:  "project",
		Usage: "Manage projects",
		Subcommands: []cli.Command{
			*newListCommand(),
			*newImportCommand(),
			*newDeleteCommand(),
			*newExportCommand(),
			*newInfoCommand(),
			*newBuildCommand(),
			*newCreateCommand(),
		},
	}
}
