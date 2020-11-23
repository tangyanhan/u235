---
title: "搭建适合日常工作的Linux桌面环境"
date: 2019-10-19T22:15:16+08:00
tags: "杂谈"
---

Linux稳定性好，Linux软件开放……不过等到决定把Linux当作日常工作用系统时，就一言难尽了……

我日常工作的需求有：

* 笔记本扩展屏幕

* Golang开发

* docker/kubernetes: 日常工作会使用Virtualbox启动1-1的K8S集群

* 输入法

* Git及文件对比

* 办公通信： Office365邮件，微信，企业微信

* 娱乐需求： 网易云音乐，播放本地音乐

* 无线投屏演示


我目前使用的笔记本是华硕灵耀，在使用不同发行版过程中遇到的坑有：

## CentOS

* 安装完毕后无法使用无线网卡，推测是内核较老缺少无线网卡
* Gnome用的是Gnome3, KDE用的还是10年前的塑料风KDE3,丑拒

## Ubuntu

* Unity3 桌面下，每天至少出现一次界面卡死，切换终端无反应，只能选择重启，另一个同事在小米笔记本上也是类似，可能是某些软件的兼容性出现问题

* Lxde桌面下，不能支持Fn系列快捷键，自带软件不支持多屏幕，必须手动添加软件并手动设定

* KDE桌面下一切表现都比较良好，但是自从 18.04.2版本后的一次系统更新后，引导必定黑屏，至今未找到原因，即使将内核切换回更新前的版本也无效，累觉不爱

## Fedora

这是我目前使用的版本，自从Kubuntu发生引导黑屏后更换数个系统，发现还是这个比较好用。目前已稳定运行两个多月，日常使用无死机，不过还是有一些问题:

* Gnome3版本对一些特殊的投屏等支持不佳，因此仍旧使用KDE
* KDE桌面环境下，当使用IBus时，将使VSCode内无法鼠标选中多行代码, 目前使用Fcitx+pinyin
* 当使用多屏幕时，如果已经设定只显示外接显示器，那么拔下视频线后，KDE桌面将不会自动切换，笔记本屏幕保持黑屏，这时需要Fn盲操作才能让桌面切换回来

# 目前的软件解决方案

* 扩展屏幕及Fn键支持： KDE5桌面

* Golang IDE: VSCode/LiteIDE

当使用go modules模式进行Golang开发时, VSCode 在Linux下存在gocode-mod占用超多内存的问题,因此在开了K8S集群后,有时候会比较卡,不得不用LiteIDE. LiteIDE作为一个Golang IDE对资源占用较小,勉强够用,但是代码提示并不进行缓存,因此经常出现跳转需要等几秒钟才能跳的问题.

* kubernetes: Virtualbox+kubeadm

Ubuntu的microk8s文档有点少，用了一段时间之后放弃了。minikube 号称支持kvm或virtualbox，曾经试过kvm+minikube的组合，代替virtualbox+minikube，很坑，各种莫名其妙的错误防不胜防，已放弃。
kubeadm是正式kubernetes的简化版，不像microk8s/minikube能够自动适应笔记本换IP的问题，不适合在家/在办公室随便玩,因此我的办法是使用两台虚拟机构成局域网搭建了一个1 Master 1 Node集群. 对内存和CPU消耗并不是很高.

* Git/文件对比： gitg, meld  (gitg有时会出现崩溃,并不是很靠谱)

* 输入法: fcitx + pinyin

经过实验, ibus和scim会在某些输入场合出现奇奇怪怪的bug,例如无法在vscode选中多行,这可能和GTK/KDE兼容性有关.

* 办公通信： 邮件，微信，企业微信通过 Virtualbox+Win7+无缝模式的方式， 32位Win7对内存占用还没有Goland高，Virtualbox提供的无缝模式可以让微信的窗口像是Linux自己的一样

* 截图: 使用KDE自带的 Spectacle , 曾经尝试过flameshot, 但flameshot会出现假死问题.

* 网易云音乐通过一个deb转yum的工具直接转换后安装， 本地音乐使用网易云或VLC

网易云音乐在Linux下吃CPU/内存比较多,有时候 K8S集群+VSCode+网易云音乐 一起开系统就会卡,这时候关掉网易云就可以略微改善情况.

* 无线投屏可以使用一个叫做airplay.jar的工具，不过需要得到AppleTV的IP地址才能用，目前还没尝试

* 其他Office等工具，LibreOffice/OpenOffice等效果并不理想，在我目前的软件组合下，WPS for Linux无法输入中文，不如直接虚拟机里用WPS或Office. 目前的方案是虚拟机开32位Win7,分配2核1G,里面安装WPS.
