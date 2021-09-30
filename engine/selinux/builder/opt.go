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

package builder

import (
	"encoding/json"
	"fmt"
	"secpaver/common/utils"
)

type resourceOpt struct {
	IsSysFile     bool
	IsPrivateFile bool
	Type          string
	Domain        string
}

func parseOpt(opts map[string][]byte) (*resourceOpt, error) {
	o := &resourceOpt{}
	minLen := 2

	for k, v := range opts {
		if k == "selinux" && len(v) > minLen {
			if err := json.Unmarshal(v, o); err != nil {
				return nil, err
			}

			return o, nil
		}
	}

	if err := checkOpt(o); err != nil {
		return nil, err
	}

	return nil, nil
}

func checkOpt(opt *resourceOpt) error {
	if opt == nil {
		return fmt.Errorf("nil data")
	}

	if opt.getDomain() != "" {
		if err := utils.CheckValidSelinuxType(opt.getDomain()); err != nil {
			return err
		}
	}

	if opt.getType() != "" {
		if err := utils.CheckValidSelinuxType(opt.getType()); err != nil {
			return err
		}
	}

	return nil
}

func (m *resourceOpt) getIsSysFile() bool {
	if m != nil {
		return m.IsSysFile
	}
	return false
}

func (m *resourceOpt) getIsPrivateFile() bool {
	if m != nil {
		return m.IsPrivateFile
	}
	return false
}

func (m *resourceOpt) getType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *resourceOpt) getDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}
