---
title: "使用Virtualbox+kubeadm搭建k8s集群"
date: 2019-11-10T20:11:43+08:00
tags: "Kubernetes"
---

搭建K8S集群,动不动就是 1 Master 3节点, 或者更多这样子. Minikube只是个单节点K8S,不能用来学习更多多节点调度之类的知识. 现在的机器性能这么强劲, 为什么不在自己的笔记本或台式电脑上, 用Virtualbox虚拟机做出自己的多节点集群呢?

第一步,安装好VirtualBox. 某些主机可能需要在BIOS中开启虚拟化支持,才能在Virtualbox中安装64位操作系统.

第二步, 下载一个系统镜像. 我使用的是[CentOS7 Minimal](https://mirrors.aliyun.com/centos/7.7.1908/isos/x86_64/CentOS-7-x86_64-Minimal-1908.iso), 几乎啥软件都没有, 能够将折腾发挥到极致.

第三步, 创建一个虚拟机. K8S集群每个节点至少2核CPU, 2G内存. 不用担心, 这并不意味着创建3个这样的节点一定会使你的6个CPU核心飙升100%, 内存占用达到6G. 硬盘至少40G, 使用动态分配可以有效降低实际磁盘占用.

**重点** 网络除了默认的NAT网络外,另外添加一个HostOnly的网络, 这样主机和所有添加了该网络的虚拟机可以形成一个局域网,能够互通:

![网络设置](/images/post/virtualbox-network.png)

安装过程略略略, 使用默认设置即可.

第四步, 更换国内软件源,并安装必要的软件.

小插曲: 通过上面的网络设置, 你的虚拟机拥有了两块网卡, 如果你发现自己不能联网,可能是因为其中一块网卡未启动. 通过 ```ip route``` 查看网卡名字,然后 ifup 网卡名 即可.

1. 移除原来的所有软件源

```
cd /etc/yum.repos.d/
rm -rf *.repo   # 这里是为了删掉不必要的软件更新等,加快更新缓存速度
curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo

```
2. 添加普通软件源, docker-ce软件源及Kubernetes软件源
```
curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
yum makecache
```

3. 安装docker/kubernetes必要组件等

```
yum install -y docker-ce kubelet kubeadm kubectl
systemctl start docker.service
systemctl enable docker.service
systemctl enable kubelet && systemctl start kubelet
```

4. 打通SSH,使你能够通过SSH连接到虚拟机

```
yum install -y openssl openssh-server
systemctl start sshd.service
systemctl enable sshd.service
```

等等,我在主机上如何连接这台虚拟机? 首先,在我们的主机上, 通过```ip route``` 找到```vboxnet0```这个网络:

![主机网络](./images/posts/host-network.png)

可见该网络的网段是 192.168.99.0/24, 我们自己的主机地址是 192.168.99.1

在虚拟机上, 通过 ```ip addr``` 命令可以查看网卡及网段:

![virtual](./images/posts/virtual-network.png)

可见虚拟机中使用该网段的网卡为 enp0s8, 地址为 192.168.99.104 , 这就是接下来它在K8S集群中的地址,也是我们可以从主机上通过SSH连接的地址.

5. 复制主机作为其它节点

我们花费了很大的经历搭建单台节点,不过好在必要的软件已经装好了, 接下来只要克隆弄好的主机,然后稍作改动就可以作为子节点了.

首先关闭虚拟机,然后在菜单上右键点击虚拟机并克隆. 注意克隆时选择生成新的MAC地址, 这样我们就不用费劲去给每台主机修改IP地址啦.

然后,我们对每台机器做以下操作:

a. 修改hostname

```
hostnamectl set-hostname k8s-node1
```

b. 修改IP地址(好像不用改了?)


5. 使用kubeadm初始化Master节点

```
kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=<虚拟机IP地址> --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers
```

cidr 设置成这个地址是为了后续安装flannel时减少折腾量. ```--apiserver-advertise-address``用于与其它节点通信,因此应该是我们第4步中获取到的SSH地址.

--image-repository 是为了加快镜像拉取速度-- 好吧, 真话是如果是在大陆, 默认镜像根本就拉取不下来.

6. 添加其它的节点

