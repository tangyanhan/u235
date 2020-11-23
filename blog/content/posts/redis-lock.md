---
title: "Redis分布式锁"
date: 2019-10-20T22:25:41+08:00
tags: ["分布式"]
---
# Redis的原子性
同一个Redis实例，它只以单个进程运行，并可以确保所有请求都是在同一个序列中执行的，因此可以保证Redis执行的语句是原子性的。 对于使用EVAL，通过LUA运行的多条语句，也可以保证像数据库事务一样具有原子性。
# 单实例Redis
一个Go的实现：https://github.com/bsm/redis-lock

单实例Redis只需借助SETNX（2.6.12后续版本只需SET  key value NX也可以做到）即可：

    LOCK(lockKey, ownerID):
    SET lockKey   ownerID  NX PX  expirationInMilliseconds

    Unlock(lockKey, ownerID):
    if GET(lockKey) == ownerID:
        DEL(lockKey)

为lockKey设置一个独一无二的值ownerID，这样在Unlock时，就不会出现lockKey正好被自动Expire删除后，原拥有者误将别人的锁释放掉的情况。

如果一系列操作需要多个Redis操作，那么应当EVAL将多个操作封装到同一段LUA代码中，否则可能导致多次通讯时差中出现意外。


这种情况仅适用于同一条key存在于同一个Redis实例的情况，例如Redis只有一个，或者不使用Master-Slave的Redis集群，例如无slave的hashring集群（利用类似一致性环形哈希计算key，最终请求落到特定节点上）

如果是使用Master-Slave的Redis集群，同一个key可以存在若干个备份，写入master的数据同步到slave中需要一段时间。考虑以下情形：

1. Client A在master上获取了对key的锁： key:A
2. master短暂故障（网络故障，重启等），但key:A尚未同步到slave
3. Client B向slave请求获得锁 key:B成功
4. master恢复运行，现在Client A、B都认为自己获得了锁

# Red-lock分布式锁算法
一个Go的实现：https://github.com/go-redsync/redsync
在使用master-slave集群时，上述锁的问题在于同步过程中发生了冲突，因此一种解决方案是同时在多个节点上获取锁，当在多数节点成功时，就意味着其它client必然只能获得少数成功，该算法来自https://redis.io/topics/distlock，根据Redis文档的说法，并未在生产环境全面验证，但理论上是行的通的。算法如下：

假定我们想要锁定的时间为T，
1. 记录下当前时间start，以毫秒（ms）为单位
2. 借助多线程/协程同时向所有Redis实例请求获得锁 K:V, px=T，设置一个请求的超时时间，例如如果T为10s，那么我们可以设置请求超时为5-50ms，这样就可以避免在一个已经死掉的client上花费过多时间，但该值不应当低于网络通讯时间
3. 不论成功与否，所有过程结束之后，计算剩余锁的时间 t=now-start
4. 如果t<T，或者成功获得锁的数量不超过集群的半数，则认为失败
5. 如果获取锁失败，那么释放掉所有集群上的锁（仅限于K:V一致）。为什么不只释放自己成功获得锁的实例呢？考虑到集群中存在Master-Slave的同步机制，以及我们设置的请求超时，最终存有我们的锁的实例将不限于曾经成功的那些，因此必须对所有实例释放我们的锁。




