---
title: "基于Hyper-V/kubeadm的k8s集群搭建"
date: 2023-08-16T10:21:55+08:00
draft: true
tags: "Kubernetes"
---

本文介绍在Windows上借助Hyper-V和kubeadm搭建k8s机器群的方法及避坑。

# 前置需求

* 虚拟机应当具有固定的ip
* 虚拟机应当能够访问外部网络
* 虚拟机能够使用主机上的代理软件（如clash）

# 配置网络

1. 确保Windows已启用Hyper-V功能并重启

以admin权限通过PowerShell运行：
```
Enable-WindowsOptionalFeature -Online -FeatureName Microsoft-Hyper-V -All
```

2. 网络配置

创建一个k8s的虚拟交换机，为其分配一个固定ip，并在主机上配置一个NAT使其中的主机能够访问外部网络
```
New-VMSwitch -Name 'k8s' -SwitchType Internal
$adapter = Get-NetAdapter | Where { $_.Name -match 'k8s' }
New-NetIPAddress -IPAddress '172.16.0.1' -InterfaceIndex $adapter.InterfaceIndex -PrefixLength 24
new-netnat -Name 'k8s' -InternalIPInterfaceAddressPrefix 172.16.0.0/24
```

![网络示意图](./images/posts/k8s-network.png)

如果在使用代理服务，确保监听在了所有地址之上，例如Clash需要打开AllowLan: true

3. 安装虚拟机

虚拟机安装无需赘述，使用centos/ubuntu-server皆可。
确保节点配置为：
* CPU >= 2
* Memory >= 4G
* SecureBoot 已禁用
* 安装并启用ssh server服务
* 网卡选择我们之前创建的'k8s'
* hostname 分别改为controlplane/node

在安装过程中,可以对网卡直接配置固定IP, 设定IPv4为manual:

```
Address=172.16.0.100 # 或.101
Gateway=172.16.0.1
DNS=223.5.5.5,172.16.0.1
```

如果错过了这些配置：

### CentOS:

* 修改 /etc/sysconfig/network-scripts/ifcfg-eth0 配置固定IP
* 修改 /etc/hostname 配置hostname

### Ubuntu

* 修改 /etc/netplan/01 后通过 netplan apply 配置固定IP
* `hostnamectl set-hostname <new-name>` 修改主机名

# 节点配置

以下内容需要在每个节点上都运行一遍

## 代理

如果外部有网站，可以临时编写一个proxy.sh, 通过source proxy.sh临时使用外部代理，解决wget/curl不能顺利访问某些网站的问题
```
export http_proxy="http://172.16.0.1:7890"
export https_proxy="http://172.16.0.1:7890"
export no_proxy="127.0.0.0/8,172.16.0.0/16,10.96.0.0/16,10.98.0.0/16,192.168.0.0/16"
```
如果自己没有代理，需要考虑利用https://developer.aliyun.com/mirror/和其它资源配置软件源、镜像源。

## 安装CRI

选定containerd为容器运行时支撑我们的Pod容器，可通过官方文档安装：

https://docs.docker.com/engine/install/centos/
https://docs.docker.com/engine/install/ubuntu/

确保安装完成后，`which containerd` 能够找到它

crictl类似docker，用于管理查看containerd容器，通过`which crictl`判断是否已经拥有crictl，如果没有新装一个，然后新建文件 `/etc/crictl.yaml`, 内容如下:
```
runtime-endpoint: unix:///run/containerd/containerd.sock
image-endpoint: unix:///run/containerd/containerd.sock
timeout: 10
debug: false
```

#### 验证安装
```
systemctl start containerd
systemctl enable containerd
crictl ps -a # 确保crictl/containerd都能够正常工作
```

### 配置containerd

#### 1. 启用CRI功能

编辑 `/etc/containerd/config.toml`, 注释掉这一行：
```
#disabled_plugins = ["cri"]
```

#### 2. 为containerd配置cgroup2

查看当前系统cgroup使用的是哪一版：
```
mount |grep cgroup2
```
如果有东西，说明使用的是cgroupv2，目前版本的containerd需要添加额外配置才能让k8s正常运行，否则会不断杀死etcd容器，需要在 `/etc/containerd/config.tomal`中追加：

```
version = 2
[plugins]
  [plugins."io.containerd.grpc.v1.cri"]
   [plugins."io.containerd.grpc.v1.cri".containerd]
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes]
        [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
          runtime_type = "io.containerd.runc.v2"
          [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc.options]
            SystemdCgroup = true
```

#### 3. containerd使用代理

修改 `/lib/systemd/system/containerd.service`, 在`[Service]` 下添加：

```
Environment="HTTP_PROXY=http://172.16.0.1:7890/"
Environment="HTTPS_PROXY=http://172.16.0.1:7890/"
Environment="NO_PROXY=10.96.0.0/12,127.0.0.0/8,192.168.0.0/16,localhost,172.16.0.0/16,.svc,.cluster.local,.ewhisper.cn"
```

然后重启containerd即可：
```
systemctl daemon-reload
systemctl restart containerd
```

## 配置内核模块

加载模块br_netfilter:
```
modprobe br_netfilter
```

修改 `/etc/sysctl.conf`, 添加如下内容：
```
net.ipv4.ip_forward=1
net.bridge.bridge-nf-call-iptables=1
```

随后通过`sysctl -p` 即可生效

## 禁用swap

使用`swapoff -a`只能临时禁用swap，重启后会失效。可以通过修改 `/etc/fstab`, 将挂载的swap移除, 重启后生效

## 禁用防火墙或开放端口

简单一点可以直接禁用防火墙：

CentOS:
```
systemctl disable firewalld
systemctl stop firewalld
```

Ubuntu:
```
ufw disable
```

开放端口则麻烦一些，需要开启以下端口，包括但不限于：

Port | Usage|
-----|------|
6443 | API Server |
10250 | Kubelet 通信 |
2379 | ETCD |
2380 | ETCD |

## 安装kubeadm/kubelet/kubectl

### CentOS
```
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg
EOF
setenforce 0
yum install -y kubelet kubeadm kubectl
systemctl enable kubelet && systemctl start kubelet
```

### Ubuntu
```
apt-get update && apt-get install -y apt-transport-https
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add - 
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF
apt-get update
apt-get install -y kubelet kubeadm kubectl
```

# 安装集群

## Controlplane

通过`kubeadm config print init-defaults > init.yaml` 生成配置模板

修改init.yaml:

1. 修改localAPIEndpoint.advertiseAddress 为节点真实地址，如`172.16.0.100`

2. 修改nodeRegistration.name 为当前节点真实hostname，如controlplane

3. 添加networking.podSubnet, 为`192.16.0.0/16`, 这样安装CNI就无需额外配置了：
```
networking:
  dnsDomain: cluster.local
  serviceSubnet: 10.96.0.0/12
  podSubnet: 192.168.0.0/16
```

随后通过`kubeadm init --config init.yaml`初始化控制节点。在kubeadm初始化结束时，会提供一条命令，可供其它节点加入。

## Node

在完成之前的节点配置后，其它节点都可以通过```kubeadm join 172.16.0.100:6443 --token abcdef.0123456789abcdef 
--discovery-token-ca-cert-hash sha256:<signature>``` 来加入集群

如果不慎漏掉了之前的log，可以通过命令找回signature：

```openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //'```

## CNI

CNI 能够帮助Pod分配IP，帮助Pod之前进行跨节点通讯，完成Pod网络隔离等。以Calico为例，可以通过官方文档，两个kubectl create -f 搞定：
https://docs.tigera.io/calico/latest/getting-started/kubernetes/quickstart

# 验证集群

配置kubectl
```
mkdir ~/.kube
cp /etc/kubernetes/admin.conf ~/.kube/config
```

查看namespace：

```
kubectl get ns
```

查看Pod IP分配：

```
kubectl get po -A -o wide```