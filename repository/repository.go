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
Package repository implements the secPaver file resource manage function.
*/
package repository

import (
	"archive/zip"
	"fmt"
	"github.com/pkg/errors"
	"gitee.com/openeuler/secpaver/common/config"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/domain"
)

var repoInstance *repo

// Repo is the interface of secPaver repository manager
type Repo interface {
	GetProjectRoot() string
	FindAllProjects() ([]*domain.Project, error)
	FindProjectByName(name string) (*domain.Project, error)
	AddProjectByZip(name string, zipReader *zip.Reader) error
	UpdateProjectByZip(name string, zipReader *zip.Reader) error
	DeleteProject(project *domain.Project) error
	ExportProjectZip(project *domain.Project, zipName string) ([]byte, error)

	GetPolicyRoot() string
	FindAllPolicies() ([]*domain.Policy, error)
	FindPolicyByName(name, engine string) (*domain.Policy, error)
	DeletePolicy(policy *domain.Policy) error
	ExportPolicyZip(policy *domain.Policy, zipName string) ([]byte, error)
}

// InitRepo set the root paths
func InitRepo(info *config.RepositoryInfo) error {
	if err := checkUID(info.ProjectRoot); err != nil {
		return errors.Wrap(err, "fail to check uid of project repository root")
	}

	if err := checkUID(info.PolicyRoot); err != nil {
		return errors.Wrap(err, "fail to check uid of policy repository root")
	}

	repoInstance = &repo{
		projectRoot: info.ProjectRoot,
		policyRoot:  info.PolicyRoot,
	}

	return nil
}

// GetRepo returns the repo single instance
func GetRepo() Repo {
	if repoInstance == nil {
		return nil
	}

	return repoInstance
}

type repo struct {
	projectRoot string
	policyRoot  string
}

func checkUID(path string) error {
	uid, err := utils.GetUIDOfFile(path)
	if err != nil {
		return err
	}

	if uid != 0 {
		return fmt.Errorf("the uid must be 0")
	}

	return nil
}
