---
title: "Kubebuilder"
date: 2020-12-29T13:58:33+08:00
---

## 安装前置

* 安装好go，配置好GOPATH, GOPROXY

* 安装好kustomize

```
go get sigs.k8s.io/kustomize/kustomize/v3
```

* 安装好controller-gen

```
go get sigs.k8s.io/controller-tools/cmd/controller-gen
```

## 安装

确保机器已经安装好go等工具，运行以下命令：

```bash
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv /tmp/kubebuilder_2.3.1_${os}_${arch} /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin
```

## 创建项目

1. 在`$GOPATH`下创建目录

```
mkdir -p $GOPATH/example
cd $GOPATH/src/example
```

2. 初始化项目

```
kubebuilder init --domain my.domain
```

这里`my.domain`可以理解为自己的项目域名。K8S中可能安装很多CRD，这些CRD都需要自己的domain作为后缀来做区分，例如`knative.dev`是`knative`项目使用的domain，它的`Service`对象是属于`serving.knative.dev/v1` 这个`APIGroup`, 在CRD列表中被标记为 `service.serving.knative.dev/v1`.

3. 创建API

K8S通过声明式API操作CRD资源，创建API其实就是在创建相关的资源定义和Controller。 `kubebuilder`可以在创建过程中帮忙生成定义和Controller，当然都只是模板:

```
kubebuilder create api --group webapp --version v1 --kind Guestbook
```

![create_api](/images/post/screenshot_20201229_141153.png)

这里，创建了一种资源 Guestbook，如果我们创建这种资源，那么它的YAML开头长这样：
```yaml
kind: Guestbook
apiVersion: webapp.my.domain/v1
```

这时，我们会发现在工程目录下， `api/v1`目录中已经创建好了类型模板， `controllers`目录下也出现了相关的Controller定义，其中的`Reconcile`函数留空让我们补充逻辑。

4. 安装运行

## 本地运行Controller


```bash
# 重新生成并安装CRD定义
make install
# 运行Controller
make run
```

在`make run`时，目前默认的模板可能会报错，这是由于`k8s.io/client-go`的版本过高，与`sigs.k8s.io/controller-runtime`不协调导致的，将`sigs.k8s.io/controller-runtime`的版本调高即可。

然后就可以通过`kubectl apply -f config/samples/`尝试，可以看到controller运行日志:

![create_sample](/images/post/screenshot_20201229_150017.png)

`make run`是在本地编译运行controller，使用的是当前的默认kubeconfig设定。

## 集群中部署运行Controller

```bash
# 构建镜像
make docker-build docker-push IMG=<some-registry>/<project-name>:tag
# 部署镜像
make deploy IMG=<some-registry>/<project-name>:tag
```

5. 卸载CRD定义

```
make uninstall
```

