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

package sepolicy

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/serule"
)

// SeRequire is the selinux policy require statement model
type SeRequire struct {
	TypeRequires  []string
	RoleRequires  []string
	AttrRequires  []string
	ClassRequires map[string][]string
}

// NewSeRequire returns a blank selinux require statement model
func NewSeRequire() *SeRequire {
	return &SeRequire{
		ClassRequires: make(map[string][]string),
	}
}

// AddTypeRequire add a type require statement to SeRequire struct
func (req *SeRequire) AddTypeRequire(tp string) {
	if tp == "" {
		return
	}

	if !utils.IsExistItem(tp, req.TypeRequires) {
		req.TypeRequires = append(req.TypeRequires, tp)
	}
}

// AddRoleRequire add a role require statement to SeRequire struct
func (req *SeRequire) AddRoleRequire(role string) {
	if role == "" {
		return
	}

	if !utils.IsExistItem(role, req.RoleRequires) {
		req.RoleRequires = append(req.RoleRequires, role)
	}
}

// AddAttrRequire add a attribute require statement to SeRequire struct
func (req *SeRequire) AddAttrRequire(attr string) {
	if attr == "" {
		return
	}

	if !utils.IsExistItem(attr, req.AttrRequires) {
		req.AttrRequires = append(req.AttrRequires, attr)
	}
}

// RemoveTypeRequire remove a type from type require list
func (req *SeRequire) RemoveTypeRequire(tp string) {
	for i := range req.TypeRequires {
		if i >= len(req.TypeRequires) {
			return
		}

		if req.TypeRequires[i] == tp {
			if i < len(req.TypeRequires)-1 {
				req.TypeRequires = append(req.TypeRequires[:i], req.TypeRequires[i+1:]...)
			} else {
				req.TypeRequires = req.TypeRequires[:i]
			}
		}
	}
}

// AddClassRequire add a class require statement to SeRequire struct
func (req *SeRequire) AddClassRequire(cls string, actions []string) error {
	if cls == "" {
		return nil
	}

	if len(actions) == 0 ||
		utils.IsExistItem("*", actions) {

		perms, err := serule.GetAllPermissionsOfClass(cls)
		if err != nil {
			return errors.Wrapf(err, "fail to get permissions of %s class", cls)
		}

		req.ClassRequires[cls] = perms

		return nil
	}

	if _, ok := req.ClassRequires[cls]; !ok {
		req.ClassRequires[cls] = actions
		return nil
	}

	for _, act := range actions {
		if !utils.IsExistItem(act, req.ClassRequires[cls]) {
			req.ClassRequires[cls] = append(req.ClassRequires[cls], act)
		}
	}

	return nil
}

// Text generate a string of selinux require statement
func (req *SeRequire) Text() string {
	var buffer bytes.Buffer
	buffer.WriteString("require\n{\n")

	for _, role := range req.RoleRequires {
		buffer.WriteString(genRoleRequireStr(role))
	}

	for _, tp := range req.TypeRequires {
		buffer.WriteString(genTypeRequireStr(tp))
	}

	for _, attr := range req.AttrRequires {
		buffer.WriteString(genAttrRequireStr(attr))
	}

	buffer.WriteString(genClassRequireStr(req.ClassRequires))

	buffer.WriteString("};\n")
	return buffer.String()
}

func genRoleRequireStr(role string) string {
	return fmt.Sprintf("\trole %s;\n", role)
}

func genTypeRequireStr(tp string) string {
	return fmt.Sprintf("\ttype %s;\n", tp)
}

func genAttrRequireStr(attr string) string {
	return fmt.Sprintf("\tattribute %s;\n", attr)
}

func genClassRequireStr(cls map[string][]string) string {
	var buffer bytes.Buffer

	for cls, acts := range cls {
		if len(acts) != 0 {
			buffer.WriteString(fmt.Sprintf("\tclass %s { %s };\n",
				cls, utils.ShowStringsWithSpace(acts)))
		}
	}

	return buffer.String()
}
