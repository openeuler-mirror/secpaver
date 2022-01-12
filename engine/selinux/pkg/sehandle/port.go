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

import (
	"fmt"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/libsepol"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
)

// LookupPortContext searches a port context item by port and proto
func (h *handleImpl) LookupPortContext(port, proto uint32) *secontext.PortContext {
	for _, info := range h.portContexts {
		if info.HighPort >= port &&
			info.LowPort <= port &&
			info.Protocol == proto {

			return info
		}
	}

	return nil
}

func (h *handleImpl) parsePortContextsFromRecords(records []*libsepol.PortRecord) error {
	for _, record := range records {
		if record == nil {
			continue
		}

		con := &secontext.PortContext{
			LowPort:  record.LowPort,
			HighPort: record.HighPort,
		}

		switch record.Protocol {
		case libsepol.IPProtoTCP:
			con.Protocol = secontext.ProtoTCP
		case libsepol.IPProtoUDP:
			con.Protocol = secontext.ProtoUDP
		case libsepol.IPProtoSCTP:
			con.Protocol = secontext.ProtoSCTP
		case libsepol.IPProtoDCCP:
			con.Protocol = secontext.ProtoDCCP
		default:
			return fmt.Errorf("invalid port protocol ID %d", record.Protocol)
		}

		if record.Context.User > uint32(len(h.users)) || record.Context.User == 0 {
			return fmt.Errorf("invalid context user ID %d", record.Context.User)
		}

		if record.Context.Role > uint32(len(h.roles)) || record.Context.Role == 0 {
			return fmt.Errorf("invalid context role ID %d", record.Context.Role)
		}

		if record.Context.Type > uint32(len(h.typeInfos)) || record.Context.Type == 0 {
			return fmt.Errorf("invalid context type ID %d", record.Context.Type)
		}

		con.Context.User = h.users[record.Context.User-1]
		con.Context.Role = h.roles[record.Context.Role-1]
		con.Context.Type = h.typeInfos[record.Context.Type-1].Name

		h.portContexts = append(h.portContexts, con)
	}

	return nil
}
