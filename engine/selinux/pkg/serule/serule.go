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
Package serule provides some functions for SELinux policy module.
*/
package serule

import (
	"fmt"
	"github.com/pkg/errors"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"strings"
)

const (
	classRoot = "/sys/fs/selinux/class"
)

// Rule is the selinux rule interface
type Rule interface {
	Text() string
	Check() error
}

// CreateCommonAllowRule returns a file avc rule
func CreateCommonAllowRule(sub, obj, class string, acts []string) (Rule, error) {
	rule := NewAvcRule("allow", sub, obj, class, acts)
	if err := rule.Check(); err != nil {
		return nil, errors.Wrap(err, "invalid common allow avc rule")
	}

	if rule == nil {
		return nil, nil
	}

	return rule, nil
}

// CreateFileAllowRules returns a file allow avc rule
func CreateFileAllowRules(sub, obj string, class secontext.FileClass, acts []string) ([]Rule, error) {
	var rules []Rule

	allClass := secontext.GetFileClassByID(class)
	for _, cls := range allClass {
		rule := NewAvcRule("allow", sub, obj, cls, acts)
		if err := rule.Check(); err != nil {
			return nil, errors.Wrap(err, "invalid common allow avc rule")
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// CreateFileTypeTransitionRule returns a type transition rule
func CreateFileTypeTransitionRule(sub, obj string, class secontext.FileClass, tgt, name string) ([]Rule, error) {
	if hasMetaCharacter(name) {
		name = ""
	}

	var rules []Rule

	allClass := secontext.GetFileClassByID(class)
	for _, cls := range allClass {
		rule := NewTypeRule("type_transition", sub, obj, cls, tgt, name)
		if err := rule.Check(); err != nil {
			return nil, errors.Wrap(err, "invalid file type transition rule")
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

// CreateFilesystemAllowRule returns a filesystem avc rule
func CreateFilesystemAllowRule(sub, obj string, acts []string) (Rule, error) {
	rule := NewAvcRule("allow", sub, obj, "filesystem", acts)
	if err := rule.Check(); err != nil {
		return nil, errors.Wrap(err, "invalid filesystem allow avc rule")
	}

	if rule == nil {
		return nil, nil
	}

	return rule, nil
}

// CreateCapabilityAllowRule returns a capability avc rule
func CreateCapabilityAllowRule(sub string, caps []string) ([]Rule, error) {
	var rules []Rule
	var cap1, cap2 []string

	for _, c := range caps {
		if checkPermission("capability", []string{c}) == nil {
			cap1 = append(cap1, c)
		} else if checkPermission("capability2", []string{c}) == nil {
			cap2 = append(cap2, c)
		} else {
			return nil, fmt.Errorf("invalid capability %s", c)
		}
	}

	if len(cap1) != 0 {
		rules = append(rules,
			NewAvcRule("allow", sub, sub, "capability", cap1))
	}

	if len(cap2) != 0 {
		rules = append(rules,
			NewAvcRule("allow", sub, sub, "capability2", cap2))
	}

	return rules, nil
}

// CreateProcessAllowRule returns a capability avc rule
func CreateProcessAllowRule(sub, obj string, perms []string) ([]Rule, error) {
	var rules []Rule
	var pro1, pro2 []string

	for _, perm := range perms {
		if checkPermission("process", []string{perm}) == nil {
			pro1 = append(pro1, perm)
		} else if checkPermission("process2", []string{perm}) == nil {
			pro2 = append(pro2, perm)
		} else {
			return nil, fmt.Errorf("invalid process permission %s", perm)
		}
	}

	if len(pro1) != 0 {
		rules = append(rules,
			NewAvcRule("allow", sub, obj, "process", pro1))
	}

	if len(pro2) != 0 {
		rules = append(rules,
			NewAvcRule("allow", sub, obj, "process2", pro2))
	}

	return rules, nil
}

// CreateNetifAllowRule returns a socket avc rule
func CreateNetifAllowRule(sub, obj string, acts []string) (Rule, error) {
	rule := NewAvcRule("allow", sub, obj, "netif", acts)
	if err := rule.Check(); err != nil {
		return nil, errors.Wrap(err, "invalid filesystem allow avc rule")
	}

	if rule == nil {
		return nil, nil
	}

	return rule, nil
}

// CreateNodeAllowRule returns a socket avc rule
func CreateNodeAllowRule(sub, obj string, acts []string) (Rule, error) {
	rule := NewAvcRule("allow", sub, obj, "node", acts)
	if err := rule.Check(); err != nil {
		return nil, errors.Wrap(err, "invalid filesystem allow avc rule")
	}

	if rule == nil {
		return nil, nil
	}

	return rule, nil
}

// CreateDomainAutoTransRule returns a domain auto transition rules
func CreateDomainAutoTransRule(sub, obj, tgt string) ([]Rule, error) {
	var rules []Rule

	rules = append(rules,
		NewTypeRule("type_transition", sub, obj, "process", tgt, ""))

	rules = append(rules,
		NewAvcRule("allow", tgt, obj, "file",
			[]string{"entrypoint", "execute", "getattr", "open", "read", "map"}))

	rules = append(rules,
		NewAvcRule("allow", sub, obj, "file",
			[]string{"execute", "getattr", "open", "read", "map"}))

	rules = append(rules,
		NewAvcRule("allow", tgt, obj, "process",
			[]string{"sigchld"}))

	rules = append(rules,
		NewAvcRule("allow", sub, tgt, "process",
			[]string{"rlimitinh", "siginh", "transition"}))

	for _, rule := range rules {
		if err := rule.Check(); err != nil {
			return nil, err
		}
	}

	return rules, nil
}

// CreateBaseFileTypeRules returns the base file type rules
func CreateBaseFileTypeRules(tp string, class secontext.FileClass) ([]Rule, error) {
	var rules []Rule
	allClass := secontext.GetFileClassByID(class)

	for _, cls := range allClass {
		rules = append(rules,
			NewAvcRule("allow", "unconfined_t", tp, cls,
				[]string{"getattr", "relabelto", "relabelfrom"}))

		rules = append(rules,
			NewAvcRule("allow", "restorecond_t", tp, cls,
				[]string{"getattr", "relabelto", "relabelfrom"}))

		rules = append(rules,
			NewAvcRule("allow", tp, "fs_t", "filesystem",
				[]string{"associate"}))

		for _, rule := range rules {
			if err := rule.Check(); err != nil {
				return nil, err
			}
		}
	}

	return rules, nil
}

// CreateUnconfinedProcessRules returns unconfined process rules
func CreateUnconfinedProcessRules(subject string) []Rule {
	var rules []Rule
	rules = append(rules,
		NewAvcRule("allow", subject, "domain", "process", []string{"*"}))

	return rules
}

// CreateUnconfinedCapRules returns unconfined capability rules
func CreateUnconfinedCapRules(subject string) []Rule {
	var rules []Rule
	rules = append(rules,
		NewAvcRule("allow", subject, subject, "capability", []string{"*"}))

	if checkClass("capability2") == nil {
		rules = append(rules,
			NewAvcRule("allow", subject, subject, "capability2", []string{"*"}))
	}

	return rules
}

// ParseRule parses an selinux rule
func ParseRule(line string) (Rule, error) {
	arr := strings.Fields(line)
	if len(arr) == 0 {
		return nil, fmt.Errorf("invalid selinux rule format %s", line)
	}

	if utils.IsExistItem(arr[0], avcRulePrefixSet) {
		return ParseAvcRule(line)
	}

	if utils.IsExistItem(arr[0], typeRulePrefixSet) {
		return ParseTypeRule(line)
	}

	return nil, fmt.Errorf("invalid selinux rule prefix %s", line)
}

var metaSet = []rune{'?', '*'}

func hasMetaCharacter(str string) bool {
	for i := 0; i < len(str); i++ {
		c := str[i]

		if c == '\\' {
			i++
			continue
		}

		if utils.IsExistItem(rune(c), metaSet) {
			return true
		}
	}

	return false
}
