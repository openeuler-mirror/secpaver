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
	"archive/zip"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"gitee.com/openeuler/secpaver/common/errdefs"
	"gitee.com/openeuler/secpaver/common/log"
	"gitee.com/openeuler/secpaver/common/project"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/domain"
	"strconv"
	"time"
)

// GetProjectRoot returns project store root
func (r *repo) GetProjectRoot() string {
	return r.projectRoot
}

// FindAllProjects searches and returns a list of existed project
func (r *repo) FindAllProjects() ([]*domain.Project, error) {
	dirs, err := utils.FindAllSubDir(r.projectRoot)
	if err != nil {
		log.Errorf("fail to search project root: %v", err)
		return nil, fmt.Errorf("fail to search project root")
	}

	prjs := make([]*domain.Project, 0, len(dirs))
	for _, dir := range dirs {
		if err := checkUID(dir); err != nil {
			log.Errorf("fail to check uid of %s project directory: %v", filepath.Base(dir), err)
			return nil, fmt.Errorf("fail to check uid of %s project directory", filepath.Base(dir))
		}

		prjs = append(prjs, domain.NewProject(dir))
	}

	return prjs, nil
}

// FindProjectByName search the project by specified name
// if the project doesn't exist, returns an error
func (r *repo) FindProjectByName(prjName string) (*domain.Project, error) {
	path := filepath.Join(r.projectRoot, prjName)

	exist, err := utils.DirExist(path)
	if err != nil {
		log.Errorf("fail to search %s project directory: %v", prjName, err)
		return nil, fmt.Errorf("fail to search %s project directory", prjName)
	}
	if !exist {
		return nil, errdefs.NewDirNotFoundError(prjName)
	}

	if err := checkUID(path); err != nil {
		log.Errorf("fail to check uid of %s project directory: %v", prjName, err)
		return nil, fmt.Errorf("fail to check uid of %s projects directory", prjName)
	}

	return domain.NewProject(path), nil
}

// AddProjectByZip imports a project by a zip file
// if the project is not valid, importing fails and returns an error
func (r *repo) AddProjectByZip(prjName string, zipReader *zip.Reader) error {
	info, err := project.ParseProjectFromZip(zipReader)
	if err != nil {
		return errors.Wrap(err, "fail to parse project info from zip")
	}

	if err := project.CheckProject(info); err != nil {
		return errors.Wrap(err, "invalid project data")
	}

	path := filepath.Join(r.projectRoot, prjName)
	if err := project.WriteProjectToDir(info, path); err != nil {
		return errors.Wrap(err, "fail to write project data to directory")
	}

	log.Infof("add %s project", prjName)

	return nil
}

// UpdateProjectByZip updates an existed project by a zip file
// if the project is not valid, updating fails and returns an error
func (r *repo) UpdateProjectByZip(prjName string, zipReader *zip.Reader) error {
	// firstly, add zip to a temp dir
	tempPrjName := prjName + strconv.FormatInt(time.Now().Unix(), 10)
	if err := r.AddProjectByZip(tempPrjName, zipReader); err != nil {
		return errors.Wrap(err, "fail to create temp project")
	}

	tempPrjPath := filepath.Join(r.projectRoot, tempPrjName)
	prjPath := filepath.Join(r.projectRoot, prjName)

	// the temp dir should be removed anyway
	defer func() {
		if err := os.RemoveAll(tempPrjPath); err != nil {
			log.Errorf("fail to remove temp project directory %s: %v", tempPrjName, err)
		}
	}()

	// if successful, remove old project and rename the temp dir
	if err := os.RemoveAll(prjPath); err != nil {
		log.Errorf("fail to remove old project directory %s: %v", prjName, err)
	}

	if err := os.Rename(tempPrjPath, prjPath); err != nil {
		log.Errorf("fail to rename temp project directory %s to %s: %v", tempPrjName, prjName, err)
		return fmt.Errorf("fail to rename temp project directory %s to %s", tempPrjName, prjName)
	}

	log.Infof("update %s projects", prjName)

	return nil
}

// DeleteProject deletes an existed project
// if the project doesn't exist, returns an error
func (r *repo) DeleteProject(project *domain.Project) error {
	if err := os.RemoveAll(project.Path); err != nil {
		log.Errorf("fail to remove %s project directory: %v", project.Name, err)
		return fmt.Errorf("fail to remove %s project directory", project.Name)
	}

	log.Infof("remove %s project directory", project.Name)

	return nil
}

// ExportProjectZip exports an existed project to zip file
func (r *repo) ExportProjectZip(project *domain.Project, zipName string) ([]byte, error) {
	zipPath := filepath.Join(r.projectRoot, zipName)
	if err := utils.ZipDir(project.Path, zipPath); err != nil {
		log.Errorf("fail to zip %s project files: %v", project.Name, err)
		return nil, fmt.Errorf("fail to zip %s project files", project.Name)
	}
	defer func() {
		if err := os.Remove(zipPath); err != nil {
			log.Errorf("fail to remove %s project zip: %v", project.Name, err)
		}
	}()

	if err := utils.CheckFileSize(zipPath); err != nil {
		log.Errorf("fail to check %s project zip file: %v", project.Name, err)
		return nil, fmt.Errorf("fail to check size of %s project zip file", project.Name)
	}

	data, err := ioutil.ReadFile(zipPath)
	if err != nil {
		log.Errorf("fail to read %s project zip file", project.Name)
		return nil, fmt.Errorf("fail to read %s project zip file", project.Name)
	}

	log.Infof("export %s project", project.Name)

	return data, nil
}
