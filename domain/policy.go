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
Package domain is the domain layer define.
*/
package domain

import (
	"path/filepath"
)

// Policy is the secPaver policy model
type Policy struct {
	Name   string
	Path   string
	Engine string
}

// NewPolicy is the constructor of Policy
func NewPolicy(path string) *Policy {
	return &Policy{
		Name:   filepath.Base(path),
		Engine: filepath.Base(filepath.Dir(path)), // now the parent dir of policy is the engine dir
		Path:   path,
	}
}
