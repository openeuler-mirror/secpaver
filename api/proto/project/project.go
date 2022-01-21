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

package project

import (
	"gitee.com/openeuler/secpaver/common/utils"
)

type resourceAlias Resource

// UnmarshalJSON is the overwrite method for json unmarshal
func (m *Resource) UnmarshalJSON(b []byte) error {
	var err error
	m.Extends, err = utils.UnmarshalJSONWithExtendParams(b, (*resourceAlias)(m))
	return err
}

// MarshalJSON is the method for json marshal
func (m *Resource) MarshalJSON() ([]byte, error) {
	return utils.MarshalJSONWithExtendParams((*resourceAlias)(m), m.GetExtends(), "extends")
}
