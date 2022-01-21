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

package engine

import (
	"context"
	"github.com/urfave/cli"
	"gitee.com/openeuler/secpaver/api/proto/engine"
	"gitee.com/openeuler/secpaver/common/client"
	"gitee.com/openeuler/secpaver/common/utils"
)

func newInfoCommand() *cli.Command {
	return &cli.Command{
		Name:        "info",
		Usage:       "Show engine information",
		ArgsUsage:   "<ENGINE>",
		Description: "Show detailed information of an engine",
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

	svc := engine.NewEngineMgrClient(c.Connection())

	req := &engine.Req{
		Name: ctx.Args().First(),
	}

	info, err := svc.InfoEngine(context.Background(), req)
	if err != nil {
		return err
	}

	table := make([][]string, 0, 10)
	table = append(table, []string{"Name", info.GetBaseInfo().GetName()})
	table = append(table, []string{"Description", info.GetBaseInfo().GetDesc()})

	headers := []string{"Attribute", "Value"}
	utils.PrintTabulate(table, headers)

	return nil
}
