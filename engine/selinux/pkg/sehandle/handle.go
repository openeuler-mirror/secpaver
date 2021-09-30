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
Package sehandle implements system policy context handle
*/
package sehandle

import (
	"fmt"
	"github.com/pkg/errors"
	"path/filepath"
	"secpaver/engine/selinux/pkg/libselinux"
	"secpaver/engine/selinux/pkg/libsepol"
	"secpaver/engine/selinux/pkg/secontext"
)

const (
	fileContextRoot = "/etc/selinux/targeted/contexts/files"
)

// Handle the the interface for selinux context manager
type Handle interface {
	GetFileContext(path string) *secontext.FileContext
	LookupFileContext(path string, class secontext.FileClass) *secontext.FileContext
	AddTempFileContext(ctx *secontext.FileContext)
	TypeHasDefined(tp string) bool
	AttrHasDefined(at string) bool
	LookupPortContext(port, proto uint32) *secontext.PortContext
	GetRoles() []string
}

type handleImpl struct {
	fileContexts        []*secontext.FileContext
	fileContextsHomedir []*secontext.FileContext
	fileContextsLocal   []*secontext.FileContext
	// note: fileContextsTemp stores the temporary file context record
	// during policy developing
	fileContextsTemp []*secontext.FileContext

	portContexts []*secontext.PortContext

	typeInfos []*libsepol.TypeInfo

	users []string
	roles []string
}

// HandleCreate creates a SELinux file context handel
func HandleCreate() (Handle, error) {
	h := &handleImpl{}

	if err := h.readPolicyDb(); err != nil {
		return nil, err
	}

	if err := h.readContextFiles(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *handleImpl) readContextFiles() error {
	var err error
	var path string

	if path = libselinux.FileContextPath(); path == "" {
		path = filepath.Join(fileContextRoot, "file_contexts")
	}

	h.fileContexts, err = secontext.ParseFileContextsFromFile(path)
	if err != nil {
		return errors.Wrap(err, "fail to parse file_contexts file")
	}

	if path = libselinux.FileContextHomedirPath(); path == "" {
		path = filepath.Join(fileContextRoot, "file_contexts.homedirs")
	}

	h.fileContextsHomedir, err = secontext.ParseFileContextsFromFile(path)
	if err != nil {
		return errors.Wrap(err, "fail to parse file_contexts.homedirs file")
	}

	if path = libselinux.FileContextLocalPath(); path == "" {
		path = filepath.Join(fileContextRoot, "file_contexts.local")
	}

	h.fileContextsLocal, _ = secontext.ParseFileContextsFromFile(path)

	return nil
}

func (h *handleImpl) readPolicyDb() error {
	handle, err := libsepol.HandleCreate()
	if err != nil {
		return err
	}
	defer libsepol.HandleDestroy(handle)

	spf, err := libsepol.PolicyFileCreate()
	if err != nil {
		return err
	}
	defer libsepol.PolicyFileFree(spf)

	policyFile := libselinux.CurrentPolicyPath()
	if policyFile == "" {
		return fmt.Errorf("fail to get current policy path, please check policy file")
	}

	fp, err := libsepol.Fopen(policyFile)
	if err != nil {
		return err
	}
	defer libsepol.Fclose(fp)

	libsepol.PolicyFileSetFp(spf, fp)
	libsepol.PolicyFileSetHandle(spf, handle)

	db, err := libsepol.PolicydbCreate()
	if err != nil {
		return err
	}
	defer libsepol.PolicydbFree(db)

	if err := libsepol.PolicydbRead(db, spf); err != nil {
		return err
	}

	if h.users, err = libsepol.GetAllUsers(db); err != nil {
		return err
	}

	if h.roles, err = libsepol.GetAllRoles(db); err != nil {
		return err
	}

	if h.typeInfos, err = libsepol.GetAllTypesAndAttrs(db); err != nil {
		return err
	}

	ports, err := libsepol.GetAllPorts(db)
	if err != nil {
		return err
	}

	if err := h.parsePortContextsFromRecords(ports); err != nil {
		return errors.Wrap(err, "fail to parse port context from db record")
	}

	return nil
}
