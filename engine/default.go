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
	"fmt"
	pbPolicy "gitee.com/openeuler/secpaver/api/proto/policy"
	"gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/policy"
	"gitee.com/openeuler/secpaver/domain"
)

// UnimplementedEngine is a default implement of a policy engine
type UnimplementedEngine struct{}

// GetName is the default implement function of Engine interface
func (e *UnimplementedEngine) GetName() string {
	return "N/A"
}

// GetDescription is the default implement function of Engine interface
func (e *UnimplementedEngine) GetDescription() string {
	return "N/A"
}

// Build is the default implement function of Engine interface
func (e *UnimplementedEngine) Build(prjInfo *project.ProjectInfo, out string, msg chan *project.Ack) error {
	return fmt.Errorf("'Build' method not implemented in this engine")
}

// GetPolicyStatus is the default implement function of Engine interface
func (e *UnimplementedEngine) GetPolicyStatus(p *domain.Policy) (string, error) {
	return policy.StatusUnknown, nil
}

// Install is the default implement function of Engine interface
func (e *UnimplementedEngine) Install(policy *domain.Policy, msg chan *pbPolicy.Ack) error {
	return fmt.Errorf("'Install' method not implemented in this engine")
}

// Uninstall is the default implement function of Engine interface
func (e *UnimplementedEngine) Uninstall(policy *domain.Policy, msg chan *pbPolicy.Ack) error {
	return fmt.Errorf("'Uninstall' method not implemented in this engine")
}
