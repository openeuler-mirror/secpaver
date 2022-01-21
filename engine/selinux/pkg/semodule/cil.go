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
Package semodule provides some functions for SELinux policy module.
*/
package semodule

import (
	"fmt"
	"github.com/pkg/errors"
	"path/filepath"
	"regexp"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/libselinux"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/libsemanage"
	"strconv"
	"strings"
)

const policyLibRoot = "/var/lib/selinux"

var fileConPathRegexp = regexp.MustCompile(`\(filecon\s".*"`)
var provideTypeRegexp = regexp.MustCompile(`\(type\s.*_t\)`)
var requireTypeRegexp = regexp.MustCompile(`\(typeattributeset\scil_gen_require\s.*_t\)`)

// CilInfo stores the parsing result of a SELinux cil module
type CilInfo struct {
	Path         string
	ProvideTypes []string
	RequireTypes []string
	FileConPaths []string
}

// ParseCilInfo parses a SELinux cil module and returns the CilInfo
func ParseCilInfo(module string, pri int) (*CilInfo, error) {
	if module == "" {
		return nil, fmt.Errorf("module name must be non-empty")
	}

	policyType, err := libselinux.GetPolicyType()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get policy type")
	}

	path := filepath.Join(getModulePath(policyType, module, pri), "cil")

	return parseCilInfo(path)
}

// GetModuleCilInfo returns the cil module information by specified name
func GetModuleCilInfo(name string) (*CilInfo, error) {
	module, err := GetModuleInfo(name)
	if err != nil {
		return nil, err
	}

	info, err := ParseCilInfo(module.Name, module.Priority)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func parseCilInfo(path string) (*CilInfo, error) {
	data, err := utils.ReadBzip2(path)
	if err != nil {
		return nil, errors.Wrap(err, "fail to read module cil file")
	}

	fcPaths := fileConPathRegexp.FindAllString(data, -1)
	for i := range fcPaths {
		fcPaths[i] = getFileConPath(fcPaths[i])
	}

	pTypes := provideTypeRegexp.FindAllString(data, -1)
	for i := range pTypes {
		pTypes[i] = getProvideType(pTypes[i])
	}

	rTypes := requireTypeRegexp.FindAllString(data, -1)
	for i := range rTypes {
		rTypes[i] = getRequireType(rTypes[i])
	}

	return &CilInfo{
		Path:         path,
		ProvideTypes: pTypes,
		RequireTypes: rTypes,
		FileConPaths: fcPaths,
	}, nil
}

func getModulePath(policyType, module string, pri int) string {
	root := libsemanage.Root()
	if root == "" {
		root = policyLibRoot
	}

	return filepath.Join(
		root, policyType, "active", "modules", strconv.Itoa(pri), module)
}

func getFileConPath(str string) string {
	result := strings.TrimPrefix(str, `(filecon "`)
	return strings.TrimSuffix(result, `"`)
}

func getProvideType(str string) string {
	result := strings.TrimPrefix(str, `(type `)
	return strings.TrimSuffix(result, `)`)
}

func getRequireType(str string) string {
	result := strings.TrimPrefix(str, `(typeattributeset cil_gen_require `)
	return strings.TrimSuffix(result, `)`)
}
