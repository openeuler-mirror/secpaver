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
	"secpaver/engine/selinux/pkg/serule"
)

// PermsSetMap is the map change secPaver permission keywords to the selinux permission set
type PermsSetMap map[string][]string

var fileTypeMap = map[string]fClass{
	project.ComFile:     secontext.ComFile,
	project.SockFile:    secontext.SockFile,
	project.LinkFile:    secontext.LnkFile,
	project.FifoFile:    secontext.FifoFile,
	project.ChrFile:     secontext.ChrFile,
	project.BlkFile:     secontext.BlkFile,
	project.DirFile:     secontext.DirFile,
	project.ExecFile:    secontext.ComFile,
	project.AllFile:     secontext.UnknownFile,
	project.UnknownFile: secontext.UnknownFile,
}

// For ease of use, secPaver packaged some common permissions
// according to different file types

var baseFileCreatePermsSet = []string{
	"open", "getattr", "setattr", "write", "create",
}

var baseFileReadPermsSet = []string{
	"open", "getattr", "read",
}

var devFileReadPermsSet = []string{
	"open", "getattr", "read", "ioctl",
}

var baseFileWritePermsSet = []string{
	"open", "getattr", "setattr", "write", "create",
}

var devFileWritePermsSet = []string{
	"open", "getattr", "setattr", "write", "ioctl",
}

var dirFileWritePermsSet = []string{
	"open", "getattr", "setattr", "write", "create",
	"add_name", "remove_name", "rmdir", "search",
}

var dirFileReadPermsSet = []string{
	"open", "getattr", "read", "search",
}

var baseFileAppendPermsSet = []string{
	"open", "getattr", "setattr", "write", "create", "append",
}

var baseFileRenamePermsSet = []string{
	"getattr", "setattr", "rename", "unlink", "link",
}

var baseFileLinkPermsSet = []string{
	"getattr", "setattr", "rename", "unlink", "link",
}

var baseFileRemovePermsSet = []string{
	"getattr", "unlink",
}

var dirFileRemovePermsSet = []string{
	"getattr", "remove_name", "rmdir",
}

var baseFileExecPermsSet = []string{
	"execute", "open", "read", "getattr",
}

var comFileExecPermsSet = []string{
	"open", "read", "execute", "execute_no_trans",
}

var baseFileIoctlPermsSet = []string{
	"ioctl",
}

var baseFileMapPermsSet = []string{
	"map",
}

var baseFileLockPermsSet = []string{
	"lock",
}

var baseFileMountonPermsSet = []string{
	"mounton",
}

var blkFilePermsSetMap = PermsSetMap{
	project.ActionFileCreate:  baseFileCreatePermsSet,
	project.ActionFileWrite:   devFileWritePermsSet,
	project.ActionFileRead:    devFileReadPermsSet,
	project.ActionFileAppend:  baseFileAppendPermsSet,
	project.ActionFileRename:  baseFileRenamePermsSet,
	project.ActionFileLink:    baseFileLinkPermsSet,
	project.ActionFileRemove:  baseFileRemovePermsSet,
	project.ActionFileLock:    baseFileLockPermsSet,
	project.ActionFileMap:     baseFileMapPermsSet,
	project.ActionFileExec:    baseFileExecPermsSet,
	project.ActionFileIoctl:   baseFileIoctlPermsSet,
	project.ActionFileMounton: baseFileMountonPermsSet,
	project.ActionFileSearch:  nil,
}

var chrFilePermsSetMap = PermsSetMap{
	project.ActionFileCreate:  baseFileCreatePermsSet,
	project.ActionFileWrite:   devFileWritePermsSet,
	project.ActionFileRead:    devFileReadPermsSet,
	project.ActionFileAppend:  baseFileAppendPermsSet,
	project.ActionFileRename:  baseFileRenamePermsSet,
	project.ActionFileLink:    baseFileLinkPermsSet,
	project.ActionFileRemove:  baseFileRemovePermsSet,
	project.ActionFileLock:    baseFileLockPermsSet,
	project.ActionFileMap:     baseFileMapPermsSet,
	project.ActionFileExec:    baseFileExecPermsSet,
	project.ActionFileIoctl:   baseFileIoctlPermsSet,
	project.ActionFileMounton: baseFileMountonPermsSet,
	project.ActionFileSearch:  nil,
}

var fifoFilePermsSetMap = PermsSetMap{
	project.ActionFileCreate:  baseFileCreatePermsSet,
	project.ActionFileWrite:   baseFileWritePermsSet,
	project.ActionFileRead:    baseFileReadPermsSet,
	project.ActionFileAppend:  baseFileAppendPermsSet,
	project.ActionFileRename:  baseFileRenamePermsSet,
	project.ActionFileLink:    baseFileLinkPermsSet,
	project.ActionFileRemove:  baseFileRemovePermsSet,
	project.ActionFileLock:    baseFileLockPermsSet,
	project.ActionFileMap:     baseFileMapPermsSet,
	project.ActionFileExec:    baseFileExecPermsSet,
	project.ActionFileIoctl:   baseFileIoctlPermsSet,
	project.ActionFileMounton: baseFileMountonPermsSet,
	project.ActionFileSearch:  nil,
}

var linkFilePermsSetMap = PermsSetMap{
	project.ActionFileCreate:  baseFileCreatePermsSet,
	project.ActionFileWrite:   baseFileWritePermsSet,
	project.ActionFileRead:    baseFileReadPermsSet,
	project.ActionFileAppend:  baseFileAppendPermsSet,
	project.ActionFileRename:  baseFileRenamePermsSet,
	project.ActionFileLink:    baseFileLinkPermsSet,
	project.ActionFileRemove:  baseFileRemovePermsSet,
	project.ActionFileLock:    baseFileLockPermsSet,
	project.ActionFileMap:     baseFileMapPermsSet,
	project.ActionFileExec:    baseFileExecPermsSet,
	project.ActionFileIoctl:   baseFileIoctlPermsSet,
	project.ActionFileMounton: baseFileMountonPermsSet,
	project.ActionFileSearch:  nil,
}

var sockFilePermsSetMap = PermsSetMap{
	project.ActionFileCreate:  baseFileCreatePermsSet,
	project.ActionFileWrite:   baseFileWritePermsSet,
	project.ActionFileRead:    baseFileReadPermsSet,
	project.ActionFileAppend:  baseFileAppendPermsSet,
	project.ActionFileRename:  baseFileRenamePermsSet,
	project.ActionFileLink:    baseFileLinkPermsSet,
	project.ActionFileRemove:  baseFileRemovePermsSet,
	project.ActionFileLock:    baseFileLockPermsSet,
	project.ActionFileMap:     baseFileMapPermsSet,
	project.ActionFileExec:    baseFileExecPermsSet,
	project.ActionFileIoctl:   baseFileIoctlPermsSet,
	project.ActionFileMounton: baseFileMountonPermsSet,
	project.ActionFileSearch:  nil,
}

var comFilePermsSetMap = PermsSetMap{
	project.ActionFileCreate:  baseFileCreatePermsSet,
	project.ActionFileWrite:   baseFileWritePermsSet,
	project.ActionFileRead:    baseFileReadPermsSet,
	project.ActionFileAppend:  baseFileAppendPermsSet,
	project.ActionFileRename:  baseFileRenamePermsSet,
	project.ActionFileLink:    baseFileLinkPermsSet,
	project.ActionFileRemove:  baseFileRemovePermsSet,
	project.ActionFileLock:    baseFileLockPermsSet,
	project.ActionFileMap:     baseFileMapPermsSet,
	project.ActionFileExec:    comFileExecPermsSet,
	project.ActionFileIoctl:   baseFileIoctlPermsSet,
	project.ActionFileMounton: baseFileMountonPermsSet,
	project.ActionFileSearch:  nil,
}

var dirPermsSetMap = PermsSetMap{
	project.ActionFileCreate:  baseFileCreatePermsSet,
	project.ActionFileWrite:   dirFileWritePermsSet,
	project.ActionFileRead:    dirFileReadPermsSet,
	project.ActionFileAppend:  baseFileAppendPermsSet,
	project.ActionFileRename:  baseFileRenamePermsSet,
	project.ActionFileLink:    baseFileLinkPermsSet,
	project.ActionFileRemove:  dirFileRemovePermsSet,
	project.ActionFileLock:    baseFileLockPermsSet,
	project.ActionFileMap:     baseFileMapPermsSet,
	project.ActionFileExec:    baseFileExecPermsSet,
	project.ActionFileIoctl:   baseFileIoctlPermsSet,
	project.ActionFileMounton: baseFileMountonPermsSet,
	project.ActionFileSearch:  dirFileReadPermsSet,
}

var filePermsDict = map[fClass]PermsSetMap{
	secontext.BlkFile:  blkFilePermsSetMap,
	secontext.ChrFile:  chrFilePermsSetMap,
	secontext.DirFile:  dirPermsSetMap,
	secontext.FifoFile: fifoFilePermsSetMap,
	secontext.LnkFile:  linkFilePermsSetMap,
	secontext.SockFile: sockFilePermsSetMap,
	secontext.ComFile:  comFilePermsSetMap,
}

func getFilePermsByActions(classID fClass, actions []string) ([]string, error) {
	var perms []string

	permsSetMap, ok := filePermsDict[classID]
	if !ok {
		return nil, fmt.Errorf("unkown resource type id: %d", classID)
	}

	// find in permsSet
	for _, action := range actions {
		if action == project.ActionFilesystemMount { // mount is a filesystem permission
			continue
		}

		if permsSet, ok := permsSetMap[action]; ok {
			perms = append(perms, permsSet...)
			continue
		}

		perms = append(perms, action) // for other perms not list above, directly append and will be checked in sepolicy
	}

	return utils.RemoveRepeatedElement(perms), nil
}

// getFileRules creates selinux file rules
func (b *Builder) getFileRules(subject *applicationItem, perm *pb.Permission) ([]serule.Rule, error) {
	if len(perm.GetResources()) == 0 || len(perm.GetActions()) == 0 {
		return nil, fmt.Errorf("invalid file permission define")
	}

	var rules []serule.Rule
	var resourceItems []*fileItem

	for _, resource := range perm.GetResources() {
		item := b.pcHandle.getFileItemByPath(resource)
		if item == nil {
			return nil, fmt.Errorf("undefined file resource %s", resource)
		}

		resourceItems = append(resourceItems, item)
	}

	for _, item := range resourceItems {
		rs, err := getFileCommonRule(subject, item, perm.GetActions())
		if err != nil {
			return nil, errors.Wrapf(err, "fail to generate file common rule for %s", item.path)
		}

		rules = append(rules, rs...)

		if utils.IsExistItem(project.ActionFileExec, perm.GetActions()) {
			rs, err := getFileExecuteRules(subject, item)
			if err != nil {
				return nil, errors.Wrapf(err, "fail to generate file execute rule for %s", item.path)
			}
			rules = append(rules, rs...)
		}

		if utils.IsExistItem(project.ActionFileCreate, perm.GetActions()) {
			rs, err := getFileCreateRules(subject, item)
			if err != nil {
				return nil, errors.Wrapf(err, "fail to generate file create rule for %s", item.path)
			}
			rules = append(rules, rs...)
		}

		if utils.IsExistItem(project.ActionFileRemove, perm.GetActions()) {
			rs, err := getFileRemoveRule(subject, item)
			if err != nil {
				return nil, errors.Wrapf(err, "fail to generate file remove rule for %s", item.path)
			}
			rules = append(rules, rs...)
		}

		if utils.IsExistItem(project.ActionFilesystemMount, perm.GetActions()) {
			rule, err := getFileMountRule(subject, item)
			if err != nil {
				return nil, errors.Wrap(err, "fail to generate filesystem mount avc rule")
			}
			rules = append(rules, rule)
		}
	}

	return rules, nil
}

func getFileCommonRule(subject *applicationItem, file *fileItem, actions []string) ([]serule.Rule, error) {
	var rules []serule.Rule

	if file.class == secontext.UnknownFile {
		for _, class := range secontext.FileClassSet {
			perms, err := getFilePermsByActions(class, actions)
			if err != nil {
				return nil, errors.Wrap(err, "fail to parse file actions")
			}

			rs, err := serule.CreateFileAllowRules(subject.domain, file.context.Type, class, perms)
			if err != nil {
				return nil, errors.Wrap(err, "fail to generate file avc rule")
			}

			rules = append(rules, rs...)
		}
	} else {
		perms, err := getFilePermsByActions(file.class, actions)
		if err != nil {
			return nil, errors.Wrap(err, "fail to parse file actions")
		}

		rs, err := serule.CreateFileAllowRules(subject.domain, file.context.Type, file.class, perms)
		if err != nil {
			return nil, errors.Wrap(err, "fail to generate file avc rule")
		}

		rules = append(rules, rs...)
	}

	return rules, nil
}

func getFileExecuteRules(sub *applicationItem, file *fileItem) ([]serule.Rule, error) {
	var rules []serule.Rule
	domain := file.execDomain

	if domain == "" {
		domain = "unconfined_t" // TODO
	}

	rs, err := serule.CreateDomainAutoTransRule(sub.domain, file.context.Type, domain)
	if err != nil {
		return nil, err
	}

	return append(rules, rs...), nil
}

func getFileCreateRules(subject *applicationItem, file *fileItem) ([]serule.Rule, error) {
	var rules []serule.Rule

	// add the write permission for dir
	avcRules, err := serule.CreateFileAllowRules(
		subject.domain, file.contextInherit.Type, secontext.DirFile, dirFileWritePermsSet)
	if err != nil {
		return nil, err
	}

	rules = append(rules, avcRules...)

	// check if the type_transition rules should be add
	if file.context.Type == file.contextInherit.Type {
		return rules, nil
	}

	typeRules, err := serule.CreateFileTypeTransitionRule(
		subject.domain, file.contextInherit.Type, file.class, file.context.Type, getBase(file.path))

	return append(rules, typeRules...), nil
}

func getFileRemoveRule(subject *applicationItem, file *fileItem) ([]serule.Rule, error) {
	// add the write permission for dir
	return serule.CreateFileAllowRules(
		subject.domain, file.contextInherit.Type, secontext.DirFile, dirFileWritePermsSet)
}

func getFileMountRule(subject *applicationItem, file *fileItem) (serule.Rule, error) {
	mountonPerms := []string{"mount", "remount", "unmount", "getattr"}

	return serule.CreateFilesystemAllowRule(
		subject.domain, file.context.Type, mountonPerms)
}
