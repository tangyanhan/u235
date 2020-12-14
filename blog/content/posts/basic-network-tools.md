---
title: "有关端口的基本网络工具"
date: 2020-12-03T14:00:51+08:00
tags: "网络"
---

# netstat

netstat在查看网络端口占用时速度比较快，也是比较常用的工具。

netstat 在显示网络连接时有一堆的参数可以叠加:

```
netstat  [address_family_options] [--tcp|-t] [--udp|-u] [--udplite|-U] [--sctp|-S] [--raw|-w] [--l2cap|-2] [--rfcomm|-f] [--listening|-l] [--all|-a]
[--numeric|-n] [--numeric-hosts] [--numeric-ports] [--numeric-users] [--symbolic|-N] [--extend|-e[--extend|-e]] [--timers|-o] [--program|-p] [--ver‐
bose|-v] [--continuous|-c] [--wide|-W]

address_family_options:

[-4|--inet]  [-6|--inet6]  [--protocol={inet,inet6,unix,ipx,ax25,netrom,ddp,bluetooth,  ...  }  ] [--unix|-x] [--inet|--ip|--tcpip] [--ax25] [--x25]
[--rose] [--ash] [--bluetooth] [--ipx] [--netrom] [--ddp|--appletalk] [--econet|--ec]

```

显示所有tcp/udp监听端口情况及进程信息，以数字形式显示IP:

```
netstat -tunlp
```

显示系统所有端口占用情况:

```
netstat -anp
```

显示所有TCP IPv4端口占用情况:

```
netstat -anpt4
```

按状态统计所有TCP IPv4使用情况:

```
netstat -nat4 | awk '{print $6}' | sort | uniq -c | sort -r
```

# lsof(list open files)

由于linux下一切东西都是文件，因此这个命令能够做各种事情。

详细版说明:https://www.cnblogs.com/sparkbj/p/7161669.html

常见的网络相关操作：

## 列出所有网络连接

```lsof -i```

## 查看目前端口使用情况

`lsof -i [46][proto][@host|addr][:svc_list|port_list]`

* 查看TCP端口使用情况: `lsof -i tcp`
* 查看TCP IPv4端口使用情况: `lsof -i 4tcp`
* 查看UDP端口使用情况: `lsof -i udp`
* 查看TCP端口3306被谁占用: `lsof -i tcp:3306`

## 查看进程文件打开情况

* 查看指定进程号打开的文件: `lsof -p [PID]`
* 查看多个进程号打开的文件： `lsof -p [PID1],[PID2]`
* 列出名字为NAME开头的进程打开的文件: `lsof -c [NAME]`
* 列出名字开头为名字A或B开头的进程打开的文件： `lsof -c A -c B`

