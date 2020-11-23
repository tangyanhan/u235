---
title: "Istio安装"
date: 2020-01-08T14:58:11+08:00
tags: "Istio"
---

# 安装前置条件

* 已安装helm 2.x, 官方暂不支持 helm 3.0及以上

## 快速安装helm

```
VER=v2.16.1
wget https://mirror.azure.cn/kubernetes/helm/helm-$VER-linux-amd64.tar.gz
tar -xvf helm-$VER-linux-amd64.tar.gz
sudo mv linux-amd64/helm /usr/local/bin 
helm init --upgrade --tiller-image gcr.azk8s.cn/kubernetes-helm/tiller:$VER --stable-repo-url https://mirror.azure.cn/kubernetes/charts/
```

# 下载源码

下面是一个国内的同步仓库, 可以避免因为网速无法从Github下载的问题.

```
git clone https://gitee.com/mirrors/istio.git
```

然后, ```git tag``` 列出可用tag, 选择一个想要的tag并checkout:
```
git checkout -b v1.4.1 1.4.1
```

# 初始化CRD

Istio创建了大量的CRD定义, Istio将它们的安装工作单独放到一个目录中, 其实就相当于我们直接 ```kubectl apply -f .```, 不过借助了helm, 可以指定一些变量.

```
cd install/kubernetes/helm/istio-init
```

这期间由于使用了```gcr.io```的镜像, 国内访问困难, 需要修改镜像地址. 打开 ```values.yaml```, 修改其中的 ```global.hub```字段, 将其改为使用 ```gcr.azk8s.cn```:

```yaml
global:
  # Default hub for Istio images.
  # Releases are published to docker hub under 'istio' project.
  # Dev builds from prow are on gcr.io
  hub: gcr.azk8s.cn/istio-testing
```

azk8s.cn是微软提供的一个镜像加速地址, 可以提供gcr.io, docker.io 这些镜像的国内加速, 比较安全可靠.

```bash
kubectl create ns istio-system
helm install . --name istio-init --namespace istio-system
```

# 安装Controller

一般的, CRD定义包括了CRD定义和Controller两部分.

在安装CRD定义之后, ```kube-api-controller```仅能提供它们的基础CRUD, 并不能完成实际功能. Controller负责这些CRD资源的实际功能实现. 它们的 helm chart 存放在```install/kubernetes/helm/istio``` 目录下.

```
cd ../istio
```

## 启用周边

在```values.yaml```中包含了大量的配置项, 其中就包括了kiali/opentracing等, 如果你需要启用他们, 需要手动修改 ```enabled: true```:

```yaml
grafana:
  enabled: true
#
# addon prometheus configuration
#
prometheus:
  enabled: true
#
# addon jaeger tracing configuration
#
tracing:
  enabled: true
#
# addon kiali tracing configuration
#
kiali:
  enabled: true
```


默认的helm chart包含了Istio主要组件及```prometheus```等, 我们需要修改以下几个地方的镜像, 以提供加速:

./values.yaml:

```yaml
global:
  # Default hub for Istio images.
  # Releases are published to docker hub under 'istio' project.
  # Dev builds from prow are on gcr.io
  hub: gcr.azk8s.cn/istio-testing
```

./charts/prometheus/values.yaml:

```
hub: dockerhub.azk8s.cn/prom
```

./charts/tracing/values.yaml:

```yaml
jaeger:
  hub: dockerhub.azk8s.cn/jaegertracing
zipkin:
  hub: dockerhub.azk8s.cn/openzipkin
```

charts/kiali/values.yaml:

```yaml
hub: quay.azk8s.cn/kiali
```

charts/grafana/values.yaml:

```yaml
image:
  repository: dockerhub.azk8s.cn/grafana/grafana
```

修改完成后, 与Istio-init 一样, 安装即可:

```
helm install . --name istio --namespace istio-system
```

然后,开始观察Pod的启动情况:
```
watch kubectl -n istio-system get po
```

如果有Pod出现问题, 针对排错即可. 如果由于上面的镜像加速没配置好,导致安装阶段出现镜像拉取失败, 不必着急重装. 通过 ```kubectl -n istio-system get deploy``` 可以发现这一步安装的其实都是一些image, 如果哪一个镜像拉取失败, 我们可以直接 ```kubectl -n istio-system edit deploy xxxx```

# 周边配置

默认情况下Kiali不会对外暴露端口, 可以简单点开放NodePort. 访问时会报告缺少secret, 可以这样创建一个 ```admin/admin```的用户名/密码:

```
kubectl create secret generic kiali -n istio-system --from-literal "username=admin" --from-literal "passphrase=admin"
```

然后删除 Kiali 的Pod, 等待重新创建即可.
