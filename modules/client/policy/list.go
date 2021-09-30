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

package policy

import (
	"context"
	"github.com/urfave/cli"
	"secpaver/api/proto/policy"
	"secpaver/common/client"
	"secpaver/common/utils"
)

func newListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "List policies",
		ArgsUsage:   "",
		Description: "List policies in secPaver server\n",
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

	svc := policy.NewPolicyMgrClient(c.Connection())
	ack, err := svc.ListPolicy(context.Background(), &policy.Req{})
	if err != nil {
		return err
	}

	table := make([][]string, 0, 10)

	for _, info := range ack.PolicyInfos {
		row := make([]string, 0, 2)
		row = append(row, info.GetName())
		row = append(row, info.GetStatus())
		table = append(table, row)
	}

	utils.PrintTabulate(table, []string{"Name", "Status"})
	return nil
}
