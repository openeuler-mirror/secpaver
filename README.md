# secPaver

## 1 概述
secPaver是一个帮助管理员定义应用程序安全策略的工具。用户可以使用secPaver为应用程序生成不同安全机制下的策略文件。

目前secPaver支持的安全机制为SELinux，AppArmor即将支持，不同的安全机制生成模块以插件（plugin）形式加载。

## 2 安装secPaver
**(1) 安装依赖软件包**

编译secPaver需要的软件有：
* make
* golang 1.11+

编译SELinux策略插件需要的软件有：
* libselinux-devel 2.9+
* libsepol-devel 2.9+
* libsemanage-devel 2.9+

运行SELinux策略插件需要的软件有：
* libselinux 2.9+
* libsepol 2.9+
* libsemanage 2.9+
* checkpolicy 2.8+
* policycoreutils 2.8+

**(2) 下载源码**
```
git clone https://gitee.com/openeuler/secpaver.git
```

**(3) 编译安装**
```
cd secpaver
make
```
编译SELinux策略插件：
```
make selinux
```
安装软件（至少一个策略插件需要被编译）：
```
make install
```

## 3 术语名词
**engine**：secPaver生成不同安全机制下的策略文件的功能模块，以插件的形式进行加载。

## 4 如何使用
secPaver为客户端/服务端架构，服务端程序为pavd，客户端程序为pav。

**(1) 启动服务端进程**
```
systemctl start pavd
```

**(2) 查看加载的策略生成插件**
```
# pav engine list

Name        Description                  
selinux     SELinux policy generator
```

**(3) 生成工程模板并在工程中配置策略**
```
pav project create my_demo .
```

**(4) 编译工程生成策略文件**
```
pav project build -d ./my_demo --engine selinux
```

**(5) 查看策略状态**

```
# pav policy list

Name                           Status
my_demo_selinux                disable
```

**(6) 加载策略**
```
# pav policy install my_demo_selinux
[info]: install SELinux policy module
[info]: start to restore file context
[info]: Finish installing policy
```

**(7) 卸载策略**
```
# pav policy uninstall my_demo_selinux
[info]: uninstall SELinux policy module
[info]: restore file context
[info]: Finish uninstalling policy uninstalling
```

**(8) 导出策略包**
```
# pav policy export my_demo_selinux .
Finish exporting: export_my_demo_selinux.zip
```

## 5 文档资料

命令行手册：[secPaver命令行手册](doc/cmd.md)

用户手册：[secPaver用户手册](doc/manual.md)

## 6 如何贡献

我们非常欢迎新贡献者加入到项目中来，也非常高兴能为新加入贡献者提供指导和帮助。在您贡献代码前，需要先签署[CLA](https://openeuler.org/en/cla.html)。
