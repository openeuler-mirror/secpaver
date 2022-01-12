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
Package project implements the project parsing.
*/
package project

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	pb "gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/global"
	"gitee.com/openeuler/secpaver/common/log"
	"gitee.com/openeuler/secpaver/common/utils"
)

// ParseProjectFromDir parses the project file and returns the project build information
func ParseProjectFromDir(dir string) (*pb.ProjectInfo, error) {
	// parse meta file
	metaInfo, err := ReadProjectMetaFromDir(dir)
	if err != nil {
		return nil, err
	}

	// parse resources definition file
	data, err := ioutil.ReadFile(filepath.Join(dir, metaInfo.GetResources()))
	if err != nil {
		return nil, fmt.Errorf("fail to read %s file", metaInfo.GetResources())
	}

	resInfo, err := parseProjectResources(data)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to parse %s file", metaInfo.GetResources())
	}

	// parse spec files
	var specInfoList []*pb.SpecInfo
	for _, spec := range metaInfo.GetSpecs() {
		data, err := ioutil.ReadFile(filepath.Join(dir, spec))
		if err != nil {
			return nil, fmt.Errorf("fail to read %s file", spec)
		}

		specInfo, err := parseProjectSpec(data)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to parse %s file", spec)
		}

		specInfo.Name = spec
		specInfoList = append(specInfoList, specInfo)
	}

	// read extend config files
	extFiles := map[string][]byte{}
	seConf := metaInfo.GetSelinux().GetConfig()
	if seConf != "" {
		data, err = ioutil.ReadFile(filepath.Join(dir, seConf))
		if err != nil {
			return nil, fmt.Errorf("fail to read %s file", seConf)
		}

		extFiles[seConf] = data
	}

	return &pb.ProjectInfo{
		Meta:     metaInfo,
		Resource: resInfo,
		Specs:    specInfoList,
		Extends:  extFiles,
	}, nil
}

// ParseProjectFromZip parses the project file and returns the project build information
func ParseProjectFromZip(zip *zip.Reader) (*pb.ProjectInfo, error) {
	if zip == nil {
		return nil, fmt.Errorf("nil zip reader")
	}

	// parse project meta
	metaInfo, err := ReadProjectMetaFromZip(zip)
	if err != nil {
		return nil, err
	}

	// parse resources definition file
	data, err := utils.ExtractFileFromZip(zip, filepath.Join(metaInfo.Name, metaInfo.GetResources()))
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read %s file", metaInfo.GetResources())
	}

	resInfo, err := parseProjectResources(data)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to parse %s file", metaInfo.GetResources())
	}

	// parse spec files
	var specInfoList []*pb.SpecInfo
	for _, spec := range metaInfo.GetSpecs() {
		data, err := utils.ExtractFileFromZip(zip, filepath.Join(metaInfo.Name, spec))
		if err != nil {
			return nil, errors.Wrapf(err, "fail to read %s file", spec)
		}

		specInfo, err := parseProjectSpec(data)
		if err != nil {
			return nil, errors.Wrapf(err, "fail to parse %s file", spec)
		}

		specInfo.Name = spec
		specInfoList = append(specInfoList, specInfo)
	}

	// read extend files
	extFiles := map[string][]byte{}
	seConf := metaInfo.GetSelinux().GetConfig()
	if seConf != "" {
		data, err := utils.ExtractFileFromZip(zip, filepath.Join(metaInfo.Name, seConf))
		if err != nil {
			return nil, errors.Wrapf(err, "fail to read %s file", seConf)
		}

		extFiles[seConf] = data
	}

	return &pb.ProjectInfo{
		Meta:     metaInfo,
		Resource: resInfo,
		Specs:    specInfoList,
		Extends:  extFiles,
	}, nil
}

// RegularProject makes a project regular, including macro expanding and group expanding
func RegularProject(project *pb.ProjectInfo) error {
	if project == nil {
		return fmt.Errorf("nil project data")
	}

	if err := regularProjectResources(project.GetResource()); err != nil {
		return err
	}

	for _, specInfo := range project.GetSpecs() {
		if err := regularProjectSpec(specInfo, project.GetResource()); err != nil {
			return err
		}
	}

	return nil
}

// ReadProjectMetaFromDir get project meta info from dir
func ReadProjectMetaFromDir(dir string) (*pb.MetaInfo, error) {
	data, err := ioutil.ReadFile(filepath.Join(dir, metaFile))
	if err != nil {
		return nil, fmt.Errorf("fail to read %s file", metaFile)
	}

	return parseProjectMeta(data)
}

// ReadProjectMetaFromZip get project meta info from a zip file
func ReadProjectMetaFromZip(zip *zip.Reader) (*pb.MetaInfo, error) {
	if zip == nil {
		return nil, fmt.Errorf("nil zip reader")
	}

	path, err := utils.SearchFileFromZipByName(zip, metaFile)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to search %s in zip file", metaFile)
	}

	data, err := utils.ExtractFileFromZip(zip, path)
	if err != nil {
		return nil, errors.Wrapf(err, "fail to read %s file", metaFile)
	}

	meta, err := parseProjectMeta(data)
	if err != nil {
		return nil, err
	}

	if filepath.Dir(path) != meta.Name {
		return nil, fmt.Errorf("the name root directory in zip must be same as project name")
	}

	return meta, nil
}

// WriteProjectToDir export a project struct to directory
func WriteProjectToDir(info *pb.ProjectInfo, dir string) error {
	if info == nil {
		return fmt.Errorf("nil project data")
	}

	if err := marshalJSONFile(info.GetMeta(), filepath.Join(dir, metaFile)); err != nil {
		return errors.Wrapf(err, "fail to create project meta file")
	}

	path := filepath.Join(dir, info.GetMeta().GetResources())
	if err := marshalJSONFile(info.GetResource(), path); err != nil {
		return errors.Wrapf(err, "fail to create project resource.json file")
	}

	for _, spec := range info.GetMeta().GetSpecs() {
		specPath := filepath.Join(dir, spec)

		for _, specInfo := range info.GetSpecs() {
			if specInfo.GetName() != spec {
				continue
			}

			if err := marshalJSONFile(specInfo, specPath); err != nil {
				return errors.Wrapf(err, "fail to create project spec file")
			}
		}
	}

	// write extend files
	for k, v := range info.GetExtends() {
		if len(v) == 0 {
			continue
		}

		path := filepath.Join(dir, k)
		if err := utils.WriteFile(path, v, global.DefaultFilePerm); err != nil {
			log.Errorf("fail to write %s file", filepath.Base(path))
			return errors.Wrapf(err, "fail to write %s file", filepath.Base(path))
		}
	}

	return nil
}

func marshalJSONFile(obj interface{}, path string) error {
	data, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return errors.Wrap(err, "fail to marshal object")
	}

	if err := utils.WriteFile(path, data, global.DefaultFilePerm); err != nil {
		log.Errorf("fail to write %s file", filepath.Base(path))
		return fmt.Errorf("fail to write %s file", filepath.Base(path))
	}

	return nil
}
