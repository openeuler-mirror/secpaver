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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"strings"
)

const ctxMaxShaLength = 4

func autoGenFileTypeByPathAndClass(path string, class fClass, isExec bool) string {
	var tp string
	if isExec {
		tp = fmt.Sprintf("auto_%s_%s_exec_t",
			getFilenameSymbol(path),
			getFilePathHashSymbol(path))
	} else {
		tp = fmt.Sprintf("auto_%s_%s%s_t",
			getFilenameSymbol(path),
			getFilePathHashSymbol(path),
			getFileClassSymbol(class))
	}

	return strings.ToLower(tp)
}

func getFilenameSymbol(path string) string {
	if path == string(os.PathSeparator) {
		return "root"
	}

	for s := range endWildcardChangeMap {
		if strings.HasSuffix(path, s) {
			path = strings.TrimSuffix(path, s)
			break
		}
	}

	dir := filepath.Dir(path)

	filename := strings.ToLower(filepath.Base(path))
	filename = strings.Replace(filename, ".", "_", -1)
	r := regexp.MustCompile("[^\\w]")

	if strings.ContainsAny(filename, "*?") {
		filename = r.ReplaceAllString(filename, "")
		return getFilenameSymbol(dir) + "_" + strings.Trim(filename, "*?")
	}

	return r.ReplaceAllString(filename, "")
}

func getFileClassSymbol(class fClass) string {
	switch class {
	case secontext.BlkFile:
		return "_blk"
	case secontext.ChrFile:
		return "_chr"
	case secontext.DirFile:
		return "_dir"
	case secontext.FifoFile:
		return "_fifo"
	case secontext.LnkFile:
		return "_lnk"
	case secontext.SockFile:
		return "_sock"
	case secontext.ComFile:
		return "_file"
	default:
		return ""
	}
}

func getFilePathHashSymbol(path string) string {
	h := sha256.New()
	_, err := h.Write([]byte(path))
	if err != nil {
		return "0000"
	}

	return hex.EncodeToString(h.Sum(nil))[0:ctxMaxShaLength]
}

// GetTransLabel gen process type trans from an exec file
// e.g.
// xxx_exec_t -> xxx_t
// xxx_t -> xxx_trans_t
// xxx -> xxx_trans_t
func getTransProcessType(label string) string {
	if strings.HasSuffix(label, "_exec_t") {
		return strings.TrimSuffix(label, "_exec_t") + "_t"
	} else if strings.HasSuffix(label, "_t") {
		return label + "rans_t"
	} else {
		return label + "_trans_t"
	}
}

// getFileClassByName returns the sepolicy type id according to the type name
func getFileClassByName(name string) (fClass, error) {
	id, ok := fileTypeMap[name]
	if !ok {
		return secontext.UnknownFile, fmt.Errorf("invalid type %s", name)
	}
	return id, nil
}
