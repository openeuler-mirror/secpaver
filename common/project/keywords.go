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

// the list of user input keywords

// the flag of rule type
const (
	RuleFileSystem = "filesystem"
	RuleCapability = "capability"
	RuleNetWork    = "network"
)

// the secPaver keywords of filesystem
const (
	ComFile     = "file"
	SockFile    = "sock_file"
	LinkFile    = "lnk_file"
	FifoFile    = "fifo_file"
	ChrFile     = "chr_file"
	BlkFile     = "blk_file"
	DirFile     = "dir"
	ExecFile    = "exec_file"
	AllFile     = "all"
	UnknownFile = ""

	ActionFileCreate  = "create"
	ActionFileRead    = "read"
	ActionFileWrite   = "write"
	ActionFileAppend  = "append"
	ActionFileRename  = "rename"
	ActionFileLink    = "link"
	ActionFileRemove  = "remove"
	ActionFileLock    = "lock"
	ActionFileMap     = "map"
	ActionFileExec    = "exec"
	ActionFileSearch  = "search"
	ActionFileIoctl   = "ioctl"
	ActionFileMounton = "mounton"

	ActionFilesystemMount = "mount"
)

// AllFileActions is selinux file action list
var AllFileActions = []string{
	ActionFileCreate,
	ActionFileRead,
	ActionFileWrite,
	ActionFileAppend,
	ActionFileRename,
	ActionFileLink,
	ActionFileRemove,
	ActionFileLock,
	ActionFileMap,
	ActionFileExec,
	ActionFileSearch,
	ActionFileIoctl,
	ActionFileMounton,
}

// the secPaver keywords of capabilities
const (
	CapAuditControl   = "audit_control"
	CapAuditRead      = "audit_read"
	CapAuditWrite     = "audit_write"
	CapAlockSuspend   = "block_suspend"
	CapChown          = "chown"
	CapDacOverride    = "dac_override"
	CapDacReadSearch  = "dac_read_search"
	CapFowner         = "fowner"
	CapFsetid         = "fsetid"
	CapIpcLock        = "ipc_lock"
	CapIpcOwner       = "ipc_owner"
	CapKill           = "kill"
	CapLease          = "lease"
	CapLinuxImmutable = "linux_immutable"
	CapMacAdmin       = "mac_admin"
	CapMacOverride    = "mac_override"
	CapMknod          = "mknod"
	CapNetAdmin       = "net_admin"
	CapNetBindService = "net_bind_service"
	CapNetBroadcast   = "net_broadcast"
	CapNetRaw         = "net_raw"
	CapSetGid         = "setgid"
	CapSetFcap        = "setfcap"
	CapSetPcap        = "setpcap"
	CapSetUID         = "setuid"
	CapSysAdmin       = "sys_admin"
	CapSysBoot        = "sys_boot"
	CapSysChroot      = "sys_chroot"
	CapSysModule      = "sys_module"
	CapSysNice        = "sys_nice"
	CapSysPacct       = "sys_pacct"
	CapSysPtrace      = "sys_ptrace"
	CapSysRawio       = "sys_rawio"
	CapSysResource    = "sys_resource"
	CapSysTime        = "sys_time"
	CapSysTtyConfig   = "sys_tty_config"
	CapSysLog         = "syslog"
	CapWakeAlarm      = "wake_alarm"
)

// the secPaver keywords of network
const (
	NetDomain   = "domain"
	NetType     = "type"
	NetProtocol = "protocol"
	NetPort     = "port"

	NetDomainInet       = "inet"
	NetDomainAx25       = "ax25"
	NetDomainIpx        = "ipx"
	NetDomainAppletalk  = "appletalk"
	NetDomainNetrom     = "netrom"
	NetDomainBridge     = "bridge"
	NetDomainAtmpvc     = "atmpvc"
	NetDomainX25        = "x25"
	NetDomainInet6      = "inet6"
	NetDomainRose       = "rose"
	NetDomainDecnet     = ""
	NetDomainNetbeui    = "netbeui"
	NetDomainSecurity   = "security"
	NetDomainKey        = "key"
	NetDomainPacket     = "packet"
	NetDomainAsh        = "ash"
	NetDomainEconet     = "econet"
	NetDomainAtmsvc     = "atmsvc"
	NetDomainSna        = "sna"
	NetDomainIrda       = "irda"
	NetDomainPppox      = "pppox"
	NetDomainWanpipe    = "wanpipe"
	NetDomainBluetooth  = "bluetooth"
	NetDomainNetlink    = "netlink"
	NetDomainUnix       = "unix"
	NetDomainRds        = "rds"
	NetDomainLlc        = "llc"
	NetDomainCan        = "can"
	NetDomainTipc       = "tipc"
	NetDomainIucv       = "iucv"
	NetDomainRxrpc      = "rxrpc"
	NetDomainIsdn       = "isdn"
	NetDomainPhonet     = "phonet"
	NetDomainIeee802154 = "ieee802154"
	NetDomainCaif       = "caif"
	NetDomainAlg        = "alg"
	NetDomainNfc        = "nfc"
	NetDomainVsock      = "vsock"
	NetDomainMpls       = "mpls"
	NetDomainIb         = "ib"
	NetDomainSmc        = "smc"

	NetTypeStream    = "stream"
	NetTypeDgram     = "dgram"
	NetTypeSeqpacket = "seqpacket"
	NetTypeRaw       = "raw"
	NetTypeRdm       = "rdm"
	NetTypePacket    = "packet"

	NetProtocolTCP  = "tcp"
	NetProtocolUDP  = "udp"
	NetProtocolICMP = "icmp"
	NetProtocolDCCP = "dccp"
	NetProtocolSCTP = "sctp"

	NetProtocolNetlinkRoute         = "route"
	NetProtocolNetlinkFirewall      = "firewall"
	NetProtocolNetlinkTcpdiag       = "tcpdiag"
	NetProtocolNetlinkNflog         = "nflog"
	NetProtocolNetlinkXfrm          = "xfrm"
	NetProtocolNetlinkSelinux       = "selinux"
	NetProtocolNetlinkAudit         = "audit"
	NetProtocolNetlinkIP6fw         = "ip6fw"
	NetProtocolNetlinkDnrt          = "dnrt"
	NetProtocolNetlinkKobjectUevent = "kobject_uevent"
	NetProtocolNetlinkIscsi         = "iscsi"
	NetProtocolNetlinkFibLookup     = "fib_lookup"
	NetProtocolNetlinkConnector     = "connector"
	NetProtocolNetlinkNetfilter     = "netfilter"
	NetProtocolNetlinkGeneric       = "generic"
	NetProtocolNetlinkScsitransport = "scsitransport"
	NetProtocolNetlinkRdma          = "rdma"
	NetProtocolNetlinkCrypto        = "crypto"

	ActionSocketRecv    = "receive"
	ActionSocketSend    = "send"
	ActionSocketConnect = "connect"
	ActionSocketBind    = "bind"
	ActionSocketAccept  = "accept"
	ActionSocketListen  = "listen"
)
