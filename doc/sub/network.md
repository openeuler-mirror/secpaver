## 网络资源权限说明

### 1. 网络资源

secPaver将网络资源划分为domain、type、protocol和port进行管理，定义的格式为：

```
domain:<domain>,type:<type>,protocol:<protocol>,port:<port>
```

支持的domain、type和protocol关键字见附表。

### 2. 网络权限

secPaver支持的网络资源权限为：

| 权限    | 说明                             |
| ------- | -------------------------------- |
| receive | 从socket中接收数据               |
| send    | 向socket发送数据                 |
| connect | 建立socket连接                   |
| bind    | 将socket绑定网络节点或socket文件 |
| accept  | 从socket队列中取出已完成的连接   |
| listen  | 收听socket                       |

### 3. 注意事项

- 目前端口的访问控制仅在SELinux策略生成中支持；
- domain、type、protocol的组合存在一定的约束限制，对于无效的组合，生成策略时会告警并忽略，不影响其他策略生成。

### 4. 附表

|   domain   |   type    |              protocol               |
| :--------: | :-------: | :---------------------------------: |
|    inet    |  stream   |     tcp(domain为inet/inet6有效)     |
|    ax25    |   dgram   |     udp(domain为inet/inet6有效)     |
|    ipx     | seqpacket |       icmp(domain为inet有效)        |
| appletalk  |    raw    |    dccp(domain为inet/inet6有效)     |
|   netrom   |    rdm    |    sctp(domain为inet/inet6有效)     |
|   bridge   |  packet   |     route(domain为netlink有效)      |
|   atmpvc   |           |    firewall(domain为netlink有效)    |
|    x25     |           |    tcpdiag(domain为netlink有效)     |
|   inet6    |           |     nflog(domain为netlink有效)      |
|    rose    |           |      xfrm(domain为netlink有效)      |
|  netbeui   |           |    selinux(domain为netlink有效)     |
|  security  |           |     audit(domain为netlink有效)      |
|    key     |           |     ip6fw(domain为netlink有效)      |
|    ash     |           |      dnrt(domain为netlink有效)      |
|   econet   |           | kobject_uevent(domain为netlink有效) |
|   atmsvc   |           |     iscsi(domain为netlink有效)      |
|    sna     |           |   fib_lookup(domain为netlink有效)   |
|    irda    |           |   connector(domain为netlink有效)    |
|   pppox    |           |   netfilter(domain为netlink有效)    |
|  wanpipe   |           |    generic(domain为netlink有效)     |
| bluetooth  |           | scsitransport(domain为netlink有效)  |
|  netlink   |           |      rdma(domain为netlink有效)      |
|    unix    |           |     crypto(domain为netlink有效)     |
|    rds     |           |                                     |
|    llc     |           |                                     |
|    can     |           |                                     |
|    tipc    |           |                                     |
|    iucv    |           |                                     |
|   rxrpc    |           |                                     |
|    isdn    |           |                                     |
|   phonet   |           |                                     |
| ieee802154 |           |                                     |
|    caif    |           |                                     |
|    alg     |           |                                     |
|    nfc     |           |                                     |
|   vsock    |           |                                     |
|    mpls    |           |                                     |
|     ib     |           |                                     |
|    nfc     |           |                                     |