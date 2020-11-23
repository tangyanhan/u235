---
title: "二进制安装Kubernetes"
date: 2019-11-04T21:43:03+08:00
tags: "Kubernetes"
---

# 准备二进制文件

uo

```
git clone https://gitee.com/mirrors/Kubernetes
mkdir -p $GOPATH/src/k8s.io
mv Kubernetes/ $GOPATH/src/k8s.io/kubernetes
cd $GOPATH/src/k8s.io/kubernetes
git checkout -b 1.16.2 v1.16.2
```

```
docker pull docker.io/gcrcontainer/kube-cross:v1.12.10-1
docker tag docker.io/gcrcontainer/kube-cross:v1.12.10-1 k8s.gcr.io/kube-cross:v1.12.10-1
make release
docker cp 3c09eb833064:/go/src/k8s.io/kubernetes/_output/dockerized/go/bin ./
```

下载Kubernetes:

进入Kubernetes release页, 选择一个版本, 根据Download Link跳转到二进制下载链接:

https://github.com/kubernetes/kubernetes/releases

```
wget https://dl.k8s.io/v1.16.2/kubernetes.tar.gz
```
解压缩之后,会发现它非常小,明显不是二进制包. 进入 cluster 目录, 运行get-kube-binaries.sh下载二进制包:
```
./get-kube-binaries.sh
```
