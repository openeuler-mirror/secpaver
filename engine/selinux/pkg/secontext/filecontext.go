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
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"regexp"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/sepath"
	"strings"
)

var metaSet = []rune{'.', '^', '$', '?', '*', '+', '|', '[', ']', '(', ')', '{', '}'}

// FileContext is the selinux file context model
type FileContext struct {
	Reg *regexp.Regexp

	Path        string
	Context     Context
	Class       FileClass
	Prefix      string
	HasMetaChar bool
	IsValid     bool
}

// NewFileContext creates a selinux file context
func NewFileContext(path string, classID FileClass, ctx Context) *FileContext {
	fc := &FileContext{
		Path:    path,
		Context: ctx,
		Class:   classID,
		IsValid: !(ctx == Context{}), // if a blank Context is set, means this fileContext is invalid
	}

	prefix := sepath.GetFixedPrefix(path)
	if len(prefix) != len(path) {
		fc.HasMetaChar = true
		fc.Prefix = prefix
	}

	return fc
}

// Text generate a string of selinux context
func (fc *FileContext) Text() string {
	if fc.Class == UnknownFile {
		return fmt.Sprintf("%s\t  \t%s\n", fc.Path, fc.Context.Text())
	}

	return fmt.Sprintf("%s\t%s\t%s\n",
		fc.Path, GetFileSymbolbyID(fc.Class), fc.Context.Text())
}

// Match test if the path is matched sepath regexp
func (fc *FileContext) Match(path string, class FileClass) bool {
	if !fc.IsValid {
		return false
	}

	var err error

	if fc.Reg == nil {
		fc.Reg, err = regexp.Compile(fmt.Sprintf("^%s$", fc.Path))
		if err != nil {
			fc.IsValid = false // if the regexp is invalid, set IsValid flag to false
			return false
		}
	}

	if (class == UnknownFile) ||
		(fc.Class == UnknownFile) ||
		(fc.Class == class) {

		return fc.Reg.MatchString(path)
	}

	return false
}

// NewFileContextFromString parses selinux file context from a file context string
func NewFileContextFromString(line string) (*FileContext, error) {
	var contextStr, path string
	var class FileClass
	var err error

	arr := strings.Fields(line)
	switch len(arr) {
	case 2:
		path = arr[0]
		class = UnknownFile
		contextStr = arr[1]

	case 3:
		path = arr[0]
		class, err = GetFileClassBySymbol(arr[1])
		if err != nil {
			return nil, errors.Wrap(err, "fail to parse file class")
		}
		contextStr = arr[2]

	default:
		return nil, fmt.Errorf("invalid format of file context define")
	}

	ctx, _ := ParseContextFromLine(contextStr)
	if ctx == nil {
		ctx = &Context{}
	}

	return NewFileContext(path, class, *ctx), nil
}

// ParseFileContextsFromFile parse file context file, invalid lines will be ignored
func ParseFileContextsFromFile(file string) ([]*FileContext, error) {
	var fctxs []*FileContext

	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("fail to open file_context file")
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		b, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("fail to read file_context file")
		}

		line := strings.TrimSuffix(string(b), "\n")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fctx, err := NewFileContextFromString(line)
		if err != nil {
			continue
		}

		fctxs = append(fctxs, fctx)
	}

	return fctxs, nil
}

// FileContextCompare compares two file context
func FileContextCompare(fc1, fc2 *FileContext) bool {
	if fc1 == nil || fc2 == nil {
		return false
	}

	if fc1.HasMetaChar == false && fc2.HasMetaChar == false {
		return strings.Compare(fc1.Path, fc2.Path) < 0
	}

	if fc1.HasMetaChar == false {
		return true
	}

	if fc2.HasMetaChar == false {
		return false
	}

	return len(fc1.Prefix) > len(fc2.Prefix)
}
