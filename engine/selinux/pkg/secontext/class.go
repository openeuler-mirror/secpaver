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

package secontext

import (
	"fmt"
	"gitee.com/openeuler/secpaver/common/utils"
)

// FileClass is the file type flag
type FileClass uint32

// file class flag
const (
	UnknownFile FileClass = 0 // all file types
	BlkFile     FileClass = 1 // block device file
	ChrFile     FileClass = 2 // character device file
	DirFile     FileClass = 3 // directory file
	FifoFile    FileClass = 4 // pipe file
	LnkFile     FileClass = 5 // link file
	SockFile    FileClass = 6 // socket file
	ComFile     FileClass = 7 // common file
)

// FileClassSet is the set of all valid file class
var FileClassSet = []FileClass{BlkFile, ChrFile, DirFile, FifoFile, LnkFile, SockFile, ComFile}

type fileClassDefine struct {
	id      FileClass
	symbol  string // for gen file context
	seClass string // for gen rule
}

// the dict of selinux file type
var fileClassDict = map[FileClass]fileClassDefine{
	UnknownFile: {
		id:      UnknownFile,
		symbol:  "",
		seClass: "",
	},
	BlkFile: {
		id:      BlkFile,
		symbol:  "-b",
		seClass: "blk_file",
	},
	ChrFile: {
		id:      ChrFile,
		symbol:  "-c",
		seClass: "chr_file",
	},
	DirFile: {
		id:      DirFile,
		symbol:  "-d",
		seClass: "dir",
	},
	FifoFile: {
		id:      FifoFile,
		symbol:  "-p",
		seClass: "fifo_file",
	},
	LnkFile: {
		id:      LnkFile,
		symbol:  "-l",
		seClass: "lnk_file",
	},
	SockFile: {
		id:      SockFile,
		symbol:  "-s",
		seClass: "sock_file",
	},
	ComFile: {
		id:      ComFile,
		symbol:  "--",
		seClass: "file",
	},
}

// GetFileClassByID finds the selinux class corresponding to type id
func GetFileClassByID(classID FileClass) []string {
	if !utils.IsExistItem(classID, FileClassSet) {
		classID = UnknownFile
	}

	if classID == UnknownFile {
		return []string{"file", "chr_file", "blk_file", "dir", "fifo_file", "lnk_file", "sock_file"}
	}

	return []string{fileClassDict[classID].seClass}
}

// GetFileSymbolbyID finds the selinux file type symbol corresponding to type id
func GetFileSymbolbyID(classID FileClass) string {
	if !utils.IsExistItem(classID, FileClassSet) {
		classID = UnknownFile
	}

	return fileClassDict[classID].symbol
}

// GetFileClassBySymbol finds the type id corresponding to selinux file type symbol
func GetFileClassBySymbol(symbol string) (FileClass, error) {
	for id, typeDef := range fileClassDict {
		if typeDef.symbol == symbol {

			return id, nil
		}
	}

	return UnknownFile, fmt.Errorf("unkown resource type symbol: %s", symbol)
}
