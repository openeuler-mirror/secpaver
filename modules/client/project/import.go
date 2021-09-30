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
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"io/ioutil"
	"path/filepath"
	"secpaver/api/proto/project"
	"secpaver/common/client"
	"secpaver/common/errdefs"
	"secpaver/common/utils"
)

func newImportCommand() *cli.Command {
	return &cli.Command{
		Name:        "import",
		Usage:       "Import a project zip file",
		UsageText:   "Pav project import [-f] <FILE>",
		Description: "Import a project zip file",
		Action:      importAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:   "force, f",
				Usage:  "force import",
				Hidden: false,
			},
		},
	}
}

func importCheck(ctx *cli.Context) error {
	if err := utils.CheckCommandArgs(ctx, 1, utils.CheckExactArgs); err != nil {
		return err
	}

	file := ctx.Args().First()
	exist, err := utils.PathExist(file)
	if err != nil {
		return err
	} else if !exist {
		return errdefs.NewFileNotFoundError(file)
	}

	if err := utils.CheckFileSize(file); err != nil {
		return err
	}

	if err := utils.CheckValidZipFile(file); err != nil {
		return errors.Wrapf(err, "%s is not a valid zip file", file)
	}

	return nil
}

func importAction(ctx *cli.Context) error {
	if err := importCheck(ctx); err != nil {
		return err
	}

	file := ctx.Args().First()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "fail to read zip file")
	}

	c, err := client.NewClientFromContext(ctx)
	if err != nil {
		return err
	}
	defer c.Close()

	svc := project.NewProjectMgrClient(c.Connection())

	req := &project.ImportProjectReq{
		File: &project.ProjectZipFile{
			Filename: filepath.Base(file),
			Data:     data,
		},
		Force: ctx.Bool("force"),
	}

	ack, err := svc.ImportProject(context.Background(), req)
	if err != nil {
		return err
	}

	utils.PrintAck(ack.GetLevel(), ack.GetStatus())

	return nil
}
