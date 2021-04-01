---
title: "StatefulSet使用及运行原理"
date: 2021-04-01T13:38:05+08:00
draft: true
tags: "Kubernetes"
---

StatefulSet 字面意思即“有状态的集合”。相比于ReplicaSet这种Pod用完即抛，每次名字都不一样的集合而言，StatefulSet可以保持Pod名字和PVC的一一对应，从而保持状态。

StatefulSet所创建的Pod名称固定为 "<StatefulSet名字>-<Pod序列号，从0开始>", 其使用PVC模板创建出的PVC模板，格式固定为"<PVC模板名>-<StatefulSet名>-<Pod序号>"

```go
// getPodName gets the name of set's child Pod with an ordinal index of ordinal
func getPodName(set *apps.StatefulSet, ordinal int) string {
	return fmt.Sprintf("%s-%d", set.Name, ordinal)
}

// getPersistentVolumeClaimName gets the name of PersistentVolumeClaim for a Pod with an ordinal index of ordinal. claim
// must be a PersistentVolumeClaim from set's VolumeClaims template.
func getPersistentVolumeClaimName(set *apps.StatefulSet, claim *v1.PersistentVolumeClaim, ordinal int) string {
	// NOTE: This name format is used by the heuristics for zone spreading in ChooseZoneForVolume
	return fmt.Sprintf("%s-%s-%d", claim.Name, set.Name, ordinal)
}
```

名称固定方便了与PVC的一一绑定，以及其它应用根据名称访问不同的Pod。 但也带来一个问题，同一个namespace下Pod名字肯定是不能重的，这就意味着StatefulSet中的某个Pod，必须先干掉原来的，才能创建出新的。当删除一个StatefulSet时，它的Pod可能还没有被删除，这时如果马上重新创建一个同名StatefulSet，那么这些孤儿Pod将会被重新纳入新的StatefulSet管控。

在为StatefulSet创建Pod时（创建，扩容，缩容），有个字段`PodManagementPolicy`控制创建Pod的行为，默认`OrderedReady`，表示按照序号从低到高的顺序，Pod进入Ready状态后才能创建下一个，缩容时则反之，从高到低。另一种策略为`Parallel`，表示可以像ReplicaSet那样，同时创建多个Pod。

`Parallel`在有多个副本时创建速度更快，但显然不适合特定场景，例如主从复制的Redis。假设我们使用 `redis-0`作为redis master， `redis-1`作为redis-slave，在使用Parallel策略时，有很大可能会出现`redis-1`先于`redis-0`进入Ready状态，这时请求就会优先落入`redis-1`，可能与我们的预期不符。它控制了创建和更新Replica时的行为，删除时不受影响。

StatefulSet经常被用于部署如Redis、Kafka、MySQL这类数据服务，往往都被多个业务依赖，其升级策略并不像无状态服务那样随意。

目前StatefulSet默认更新策略为`RollingUpdate`,Pod将会按照序号从高到低的顺序更新，与创建时的行为基本一致，也是要等到上一个Pod进入Ready状态后，才能继续下一个。

默认情况下，`RollingUpdate`会更新全部Pod，这时可以指定`partition`值，只有序号大于等于该值的Pod才会被删除重建。

另一种策略为`OnDelete`，当我们更新了StatefulSet里的设定时，所有的Pod都不会被删除更新，只有手动删除一个Pod之后，才会按照新的定义重新创建一个Pod。这种策略主要是为了实现1.6以前的默认行为。

