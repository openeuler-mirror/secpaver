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
Package selinux implements selinux policy engine
*/
package selinux

import (
	"secpaver/api/proto/policy"
	"secpaver/api/proto/project"
	"secpaver/domain"
	"secpaver/engine"
	"secpaver/engine/selinux/builder"
	"secpaver/engine/selinux/manager"
)

// Engine is the SELinux engine instance
type Engine struct {
	engine.UnimplementedEngine
	Manager *manager.Manager
	Builder *builder.Builder
}

// NewEngine create a SELinux engine
func NewEngine() *Engine {
	return &Engine{
		Manager: manager.NewManager(),
		Builder: builder.NewBuilder(),
	}
}

// GetName returns engine name
func (e *Engine) GetName() string {
	return "selinux"
}

// GetDescription returns engine description
func (e *Engine) GetDescription() string {
	return "SELinux policy generator"
}

// Build does build work
func (e *Engine) Build(prjInfo *project.ProjectInfo, out string, msg chan *project.Ack) error {
	return e.Builder.Build(prjInfo, out, msg)
}

// GetPolicyStatus returns the status of policy
func (e *Engine) GetPolicyStatus(policy *domain.Policy) (string, error) {
	return e.Manager.GetPolicyStatus(policy)
}

// Install installs a policy to system
func (e *Engine) Install(policy *domain.Policy, msg chan *policy.Ack) error {
	return e.Manager.Install(policy, msg)
}

// Uninstall uninstalls a policy in system
func (e *Engine) Uninstall(policy *domain.Policy, msg chan *policy.Ack) error {
	return e.Manager.Uninstall(policy, msg)
}
