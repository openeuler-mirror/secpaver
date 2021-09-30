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

package errdefs

import "fmt"

type fileNotFoundError struct {
	File string
}

// Error returns the error message of fileNotFoundError
func (e *fileNotFoundError) Error() string {
	return fmt.Sprintf("%s file not found", e.File)
}

// NotFound error type
func (e *fileNotFoundError) NotFound() {}

// NewFileNotFoundError creates a fileNotFoundError error
func NewFileNotFoundError(file string) error {
	return &fileNotFoundError{
		File: file,
	}
}

type dirNotFoundError struct {
	File string
}

// Error returns the error message of dirNotFoundError
func (e *dirNotFoundError) Error() string {
	return fmt.Sprintf("%s directory not found", e.File)
}

// NotFound error type
func (e *dirNotFoundError) NotFound() {}

// NewDirNotFoundError creates a dirNotFoundError error
func NewDirNotFoundError(file string) error {
	return &dirNotFoundError{
		File: file,
	}
}
