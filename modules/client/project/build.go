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
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io"
	pb "secpaver/api/proto/project"
	"secpaver/common/client"
	"secpaver/common/errdefs"
	"secpaver/common/project"
	"secpaver/common/utils"
)

func newBuildCommand() *cli.Command {
	return &cli.Command{
		Name:        "build",
		Usage:       "Generate builder files",
		UsageText:   "pav project build --engine <ENGINE> [-r <PROJECT> | -d <PATH>]",
		Description: "Build project",
		Action:      buildAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "engine",
				Usage: `Select the engine`,
				Value: "selinux",
			},
			&cli.StringFlag{
				Name:  "r, remote",
				Usage: `Build project in the server repository`,
			},
			&cli.StringFlag{
				Name:  "d, dir",
				Usage: `Build local project`,
			},
		},
	}
}

func buildCheck(ctx *cli.Context) error {
	if err := utils.CheckCommandArgs(ctx, 0, utils.CheckExactArgs); err != nil {
		return err
	}

	if (ctx.String("remote") != "") &&
		(ctx.String("dir") != "") {

		return fmt.Errorf("the 'remote' and 'dir' parameters can be set both")
	}

	if (ctx.String("remote") == "") &&
		(ctx.String("dir") == "") {

		return fmt.Errorf("please set 'remote' or 'dir' parameter")
	}

	if ctx.String("dir") == "" {
		return nil
	}

	dirExists, err := utils.DirExist(ctx.String("dir"))
	if err != nil {
		return err
	}
	if !dirExists {
		return errdefs.NewDirNotFoundError(ctx.String("dir"))
	}

	return nil
}

func buildAction(ctx *cli.Context) error {
	if err := buildCheck(ctx); err != nil {
		return err
	}

	var req *pb.BuildProjectReq

	if ctx.String("remote") != "" {
		req = &pb.BuildProjectReq{
			Remote:        true,
			RemoteProject: ctx.String("remote"),
			Engine:        ctx.String("engine"),
		}
	} else {
		prjInfo, err := project.ParseProjectFromDir(ctx.String("dir"))
		if err != nil {
			return errors.Wrap(err, "fail to parse project from the directory")
		}

		req = &pb.BuildProjectReq{
			Remote:  false,
			Project: prjInfo,
			Engine:  ctx.String("engine"),
		}
	}

	c, err := client.NewClientFromContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	svc := pb.NewProjectMgrClient(c.Connection())

	stream, err := svc.BuildProject(context.Background(), req)
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
