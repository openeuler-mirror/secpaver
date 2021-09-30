# secPaver命令行手册

[TOC]

## 前言

本文档介绍secPaver客户端程序pav和服务端程序pavd的命令行使用。

## 1 概述

pav是secPaver的命令行客户端程序，主要负责使用命令行的方式和服务端进行信息交互。

pavd是secPaver的服务端程序，主要负责管理服务端主机存储的secPaver资源，同时处理客户端发送的请求，并反馈处理结果。

## 2 命令格式说明

secPaver的命令行格式为：
```
pav [options] OBJECT SUMCOMMAND [SUBCOMMAND_ARGUMENTS]
pavd [options]
```

其中options表示全局参数，OBJECT表示命令行作用的管理对象，SUMCOMMAND表示管理对象支持的子命令，SUBCOMMAND_ARGUMENTS表示子命令参数。

* pav中的全局参数可以和管理类命令同时使用，但全局参数必须放置在管理类命令前，如：
```
pav -s /var/run/secpaver/pavd.sock project list
```

* 命令格式中，[] 表示参数可选，<> 表示参数必选。

## 3 pav命令

### 3.1 全局参数

pav程序的全局参数包含帮助信息、版本信息查询选项以及一些连接配置。

#### 3.1.1 --help, -h

**功能描述：**

打印pav命令行帮助信息。

**命令格式：**

```
pav --help
```

#### 3.1.2	--version, -v

**功能描述：**

打印pav版本号。

**命令格式：**

```
pav --version
```

#### 3.1.3	--socket, -s
**命令描述：**

指定服务端grpc连接使用的unix socket文件路径，默认为/var/run/secpaver/pavd.sock。

**命令格式：**

```
pav --socket <PATH>
```

### 3.2 engine管理子命令

**命令描述：**

管理服务端的策略生成插件。

**命令格式：**

```
pav engine SUBCOMMAND [SUBCOMMAND_ARGUMENTS]
```

支持的SUBCOMMAND如下：

#### 3.2.1	list

**命令描述：**

列出加载的策略生成插件。

**命令格式：**

```
pav engine list
```

**命令参数：**

无

**使用示例：**

```
# pav engine list

Name        Desc                               
selinux     selinux policy generate engine     
```

#### 3.2.2	info

**命令描述：**

列出服务端策略生成插件的详细信息。

**命令格式：**

```
pav engine info <ENGINE>
```

**命令参数：**

|  参数  |      说明      |
| :----: | :------------: |
| ENGINE | 指定策略生成插件名称 |

**使用示例：**

```
# pav engine info selinux

Attribute     Value                             
Name          selinux                           
Desc          selinux policy generate engine    
```

### 3.3 project管理子命令

**命令描述：**

管理服务端工程文件。

**命令格式：**

```
pav project SUBCOMMAND [SUBCOMMAND_ARGUMENTS]
```

支持的SUBCOMMAND如下：

#### 3.3.1	create

**命令描述：**

生成模板工程。

**命令格式：**

```
pav project create <NAME> <PATH>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| NAME | 工程名称 |
| PATH | 生成路径，需为有效目录 |

**使用示例：**

```
# pav project create demo /root
create demo template project at /root/demo
# tree /root/demo/
/root/demo/
├── resources.json
├── selinux.json
├── specs
│   └── module_demo.json
└── pav.proj

1 directory, 4 files
```

#### 3.3.2	list

**命令描述：**

列出服务端管理的工程。

**命令格式：**

```
pav project list
```

**命令参数：**

无

**使用示例：**

```
# pav project list

Name       Version  
demo       1.0
```

#### 3.3.3	info

**命令描述：**

列出服务端管理的指定工程的详细信息

**命令格式：**

```
pav project info <PROJECT>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| PROJECT | 服务端管理的工程名称 |

**使用示例：**
```
# pav project info demo

Attribute        Value                     
name             demo                      
resource file    resources.json            
spec files       specs/module_demo.json
```

#### 3.3.4	import

**命令描述：**

导入工程文件，工程文件必须是zip格式的。

**命令格式：**
```
pav project import [-f] <FILE>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| --force, -f | 若服务端存在同名的工程，允许覆盖导入 |
| FILE | 工程zip文件 |

**使用示例：**
```
# pav project import demo.zip 
[info]: Finish importing demo project
```

#### 3.3.5	export

**命令描述：**

从服务端导出指定的工程文件，导出结果是一个zip文件。

**命令格式：**
```
pav project export [-f] <PROJECT> <PATH>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| --force, -f | 若导出路径存在同名的文件，允许覆盖导出 |
| PROJECT | 服务端的工程名称 |
| PATH | 导出路径，需为有效目录 |

**使用示例：**

```
# pav project export demo .
Finish exporting: export_demo.zip
```

#### 3.3.6	delete

**命令描述：**

从服务端删除指定的工程。

**命令格式：**
```
pav project delete <PROJECT>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| PROJECT | 服务端的工程名称 |

**使用示例：**
```
# pav project delete demo
[info]: Finish deleting demo project
```

#### 3.3.7	build

**命令描述：**

编译指定工程，生成并导出指定安全机制的策略文件，被编译的工程可以是一个本地工程，也可以是服务端保存的工程。

**命令格式：**
```
pav project build --engine <ENGINE> <-d PATH | -r PROJECT>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| --engine ENGINE | 指定生成策略文件的安全机制 |
| -d PATH | 本地工程所在的目录 |
| -r PROJECT | 服务端管理的工程名 |

**使用示例：**

```
# pav project build -d ./demo --engine selinux
[info]: Finish building demo project

# pav project build -r demo --engine selinux
[info]: Finish building demo project
```

### 3.4 policy管理子命令

**命令描述：**

管理服务端生成的策略。

**命令格式：**

```
pav policy SUBCOMMAND
```

支持的SUBCOMMAND如下：

#### 3.4.1	list

**命令描述：**

列出服务的生成的策略。

**命令格式：**
```
pav policy list
```

**命令参数：**

无

**使用示例：**

```
# pav policy list

Name                       Status     
demo_selinux               disable
```

#### 3.4.2	install

**命令描述：**

安装并加载生成的策略。

**命令格式：**

```
pav policy install <POLICY>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| POLICY | 指定的服务端策略名称 |

**使用示例：**

```
# pav policy install demo_selinux
[info]: install SELinux policy module
[info]: start to restore file context
[info]: Finish installing policy
```

#### 3.4.3	uninstall

**命令描述：**

从系统中卸载生成的策略。

**命令格式：**

```
pav policy uninstall <POLICY>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| POLICY | 指定的服务端策略名称 |

**使用示例：**

```
# pav policy uninstall demo_selinux
[info]: uninstall SELinux policy module
[info]: restore file context
[info]: Finish uninstalling policy uninstalling
```

#### 3.4.4	export

**命令描述：**

导出服务端生成的策略文件，以zip格式导出。

**命令格式：**
```
pav policy export [-f] <POLICY> <PATH>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| --force, -f | 若导出路径存在同名的文件，允许覆盖导出 |
| POLICY | 指定的服务端策略名称 |
| PATH | 导出路径，需为有效目 |

**使用示例：**

```
# pav policy export demo_selinux .
Finish exporting: export_demo_selinux.zip
```

#### 3.4.5	delete

**命令描述：**

删除指定的服务端策略。

**命令格式：**
```
pav policy delete <POLICY>
```

**命令参数：**

| 参数 | 说明 |
| :---: | :---: |
| POLICY | 指定的服务端策略名称 |

**使用示例：**
```
# pav policy delete demo_selinux
[info]: Finish deleting policy
```

## 4 pavd程序命令

### 4.1 全局命令

pavd命令中的全局命令包含帮助信息、版本信息查询选项以及一些基本配置。

#### 4.1.1 --help, -h

**功能描述：**

打印pavd命令行帮助信息。

**命令格式：**

```
pavd --help
```

#### 4.1.2	--version, -v

**功能描述：**

打印pavd版本号。

**命令格式：**

```
pavd --version
```

#### 4.1.3	--config, -c

**命令描述：**

指定使用的配置文件路径，默认为/etc/secpaver/pavd/config.json。

**命令格式：**

```
pavd --config <FILE>
```

#### 4.1.4	--socket, -s
**命令描述：**

客户端grpc连接服务端所使用的socket文件路径，默认为/var/run/secpaver/pavd.sock。

**命令格式：**

```
pavd --socket <PATH>
```

#### 4.1.5	--log-level, -l

**命令描述：**

设置日志等级。日志等级可以是debug, info, warn, error, fatal, panic。日志等级可以在配置文件中指定；若配置文件未指定，则默认为info。

**命令格式：**

```
pavd --log-level <LEVEL>
```

### 4.2 systemd命令

**命令描述：**

secPaver安装后，可以通过systemd命令对pavd服务进程管理。

**命令格式：**

启动pavd服务：

```
systemctl start pavd
```

停止pavd服务：

```
systemctl stop pavd
```

查询pavd服务状态：

```
systemctl status pavd
```