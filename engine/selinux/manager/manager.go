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
Package manager implements the SELinux policy manager function.
*/
package manager

import (
	"fmt"
	"github.com/pkg/errors"
	"path/filepath"
	pb "gitee.com/openeuler/secpaver/api/proto/policy"
	"gitee.com/openeuler/secpaver/common/ack"
	"gitee.com/openeuler/secpaver/common/log"
	"gitee.com/openeuler/secpaver/common/policy"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/domain"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/semodule"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/sepath"
)

const systemPolicyPriority = 100

// Manager is the SELinux policy manager model
type Manager struct{}

// NewManager returns a new instance of policy manager
func NewManager() *Manager {
	return &Manager{}
}

// GetPolicyStatus returns the status of a SELinux policy
func (m *Manager) GetPolicyStatus(p *domain.Policy) (string, error) {
	if p == nil {
		return "", fmt.Errorf("nil policy request")
	}

	module, err := semodule.GetModuleInfo(p.Name)
	if err != nil {
		return "", errors.Wrap(err, "fail to get module info")
	}

	if (module.Status == semodule.ModuleDisabled) ||
		(module.Status == semodule.ModuleNotExist) {

		return policy.StatusDisable, nil
	}

	return policy.StatusActive, nil
}

// Install installs a SELinux policy to system
func (m *Manager) Install(p *domain.Policy, msg chan *pb.Ack) error {
	if p == nil {
		return fmt.Errorf("nil policy request")
	}

	sendMsg(msg, ack.LevelInfo, "install SELinux policy module")

	module, err := semodule.GetModuleInfo(p.Name)
	if err != nil {
		return errors.Wrapf(err, "fail to get %s policy module information", p.Name)
	}

	if module.Priority == systemPolicyPriority {
		return fmt.Errorf("%s is a system policy module, can't be modified", p.Name)
	}

	oldInfo, _ := semodule.ParseCilInfo(module.Name, module.Priority)

	policyFile, err := getAndCheckPpFilepath(p)
	if err != nil {
		return errors.Wrap(err, "fail to check policy binary file")
	}

	err = semodule.InstallModuleFile(policyFile)
	if err != nil {
		return errors.Wrap(err, "fail to install policy module")
	}

	sendMsg(msg, ack.LevelInfo, "start to restore file context")

	newInfo, err := semodule.GetModuleCilInfo(p.Name)
	if err != nil {
		return err
	}

	allPaths := newInfo.FileConPaths
	if oldInfo != nil {
		allPaths = utils.RemoveRepeatedElement(append(allPaths, oldInfo.FileConPaths...))
	}

	for _, path := range allPaths {
		restoreContextWithRegPath(path)
	}

	sendMsg(msg, ack.LevelInfo, "Finish installing policy")

	return nil
}

// Uninstall uninstalls a SELinux policy in system
func (m *Manager) Uninstall(p *domain.Policy, msg chan *pb.Ack) error {
	sendMsg(msg, ack.LevelInfo, "uninstall SELinux policy module")

	module, err := semodule.GetModuleInfo(p.Name)
	if err != nil {
		return errors.Wrapf(err, "fail to get %s policy module information", p.Name)
	}

	if module.Priority == systemPolicyPriority {
		return fmt.Errorf("%s is a system policy module, can't be modified", p.Name)
	}

	if module.Status != semodule.ModuleEnabled {
		return fmt.Errorf("%s module is not installed", p.Name)
	}

	info, err := semodule.ParseCilInfo(module.Name, module.Priority)
	if err != nil {
		return errors.Wrapf(err, "fail to parse %s policy cil module information", p.Name)
	}

	if err := semodule.RemoveModule(p.Name); err != nil {
		return errors.Wrap(err, "fail to remove module")
	}

	sendMsg(msg, ack.LevelInfo, "restore file context")

	for _, path := range info.FileConPaths {
		restoreContextWithRegPath(path)
	}

	sendMsg(msg, ack.LevelInfo, "Finish uninstalling policy uninstalling")

	return nil
}

func getAndCheckPpFilepath(policy *domain.Policy) (string, error) {
	path := filepath.Join(policy.Path, policy.Name+".pp")
	uid, err := utils.GetUIDOfFile(path)
	if err != nil {
		log.Errorf("fail to get uid of policy file %s", filepath.Base(path))
		return "", fmt.Errorf("fail to get uid of policy file %s", filepath.Base(path))
	}

	if uid != 0 {
		return "", fmt.Errorf("the uid of policy file %s must be 0", filepath.Base(path))
	}

	return path, nil
}

func sendMsg(ch chan *pb.Ack, level, msg string) {
	if ch == nil {
		return
	}

	ch <- &pb.Ack{
		Level:  level,
		Status: msg,
	}
}

func restoreContextWithRegPath(path string) {
	fPath := path

	prefix := sepath.GetFixedPrefix(path)
	if path != prefix {
		if prefix[len(prefix)-1] == '/' {
			fPath = prefix
		} else {
			fPath = prefix + "*"
		}
	}

	log.Debugf("restore context of %s", fPath)
	if r, _ := secontext.RestoreconPath(fPath, true); r != "" {
		log.Debugf("restorecon log: %s", r)
	}
}
