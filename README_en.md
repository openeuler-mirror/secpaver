# secPaver

## 1 Introduction
secPaver is a tool to help OS administrators define security policies of different security mechanisms for applications.

Now secPaver supports SELinux security mechanisms. And AppArmor will be supported soon.

## 2 Installation
**(1) Dependency packages**
   
   To build secPaver, the following packages are required:
   * make
   * golang 1.11+
   
   To build SELinux engine plugin, the following packages are also required:
   * libselinux-devel 2.9+
   * libsepol-devel 2.9+
   * libsemanage-devel 2.9+
   
   To run SELinux engine plugin, the following packages are required:
   * libselinux 2.9+
   * libsepol 2.9+
   * libsemanage 2.9+
   * checkpolicy 2.8+
   * policycoreutils 2.8+
   
   **(2) Download source code**
   ```
   git clone https://gitee.com/openeuler/secpaver.git
   ```
   
   **(3) Build and install**
   ```
   cd secpaver
   make
   ```
   Build SELinux engine plugin:
   ```
   make selinux
   ```
   Install (after at least one engine plugin is built):
   ```
   make install
   ```

## 3 Terms
**engine**：A software module for generating security policy based on a given security mechanism. An engine in secPaver is loaded as a software plugin.

## 4 How to Use
secPaver uses client/server architecture; the server process is pavd, and the client process is pav.

**(1) Start pavd process**

```
systemctl start pavd
```

**(2) Check loaded policy generator plugins**

```
# pav engine list

Name        Description                  
selinux     SELinux policy generator
```

**(3) Create a template project and modify it**

```
pav project create my_demo .
```

**(4) Build project to policy**

```
pav project build -d ./my_demo --engine selinux
```

**(5) List generated policies**

```
# pav policy list

Name                           Status     
my_demo_selinux                disable
```

**(6) Install policy**

```
# pav policy install my_demo_selinux
[info]: install SELinux policy module
[info]: start to restore file context
[info]: Finish installing policy
```

**(7) Uninstall policy**

```
# pav policy uninstall my_demo_selinux
[info]: uninstall SELinux policy module
[info]: restore file context
[info]: Finish uninstalling policy uninstalling
```

**(8) Export policy package**

```
# pav policy export my_demo_selinux .
Finish exporting: export_my_demo_selinux.zip
```

## 5 Document

Command manual: [secPaver Command Manual](doc/cmd_en.md)

User manual：[secPaver User Manual](doc/manual_en.md)

## 6 How to Contribute

We welcome new contributors to the project, and are pleased to provide guidance and assistance to new contributors. Before you contribute code, you need to sign [CLA](https://openeuler.org/en/cla.html)。
