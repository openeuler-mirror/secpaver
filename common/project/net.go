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
	"fmt"
	"github.com/pkg/errors"
	"secpaver/common/utils"
	"strings"
)

const netParamMax = 4    // domain type protocol port
const netParamFields = 2 // key value

// NetInfo is the net resource struct
type NetInfo struct {
	Domain   string
	Type     string
	Protocol string
	Port     uint32
}

// ParseNetwork parse user net resource input as 'domain:XXX, type:XXX, protol:XXX'
func ParseNetwork(line string) (*NetInfo, error) {
	info := &NetInfo{}
	attrs := strings.Split(utils.TrimSpaceAndTab(line), ",")
	if len(attrs) == 0 || len(attrs) > netParamMax {
		return nil, fmt.Errorf("invalid parameter number of network resource %s", line)
	}

	for _, attr := range attrs {
		if attr == "" {
			continue
		}

		items := strings.Split(attr, ":")
		if len(items) != netParamFields {
			return nil, fmt.Errorf("invalid network resource %s", line)
		}
		items[0] = strings.ToLower(strings.TrimSpace(items[0]))
		items[1] = strings.ToLower(strings.TrimSpace(items[1]))

		switch items[0] {
		case NetDomain:
			info.Domain = items[1]

		case NetType:
			info.Type = items[1]

		case NetProtocol:
			info.Protocol = items[1]

		case NetPort:
			if items[1] == "" {
				continue
			}

			port, err := utils.CheckPort(items[1])
			if err != nil {
				return nil, errors.Wrap(err, "invalid port param")
			}

			info.Port = uint32(port)

		default:
			return nil, fmt.Errorf("invalid network field %s in %s", items[0], line)
		}
	}

	if info.Type == "" && info.Protocol == "" && info.Domain == "" {
		return nil, fmt.Errorf("invalid network resource %s, all network fields are blank", line)
	}

	return info, nil
}
