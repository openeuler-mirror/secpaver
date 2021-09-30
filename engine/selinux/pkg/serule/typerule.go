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
	"strings"
)

// TypeRule is selinux type rule model
type TypeRule struct {
	Prefix     string
	Subject    string
	Object     string
	Class      string
	Target     string
	ObjectName string
}

// Text returns the expression for a SELinux type_transition rule
func (r *TypeRule) Text() string {
	if r.ObjectName == "" {
		return fmt.Sprintf("%s %s %s : %s %s;\n",
			r.Prefix, r.Subject, r.Object, r.Class, r.Target)
	}

	return fmt.Sprintf("%s %s %s : %s %s \"%s\";\n",
		r.Prefix, r.Subject, r.Object, r.Class, r.Target, r.ObjectName)
}

// Check checks if the rule is valid
func (r *TypeRule) Check() error {
	return checkTypeRule(r)
}

// CheckConflict check if two rules are conflict
func (r *TypeRule) CheckConflict(rule *TypeRule) error {
	if (r.Prefix == rule.Prefix) &&
		(r.Subject == rule.Subject) &&
		(r.Object == rule.Object) &&
		(r.Class == rule.Class) &&
		(r.ObjectName == rule.ObjectName) {

		if r.Target != rule.Target {
			return fmt.Errorf("type rule conflict")
		}
	}

	return nil
}

// NewTypeRule returns a SELinux type_transition rule
func NewTypeRule(prefix, sub, obj, class, tg, name string) *TypeRule {
	return &TypeRule{
		Prefix:     prefix,
		Subject:    sub,
		Object:     obj,
		Class:      class,
		Target:     tg,
		ObjectName: name,
	}
}

const (
	numTypeDomain    = 2
	numTypePreFields = 3
	numTypeSufFields = 2
)

// ParseTypeRule parse a string to type rule
func ParseTypeRule(line string) (*TypeRule, error) {
	rule := &TypeRule{}

	lineTrim := strings.TrimSpace(line)
	lineTrim = strings.TrimSuffix(lineTrim, ";")

	arr := strings.Split(lineTrim, ":")
	if len(arr) != numTypeDomain {
		return nil, fmt.Errorf("invalid type rule format: %s", line)
	}

	// get prefix, subject and object
	fields := strings.Fields(arr[0])
	if len(fields) != numTypePreFields {
		return nil, fmt.Errorf("invalid type rule format: %s", line)
	}

	rule.Prefix, rule.Subject, rule.Object = fields[0], fields[1], fields[2]

	// get class and target
	fields = strings.Fields(arr[1])
	if len(fields) < numTypeSufFields {
		return nil, fmt.Errorf("invalid type rule format: %s", line)
	}

	rule.Class = fields[0]
	rule.Target = fields[1]

	if len(fields) > numTypeSufFields {
		name := strings.TrimPrefix(fields[2], "\"")
		name = strings.TrimSuffix(name, "\"")
		rule.ObjectName = name
	}

	return rule, nil
}
