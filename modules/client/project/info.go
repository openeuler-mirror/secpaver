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
	"secpaver/api/proto/project"
	"secpaver/common/client"
	"secpaver/common/utils"
)

func newInfoCommand() *cli.Command {
	return &cli.Command{
		Name:        "info",
		Usage:       "Show project information",
		ArgsUsage:   "<PROJECT>",
		Description: "Show detailed information of a project",
		Action:      infoAction,
	}
}

func infoCheck(ctx *cli.Context) error {
	return utils.CheckCommandArgs(ctx, 1, utils.CheckExactArgs)
}

func infoAction(ctx *cli.Context) error {
	if err := infoCheck(ctx); err != nil {
		return err
	}

	c, err := client.NewClientFromContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	svc := project.NewProjectMgrClient(c.Connection())

	req := &project.Req{
		Name: ctx.Args().First(),
	}

	info, err := svc.InfoProject(context.Background(), req)
	if err != nil {
		return err
	}

	// print info table
	table := make([][]string, 0, 10)
	table = append(table, []string{"name", info.GetMetaInfo().GetName()})
	table = append(table, []string{"resource file", info.GetMetaInfo().GetResources()})

	strSpecs := utils.ShowStringsWithSpace(info.GetMetaInfo().GetSpecs())
	table = append(table, []string{"spec files", strSpecs})

	headers := []string{"Attribute", "Value"}
	utils.PrintTabulate(table, headers)

	return nil
}
