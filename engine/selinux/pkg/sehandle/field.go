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

// GetRoles return all roles of policy
func (h *handleImpl) GetRoles() []string {
	var rs []string
	for _, role := range h.roles {
		rs = append(rs, role)
	}

	return rs
}

// TypeHasDefined checks if a SELinux type has been defined in system
func (h *handleImpl) TypeHasDefined(tp string) bool {
	for _, info := range h.typeInfos {
		if info.Name == tp {
			return info.IsAttr == false
		}
	}

	return false
}

// AttrHasDefined checks if a SELinux attribute has been defined in system
func (h *handleImpl) AttrHasDefined(attr string) bool {
	for _, info := range h.typeInfos {
		if info.Name == attr {
			return info.IsAttr == true
		}
	}

	return false
}
