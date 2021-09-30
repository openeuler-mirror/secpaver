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
	"secpaver/engine/selinux/pkg/serule"
)

var defaultAvcRulesText = []string{
	"allow SUBJECT user_devpts_t:chr_file { append read write getattr ioctl };",
	"allow SUBJECT file_type:dir { search getattr };",
}

var defaultNetworkInetAvcRulesText = []string{
	"allow SUBJECT sysctl_net_t:dir { search };",
	"allow SUBJECT sysctl_net_t:file { open read };",
}

var defaultAvcRules []serule.AvcRule
var defaultNetworkInetAvcRules []serule.AvcRule

func init() {
	for _, line := range defaultAvcRulesText {
		if rule, _ := serule.ParseAvcRule(line); rule != nil {
			defaultAvcRules = append(defaultAvcRules, *rule)
		}
	}

	for _, line := range defaultNetworkInetAvcRulesText {
		if rule, _ := serule.ParseAvcRule(line); rule != nil {
			defaultNetworkInetAvcRules = append(defaultNetworkInetAvcRules, *rule)
		}
	}
}

func getDefaultRules(subject string) []serule.Rule {
	var rules []serule.Rule
	for _, rule := range defaultAvcRules {
		tmpRule := rule
		tmpRule.Subject = subject
		rules = append(rules, &tmpRule)
	}

	return rules
}

func getUnconfinedRules(subject string) []serule.Rule {
	var rules []serule.Rule
	rules = append(rules, serule.CreateUnconfinedCapRules(subject)...)
	rules = append(rules, serule.CreateUnconfinedProcessRules(subject)...)
	return rules
}
