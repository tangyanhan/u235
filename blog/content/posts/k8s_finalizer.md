---
title: "K8S中的删除：垃圾回收，资源清理以及Finalizer实现"
date: 2023-11-24T12:14:07+08:00
tags: "Kubernetes"
---

现代大多数高级语言都会自动进行垃圾回收(Garbage Collection)，但一般也会强调对于特殊的资源，例如fd, net connection在删除前进行手动回收, 这些一般会放在析构函数之类的地方，对应K8S中就是Finalizer。

# Garbage Collection

Garbage Collection默认由Controller Manager自带实现，具体使用需要进行参数调整。

## 以Succeeded/Failed结束的Pod

对于已经终止的Pod，由Controller Manager通过`terminated-pod-gc-threshold`这一阈值进行管理，目前其默认值为12500，如果集群中需要更好的清理他们，可能需要调低数值：

```
func RecommendedDefaultPodGCControllerConfiguration(obj *kubectrlmgrconfigv1alpha1.PodGCControllerConfiguration) {
	if obj.TerminatedPodGCThreshold == 0 {
		obj.TerminatedPodGCThreshold = 12500
	}
}
```

## 已完成的Job

适用于在Job进入Completed状态一段时间后自动清理： `.spec.ttlSecondsAfterFinished`

## Owner Reference失效的孤儿资源

### 级联删除

删除策略可以通过在client-go中修改Pro0pagationPolicy来改变，或者在kubectl中通过`--cascade`参数改变。

#### 同步级联删除(Foreground cascading deletion)

```kubectl delete deploy nginx --cascade=foreground```

* API Server修改`metadata.deletionTimestamp`作为标记删除时间
* API Server 更新`metadata.finalizers`为`forgroundDeletion`
* 资源直至整个删除结束才会消失

#### 异步级联删除

默认通过kubectl删除资源时，就会通过此选项。在使用client-go时，需要额外制定。

```kubectl delete deploy nginx --cascade=background```

#### 单独删除，故意让子级资源成为孤儿

```kubectl delete deploy nginx --cascade=orphan```

这种情况下，孤儿资源过一段时间就会被Controller进行回收。

## 节点容器资源回收

节点资源由kubelet进行回收。

### 镜像回收

当磁盘资源达到一定程度时，开始回收下载镜像所占用的磁盘资源，主要控制参数：

```
ImageGCHighThresholdPercent // 磁盘占用超过该值，必定开始回收
ImageGCLowThresholdPercent  // 磁盘占用低于该值，必定不会回收
ImageMinimumGCAge           // 镜像未被使用一段时间后应该被回收，这里规定这个最低时长
```

### 容器回收

一些容器构建(docker build)等会在容器构建打包过程中创建出一些中间容器，他们往往都不活跃，kubelet可以自动回收这些容器资源。

# Pod足迹清除

有些Pod在运行过程中会在磁盘上写入一些文件，要清除这些文件，可以通过`.spec.containers.lifecycle.preStop`来指定容器终止前的动作

# Finalizer

API Server一般只负责资源的CRUD，在删除资源时，它所能做的就是修改标记删除时间，以及从etcd中删除数据。具体的额外动作需要通过 `.metadata.finalizers`来指定。

## Finalizer工作机制

在请求删除资源时，如果目标资源的`.metadata.finalizers`字段不为空，那么API Server只会对资源进行标记，直到finalizers字段为空，才会进行etcd资源删除。

通过观察Controller Manager自带的几个finalizer可知finalizer的工作机制与其它的Resource controller非常接近：

1. 监视资源的更新和删除事件
2. 当有符合条件的资源进行更新时，将其finalizers字段中加入我们要使用的finalizer
3. 当有带有该finalizer字段的资源进行删除时，进行相关动作后，将我们的finalizer字段从中移除
4. 等待所有finalizer字段被移除，API Server即可进一步实际删除资源

因此在资源卡住时，可以通过强制清除finalizer字段来进行强行删除. 那么这些finalizer实际做了什么，如果强行删除会有什么影响呢？

### 删除namespace时，背后都在发生什么？

在操作自己搭建的小集群时，偶尔会尝试删除整个namespace，经常会出现这个namespace迟迟删不掉的情况，这是因为namespace创建时就已经在finalizer中加入了'kubernetes'这个字段：
```json
"spec": {
    "finalizers": [
        "kubernetes"
    ]
},
```

当删除卡死时，可以尝试通过PUT请求直接清空finalizer，那么这可能会导致什么后果？

finalizer相关逻辑与namespace controller共同运行于controller manager中，在namespace已经被标记删除且finalizer尚未清空时，会通过 [deleteAllContents](https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/namespace/deletion/namespaced_resources_deleter.go#L502) 发现所有GroupVersionResources，并删除它们，提供一个预估时间。

预估时间主要通过GracePeriodSeconds来进行粗略估计，对于部分支持Graceful Deletion的资源，也对应的进行Graceful Delete（目前主要原生资源是Pod支持Graceful Delete），其它自定义资源需要支持RESTGracefulDeleteStrategy：

```go
// RESTGracefulDeleteStrategy must be implemented by the registry that supports
// graceful deletion.
type RESTGracefulDeleteStrategy interface {
	// CheckGracefulDelete should return true if the object can be gracefully deleted and set
	// any default values on the DeleteOptions.
	// NOTE: if return true, `options.GracePeriodSeconds` must be non-nil (nil will fail),
	// that's what tells the deletion how "graceful" to be.
	CheckGracefulDelete(ctx context.Context, obj runtime.Object, options *metav1.DeleteOptions) bool
}
```

因此清除finalizer的主要问题，就是会导致其中一些需要Graceful Deletion的资源不能得到正确删除，可能会出现一些资源泄露、状态不一致之类的问题。另外，这些资源的删除过程也将不再由k8s显式控制，无法观察他们的状态。

