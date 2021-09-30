## Network Permission

### 1. Network Resources

The secPaver classifies network resources into four domains: domain, type, protocol, and port. The format is as follows:

```
domain:<domain>,type:<type>,protocol:<protocol>,port:<port>
```

For details about the domain, type, and protocol keywords, see the appendix table.

### 2. Network Actions

| Action  | Description                                          |
| ------- | ---------------------------------------------------- |
| receive | Receives data from the socket.                       |
| send    | Sending data to the socket                           |
| connect | Establishing a socket connection                     |
| bind    | Binding a socket to a network node or a socket file  |
| accept  | Obtains completed connections from the socket queue. |
| listen  | Listening socket                                     |

### 3. Note

- Currently, port access control is supported only in SELinux policy generation.
- The combinations of domain, type, and protocol have certain restrictions. For invalid combinations, a warning is generated and ignored during policy generation, which does not affect the generation of other policies.

### 4. Appendix

|   domain   |   type    |                 protocol                  |
| :--------: | :-------: | :---------------------------------------: |
|    inet    |  stream   |     tcp (valid for inet/inet6 domain)     |
|    ax25    |   dgram   |     udp (valid for inet/inet6 domain)     |
|    ipx     | seqpacket |       icmp (valid for inet domain)        |
| appletalk  |    raw    |    dccp (valid for inet/inet6 domain)     |
|   netrom   |    rdm    |    sctp (valid for inet/inet6 domain)     |
|   bridge   |  packet   |     route (valid for netlink domain)      |
|   atmpvc   |           |    firewall (valid for netlink domain)    |
|    x25     |           |    tcpdiag (valid for netlink domain)     |
|   inet6    |           |     nflog (valid for netlink domain)      |
|    rose    |           |      xfrm (valid for netlink domain)      |
|  netbeui   |           |    selinux (valid for netlink domain)     |
|  security  |           |     audit (valid for netlink domain)      |
|    key     |           |     ip6fw (valid for netlink domain)      |
|    ash     |           |      dnrt (valid for netlink domain)      |
|   econet   |           | kobject_uevent (valid for netlink domain) |
|   atmsvc   |           |     iscsi (valid for netlink domain)      |
|    sna     |           |   fib_lookup (valid for netlink domain)   |
|    irda    |           |   connector (valid for netlink domain)    |
|   pppox    |           |   netfilter (valid for netlink domain)    |
|  wanpipe   |           |    generic (valid for netlink domain)     |
| bluetooth  |           | scsitransport (valid for netlink domain)  |
|  netlink   |           |      rdma (valid for netlink domain)      |
|    unix    |           |     crypto (valid for netlink domain)     |
|    rds     |           |                                           |
|    llc     |           |                                           |
|    can     |           |                                           |
|    tipc    |           |                                           |
|    iucv    |           |                                           |
|   rxrpc    |           |                                           |
|    isdn    |           |                                           |
|   phonet   |           |                                           |
| ieee802154 |           |                                           |
|    caif    |           |                                           |
|    alg     |           |                                           |
|    nfc     |           |                                           |
|   vsock    |           |                                           |
|    mpls    |           |                                           |
|     ib     |           |                                           |
|    nfc     |           |                                           |