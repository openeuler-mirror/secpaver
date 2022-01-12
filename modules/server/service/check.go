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

package service

import (
	"github.com/pkg/errors"
	pbEngine "gitee.com/openeuler/secpaver/api/proto/engine"
	pbPolicy "gitee.com/openeuler/secpaver/api/proto/policy"
	pbProject "gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/utils"
)

func checkEngineReq(req *pbEngine.Req) error {
	return checkEngineName(req.GetName())
}

func checkProjectReq(req *pbProject.Req) error {
	return checkProjectName(req.GetName())
}

func checkImportProjectReq(req *pbProject.ImportProjectReq) error {
	return utils.CheckZipFileName(req.GetFile().GetFilename())
}

func checkBuildProjectReq(req *pbProject.BuildProjectReq) error {
	if err := checkEngineName(req.GetEngine()); err != nil {
		return err
	}

	if req.GetRemote() {
		return checkProjectName(req.GetRemoteProject())
	}

	return checkProjectName(req.GetProject().GetMeta().GetName())
}

func checkPolicyReq(req *pbPolicy.Req) error {
	if err := checkPolicyName(req.GetName()); err != nil {
		return err
	}

	return checkEngineName(req.GetEngine())
}

func checkEngineName(name string) error {
	if err := utils.CheckUnsafeArg(name); err != nil {
		return errors.Wrapf(err, "invalid engine name %s", name)
	}

	return nil
}

func checkProjectName(name string) error {
	if err := utils.CheckUnsafeArg(name); err != nil {
		return errors.Wrapf(err, "invalid project name %s", name)
	}

	return nil
}

func checkPolicyName(name string) error {
	if err := utils.CheckUnsafeArg(name); err != nil {
		return errors.Wrapf(err, "invalid policy name %s", name)
	}

	return nil
}
