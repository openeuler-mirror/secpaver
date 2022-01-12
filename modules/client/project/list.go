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
	"context"
	"github.com/urfave/cli"
	"gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/client"
	"gitee.com/openeuler/secpaver/common/utils"
)

func newListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List projects",
		ArgsUsage:   "",
		Description: "List projects in secPaver server\n",
		Action:      listAction,
	}
}

func listCheck(ctx *cli.Context) error {
	// list command should not have any args
	return utils.CheckCommandArgs(ctx, 0, utils.CheckExactArgs)
}

func listAction(ctx *cli.Context) error {
	if err := listCheck(ctx); err != nil {
		return err
	}

	c, err := client.NewClientFromContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	svc := project.NewProjectMgrClient(c.Connection())
	infos, err := svc.ListProject(context.Background(), &project.Req{})
	if err != nil {
		return err
	}

	table := make([][]string, 0, 10)

	for _, info := range infos.ProjectInfos {
		row := make([]string, 0, 2)
		row = append(row, info.GetName())
		row = append(row, info.GetVersion())
		table = append(table, row)
	}

	utils.PrintTabulate(table, []string{"Name", "Version"})

	return nil
}
