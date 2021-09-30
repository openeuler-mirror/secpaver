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

package serule

import (
	"fmt"
	"secpaver/common/utils"
	"sort"
	"strings"
)

// AvcRule is the selinux avc rule model
type AvcRule struct {
	Prefix  string
	Subject string
	Object  string
	Class   string
	Actions []string
}

// Text generate a string of selinux rule
func (r *AvcRule) Text() string {
	if utils.IsExistItem("*", r.Actions) {
		return fmt.Sprintf("%s %s %s : %s *;\n",
			r.Prefix, r.Subject, r.Object, r.Class)
	}

	sort.Strings(r.Actions)
	actStr := utils.ShowStringsWithSpace(r.Actions)

	return fmt.Sprintf("%s %s %s : %s { %s };\n",
		r.Prefix, r.Subject, r.Object, r.Class, actStr)
}

// Check checks if the rule is valid
func (r *AvcRule) Check() error {
	return checkAvcRule(r)
}

// Merge merges 2 avc rules, if merge success, return true
func (r *AvcRule) Merge(r1 *AvcRule) bool {
	if (r.Prefix == r1.Prefix) &&
		(r.Subject == r1.Subject) &&
		(r.Object == r1.Object) &&
		(r.Class == r1.Class) {

		for _, act := range r1.Actions {
			if utils.IsExistItem("*", r.Actions) ||
				utils.IsExistItem("*", r1.Actions) {

				r.Actions = []string{"*"}

				return true
			}

			if !utils.IsExistItem(act, r.Actions) {
				r.Actions = append(r.Actions, act)
			}
		}

		return true
	}

	return false
}

// NewAvcRule returns a selinux policy rule
func NewAvcRule(prefix string, sub, obj string, class string, acts []string) *AvcRule {
	if utils.IsExistItem("*", acts) {
		return &AvcRule{
			Prefix:  prefix,
			Subject: sub,
			Object:  obj,
			Class:   class,
			Actions: []string{"*"},
		}
	}

	return &AvcRule{
		Prefix:  prefix,
		Subject: sub,
		Object:  obj,
		Class:   class,
		Actions: acts,
	}
}

const (
	numAvcDomain    = 2
	numAvcPreFields = 3
	numAvcSufFields = 2
)

// ParseAvcRule parse a string to avc rule
func ParseAvcRule(line string) (*AvcRule, error) {
	rule := &AvcRule{}

	lineTrim := strings.TrimSpace(line)
	lineTrim = strings.TrimSuffix(lineTrim, ";")

	arr := strings.Split(lineTrim, ":")
	if len(arr) != numAvcDomain {
		return nil, fmt.Errorf("invalid avc rule format: %s", line)
	}

	// get prefix, subject and object
	fields := strings.Fields(arr[0])
	if len(fields) != numAvcPreFields {
		return nil, fmt.Errorf("invalid avc rule format: %s", line)
	}

	rule.Prefix, rule.Subject, rule.Object = fields[0], fields[1], fields[2]

	// get class and permissions
	fields = strings.Fields(arr[1])
	if len(fields) < numAvcSufFields {
		return nil, fmt.Errorf("invalid avc rule format: %s", line)
	}

	for i, field := range fields {
		if i == 0 {
			rule.Class = fields[0]
			continue
		}

		action := strings.TrimSuffix(field, "}")
		action = strings.TrimPrefix(action, "{")

		if action != "" {
			rule.Actions = append(rule.Actions, action)
		}
	}

	return rule, nil
}
