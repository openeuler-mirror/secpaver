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
	pb "gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/utils"
)

const (
	maxApplicationListLength = 100
	maxPermissionListLength  = 500
)

func parseProjectSpec(data []byte) (*pb.SpecInfo, error) {
	spec := &pb.SpecInfo{}

	if err := json.Unmarshal(data, spec); err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal json file")
	}

	if err := checkProjectSpec(spec); err != nil {
		return nil, errors.Wrap(err, "fail to check spec file")
	}

	return spec, nil
}

func checkProjectSpec(spec *pb.SpecInfo) error {
	if len(spec.GetApplicationList()) == 0 {
		return fmt.Errorf("the length of applicationList can't be zero")
	}

	if len(spec.GetApplicationList()) > maxApplicationListLength {
		return fmt.Errorf(
			"the length of applicationList should be less than %d", maxApplicationListLength)
	}

	for _, policy := range spec.GetApplicationList() {
		if policy.GetApplication() == nil {
			return fmt.Errorf("invalid application definition")
		}
	}

	return nil
}

func regularProjectSpec(specInfo *pb.SpecInfo, resInfo *pb.ResourceInfo) error {
	for _, policy := range specInfo.GetApplicationList() {
		regPath, err := MacroReplace(policy.GetApplication().GetPath(), resInfo.GetMacroList())
		if err != nil {
			return err
		}

		if findResourceByPath(regPath, resInfo.GetResourceList()) == nil {
			return fmt.Errorf("undefined application path %s", regPath)
		}

		if policy.GetApplication() == nil {
			return fmt.Errorf("invalid application definition")
		}

		policy.Application.Path = regPath

		if len(policy.GetPermissionList()) > maxPermissionListLength {
			return fmt.Errorf(
				"the length of permission list should be less than %d", maxPermissionListLength)
		}

		for _, permission := range policy.GetPermissionList() {
			if permission.GetType() != RuleFileSystem {
				continue
			}

			regPaths, err := getRegularFileResources(permission.GetResources(), resInfo)
			if err != nil {
				return err
			}

			permission.Resources = regPaths
		}
	}

	return nil
}

func getRegularFileResources(resources []string, resInfo *pb.ResourceInfo) ([]string, error) {
	var regPaths []string

	for _, res := range resources {
		// check if resource is a group expression
		if name := getGroupName(res); name != "" {
			if group := findGroupByName(name, resInfo.GetGroupList()); group != nil {
				regPaths = append(regPaths, group.GetResources()...)
				continue
			}
		}

		// resource is not a group expression
		regPath, err := MacroReplace(res, resInfo.GetMacroList())
		if err != nil {
			return nil, err
		}

		if findResourceByPath(regPath, resInfo.GetResourceList()) == nil {
			return nil, fmt.Errorf("undefined file resource path %s", regPath)
		}

		if !utils.IsExistItem(regPath, regPaths) {
			regPaths = append(regPaths, regPath)
		}
	}

	return regPaths, nil
}
