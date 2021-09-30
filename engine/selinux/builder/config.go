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

package builder

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

const maxExtraRuleNumber = 5000

type config struct {
	Policy     *policyOpt `json:"policy,omitempty"`
	ExtraRules []string   `json:"extraRules,omitempty"`
}

type policyOpt struct {
	Monolithic bool `json:"monolithic,omitempty"`
}

func parseProjectSelinuxConfigFile(data []byte) (*config, error) {
	info := &config{}

	if err := json.Unmarshal(data, info); err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal json file")
	}

	if err := checkProjectSelinuxConfigFile(info); err != nil {
		return nil, errors.Wrap(err, "fail to check project selinux config file")
	}

	return info, nil
}

func checkProjectSelinuxConfigFile(conf *config) error {
	if len(conf.getExtraRules()) > maxExtraRuleNumber {
		return fmt.Errorf("the number of extra rules should be less than %d", maxExtraRuleNumber)
	}

	return nil
}

func (m *config) getPolicy() *policyOpt {
	if m != nil {
		return m.Policy
	}
	return nil
}

func (m *config) getExtraRules() []string {
	if m != nil {
		return m.ExtraRules
	}
	return nil
}

func (m *policyOpt) getMonolithic() bool {
	if m != nil {
		return m.Monolithic
	}
	return false
}
