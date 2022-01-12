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

package sehandle

import (
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"sort"
)

type fcLookupFunc func(fc *secontext.FileContext) bool

// GetFileContext searches a file context item by path exactly match
func (h *handleImpl) GetFileContext(path string) *secontext.FileContext {
	f := func(fc *secontext.FileContext) bool {
		return fc.Path == path
	}

	return h.lookup(f)
}

// Lookup searches a file context item by path and class
func (h *handleImpl) LookupFileContext(path string, class secontext.FileClass) *secontext.FileContext {
	sort.Slice(h.fileContextsTemp, func(i, j int) bool {
		return secontext.FileContextCompare(h.fileContextsTemp[i], h.fileContextsTemp[j])
	})

	f := func(fc *secontext.FileContext) bool {
		return fc.Match(path, class)
	}

	return h.lookup(f)
}

// AddTempFileContext adds a temporary file context
func (h *handleImpl) AddTempFileContext(ctx *secontext.FileContext) {
	h.fileContextsTemp = append(h.fileContextsTemp, ctx)
}

func (h *handleImpl) lookup(f fcLookupFunc) *secontext.FileContext {
	for i := len(h.fileContextsTemp) - 1; i >= 0; i-- {
		if f(h.fileContextsTemp[i]) {
			return h.fileContextsTemp[i]
		}
	}

	for i := len(h.fileContextsLocal) - 1; i >= 0; i-- {
		if f(h.fileContextsLocal[i]) {
			return h.fileContextsLocal[i]
		}
	}

	for i := len(h.fileContextsHomedir) - 1; i >= 0; i-- {
		if f(h.fileContextsHomedir[i]) {
			return h.fileContextsHomedir[i]
		}
	}

	for i := len(h.fileContexts) - 1; i >= 0; i-- {
		if f(h.fileContexts[i]) {
			return h.fileContexts[i]
		}
	}

	return nil
}
