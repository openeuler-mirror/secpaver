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
main package of pav application
*/
package main

import (
	"fmt"
	"github.com/urfave/cli"
	"os"
	"gitee.com/openeuler/secpaver/common/global"
	"gitee.com/openeuler/secpaver/modules/client"
	"sort"
)

func newPavApp() *cli.App {
	app := cli.NewApp()
	app.Name = "pav"
	app.Usage = "pav is a command line client for secPaver"
	app.Version = global.PavVersion
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "socket, s",
			Usage: "Set socket file path of grpc connection",
			Value: global.DefaultGrpcSocket,
		},
	}

	app.Commands = client.NewClientCmd()
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}

func main() {
	if err := newPavApp().Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
