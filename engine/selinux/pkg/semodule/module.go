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

package semodule

import (
	"fmt"
	"github.com/pkg/errors"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/libsemanage"
)

// ModuleStatus is the flag for SELinux policy module status
type ModuleStatus int

const (
	// ModuleEnabled means a module has been installed and is enabled
	ModuleEnabled ModuleStatus = 1
	// ModuleDisabled means a module has been installed but is disabled
	ModuleDisabled ModuleStatus = 2
	// ModuleNotExist means a module has not been installed
	ModuleNotExist ModuleStatus = 3
)

// Module stores the information of an SELinux policy module
type Module struct {
	Name     string
	Priority int
	Status   ModuleStatus
}

// GetModuleInfo returns a module's information by specified name
func GetModuleInfo(name string) (*Module, error) {
	if name == "" {
		return nil, fmt.Errorf("module name must be non-empty")
	}

	sh, err := libsemanage.HandleCreate()
	if err != nil {
		return nil, errors.Wrap(err, "fail to create semanage handle")
	}
	defer libsemanage.HandleDestroy(sh)

	if err := libsemanage.Connect(sh); err != nil {
		return nil, errors.Wrap(err, "fail to connect semanage handle")
	}
	defer libsemanage.Disconnect(sh)

	key, err := libsemanage.ModuleKeyCreate(sh)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create semanage module key")
	}
	defer libsemanage.ModuleKeyDestroy(sh, key)

	if err := libsemanage.ModuleKeySetName(sh, key, name); err != nil {
		return nil, errors.Wrap(err, "fail to set semanage module key name")
	}

	info, err := libsemanage.ModuleGetModuleInfo(sh, key)
	if err != nil {
		return &Module{
			Name:   name,
			Status: ModuleNotExist,
		}, nil
	}
	defer libsemanage.ModuleInfoDestroy(sh, info)

	pri, err := libsemanage.ModuleInfoGetPriority(sh, info)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get semanage module priority")
	}

	var status = ModuleDisabled
	enabled, err := libsemanage.ModuleInfoGetEnabled(sh, info)
	if err != nil {
		return nil, errors.Wrap(err, "fail to get semanage module enabled")
	}

	if enabled {
		status = ModuleEnabled
	}

	return &Module{
		Name:     name,
		Priority: pri,
		Status:   status,
	}, nil
}

// RemoveModule removes an installed SELinux policy module
func RemoveModule(name string) error {
	if name == "" {
		return fmt.Errorf("module name must be non-empty")
	}

	sh, err := libsemanage.HandleCreate()
	if err != nil {
		return errors.Wrap(err, "fail to create semanage handle")
	}
	defer libsemanage.HandleDestroy(sh)

	if err := libsemanage.Connect(sh); err != nil {
		return errors.Wrap(err, "fail to connect semanage handle")
	}
	defer libsemanage.Disconnect(sh)

	if err := libsemanage.BeginTransaction(sh); err != nil {
		return errors.Wrap(err, "fail to begin semanage transaction")
	}

	if err := libsemanage.ModuleRemove(sh, name); err != nil {
		return errors.Wrap(err, "fail to remove semanage module")
	}

	libsemanage.SetRebuild(sh)

	if err := libsemanage.Commit(sh); err != nil {
		return errors.Wrap(err, "fail to commit semanage transaction")
	}

	return nil
}

// InstallModuleFile installs a policy module
func InstallModuleFile(file string) error {
	if file == "" {
		return fmt.Errorf("module file must be non-empty")
	}

	sh, err := libsemanage.HandleCreate()
	if err != nil {
		return errors.Wrap(err, "fail to create semanage handle")
	}
	defer libsemanage.HandleDestroy(sh)

	if err := libsemanage.Connect(sh); err != nil {
		return errors.Wrap(err, "fail to connect semanage handle")
	}
	defer libsemanage.Disconnect(sh)

	if err := libsemanage.BeginTransaction(sh); err != nil {
		return errors.Wrap(err, "fail to begin semanage transaction")
	}

	if err := libsemanage.ModuleInstallFile(sh, file); err != nil {
		return errors.Wrap(err, "fail to install module file")
	}

	libsemanage.SetRebuild(sh)

	if err := libsemanage.Commit(sh); err != nil {
		return errors.Wrap(err, "fail to commit semanage transaction")
	}

	return nil
}
