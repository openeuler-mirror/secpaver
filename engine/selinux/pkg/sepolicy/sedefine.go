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
	"secpaver/common/utils"
)

// SeDefine is the selinux policy define statement model
type SeDefine struct {
	RoleTypeDefine map[string][]string
	TypeAttrDefine map[string][]string
}

// NewSeDefine returns a blank SeDefine struct
func NewSeDefine() *SeDefine {
	return &SeDefine{
		RoleTypeDefine: make(map[string][]string),
		TypeAttrDefine: make(map[string][]string),
	}
}

// Text generate a string of selinux policy define statement
func (def *SeDefine) Text() string {
	var buffer bytes.Buffer
	for role, tps := range def.RoleTypeDefine {
		buffer.WriteString(genRoleTypeDefine(role, tps))
	}

	for tp, attrs := range def.TypeAttrDefine {
		buffer.WriteString(genTypeAttrDefine(tp, attrs))
	}

	return buffer.String()
}

// AddRoleTypeDefine adds a role type definition to SeDefine struct
func (def *SeDefine) AddRoleTypeDefine(role, tp string) {
	if role == "" || tp == "" {
		return
	}

	if !utils.IsExistItem(tp, def.RoleTypeDefine[role]) {
		def.RoleTypeDefine[role] = append(def.RoleTypeDefine[role], tp)
	}
}

// AddTypeDefine adds a type definition to SeDefine struct
func (def *SeDefine) AddTypeDefine(tp string) {
	if tp == "" {
		return
	}

	def.TypeAttrDefine[tp] = append(def.TypeAttrDefine[tp], []string{}...)
}

// AddTypeAttrDefine adds a type attribute definition to SeDefine struct
func (def *SeDefine) AddTypeAttrDefine(tp, attr string) {
	if tp == "" || attr == "" {
		return
	}

	if !utils.IsExistItem(attr, def.TypeAttrDefine[tp]) {
		def.TypeAttrDefine[tp] = append(def.TypeAttrDefine[tp], attr)
	}
}

func genRoleTypeDefine(role string, tps []string) string {
	var buffer bytes.Buffer
	for _, tp := range tps {
		buffer.WriteString(fmt.Sprintf("role %s types %s;\n", role, tp))
	}

	return buffer.String()
}

func genTypeAttrDefine(tp string, attrs []string) string {
	var buffer bytes.Buffer

	if len(attrs) == 0 {
		buffer.WriteString(fmt.Sprintf("type %s;\n", tp))
	} else {
		buffer.WriteString(fmt.Sprintf("type %s", tp))
		for _, attr := range attrs {
			if attr != "" {
				buffer.WriteString(fmt.Sprintf(", %s", attr))
			}
		}

		buffer.WriteString(";\n")
	}

	return buffer.String()
}
