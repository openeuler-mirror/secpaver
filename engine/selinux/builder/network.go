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
	"gitee.com/openeuler/secpaver/common/log"
	"gitee.com/openeuler/secpaver/common/project"
	"gitee.com/openeuler/secpaver/common/utils"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/secontext"
	"gitee.com/openeuler/secpaver/engine/selinux/pkg/serule"
)

// the socket can specified by domain
var commonDomainSocketMap = map[string]string{
	project.NetDomainAx25:       "ax25_socket",
	project.NetDomainIpx:        "ipx_socket",
	project.NetDomainAppletalk:  "appletalk_socket",
	project.NetDomainNetrom:     "netrom_socket",
	project.NetDomainBridge:     "bridge_socket",
	project.NetDomainAtmpvc:     "atmpvc_socket",
	project.NetDomainX25:        "x25_socket",
	project.NetDomainRose:       "rose_socket",
	project.NetDomainDecnet:     "decnet_socket",
	project.NetDomainNetbeui:    "",
	project.NetDomainSecurity:   "",
	project.NetDomainKey:        "",
	project.NetDomainPacket:     "packet_socket",
	project.NetDomainAsh:        "",
	project.NetDomainEconet:     "",
	project.NetDomainAtmsvc:     "atmsvc_socket",
	project.NetDomainSna:        "",
	project.NetDomainIrda:       "irda_socket",
	project.NetDomainPppox:      "pppox_socket",
	project.NetDomainWanpipe:    "",
	project.NetDomainBluetooth:  "bluetooth_socket",
	project.NetDomainRds:        "rds_socket",
	project.NetDomainLlc:        "llc_socket",
	project.NetDomainCan:        "can_socket",
	project.NetDomainTipc:       "tipc_socket",
	project.NetDomainIucv:       "iucv_socket",
	project.NetDomainRxrpc:      "rxrpc_socket",
	project.NetDomainIsdn:       "isdn_socket",
	project.NetDomainPhonet:     "phonet_socket",
	project.NetDomainIeee802154: "ieee802154_socket",
	project.NetDomainCaif:       "caif_socket",
	project.NetDomainAlg:        "alg_socket",
	project.NetDomainNfc:        "nfc_socket",
	project.NetDomainVsock:      "vsock_socket",
	project.NetDomainMpls:       "mpls_socket",
	project.NetDomainIb:         "ib_socket",
	project.NetDomainSmc:        "smc_socket",
}

var inetSocketSet = []string{
	"tcp_socket",
	"udp_socket",
	"dccp_socket",
	"sctp_socket",
	"icmp_socket",
	"rawip_socket",
}

var inetTypeSocketMap = map[string][]string{
	project.NetTypeStream:    {"tcp_socket", "sctp_socket"},
	project.NetTypeDgram:     {"udp_socket"},
	project.NetTypeSeqpacket: {"sctp_socket"},
	project.NetTypeRaw:       {"rawip_socket", "icmp_socket", "sctp_socket"},
	project.NetTypeRdm:       {},
	project.NetTypePacket:    {},
}

var inetProtocolSocketMap = map[string][]string{
	project.NetProtocolTCP:  {"tcp_socket"},
	project.NetProtocolUDP:  {"udp_socket"},
	project.NetProtocolDCCP: {"dccp_socket"},
	project.NetProtocolSCTP: {"sctp_socket", "rawip_socket"},
	project.NetProtocolICMP: {"icmp_socket", "rawip_socket"},
}

var netlinkSocketSet = []string{
	"netlink_route_socket",
	"netlink_firewall_socket",
	"netlink_tcpdiag_socket",
	"netlink_nflog_socket",
	"netlink_selinux_socket",
	"netlink_audit_socket",
	"netlink_ip6fw_socket",
	"netlink_dnrt_socket",
	"netlink_kobject_uevent_socket",
	"netlink_iscsi_socket",
	"netlink_fib_lookup_socket",
	"netlink_connector_socket",
	"netlink_netfilter_socket",
	"netlink_generic_socket",
	"netlink_scsitransport_socket",
	"netlink_rdma_socket",
	"netlink_crypto_socket",
	"netlink_socket",
}

var netlinkProtocolSocketMap = map[string]string{
	project.NetProtocolNetlinkRoute:         "netlink_route_socket",
	project.NetProtocolNetlinkFirewall:      "netlink_firewall_socket",
	project.NetProtocolNetlinkTcpdiag:       "netlink_tcpdiag_socket",
	project.NetProtocolNetlinkNflog:         "netlink_nflog_socket",
	project.NetProtocolNetlinkSelinux:       "netlink_selinux_socket",
	project.NetProtocolNetlinkAudit:         "netlink_audit_socket",
	project.NetProtocolNetlinkIP6fw:         "netlink_ip6fw_socket",
	project.NetProtocolNetlinkDnrt:          "netlink_dnrt_socket",
	project.NetProtocolNetlinkKobjectUevent: "netlink_kobject_uevent_socket",
	project.NetProtocolNetlinkIscsi:         "netlink_iscsi_socket",
	project.NetProtocolNetlinkFibLookup:     "netlink_fib_lookup_socket",
	project.NetProtocolNetlinkConnector:     "netlink_connector_socket",
	project.NetProtocolNetlinkNetfilter:     "netlink_netfilter_socket",
	project.NetProtocolNetlinkGeneric:       "netlink_generic_socket",
	project.NetProtocolNetlinkScsitransport: "netlink_scsitransport_socket",
	project.NetProtocolNetlinkRdma:          "netlink_rdma_socket",
	project.NetProtocolNetlinkCrypto:        "netlink_crypto_socket",
}

var unixSocketSet = []string{
	"unix_stream_socket",
	"unix_dgram_socket",
}

var unixTypeSocketMap = map[string][]string{
	project.NetTypeStream: {"unix_stream_socket"},
	project.NetTypeDgram:  {"unix_dgram_socket"},
}

func getInetSocketClasses(tp, protocol string) []string {
	if tp == "" && protocol == "" { // return all inet sockets
		return inetSocketSet
	}

	if protocol == "" { // only specify the socket type
		return inetTypeSocketMap[tp]
	}

	return inetProtocolSocketMap[protocol]
}

func getNetlinkSocketClasses(protocol string) []string {
	if protocol == "" {
		return netlinkSocketSet
	}

	return []string{
		netlinkProtocolSocketMap[protocol],
	}
}

func getUnixSocketClasses(tp string) []string {
	if tp == "" {
		return unixSocketSet
	}

	return unixTypeSocketMap[tp]
}

func getCommonSocketClasses(domain string) []string {
	s := commonDomainSocketMap[domain]
	if s == "" {
		return nil
	}

	return []string{s}
}

func getSocketClasses(info *project.NetInfo) []string {
	if info.Domain != "" {
		switch info.Domain {
		case project.NetDomainInet, project.NetDomainInet6:
			return getInetSocketClasses(info.Type, info.Protocol)

		case project.NetDomainNetlink:
			return getNetlinkSocketClasses(info.Protocol)

		case project.NetDomainUnix:
			return getUnixSocketClasses(info.Type)

		default:
			return getCommonSocketClasses(info.Domain)
		}
	}

	if info.Type != "" {
		var sockets []string
		sockets = append(sockets, getInetSocketClasses(info.Type, info.Protocol)...)
		sockets = append(sockets, getUnixSocketClasses(info.Type)...)

		return sockets
	}

	if info.Protocol != "" {
		var sockets []string
		sockets = append(sockets, getInetSocketClasses(info.Type, info.Protocol)...)
		sockets = append(sockets, getNetlinkSocketClasses(info.Protocol)...)

		return sockets
	}

	return nil
}

var netifPermMap = map[string]map[string][]string{
	"tcp_socket": {
		project.ActionSocketRecv: {"egress", "ingress", "tcp_recv"},
		project.ActionSocketSend: {"egress", "ingress", "tcp_send"},
	},
	"udp_socket": {
		project.ActionSocketRecv: {"egress", "ingress", "udp_recv"},
		project.ActionSocketSend: {"egress", "ingress", "udp_send"},
	},
	"dccp_socket": {
		project.ActionSocketRecv: {"egress", "ingress", "dccp_recv"},
		project.ActionSocketSend: {"egress", "ingress", "dccp_send"},
	},
	"rawip_socket": {
		project.ActionSocketRecv: {"egress", "ingress", "rawip_recv"},
		project.ActionSocketSend: {"egress", "ingress", "rawip_send"},
	},
}

func getNetInterfaceRule(subject, socket string, actions []string) (serule.Rule, error) {
	var perms []string

	pMap, ok := netifPermMap[socket]
	if !ok {
		return nil, nil
	}

	for _, action := range actions {
		if p, ok := pMap[action]; ok {
			perms = append(perms, p...)
		}
	}

	if len(perms) != 0 {
		return serule.CreateNetifAllowRule(subject, "netif_type", perms)
	}

	return nil, nil
}

var nodePermMap = map[string]map[string][]string{
	"tcp_socket": {
		project.ActionSocketRecv: {"enforce_dest", "recvfrom", "tcp_recv"},
		project.ActionSocketSend: {"enforce_dest", "sendto", "tcp_send"},
	},
	"udp_socket": {
		project.ActionSocketRecv: {"enforce_dest", "recvfrom", "udp_recv"},
		project.ActionSocketSend: {"enforce_dest", "sendto", "udp_send"},
	},
	"dccp_socket": {
		project.ActionSocketRecv: {"enforce_dest", "recvfrom", "dccp_recv"},
		project.ActionSocketSend: {"enforce_dest", "sendto", "dccp_send"},
	},
	"rawip_socket": {
		project.ActionSocketRecv: {"enforce_dest", "recvfrom", "rawip_recv"},
		project.ActionSocketSend: {"enforce_dest", "sendto", "rawip_send"},
	},
}

func genNetNodeRule(sub, socket string, actions []string) ([]serule.Rule, error) {
	var rules []serule.Rule
	var perms []string

	pMap, ok := nodePermMap[socket]
	if ok {
		for _, action := range actions {
			if p, ok := pMap[action]; ok {
				perms = append(perms, p...)
			}
		}
	}

	if len(perms) != 0 {
		rule, err := serule.CreateNodeAllowRule(sub, "node_type", perms)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	if utils.IsExistItem(socket, inetSocketSet) {
		rule, err := serule.CreateCommonAllowRule(sub, "node_type", socket, []string{"node_bind"})
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	return rules, nil
}

var portPermMap = map[string]map[string][]string{
	"tcp_socket": {
		project.ActionSocketBind:    {"name_bind", "node_bind"},
		project.ActionSocketConnect: {"name_connect"},
		project.ActionSocketRecv:    {"recv_msg"},
		project.ActionSocketSend:    {"send_msg"},
	},
	"dccp_socket": {
		project.ActionSocketBind:    {"name_bind", "node_bind"},
		project.ActionSocketConnect: {"name_connect"},
		project.ActionSocketRecv:    {"recv_msg"},
		project.ActionSocketSend:    {"send_msg"},
	},
	"sctp_socket": {
		project.ActionSocketBind:    {"name_bind", "node_bind"},
		project.ActionSocketConnect: {"name_connect"},
		project.ActionSocketRecv:    {"recv_msg"},
		project.ActionSocketSend:    {"send_msg"},
	},
	"udp_socket": {
		project.ActionSocketBind: {"name_bind", "node_bind"},
		project.ActionSocketRecv: {"recv_msg"},
		project.ActionSocketSend: {"send_msg"},
	},
	"rawip_socket": {
		project.ActionSocketBind: {"name_bind", "node_bind"},
		project.ActionSocketRecv: {"recv_msg"},
		project.ActionSocketSend: {"send_msg"},
	},
	"icmp_socket": {
		project.ActionSocketBind: {"name_bind", "node_bind"},
		project.ActionSocketRecv: {"recv_msg"},
		project.ActionSocketSend: {"send_msg"},
	},
}

func getNetPortRule(subject, portType, socket string, actions []string) (serule.Rule, error) {
	var perms []string

	pMap, ok := portPermMap[socket]
	if !ok {
		return nil, nil
	}

	for _, action := range actions {
		if p, ok := pMap[action]; ok {
			perms = append(perms, p...)
		}
	}

	if len(perms) != 0 {
		return serule.CreateCommonAllowRule(subject, portType, socket, perms)
	}

	return nil, nil
}

func getNetProcessRule(subject, socket string, actions []string) (serule.Rule, error) {
	if socket == "unix_stream_socket" {
		return serule.CreateCommonAllowRule(
			subject, "domain", socket, []string{"connectto"})
	}

	return nil, nil
}

var baseSocketPerms = []string{
	"create", "lock", "read", "write", "getattr", "setattr",
	"getopt", "setopt", "append", "shutdown", "ioctl",
}

var socketActionSet = []string{
	project.ActionSocketRecv,
	project.ActionSocketSend,
	project.ActionSocketConnect,
	project.ActionSocketBind,
	project.ActionSocketAccept,
	project.ActionSocketListen,
}

var commonSocketPermsSetMap = map[string][]string{
	project.ActionSocketAccept:  {"accept"},
	project.ActionSocketConnect: {"connect"},
	project.ActionSocketBind:    {"bind"},
	project.ActionSocketRecv:    {"recvfrom"},
	project.ActionSocketSend:    {"sendto"},
	project.ActionSocketListen:  {"listen"},
}

var extendNetlinkSocketSet = []string{
	"netlink_route_socket",
	"netlink_firewall_socket",
	"netlink_tcpdiag_socket",
	"netlink_xfrm_socket",
	"netlink_audit_socket",
	"netlink_ip6fw_socket",
}

var extendNetlinkSocketPermsSetMap = map[string][]string{
	project.ActionSocketAccept:  {"accept"},
	project.ActionSocketConnect: {"connect"},
	project.ActionSocketBind:    {"bind"},
	project.ActionSocketRecv:    {"recvfrom", "nlmsg_read"},
	project.ActionSocketSend:    {"sendto", "nlmsg_read"},
	project.ActionSocketListen:  {"listen"},
}

func getNetSocketRule(subject, socket string, actions []string) (serule.Rule, error) {
	var perms []string
	var permMap map[string][]string

	if utils.IsExistItem(socket, extendNetlinkSocketSet) {
		permMap = extendNetlinkSocketPermsSetMap
	} else {
		permMap = commonSocketPermsSetMap
	}

	for _, act := range actions {
		if p, ok := permMap[act]; ok {
			perms = append(perms, p...)
		}
	}

	perms = append(perms, baseSocketPerms...)

	if len(perms) == 0 {
		return nil, nil
	}

	return serule.CreateCommonAllowRule(subject, subject, socket, perms)
}

var portProtoMap = map[string]uint32{
	"tcp_socket":  secontext.ProtoTCP,
	"udp_socket":  secontext.ProtoUDP,
	"sctp_socket": secontext.ProtoSCTP,
	"dccp_socket": secontext.ProtoDCCP,
}

func (b *Builder) getNetworkRules(subject *applicationItem, perm *pb.Permission) ([]serule.Rule, error) {
	if len(perm.GetResources()) == 0 || len(perm.GetActions()) == 0 {
		return nil, fmt.Errorf("invalid network permission define")
	}

	for _, act := range perm.GetActions() {
		if !utils.IsExistItem(act, socketActionSet) {
			return nil, fmt.Errorf("invalid socket action %s", act)
		}
	}

	var rules []serule.Rule
	var sockets []string

	for _, res := range perm.GetResources() {
		info, err := project.ParseNetwork(res)
		if err != nil {
			return nil, errors.Wrap(err, "fail to parse network resource")
		}

		// get socket classes by network resource
		s := getSocketClasses(info)
		if len(s) == 0 {
			log.Warnf("invalid network resource define %s", res)
			continue
		}

		for _, socket := range s {
			pt := b.getPortType(info.Port, socket)
			if pt == "" {
				continue
			}

			rule, err := getNetPortRule(subject.domain, pt, socket, perm.GetActions())
			if err != nil {
				return nil, errors.Wrap(err, "fail to gen net port rule")
			}

			rules = append(rules, rule)
		}

		sockets = append(sockets, s...)
	}

	for _, socket := range sockets {
		if socket == "" {
			continue
		}

		socketRules, err := getAllSocketRules(subject.domain, socket, perm.GetActions())
		if err != nil {
			return nil, err
		}

		rules = append(rules, socketRules...)
	}

	return rules, nil
}

func getAllSocketRules(subject, socket string, actions []string) ([]serule.Rule, error) {
	var rules []serule.Rule

	if utils.IsExistItem(socket, inetSocketSet) {
		rules = append(rules, createDefaultInetRules(subject)...)
	}

	rule, err := getNetSocketRule(subject, socket, actions)
	if err != nil {
		return nil, errors.Wrap(err, "fail to gen net socket rule")
	}
	rules = append(rules, rule)

	rule, err = getNetInterfaceRule(subject, socket, actions)
	if err != nil {
		return nil, errors.Wrap(err, "fail to gen net interface rule")
	}
	rules = append(rules, rule)

	rs, err := genNetNodeRule(subject, socket, actions)
	if err != nil {
		return nil, errors.Wrap(err, "fail to gen net node rule")
	}
	rules = append(rules, rs...)

	rule, err = getNetProcessRule(subject, socket, actions)
	if err != nil {
		return nil, errors.Wrap(err, "fail to gen net process rule")
	}
	rules = append(rules, rule)

	return rules, nil
}

func (b *Builder) getPortType(port uint32, socket string) string {
	proto, ok := portProtoMap[socket]
	if !ok {
		return ""
	}

	pt := "port_type"

	if port != 0 {
		pc := b.pcHandle.systemContextHandle.LookupPortContext(port, proto)
		if pc != nil {
			pt = pc.Context.Type
		}
	}

	return pt
}

func createDefaultInetRules(subject string) []serule.Rule {
	var rules []serule.Rule
	r1, _ := serule.CreateFileAllowRules(
		subject, "sysctl_net_t", secontext.DirFile, []string{"search"})

	r2, _ := serule.CreateFileAllowRules(
		subject, "sysctl_net_t", secontext.ComFile, []string{"getattr", "read", "open"})

	rules = append(rules, r1...)
	rules = append(rules, r2...)

	return rules
}
