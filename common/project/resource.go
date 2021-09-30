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

package project

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"path/filepath"
	"regexp"
	pb "secpaver/api/proto/project"
	"secpaver/common/utils"
	"strings"
)

const (
	maxMacroDepth         = 10
	maxMacroListLength    = 50
	maxMacroNameLength    = 50
	maxGroupListLength    = 50
	maxResourceListLength = 5000
)

var (
	macroRegexp     = regexp.MustCompile("\\$\\([a-zA-Z_][a-zA-Z0-9_]*\\)")
	macroNameRegexp = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
)

func parseProjectResources(data []byte) (*pb.ResourceInfo, error) {
	info := &pb.ResourceInfo{}

	err := json.Unmarshal(data, info)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal resource json file")
	}

	if err := checkProjectResources(info); err != nil {
		return nil, errors.Wrap(err, "fail to check resource file")
	}

	return info, nil
}

func checkProjectResources(resInfo *pb.ResourceInfo) error {
	if resInfo == nil {
		return fmt.Errorf("nil resource info")
	}

	if len(resInfo.GetResourceList()) == 0 {
		return fmt.Errorf("the length of resourceList can't be zero")
	}

	if len(resInfo.GetResourceList()) > maxResourceListLength {
		return fmt.Errorf("the length of resource list should be less than %d", maxResourceListLength)
	}

	return checkMacroAndGroup(resInfo.GetMacroList(), resInfo.GetGroupList())
}

func checkMacroAndGroup(macros []*pb.Macro, groups []*pb.ResGroup) error {
	if len(macros) > maxMacroListLength {
		return fmt.Errorf("the length of macro list needs to be less than %d", maxMacroListLength)
	}

	if len(groups) > maxGroupListLength {
		return fmt.Errorf("the length of group list needs to be less than %d", maxGroupListLength)
	}

	macroMap := make(map[string]string, len(macros))
	groupMap := make(map[string][]string, len(groups))

	for _, macro := range macros {
		if !isValidMacroName(macro.GetName()) {
			return fmt.Errorf("invalid macro name %s", macro.GetName())
		}

		if _, ok := macroMap[macro.GetName()]; ok {
			return fmt.Errorf("redefine macro %s", macro.GetName())
		}

		macroMap[macro.GetName()] = macro.GetValue()
	}

	for _, group := range groups {
		if !isValidMacroName(group.GetName()) {
			return fmt.Errorf("invalid group name %s", group.GetName())
		}

		if _, ok := groupMap[group.GetName()]; ok {
			return fmt.Errorf("redefine group %s", group.GetName())
		}

		if _, ok := macroMap[group.GetName()]; ok {
			return fmt.Errorf("group name %s has been used in macros", group.GetName())
		}
	}

	return nil
}

func regularProjectResources(resInfo *pb.ResourceInfo) error {
	for _, resource := range resInfo.GetResourceList() {
		regPath, err := MacroReplace(resource.GetPath(), resInfo.GetMacroList())
		if err != nil {
			return errors.Wrapf(err, "fail to make resource path regular")
		}

		if !utils.IsAbsolutePath(regPath) {
			return fmt.Errorf("invaild file path %s in resource define, path must be absolute", regPath)
		}

		resource.Path = regPath
	}

	for _, group := range resInfo.GetGroupList() {
		for i, path := range group.GetResources() {
			regPath, err := MacroReplace(path, resInfo.GetMacroList())
			if err != nil {
				return errors.Wrapf(err, "fail to make resource path regular in group %s", group.Name)
			}

			group.Resources[i] = regPath
		}
	}

	return nil
}

func findResourceByPath(path string, resources []*pb.Resource) *pb.Resource {
	for _, resource := range resources {
		if resource.GetPath() == path {
			return resource
		}
	}

	return nil
}

func findGroupByName(name string, groups []*pb.ResGroup) *pb.ResGroup {
	for _, group := range groups {
		if group.GetName() == name {
			return group
		}
	}

	return nil
}

// MacroReplace expand macro of path
func MacroReplace(path string, macros []*pb.Macro) (string, error) {
	if len(macros) > maxMacroListLength {
		return "", fmt.Errorf("the length of macro list needs to be less than %d", maxMacroListLength)
	}

	for i := 0; i < maxMacroDepth; i++ {
		p, err := macroReplace(path, macros)
		if err != nil {
			return "", err
		}

		path = p
	}

	if strings.Count(path, "$(") > 0 {
		return "", fmt.Errorf("the depth of macro in %s needs to be less than %d", path, maxMacroDepth)
	}

	return filepath.Clean(path), nil
}

func macroReplace(path string, macros []*pb.Macro) (string, error) {
	oriPath := path
	macroStrArr := macroRegexp.FindAllString(path, -1)

	for _, macroStr := range macroStrArr {
		name := getMacroName(macroStr)
		val := getMacroValue(name, macros)
		if val == "" {
			return "", fmt.Errorf("undefined macro %s in path %s", name, oriPath)
		}

		path = strings.Replace(path, macroStr, val, 1)
	}

	return path, nil
}

func isValidMacroName(name string) bool {
	if len(name) > maxMacroNameLength {
		return false
	}

	return macroNameRegexp.MatchString(name)
}

func getMacroValue(name string, macros []*pb.Macro) string {
	for _, macro := range macros {
		if macro.GetName() == name {
			return macro.GetValue()
		}
	}

	return ""
}

func getMacroName(str string) string {
	str = strings.TrimSpace(str)
	if !isMacro(str) {
		return ""
	}

	str = strings.TrimSuffix(str, ")")
	str = strings.TrimPrefix(str, "$(")
	return str
}

func isMacro(str string) bool {
	return macroRegexp.MatchString(str)
}

func getGroupName(str string) string {
	return getMacroName(str)
}
