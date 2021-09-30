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
	"fmt"
	"github.com/pkg/errors"
	pb "secpaver/api/proto/project"
	"secpaver/common/utils"
)

// CheckProject checks a project is valid
func CheckProject(project *pb.ProjectInfo) error {
	if project == nil {
		return fmt.Errorf("nil project data")
	}

	if err := checkProjectMeta(project.GetMeta()); err != nil {
		return errors.Wrap(err, "fail to check project meta data")
	}

	if err := checkProjectResources(project.GetResource()); err != nil {
		return errors.Wrap(err, "fail to check project resource data")
	}

	if len(project.GetSpecs()) == 0 {
		return fmt.Errorf("can't find any spec data")
	}

	for _, spec := range project.GetSpecs() {
		if spec.GetName() == "" {
			return fmt.Errorf("spec name not set")
		}

		if err := utils.CheckUnsafePath(spec.GetName()); err != nil {
			return errors.Wrapf(err, "fail to check spec file name %s", spec.GetName())
		}

		if err := checkProjectSpec(spec); err != nil {
			return errors.Wrapf(err, "fail to check %s spec file", spec.GetName())
		}
	}

	return nil
}
