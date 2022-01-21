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
main package of Pavd application
*/
package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"gitee.com/openeuler/secpaver/common/errdefs"
	"gitee.com/openeuler/secpaver/common/global"
	"gitee.com/openeuler/secpaver/common/utils"
	"sort"
)

func checkArgs(ctx *cli.Context) error {
	if err := utils.CheckCommandArgs(ctx, 0, utils.CheckExactArgs); err != nil {
		return err
	}

	file := ctx.String("config")
	if file == "" {
		return fmt.Errorf("please specified config file")
	}

	exists, err := utils.PathExist(file)
	if err != nil {
		return errors.Wrap(err, "fail to search config file")
	}
	if !exists {
		return errdefs.NewFileNotFoundError(filepath.Base(file))
	}

	return nil
}

func doBeforeJob(ctx *cli.Context) error {
	return checkArgs(ctx)
}

func newPavdApp() *cli.App {
	app := cli.NewApp()
	app.Name = "pavd"
	app.Usage = "Pavd daemon is the security modeling system daemon"
	app.Version = global.PavdVerison

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "socket, s",
			Usage: "Set socket file path of grpc connection",
			Value: global.DefaultGrpcSocket,
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Set config file",
			Value: global.DefaultConfigFile,
		},
		cli.StringFlag{
			Name: "log-level, l",
			Usage: fmt.Sprintf("Set the logging level (default: \"%s\")",
				global.DefaultLogLevel),
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	app.Before = doBeforeJob
	app.Action = runPavd
	return app
}

func main() {
	if err := newPavdApp().Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
