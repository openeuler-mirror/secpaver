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
Package client implements the Pav cmd application.
*/
package client

import (
	"github.com/urfave/cli"
	"secpaver/modules/client/engine"
	"secpaver/modules/client/policy"
	"secpaver/modules/client/project"
)

// NewClientCmd returns all client commands of Pav
func NewClientCmd() []cli.Command {
	return []cli.Command{
		*project.NewProjectCmd(),
		*engine.NewEngineCmd(),
		*policy.NewPolicyCmd(),
	}
}
