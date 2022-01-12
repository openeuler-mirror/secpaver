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
	"fmt"
	"github.com/pkg/errors"
	pb "gitee.com/openeuler/secpaver/api/proto/project"
	"gitee.com/openeuler/secpaver/common/project"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/serule"
)

var capPermsDict = map[string]string{
	project.CapAuditControl:   "audit_control",
	project.CapAuditRead:      "audit_read",
	project.CapAuditWrite:     "audit_write",
	project.CapAlockSuspend:   "block_suspend",
	project.CapChown:          "chown",
	project.CapDacOverride:    "dac_override",
	project.CapDacReadSearch:  "dac_read_search",
	project.CapFowner:         "fowner",
	project.CapFsetid:         "fsetid",
	project.CapIpcLock:        "ipc_lock",
	project.CapIpcOwner:       "ipc_owner",
	project.CapKill:           "kill",
	project.CapLease:          "lease",
	project.CapLinuxImmutable: "linux_immutable",
	project.CapMacAdmin:       "mac_admin",
	project.CapMacOverride:    "mac_override",
	project.CapMknod:          "mknod",
	project.CapNetAdmin:       "net_admin",
	project.CapNetBindService: "net_bind_service",
	project.CapNetBroadcast:   "net_broadcast",
	project.CapNetRaw:         "net_raw",
	project.CapSetGid:         "setgid",
	project.CapSetFcap:        "setfcap",
	project.CapSetPcap:        "setpcap",
	project.CapSetUID:         "setuid",
	project.CapSysAdmin:       "sys_admin",
	project.CapSysBoot:        "sys_boot",
	project.CapSysChroot:      "sys_chroot",
	project.CapSysModule:      "sys_module",
	project.CapSysNice:        "sys_nice",
	project.CapSysPacct:       "sys_pacct",
	project.CapSysPtrace:      "sys_ptrace",
	project.CapSysRawio:       "sys_rawio",
	project.CapSysResource:    "sys_resource",
	project.CapSysTime:        "sys_time",
	project.CapSysTtyConfig:   "sys_tty_config",
	project.CapSysLog:         "syslog",
	project.CapWakeAlarm:      "wake_alarm",
}

// getCapabilityByActions returns the selinux capabilities corresponding to the user input
func getCapabilityByActions(actions []string) ([]string, error) {
	var perms []string

	for _, action := range actions {
		if perm, ok := capPermsDict[action]; ok {
			perms = append(perms, perm)
		} else {
			return nil, fmt.Errorf("invalid capability %s", action)
		}
	}

	return utils.RemoveRepeatedElement(perms), nil
}

// genCapabilityRules creates selinux capability rules
func (b *Builder) getCapabilityRules(subject *applicationItem, perm *pb.Permission) ([]serule.Rule, error) {
	if len(perm.GetActions()) == 0 {
		return nil, fmt.Errorf("invalid capability permission define")
	}

	caps, err := getCapabilityByActions(perm.GetActions())
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse capability actions")
	}

	rules, _ := serule.CreateCapabilityAllowRule(subject.domain, caps)

	return append(rules, createProcessCapabilityRules(subject)...), nil
}

func createProcessCapabilityRules(subject *applicationItem) []serule.Rule {
	var rules []serule.Rule
	rs, _ := serule.CreateProcessAllowRule(
		subject.domain, subject.domain, []string{"setcap", "getcap"})

	return append(rules, rs...)
}
