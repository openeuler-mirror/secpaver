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
	"path/filepath"
	"gitee.com/openeuler/secpaver/common/utils"
)

var avcRulePrefixSet = []string{
	"allow",
	"dontaudit",
	"auditallow",
	"neverallow",
}

var typeRulePrefixSet = []string{
	"type_transition",
}

func checkAvcRule(rule *AvcRule) error {
	if err := checkAvcRulePrefix(rule.Prefix); err != nil {
		return err
	}

	if err := checkClass(rule.Class); err != nil {
		return err
	}

	if err := checkPermission(rule.Class, rule.Actions); err != nil {
		return err
	}

	return nil
}

func checkTypeRule(rule *TypeRule) error {
	if err := checkTypeRulePrefix(rule.Prefix); err != nil {
		return err
	}

	if err := checkClass(rule.Class); err != nil {
		return err
	}

	return nil
}

func checkAvcRulePrefix(prefix string) error {
	if !utils.IsExistItem(prefix, avcRulePrefixSet) {
		return fmt.Errorf("invalid avc rule prefix %s", prefix)
	}

	return nil
}

func checkTypeRulePrefix(prefix string) error {
	if !utils.IsExistItem(prefix, typeRulePrefixSet) {
		return fmt.Errorf("invalid type rule prefix %s", prefix)
	}

	return nil
}

func checkClass(class string) error {
	exist, err := utils.PathExist(filepath.Join(classRoot, class))
	if err != nil {
		return fmt.Errorf("please check selinux mnt directory")
	}

	if !exist {
		return fmt.Errorf("invalid selinux class %s", class)
	}

	return nil
}

func checkPermission(class string, perms []string) error {
	if len(perms) == 0 {
		return fmt.Errorf("at least one permission should be set")
	}

	for _, perm := range perms {
		if perm == "*" {
			continue
		}

		exist, err := utils.PathExist(filepath.Join(classRoot, class, "perms", perm))
		if err != nil {
			return fmt.Errorf("please check selinux mnt directory")
		}

		if !exist {
			return fmt.Errorf("invalid selinux permission %s of class %s", perm, class)
		}
	}

	return nil
}
