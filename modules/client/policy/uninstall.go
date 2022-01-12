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
	"io"
	"gitee.com/openeuler/secpaver/api/proto/policy"
	"gitee.com/openeuler/secpaver/common/client"
	"gitee.com/openeuler/secpaver/common/utils"
)

func newUninstallCommand() *cli.Command {
	return &cli.Command{
		Name:        "uninstall",
		Usage:       "Uninstall a policy",
		ArgsUsage:   "<POLICY>",
		Description: "Uninstall a policy in secPaver server",
		Action:      uninstallAction,
	}
}

func uninstallCheck(ctx *cli.Context) error {
	return utils.CheckCommandArgs(ctx, 1, utils.CheckExactArgs)
}

func uninstallAction(ctx *cli.Context) error {
	if err := uninstallCheck(ctx); err != nil {
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

	stream, err := svc.UninstallPolicy(context.Background(), req)
	if err != nil {
		return err
	}

	for {
		ack, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		utils.PrintAck(ack.GetLevel(), ack.GetStatus())
	}

	return nil
}
