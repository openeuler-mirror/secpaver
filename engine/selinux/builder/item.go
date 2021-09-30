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

package builder

import (
	"fmt"
	"github.com/pkg/errors"
	pb "secpaver/api/proto/project"
	"secpaver/common/project"
	"secpaver/common/utils"
	"secpaver/engine/selinux/pkg/secontext"
	"secpaver/engine/selinux/pkg/sehandle"
	"secpaver/engine/selinux/pkg/semodule"
)

func (f *fileItem) fileContext() *secontext.FileContext {
	return secontext.NewFileContext(f.sePath, f.class, f.context)
}

type projectContextHandle struct {
	fileItems []*fileItem
	appItems  []*applicationItem

	systemContextHandle sehandle.Handle

	maskTypes            []string
	maskFileContextPaths []string
}

type fileItem struct {
	path                  string // origin resource filepath
	sePath                string // SELinux style filepath
	class                 fClass
	context               secontext.Context
	contextInherit        secontext.Context
	execDomain            string
	isPrivate             bool
	isExecFile            bool
	hasTypeDefined        bool
	hasFileContextDefined bool
}

type applicationItem struct {
	file *fileItem

	domain string

	isPermissive bool
	isUnconfined bool
}

func (h *projectContextHandle) initProject(project *pb.ProjectInfo) error {
	// 1. set mask selinux information, maybe same module has been installed in system
	if err := h.setMask(
		[]string{project.GetMeta().GetName(), basePolicyName(project.GetMeta().GetName())}); err != nil {
		return err
	}

	// 2. init all file resource information, alloc context for every file
	if err := h.initAllFileResources(project.GetResource().GetResourceList()); err != nil {
		return err
	}

	// 3. init all applications information
	if err := h.initAllApplications(project.Specs); err != nil {
		return err
	}

	return nil
}

func (h *projectContextHandle) setMask(modules []string) error {
	for _, module := range modules {
		moduleInfo, err := semodule.GetModuleInfo(module)
		if err != nil {
			return errors.Wrap(err, "fail to get semodule info")
		}

		if moduleInfo.Status == semodule.ModuleNotExist {
			continue
		}

		cilInfo, err := semodule.ParseCilInfo(moduleInfo.Name, moduleInfo.Priority)
		if err != nil {
			return errors.Wrap(err, "fail to get semodule info from cil file")
		}

		h.maskFileContextPaths = append(h.maskFileContextPaths, cilInfo.FileConPaths...)
		h.maskTypes = append(h.maskTypes, cilInfo.ProvideTypes...)
	}

	for _, fcPath := range h.maskFileContextPaths {
		if fc := h.systemContextHandle.GetFileContext(fcPath); fc != nil {
			fc.IsValid = false
		}
	}

	return nil
}

func (h *projectContextHandle) initAllFileResources(resources []*pb.Resource) error {
	for _, resource := range resources {
		class, err := getFileClassByName(resource.GetType())
		if err != nil {
			return errors.Wrapf(err, "invalid setting of resource %s", resource.GetPath())
		}

		opt, err := parseOpt(resource.GetExtends())
		if err != nil {
			return errors.Wrap(err, "fail to parse extend option")
		}

		item := &fileItem{
			path:       resource.GetPath(),
			sePath:     getSePath(resource.GetPath()),
			class:      class,
			isExecFile: resource.GetType() == project.ExecFile,
			isPrivate:  opt.getIsPrivateFile(),
			execDomain: opt.getDomain(),
		}

		if err := h.setContextOfFileItem(item, opt); err != nil {
			return err
		}

		if !item.hasFileContextDefined {
			h.systemContextHandle.AddTempFileContext(
				secontext.NewFileContext(item.sePath, class, item.context))
		}

		h.fileItems = append(h.fileItems, item)
	}

	// second time, set inherit context for every resource
	for _, item := range h.fileItems {
		h.setInheritContextOfFileItem(item)
	}

	return nil
}

func (h *projectContextHandle) initAllApplications(specs []*pb.SpecInfo) error {
	for _, spec := range specs {
		for _, app := range spec.GetApplicationList() {
			fItem := h.getFileItemByPath(app.GetApplication().GetPath())
			if fItem == nil {
				return fmt.Errorf("undefine application %s", app.GetApplication().GetPath())
			}

			aItem := &applicationItem{
				file:         fItem,
				isPermissive: app.GetApplication().GetIsPermissive(),
				isUnconfined: app.GetApplication().GetIsUnconfined(),
			}

			if fItem.execDomain == "" {
				fItem.execDomain = getTransProcessType(fItem.context.Type)
			}

			aItem.domain = fItem.execDomain
			h.appItems = append(h.appItems, aItem)
		}
	}

	return nil
}

func (h *projectContextHandle) setContextOfFileItem(item *fileItem, opt *resourceOpt) error {
	fcLookup := h.systemContextHandle.LookupFileContext(item.path, item.class)

	if opt.getIsSysFile() { // is system file, use its default selinux type
		if opt.getIsPrivateFile() {
			return fmt.Errorf(
				"invalid setting of %s, isSysFile and isPrivate can't be true both", item.path)
		}

		if opt.getType() != "" {
			return fmt.Errorf(
				"invalid setting of %s, isSysFile and type can't be set both", item.path)
		}

		// inherit file context
		item.context = fcLookup.Context
		item.hasTypeDefined = true
		item.hasFileContextDefined = true

	} else { // create a new file context
		if fcLookup.Path == item.sePath {
			return fmt.Errorf(
				"the same resource path has been defined, must specify isSysFile flag of %s to true", item.path)
		}

		// set hasFileContextDefined to false because the context will be create
		item.hasFileContextDefined = false

		if opt.getType() != "" { // create context by specified type

			item.context = *secontext.CreateDefaultObjectContext(opt.getType())
			if utils.IsExistItem(opt.getType(), h.maskTypes) {
				item.hasTypeDefined = false
				return nil
			}

			if h.systemContextHandle.AttrHasDefined(opt.getType()) {
				return fmt.Errorf("specified type %s of %s has been defined as an attribute in system policy",
					opt.getType(), item.path)
			}

			// check if the type has been define in system
			item.hasTypeDefined = h.systemContextHandle.TypeHasDefined(opt.getType())

		} else {

			// create selinux type by filepath and class
			item.context = *secontext.CreateDefaultObjectContext(
				autoGenFileTypeByPathAndClass(item.path, item.class, item.isExecFile))
		}
	}

	return nil
}

func (h *projectContextHandle) setInheritContextOfFileItem(item *fileItem) {
	if fc := h.systemContextHandle.LookupFileContext(getDir(item.path), secontext.DirFile); fc != nil {
		item.contextInherit = fc.Context
	}
}

func (h *projectContextHandle) getFileItemByPath(path string) *fileItem {
	for _, item := range h.fileItems {
		if item.path == path {
			return item
		}
	}

	return nil
}

func (h *projectContextHandle) getAppItemByPath(path string) *applicationItem {
	for _, item := range h.appItems {
		if item.file.path == path {
			return item
		}
	}

	return nil
}
