---
title: "Runc容器运行过程及两个逃逸漏洞原理"
date: 2021-01-28T09:13:11+08:00
draft: true
tags: "Kubernetes" "安全"
---

在每一个Kubernetes节点中，运行着kubelet，负责为Pod创建销毁容器，kubelet预定义了API接口，通过GRPC从指定的位置调用特定的API进行相关操作。而这些CRI的实现者，如cri-o, containerd等，通过调用runc创建出容器。runc功能相对单一，即针对特定的配置，构建出容器运行指定进程，它不能直接用来构建镜像，kubernetes依赖的如cri-o这类CRI，在runc基础上增加了通过API管理镜像，容器等功能。

Kubelet，Cri-O，runc，Linux大致层级示意图如下：

![container-runtime](/images/post/container-runtime.png)

在Kubernetes源码中，可以在`pkg/kubelet/cri`目录下找到相关代码，其中`remote`目录包含了常见的如镜像拉取，容器创建等操作，`streaming`目录中包含了一些需要TCP流的操作，如attach，port-forward等。

# 构建并使用runc运行一个容器

## 构建
runc的源码可以下载并通过make命令构建: 

```bash
git clone https://github.com/opencontainers/runc.git
cd runc
make
```

runc是一个Go程序，使用CGO调用了一些外部库。构建除了需要安装go之外，可能需要额外安装如pkgconfig, libseccomp-dev, libseccomp等包，视具体错误排查。

## 运行容器

runc并不负责从镜像等上下文直接创建容器，因此需要从docker等更高级的运行时直接导出CRI，会更容易一些。

```bash
mkdir /mycontainer
cd /mycontainer

# create the rootfs directory
mkdir rootfs

# export busybox via Docker into the rootfs directory
docker export $(docker create busybox) | tar -C rootfs -xvf -
```

然后使用构建好的runc创建出容器运行的具体配置config.json:

```bash
runc spec
```

切换到root，然后运行:

```bash
runc run <container-id>
```

或者将`run`命令拆解，分成多步运行:

```bash
runc create <container-id>
runc list # 列出创建状态的容器
runc start <container-id>
runc list
runc delete <container-id>
```

## 运行分析

通过`ps axef`命令，打印出所有进程及其层级关系，发现我们之前运行的容器进程关系如图:

![process-hierachy](/images/post/Screenshot_20210128_094910.png)

表面上看，在通过`runc run <container-id>`之后，进程创建了一个子进程`sh`，也就是我们进入容器后指定运行的第一个程序。

## 容器隔离机制

Linux通过namespace将不同的资源隔离，就像一个沙箱一样。被隔离到某个namespace中的内容，无法访问到其它namespace的内容。可以通过`unshare`或`clone`设置标志位来将进程放入新的命名空间。

* PID Namespace

标志位: `CLONE_NEWPID`

文档: `man pid_namespaces`

即进程命名空间，在一个隔离的空间中，PID从1开始，相同PID与主机PID不构成冲突。类似主机的init进程，PID为1的进程被终止时，该命名空间下的所有进程都会收到`SIGKILL`信号从而被终止。正因如此，一个容器的初始化进程只能是一个，而且终止后容器也就被停止了。

在不同的PID命名空间，进程互相看不到对方，不能通过PID找到对方，`/proc`目录下也只能看到自己命名空间中的进程。但是一个父进程`fork`出的子进程可以通过`set_ns`放入子命名空间，在父进程的命名空间，仍然可以看到这个子进程，只是PID不一样。进程可以被挪到子命名空间，但不能被反向挪回更高级的命名空间。

* User Namespace

标志位: `CLONE_NEWUSER`

文档: `man user_namespaces`

用户命名空间，主要隔离的是安全相关的id和属性，尤其是用户id和用户组id，root目录，密钥，以及各种进程的能力(capabilities)。

* Network Namespace

标志位: `CLONE_NEWNET`

文档: `man network_namespaces`

网络命名空间，隔离出全新的网络设备、IPv4、IPv6协议栈，路由表，iptables规则，以及sockets,端口号。隔离后，命名空间下的`/proc/net`以及`/sys/net`也会不同于上级命名空间。

一对`veth`网络设备可以实现跨命名空间的通信，也可以桥接到主机物理设备上。

* UTS Namespace

标志位: `CLONE_NEWUTS`

隔离hostname以及NIS domain name，两个都是主机对自己的网络标识，在容器中可以重新定义。

* Cgroup Namespace

标志位: `CLONE_NEWCGROUP`

文档: `man cgroup_namespaces`

Cgroup Namespace隔离出了新的cgroup分组，可以通过`/proc/[PID]/cgroup`获得进程的cgroup相关情况。Cgroup是`Control group`的缩写，不同的cgroup通过树状结构组织在一起，每个节点上挂载着不同的控制器(controllers)，可以通过他们控制cpu, memory, bulkio等资源的使用，也可以获得cpu等资源的使用情况，或者通过freezer将进程暂时冻结或恢复。

* IPC Namespace

标志位: `CLONE_NEWIPC`

隔离了System V常用的IPC通信手段，如信号量(semaphore)，共享内存，消息队列等。

* Mount Namespace

标志位： `CLONE_NEWNS`

挂载文件系统的隔离，但是一部分文件系统也可以通过共享跨命名空间共享。通过`cat /proc/self/mountinfo`可以获得挂载信息，带有`shared`标志的就是共享出来的部分。

## 逃逸漏洞分析





