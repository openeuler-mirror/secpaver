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
	defVersion    = "1.0"
	metaFile      = "pav.proj"
	maxSpecNumber = 100
)

func parseProjectMeta(data []byte) (*pb.MetaInfo, error) {
	meta := &pb.MetaInfo{}

	if err := json.Unmarshal(data, meta); err != nil {
		return nil, errors.Wrap(err, "fail to parse json file")
	}

	if meta.GetVersion() == "" {
		meta.Version = defVersion
	}

	if err := checkProjectMeta(meta); err != nil {
		return nil, errors.Wrap(err, "fail to check project meta")
	}

	return meta, nil
}

func checkProjectMeta(meta *pb.MetaInfo) error {
	if meta == nil {
		return fmt.Errorf("nil project meta")
	}

	if err := utils.CheckUnsafeArg(meta.GetName()); err != nil {
		return err
	}

	if err := utils.CheckVersion(meta.GetVersion()); err != nil {
		return err
	}

	// check resource definition file path
	if meta.GetResources() == "" {
		return fmt.Errorf("project should have a resource definition file")
	}

	if err := utils.CheckUnsafePath(meta.GetResources()); err != nil {
		return err
	}

	// check spec file path
	if len(meta.GetSpecs()) == 0 {
		return fmt.Errorf("project should have at least one application policy spec file")
	}

	if len(meta.GetSpecs()) > maxSpecNumber {
		return fmt.Errorf("specs number should be less than %d", maxSpecNumber)
	}

	for _, spec := range meta.GetSpecs() {
		if err := utils.CheckUnsafePath(spec); err != nil {
			return err
		}
	}

	// check selinux config file path
	seConf := meta.GetSelinux().GetConfig()
	if seConf != "" {
		if err := utils.CheckUnsafePath(seConf); err != nil {
			return err
		}
	}

	return nil
}
