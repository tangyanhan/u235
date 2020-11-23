---
title: "Kubernetes源码分析:选举实现"
date: 2020-06-17T15:25:01+08:00
tags: "Kubernetes"
---

Kubernetes 中有很多服务都是可以高可用部署的, 但有些需要对K8S资源进行创建/更新操作的服务, 在观察到资源变化后, 很容易同时对资源进行操作, 导致不必要的更新.

虽然K8S有 ResourceVersion 这种锁存在, 多个实例同时进行资源的创建/更新操作也会浪费资源, 使得代码逻辑变得复杂. ControllerManagers, Scheduler 两个组件目前都是使用 client-go 中的 ```leaderelection``` 组件, 实现单个 leader 的选举.

## 竞态资源

我们知道, K8S的包括List-Watch机制, 以及数据一致性都是靠 ETCD 来保证的. 有了ETCD, K8S就可以假设自己有了一块可以保证数据一致性的区域, 允许各个组件通过向它写入/读取数据, 来完成选举过程.

ETCD已经通过Raft实现了一致性的保证, 我们只需要借助K8S自带的某种资源, 所有候选人(Leader Candidate) 读写同一个 namespace/name 的某种资源即可.

目前client-go 自带的竞态资源有 ConfigMap/Endpoints/Coordination, 都是结构比较简单的数据类型. 默认推荐为 Endpoints.

围绕竞态资源的操作被封装为 ```k8s.io/client-go/tools/leaderelection/resourcelock/interface.go:Interface```

```go
type Interface interface {
	// Get returns the LeaderElectionRecord
	Get() (*LeaderElectionRecord, []byte, error)

	// Create attempts to create a LeaderElectionRecord
	Create(ler LeaderElectionRecord) error

	// Update will update and existing LeaderElectionRecord
	Update(ler LeaderElectionRecord) error

	// RecordEvent is used to record events
	RecordEvent(string)

	// Identity will return the locks Identity
	Identity() string

	// Describe is used to convert details on current resource lock
	// into a string
	Describe() string
}
```

观察具体实现, 会发现主要通过对具体竞态资源的 annotations 字段进行更新来完成. 选举成功的Leader信息, 将会被以key ```control-plane.alpha.kubernetes.io/leader``` 写入annotation.

每个候选人都持有一些信息, 一旦选举成功, 就会被写入:
```go
type LeaderElectionRecord struct {
	// HolderIdentity is the ID that owns the lease. If empty, no one owns this lease and
	// all callers may acquire. Versions of this library prior to Kubernetes 1.14 will not
	// attempt to acquire leases with empty identities and will wait for the full lease
	// interval to expire before attempting to reacquire. This value is set to empty when
    // a client voluntarily steps down.
    // 候选人的唯一ID
    HolderIdentity       string      `json:"holderIdentity"`
    // 下一轮间隔
    LeaseDurationSeconds int         `json:"leaseDurationSeconds"`
    // 竞选成功的时间
    AcquireTime          metav1.Time `json:"acquireTime"`
    // 连任时间
    RenewTime            metav1.Time `json:"renewTime"`
    // 竞态计数器, 不同的候选人竞选成功后, 会将该值+1
	LeaderTransitions    int         `json:"leaderTransitions"`
}
```

## 选举过程

选举过程主要实现在 ```k8s.io/client-go/tools/leaderelection.go``` 中.

0. 填写自己的选票, 将 RenewTime/AcquireTime 设置为当前时间

```go
now := metav1.Now()
leaderElectionRecord := rl.LeaderElectionRecord{
    HolderIdentity:       le.config.Lock.Identity(),
    LeaseDurationSeconds: int(le.config.LeaseDuration / time.Second),
    RenewTime:            now,
    AcquireTime:          now,
}
```

1. 获取竞态资源, 如果资源不存在, 创建一个.

* 创建成功, 则选举立即成功. 
* 创建失败, 选举结束, 等待下一轮重试.

```go
oldLeaderElectionRecord, oldLeaderElectionRawRecord, err := le.config.Lock.Get()
if err != nil {
    if !errors.IsNotFound(err) {
        klog.Errorf("error retrieving resource lock %v: %v", le.config.Lock.Describe(), err)
        return false
    }
    if err = le.config.Lock.Create(leaderElectionRecord); err != nil {
        klog.Errorf("error initially creating leader election record: %v", err)
        return false
    }
    le.observedRecord = leaderElectionRecord
    le.observedTime = le.clock.Now()
    return true
}
```

2. 检查竞态资源, 如果目前的Leader不是自己, 而且现在还没有到Leader卸任时间, 则认可目前的Leader, 放弃竞选.

```go
// 2. Record obtained, check the Identity & Time
if !bytes.Equal(le.observedRawRecord, oldLeaderElectionRawRecord) {
    le.observedRecord = *oldLeaderElectionRecord
    le.observedRawRecord = oldLeaderElectionRawRecord
    le.observedTime = le.clock.Now()
}
if len(oldLeaderElectionRecord.HolderIdentity) > 0 &&
    le.observedTime.Add(le.config.LeaseDuration).After(now.Time) &&
    !le.IsLeader() {
    klog.V(4).Infof("lock is held by %v and has not yet expired", oldLeaderElectionRecord.HolderIdentity)
    return false
}
```

3. 如果自己是Leader, AcquireTime/LeaderTransitions保持不变, 否则 LeaderTransitions+1

```go
// 3. We're going to try to update. The leaderElectionRecord is set to it's default
// here. Let's correct it before updating.
if le.IsLeader() {
    leaderElectionRecord.AcquireTime = oldLeaderElectionRecord.AcquireTime
    leaderElectionRecord.LeaderTransitions = oldLeaderElectionRecord.LeaderTransitions
} else {
    leaderElectionRecord.LeaderTransitions = oldLeaderElectionRecord.LeaderTransitions + 1
}
```

4. 尝试将竞态资源更新成自己的选票

* 更新失败, 则可能其它候选人已经成功, 竞选失败
* 更新成功, 连任/选举成功

```go
// update the lock itself
if err = le.config.Lock.Update(leaderElectionRecord); err != nil {
    klog.Errorf("Failed to update lock: %v", err)
    return false
}
le.observedRecord = leaderElectionRecord
le.observedTime = le.clock.Now()
return true
```

## 主动卸任

当组件退出时, 需要主动卸任. 通过将竞态资源更新为 LeaderTransitions保持不变, 其它为空的值.

## 边缘情况探讨

* 初始启动

谁成功创建了边缘资源, 就成功竞选Leader.

* Leader卸任/Leader失败

Leader卸任和失败本质都差不多, 切换到其它候选人的时间差最大为 RetryPeriod.

* 时钟不同步

假设两个Controller分别运行于不同节点, 他们通过系统调用获得的时间取自各自的操作系统.

当时钟不同步时, 将导致候选者尝试参加竞选时, 对竞选的时机产生不同意见.

1. 候选者时钟落后于Leader, 将导致Leader卸任/失败后, 候选者迟迟不能参与竞选, 直到弥补这段时间差为止.

2. 候选者时钟超前于Leader, 将导致Leader任期内, 候选者尝试竞选. 这时由于只有一个候选者发起竞选, 它将顺利成功, 而此时Leader将直到下个RetryPeriod才能发现该问题, 期间将会出现两个Leader.

