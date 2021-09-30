## capability关键字说明

用户可使用的的capability关键字如下表：

|    capability    |                             说明                             |
| :--------------: | :----------------------------------------------------------: |
|  audit_control   | 启用和禁用内核审计；修改审计过滤器规则；提取审计状态和过滤规则。 |
|    audit_read    |         允许通过一个多播netlink socket读取审计日志。         |
|   audit_write    |                  允许向内核审计日志写记录。                  |
|  block_suspend   |                      允许阻塞系统挂起。                      |
|      chown       |                 允许修改文件的uids和gids 。                  |
|   dac_override   |         绕过文件的dac权限检查（读、写和执行权限）。          |
| dac_read_search  | 绕过文件和目录的dac权限检查（文件的读权限，目录的读和执行权限）。 |
|      fowner      |            绕过对文件uid和进程uid匹配的权限检查。            |
|      fsetid      |                 允许设置文件的 setuid 标志。                 |
|     ipc_lock     |                    允许锁定共享内存片段。                    |
|    ipc_owner     |                     忽略ipc所有权检查。                      |
|       kill       |               允许对不属于自己的进程发送信号。               |
|      lease       |                允许修改文件锁的fl_lease标志。                |
| linux_immutable  |          允许修改文件的immutable和append属性标志。           |
|    mac_admin     |                   允许设置或更改mac配置。                    |
|   mac_override   |                      允许覆盖mac配置。                       |
|      mknod       |                   允许使用mknod创建文件。                    |
|    net_admin     |                    允许执行网络管理任务。                    |
| net_bind_service |                  允许绑定到小于1024的端口。                  |
|  net_broadcast   |               允许进行套接字广播，并收听多播。               |
|     net_raw      |                     允许使用原始套接字。                     |
|      setgid      |                     允许改变进程的gid。                      |
|     setfcap      |              允许为文件设置任意的capabilities。              |
|     setpcap      |  允许向其他进程转移capability以及删除其他进程的capability。  |
|      setuid      |                      允许改变进程的uid                       |
|    sys_admin     |                    允许执行系统管理任务。                    |
|     sys_boot     |                       允许使用reboot。                       |
|    sys_chroot    |                   允许使用chroot系统调用。                   |
|    sys_module    |                   允许加载和卸载内核模块。                   |
|     sys_nice     |                    允许设置进程的优先级。                    |
|    sys_pacct     |                  允许执行进程的bsd式审计。                   |
|    sys_ptrace    |                        允许跟踪进程。                        |
|    sys_rawio     |    允许直接访问/devport,/dev/mem,/dev/kmem及原始块设备。     |
|   sys_resource   |                        忽略资源限制。                        |
|     sys_time     |                      允许改变系统时钟。                      |
|  sys_tty_config  |                    允许设置tty终端设备。                     |
|      syslog      |                   允许使用syslog系统调用。                   |
|    wake_alarm    |                      允许触发唤醒系统。                      |