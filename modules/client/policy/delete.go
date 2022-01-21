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
	"gitee.com/openeuler/secpaver/api/proto/policy"
	"gitee.com/openeuler/secpaver/common/client"
	"gitee.com/openeuler/secpaver/common/utils"
)

func newDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Usage:       "Delete a policy",
		ArgsUsage:   "<POLICY>",
		Description: "Delete a policy in secPaver server",
		Action:      deleteAction,
	}
}

func deleteCheck(ctx *cli.Context) error {
	return utils.CheckCommandArgs(ctx, 1, utils.CheckExactArgs)
}

func deleteAction(ctx *cli.Context) error {
	if err := deleteCheck(ctx); err != nil {
		return err
	}

	c, err := client.NewClientFromContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	svc := policy.NewPolicyMgrClient(c.Connection())

	name, engine := utils.SplitByLastUnderscore(ctx.Args().First())
	req := &policy.Req{
		Name:   name,
		Engine: engine,
	}

	ack, err := svc.DeletePolicy(context.Background(), req)
	if err != nil {
		return err
	}

	utils.PrintAck(ack.GetLevel(), ack.GetStatus())

	return nil
}
