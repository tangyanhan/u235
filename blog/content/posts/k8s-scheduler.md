---
title: "K8S Scheduler原理及扩展"
date: 2020-05-12T10:55:30+08:00
tags: "Kubernetes"
---

# Schedule 过程

Kubernetes 默认的 Scheduler 负责调度Pod. 在由 ApiServer 创建出Pod后, Scheduler 负责写入NodeName字段, 然后由对应Node上的Kubelet负责创建Pod实际的containers等.

默认调度器主要代码实现在这个文件中: ```pkg/scheduler/core/generic_scheduler.go```

主要工作过程, 及相关函数如图:

![Schedule](/images/post/k8s-scheduler.png)

图中红色区域, 都是可以返回一个非成功状态从而中断调度过程的. 但是并非所有步骤都有官方默认实现.

目前的scheduler自带实现主要是集中在两个部分: Filter 与 Score. 

Filter 实现了Pod调度的硬性需求筛选, 例如 NodeSelector, NodeAffinity(Required...)等

Score 则对筛选出的Pod进行了评分, 主要是针对软性需求. 例如非强制的NodeAffinity, Image Locality等.

以 Image Locality 为例, 它实现的功能是对"Pod中使用的镜像是否在本地"这一点进行打分. 因为一个节点如果包含了更多Pod需要的镜像, 那么拉取时间就会降低, 创建Pod可以更快. 对应的实现函数为 ```ImageLocalityPriorityMap```,这里不再赘述.

那么假设一个 image tag设定为 latest, 而 imagePullPolicy=Always, 那么这个插件计算出的分数就可能难以正确匹配. 可见 Pod 中引用image时, 使用明确的tag更容易让Pod分配到合适的节点.

[官方文档参考](https://kubernetes.io/docs/concepts/scheduling-eviction/scheduling-framework/)

# Preempt 抢占过程

当有一个Pod无法被调度到节点时, 可以进行抢占过程(Preempt), 通过让某个节点赶走一部分较低优先级的Pod, 以便顺利让Pod调度到节点上.

其基本流程如下:

![抢占](/images/post/sched-preempt.png)

抢占的目标是使危害最小:

1. 优先级>=Pod的,不应该被驱逐

2. 应该驱逐尽可能少的Pod

3. 应该驱逐优先级尽可能低的Pod

4. 可以驱逐尽可能新的Pod, 避免影响到一些长期运行的服务

其主要实现在 genericScheduler.Preempt 中.

# 扩展Scheduler

Scheduler既可以使用源码自带的算法, 也可以通过外部独立运行的进程, 实现调度

## 1. 浅层扩展: watch Pod

首先, 设置 ```Pod.spec.schedulerName=xxx```

由于系统并没有这样的Scheduler, Pod将陷入 Pending状态.

然后, 通过 client-go watch Pod, 将观察到 schedulerName=xxx 的Pod, 修改其 NodeName, 使其分配到具体节点.

pod.spec.nodeName 是一个特殊字段, 当其被修改为非空值时, 将会不经过任何调度器, 直接分配到具体Node. 假设过程中发生资源等的异常, 那么Pod失败.

## 2. 源码扩展, 集群中增加一个自定义Scheduler

从源码从面, 修改原来的代码, 或手动增加 Scheduler/Extender. 缺点是需要维护自己的源码, 同时还要与官方scheduler保持一定的同步.

集群中允许存在多个Scheduler. 你也可以在官方源码基础上做一定修改, 然后将其作为一个Pod的可选调度器加入集群中.

官方文档提供了集群中增加另一个Scheduler的方法: [参考文档](https://kubernetes.io/docs/tasks/administer-cluster/configure-multiple-schedulers/#run-the-second-scheduler-in-the-cluster)

调度时, 修改 ```pod.spec.schedulerName``` 即可启用自定义调度器.

## 3. 外部Extender

Kubernetes 默认允许Extender在多个点对KubeScheduler的调度行为进行干涉.

下图中绿色的部分都可以通过Extender扩展.

![Extender](/images/post/scheduling-framework-extensions.png)

Extender 需要修改原 kube-scheduler 的配置文件中 ```extenders```一节. 修改配置后, Extender可以通过API形式提供对Pod调度的干涉.

### 0. 编译并启动一个extender

样例源码放在我的 github 仓库中: 

https://github.com/tangyanhan/u235/tree/master/go/advanced/cmd/scheduler-extender

在下文中, 我启动的extender 运行在 192.168.99.1:9000, 因此有相关配置.

每个Extender可以扩展的请求, 都是一种POST请求, 返回200并包含指定结构即成功. 至于具体的输入输出, 可以参考kubernetes源码中的测试样例.

### 1. 增加配置yaml

```yaml
# /etc/kubernetes/scheduler-config.yaml
apiVersion: kubescheduler.config.k8s.io/v1alpha2
kind: KubeSchedulerConfiguration
clientConnection:
  kubeconfig: "/etc/kubernetes/scheduler.conf"
extenders:
- urlPrefix: "http://192.168.99.1:9000"
  filterVerb: filter
  prioritizeVerb: prioritize
  weight: 1
```

### 2. 增加启动参数 --config

我的kube-scheduler是kubeadm通过 StaticPod 启动的, 直接使用kubelet并不能正确影响它. 修改我的kube-scheduler static pod 设定, 增加 --config 指向我们的配置文件:

```yaml
# /etc/kubernetes/manifests/kube-scheduler.yaml
spec:
  containers:
  - command:
    - kube-scheduler
    - --authentication-kubeconfig=/etc/kubernetes/scheduler.conf
    - --authorization-kubeconfig=/etc/kubernetes/scheduler.conf
    - --bind-address=127.0.0.1
    - --kubeconfig=/etc/kubernetes/scheduler.conf
    - --leader-elect=true
    - --config=/etc/kubernetes/scheduler-config.yaml
# ...
```

修改一下挂载目录, 然后 ```systemctl restart kubelet``` 重建static pod.

**小坑**
1. manifests 下的所有文件都会被视为yaml并创建, 修改配置时, 备份文件不要放在同一文件夹下
2. kubectl 不能对static pod有效操作, 应当通过docker rm删除static pod的容器, 或通过restart kubelet来重建static pod.

### 3. 创建一个Pod, 观察我们的extender日志

```bash
[ethan@ethan go]$ ./scheduler-extender 
2020/05/13 16:17:35 Started extender
2020/05/13 16:17:48 Filter: called for pod. namespace= default name= 12345
Node: node [{InternalIP 192.168.99.102} {Hostname node}]
```

胜利!