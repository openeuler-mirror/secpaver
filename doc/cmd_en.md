# secPaver Command Manual

[TOC]

## Preface

This document describes how to use the commands in secPaver package.

## 1 Introduction

pav is the client process of secPaver. It provides a command-line interface for users to use secPaver.

pavd is the server process of secPaver that manages the data resources, processes client requests, and returns processing results.

## 2 Command Format

The format of secPaver command is:
```
pav [options] OBJECT SUMCOMMAND [SUBCOMMAND_ARGUMENTS]
pavd [options]
```

where the options indicate global parameters, OBJECT indicates the managed object of the command, SUMCOMMAND indicates the subcommand supported by the managed object, and SUBCOMMAND_ARGUMENTS indicates the subcommand parameter.

* Global parameters of pav can be used for management subcommands, but global parameters must be placed before management subcommands. For example:
```
pav -s /var/run/secpaver/pavd.sock project list
```

* `[PARM]` indicates an optional parameter, and `<PARM>` indicates a required parameter.

## 3 pav Command

### 3.1 Global Parameters

The global parameters of pav include help information, version query, and connection configurations.

#### 3.1.1 --help, -h

**Description:**

Print the help information。

**Format:**

```
pav --help
```

#### 3.1.2	--version, -v

**Description:**

Print pav version.

**Format:**

```
pav --version
```

#### 3.1.3	--socket, -s
**Description:**

Specify the Unix socket file used by grpc connection. The default value is /var/run/secpaver/pavd.sock.

**Format:**

```
pav --socket <PATH>
```

### 3.2 Engine Manage Command

**Description:**

Manage the policy generator plugin.

**Format:**

```
pav engine SUBCOMMAND
```

Supported subcommands:

#### 3.2.1	list

**Description:**

List all loaded policy generator plugins.

**Format:**

```
pav engine list
```

**Parameters:**

None

**Example:**

```
# pav engine list

Name        Desc                               
selinux     selinux policy generate engine     
```

#### 3.2.2	info

**Description:**

List the details of one loaded policy generator plugin.

**Format:**

```
pav engine info <ENGINE>
```

**Parameters:**

|  Parameter  |      Description      |
| :----: | :------------: |
| ENGINE | Name of loaded policy generator plugin |

**Example:**

```
# pav engine info selinux

Attribute     Value                             
Name          selinux                           
Desc          selinux policy generate engine    
```

### 3.3 Project Manage Command

**Description:**

Manage the projects stored in secPaver server.

**Format:**

```
pav project SUBCOMMAND
```

Supported subcommands:

#### 3.3.1	create

**Description:**

Create a template project.

**Format:**

```
pav project create <NAME> <PATH>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| NAME | Project name |
| PATH | Local path for the created project  |

**Example:**

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

**Description:**

List all projects stored in secPaver server.

**Format:**

```
pav project list
```

**Parameters:**

None

**Example:**

```
# pav project list

Name       Version  
demo       1.0
```

#### 3.3.3	info

**Description:**

List details of a specified project stored in secPaver server.

**Format:**

```
pav project info <PROJECT>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| PROJECT | Name of the project on server side |

**Example:**
```
# pav project info demo

Attribute        Value                     
name             demo                      
resource file    resources.json            
spec files       specs/module_demo.json
```

#### 3.3.4	import

**Description:**

Import a project file (zip format) to secPaver server.

**Format:**

```
pav project import [-f] <FILE>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| --force, -f | Overwrite import |
| FILE | Project zip file |

**Example:**
```
# pav project import demo.zip 
[info]: Finish importing demo project
```

#### 3.3.5	export

**Description:**

Export a given project as a .zip file from the server side.

**Format:**
```
pav project export [-f] <PROJECT> <PATH>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| --force, -f | Overwrite export |
| PROJECT | Name of the project on server side |
| PATH | Export path |

**Example:**

```
# pav project export demo .
Finish exporting: export_demo.zip
```

#### 3.3.6	delete

**Description:**

Delete a specified project stored in secPaver server.

**Format:**
```
pav project delete <PROJECT>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| PROJECT | Name of the project on server side |

**Example:**
```
# pav project delete demo
[info]: Finish deleting demo project
```

#### 3.3.7	build

**Description:**

Build a project and generate policy based on specified engine. The project can be one on the secPaver server, or one specified by a local path.

**Format:**
```
pav project build --engine <ENGINE> <-d PATH | -r PROJECT>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| --engine ENGINE | Specify the policy generator plugin |
| -d PATH | path for local project |
| -r PROJECT | Name of project stored in secPaver server |

**Example:**

```
# pav project build -d ./demo --engine selinux
[info]: Finish building demo project

# pav project build -r demo --engine selinux
[info]: Finish building demo project
```

### 3.4 Policy Manage Command

**Description:**

Manage the generated policies.

**Format:**

```
pav policy SUBCOMMAND
```

Supported subcommands:

#### 3.4.1	list

**Description:**

List all generated policies.

**Format:**
```
pav policy list
```

**Parameters:**

None

**Example:**

```
# pav policy list

Name                       Status     
demo_selinux               disable
```

#### 3.4.2	install

**Description:**

Install a generated policy.

**Format:**

```
pav policy install <POLICY>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| POLICY | Name of the policy on server side |

**Example:**

```
# pav policy install demo_selinux
[info]: install SELinux policy module
[info]: start to restore file context
[info]: Finish installing policy
```

#### 3.4.3	uninstall

**Description:**

Uninstall a generated policy.

**Format:**

```
pav policy uninstall <POLICY>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| POLICY | Name of the policy on server side |

**Example:**

```
# pav policy uninstall demo_selinux
[info]: uninstall SELinux policy module
[info]: restore file context
[info]: Finish uninstalling policy uninstalling
```

#### 3.4.4	export

**Description:**

Export a generated policy to a zip file.

**Format:**
```
pav policy export [-f] <POLICY> <PATH>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| --force, -f | Overwrite export |
| POLICY | Name of the policy on server side |
| PATH | Export path |

**Example:**

```
# pav policy export demo_selinux .
Finish exporting: export_demo_selinux.zip
```

#### 3.4.5	delete

**Description:**

Delete policy from server side.

**Format:**
```
pav policy delete <POLICY>
```

**Parameters:**

| Parameter | Description |
| :---: | :---: |
| POLICY | Name of the policy on server side |

**Example:**
```
# pav policy delete demo_selinux
[info]: Finish deleting policy
```

## 4 pavd Command

### 4.1 Global Parameters

The global parameters of pavd include help information, version query, and basic configurations.

#### 4.1.1 --help, -h

**Description:**

Print help information.

**Format:**

```
pavd --help
```

#### 4.1.2	--version, -v

**Description:**

Print pavd version.

**Format:**

```
pavd --version
```

#### 4.1.3	--config, -c

**Description:**

Specified config file, the default value is /etc/secpaver/pavd/config.json.

**Format:**

```
pavd --config <FILE>
```

#### 4.1.4	--socket, -s

**Description:**

Grpc socket file path, the default value is /var/run/secpaver/pavd.sock.

**Format:**

```
pavd --socket <PATH>
```

#### 4.1.5	--log-level, -l

**Description:**

Specifies the log level (values could be one of debug, info, warn, error, fatal, panic); It can be specified in the config file. If it is not specified in the config file, it defaults to "info" level.

**Format:**

```
pavd --log-level <LEVEL>
```

### 4.2 systemd command

**Description:**

After secPaver is installed, the pavd service process can be managed through the systemd command.

**Format:**

Start pavd service:

```
systemctl start pavd
```

Stop pavd service:

```
systemctl stop pavd
```

Query pavd service status:

```
systemctl status pavd
```