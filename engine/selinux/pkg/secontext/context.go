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
Package secontext implements the selinux context utils
*/
package secontext

import (
	"fmt"
	"gitee.com/openeuler/secpaver/common/utils"
	"strings"
)

// Context is the selinux context model
type Context struct {
	User        string
	Role        string
	Type        string
	Sensitivity string
	Categories  string
}

// Text generate a string of selinux context
func (ctx *Context) Text() string {
	if ctx.Sensitivity == "" {
		return fmt.Sprintf("%s:%s:%s",
			ctx.User, ctx.Role, ctx.Type)
	}

	if ctx.Categories == "" {
		return fmt.Sprintf("%s:%s:%s:%s",
			ctx.User, ctx.Role, ctx.Type, ctx.Sensitivity)
	}

	return fmt.Sprintf("%s:%s:%s:%s:%s",
		ctx.User, ctx.Role, ctx.Type, ctx.Sensitivity, ctx.Categories)
}

// NewContext returns a selinux context
func NewContext(user, role, tp, sen, cate string) *Context {
	return &Context{
		User:        user,
		Role:        role,
		Type:        tp,
		Sensitivity: sen,
		Categories:  cate,
	}
}

// CreateDefaultObjectContext creates a default SELinux context with specified type
func CreateDefaultObjectContext(tp string) *Context {
	return &Context{
		User:        "system_u",
		Role:        "object_r",
		Type:        tp,
		Sensitivity: "s0",
		Categories:  "",
	}
}

// ParseContextFromLine parse context string, usually in base context file
func ParseContextFromLine(line string) (*Context, error) {
	line = utils.TrimSpaceAndTab(line)
	arr := strings.Split(line, ":")

	switch len(arr) {
	case 3: // Doesn't have Sensitivity and Categories
		return NewContext(arr[0], arr[1], arr[2], "", ""), nil
	case 4: // Maybe only have Sensitivity
		return NewContext(arr[0], arr[1], arr[2], arr[3], ""), nil
	case 5: // Have Sensitivity and Categories
		return NewContext(arr[0], arr[1], arr[2], arr[3], arr[4]), nil
	default:
		return nil, fmt.Errorf("invalid file context: %s", line)
	}
}
