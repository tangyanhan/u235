---
title: "Go中通过Lua操纵Redis"
date: 2017-11-20T22:23:04+08:00
categories: ["分布式", "Go"]
---

为了在我的一个基本库中降低与Redis的通讯成本，我将一系列操作封装到LUA脚本中，借助Redis提供的EVAL命令来简化操作。
EVAL能够提供的特性：
1. 可以在LUA脚本中封装若干操作，如果有多条Redis指令，封装好之后只需向Redis一次性发送所有参数即可获得结果
2. Redis可以保证Lua脚本运行期间不会有其他命令插入执行，提供像数据库事务一样的原子性
3. Redis会根据脚本的SHA值缓存脚本，已经缓存过的脚本不需要再次传输Lua代码，减少了通信成本，此外在自己代码中改变Lua脚本，执行时Redis必定也会使用最新的代码。

导入常见的Go库如	"github.com/go-redis/redis"，就可以实现以下代码。

## 生成一段Lua脚本
    // KEYS: key for record
    // ARGV: fieldName, currentUnixTimestamp, recordTTL
    // Update expire field of record key to current timestamp, and renew key expiration
    var updateRecordExpireScript = redis.NewScript(`
    redis.call("EXPIRE", KEYS[1], ARGV[3])
    redis.call("HSET", KEYS[1], ARGV[1], ARGV[2])
    return 1
    `)
该变量创建时，Lua代码不会被执行，也不需要有已存的Redis连接。
Redis提供的Lua脚本支持，默认有KEYS、ARGV两个数组，KEYS代表脚本运行时传入的若干键值，ARGV代表传入的若干参数。由于Lua代码需要保持简洁，难免难以读懂，最好为这些参数写一些注释
注意：上面一段代码使用``跨行，`所在的行虽然空白回车，也会被认为是一行，报错时不要看错代码行号。

## 运行一段Lua脚本
        updateRecordExpireScript.Run(c.Client, []string{recordKey(key)}, 
										expireField,
										time.Now().UTC().UnixNano(), int64(c.opt.RecordTTL/time.Second)).Err()

运行时，Run将会先通过EVALSHA尝试通过缓存运行脚本。如果没有缓存，则使用EVAL运行，这时Lua脚本才会被整个传入Redis。

## Lua脚本的限制
1. Redis不提供引入额外的包，例如os等，只有redis这一个包可用。
2. Lua脚本将会在一个函数中运行，所有变量必须使用local声明
3. return返回多个值时，Redis将会只给你第一个

## 脚本中的类型限制
1. 脚本返回nil时，Go中得到的是err = redis.Nil（与Get找不到值相同）
2. 脚本返回false时，Go中得到的是nil，脚本返回true时，Go中得到的是int64类型的1
3. 脚本返回{"ok": ...}时，Go中得到的是redis的status类型（true/false)
4. 脚本返回{"err": ...}时，Go中得到的是err值，也可以通过return redis.error_reply("My Error")达成
5. 脚本返回number类型时，Go中得到的是int64类型
6. 传入脚本的KEYS/ARGV中的值一律为string类型，要转换为数字类型应当使用to_number

## 如果脚本运行了很久会发生什么？
Lua脚本运行期间，为了避免被其他操作污染数据，这期间将不能执行其它命令，一直等到执行完毕才可以继续执行其它请求。当Lua脚本执行时间超过了lua-time-limit时，其他请求将会收到Busy错误，除非这些请求是SCRIPT KILL（杀掉脚本）或者SHUTDOWN NOSAVE（不保存结果直接关闭Redis）

更多内容参考以下地址，我这里主要是根据使用Go的经验提供一些总结。
https://redis.io/commands/eval

一段更“复杂”的脚本，它要求在获取一个key值时，如果该值访问较多，就延长生存周期。此外还要比较更新时间，如果不需要更新，则直接返回取到的值，否则返回redis.Nil

```go
	// KEYS: rec:key, key
	// ARGV: currentUnixTimestamp, hotHit, recordTTL, ttl
	// When there's a hit, 
	var fetchRecordScript = redis.NewScript(`
	local value = redis.call("GET", KEYS[2])
	if(value == nil) then return nil end
	local hit = redis.call("HINCRBY", KEYS[1], "hit", 1)
	redis.call("EXPIRE", KEYS[1], ARGV[3])
	local minHotHit = tonumber(ARGV[2])
	local keyTTL = tonumber(ARGV[4])
	if(hit > minHotHit)then
		keyTTL = keyTTL * 2
	end
	redis.call("EXPIRE", KEYS[2], keyTTL)
	local expire = tonumber(redis.call("HGET", KEYS[1], "expire"))        
	local unixTime = tonumber(ARGV[1])
	if(expire == nil or expire < unixTime) then
		return nil
	else
		return value
	end
	`)
	// KEYS: key for record
	// ARGV: fieldName, currentUnixTimestamp, recordTTL
	// Update expire field of record key to current timestamp, and renew key expiration
	var updateRecordExpireScript = redis.NewScript(`
	redis.call("EXPIRE", KEYS[1], ARGV[3])
	redis.call("HSET", KEYS[1], ARGV[1], ARGV[2])
	return 1
	`)
```