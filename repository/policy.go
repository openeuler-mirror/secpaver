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

package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"gitee.com/openeuler/secpaver/common/errdefs"
	"gitee.com/openeuler/secpaver/common/log"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/domain"
)

// GetPolicyRoot returns policy store root
func (r *repo) GetPolicyRoot() string {
	return r.policyRoot
}

// FindAllProjects searches and returns a list of existed project
func (r *repo) FindAllPolicies() ([]*domain.Policy, error) {
	engDirs, err := utils.FindAllSubDir(r.policyRoot)
	if err != nil {
		log.Errorf("fail to search policy root: %v", err)
		return nil, fmt.Errorf("fail to search policy root")
	}

	var policies []*domain.Policy

	for _, dir := range engDirs {
		dirs, err := utils.FindAllSubDir(dir)
		if err != nil {
			log.Errorf("fail to search %s engine directory: %v", filepath.Base(dir), err)
			return nil, fmt.Errorf("fail to search %s engine directory", filepath.Base(dir))
		}

		for _, d := range dirs {
			if err := checkUID(d); err != nil {
				log.Errorf("fail to check uid of %s policy directory: %v", filepath.Base(d), err)
				return nil, fmt.Errorf("fail to check uid of %s policy directory", filepath.Base(d))
			}

			policies = append(policies, domain.NewPolicy(d))
		}
	}

	return policies, nil
}

// FindPolicyByName search the policy by specified name
// if the policy doesn't exist, returns an error
func (r *repo) FindPolicyByName(name, engine string) (*domain.Policy, error) {
	path := filepath.Join(r.policyRoot, engine, name)

	exist, err := utils.DirExist(path)
	if err != nil {
		log.Errorf("fail to search %s policy directory: %v", name, err)
		return nil, fmt.Errorf("fail to search %s policy directory", name)
	}
	if !exist {
		return nil, errdefs.NewDirNotFoundError(name)
	}

	if err := checkUID(path); err != nil {
		log.Errorf("fail to check uid of %s policy directory: %v", name, err)
		return nil, fmt.Errorf("fail to check uid of %s policy directory", name)
	}

	return domain.NewPolicy(path), nil
}

// DeletePolicy deletes an existed project
// if the project doesn't exist, returns an error
func (r *repo) DeletePolicy(policy *domain.Policy) error {
	if err := os.RemoveAll(policy.Path); err != nil {
		log.Errorf("fail to remove %s policy directory", policy.Name)
		return fmt.Errorf("fail to remove %s policy directory", policy.Name)
	}

	log.Infof("remove %s policy directory", policy.Name)

	return nil
}

// ExportPolicyZip exports an existed policy to zip file
func (r *repo) ExportPolicyZip(policy *domain.Policy, zipName string) ([]byte, error) {
	zipPath := filepath.Join(r.policyRoot, zipName)
	if err := utils.ZipDir(policy.Path, zipPath); err != nil {
		log.Errorf("fail to zip %s policy files: %v", policy.Name, err)
		return nil, fmt.Errorf("fail to zip %s policy files", policy.Name)
	}
	defer func() {
		if err := os.Remove(zipPath); err != nil {
			log.Errorf("fail to remove %s policy zip file: %v", policy.Name, err)
		}
	}()

	if err := utils.CheckFileSize(zipPath); err != nil {
		log.Errorf("fail to check size of %s policy zip file: %v", policy.Name, err)
		return nil, fmt.Errorf("fail to check size of %s policy zip file", policy.Name)
	}

	data, err := ioutil.ReadFile(zipPath)
	if err != nil {
		log.Errorf("fail to read %s policy zip file", policy.Name)
		return nil, fmt.Errorf("fail to read %s policy zip file", policy.Name)
	}

	log.Infof("export %s policy", policy.Name)

	return data, nil
}
