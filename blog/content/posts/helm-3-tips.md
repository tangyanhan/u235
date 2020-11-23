---
title: "Helm(3): 小技巧"
date: 2020-02-24T20:51:00+08:00
tags: "Helm"
---

# 用变量组合成新的模板

假设我们的一个组件在Chart中组合出了一个配置变量, 然后正巧有好几个文件要用到它:

```
http://{{ .Release.Name }}.{{ .Release.Namespace }}:{{ .Values.service.servicePort }}
```

Golang Template已经自带一些如局部变量来做类似的事情, 但它的作用域无疑是有限的,也不利于统一管理. 这就有了 ```define```的用武之地.

这里```define```其实是Go模板自带的语法,可以轻易组合出新的嵌套模板, 基本语法是:

```yaml
{{- define "tmpl.addr" -}}
http://{{ .Release.Name }}.{{ .Release.Namespace }}:{{ .Values.service.servicePort }}
{{- end -}}
```

如果这个模板有很多文件会用到它, 那么可以将它放到一个专门的文件中, 例如 ```templates/_helpers.tpl```, 这样你不必费尽心机避免模板被当作Kubernetes Object的一部分, 也容易寻找和维护.

我们就可以在其它文件中这样使用它:
```
{{ template "tmpl.addr" . }}
```

在templates文件夹下, 只要是 ```_```开头的文件都不会被渲染作为Kubernetes Object,因此其实是可以随便起的. 在这个文件夹下, 还有一个不会被渲染成 Kubernetes Object的文件, 叫做```NOTES.txt```, 稍后再说.

# 安装完成后的叨逼叨

事实上, 很多时候负责安装的用户往往并不会仔细去研究你提供的安装文档(有时,你甚至没有一份详尽的文档). 往往用户安装完成之后, TA也不知道安装了些什么东西, 接下来该做什么. 如果你安装过一些类似的中间件, 例如 ```bitnami/kafka```, 那么你会发现安装完后, 它会打印出一堆信息, 告诉你该如何访问它, 以及其它注意事项等.

安装完成后打印出的信息放在哪儿了? 那就是 ```templates/NOTES.txt```. 上面提到过, 这个文件不会被当成Kubernetes Object渲染并安装到集群中, 因此不必担心这个问题. NOTES.txt 中同样也是模板语法, 并没有什么区别.

如果你是一个后端组件, 你可以用它来告诉使用者该怎么调用你的API, 还可以告诉你的同事该如何去配置nginx等等.

# CRD与卸载

我们在安装某些特殊的Helm Chart时, 可能会创建CRD这种集群资源. 在过往的Helm2中, 这都是很难缠的:

假设你在```templates```下的模板中创建了一个CRD, 那么, 当卸载它时, 它将删除集群资源CRD, 如果此时集群中同时还安装了另一个使用该集群的Controller, 那么Controller将会遇到错误.

当在一个集群安装两个这样的Chart时, Helm将会报告CRD这种资源已经存在, 导致无法正常安装.

在Helm2中, 使用的是一种叫做```crd-install```的钩子来解决这个问题,可以在```prometheus-operator```中看到这种遗留痕迹:
```
{{- if and .Values.prometheusOperator.enabled .Values.prometheusOperator.createCustomResource -}}
{{- range $path, $bytes := .Files.Glob "crds/*.yaml" }}
{{ $.Files.Get $path }}
---
{{- end }}
{{- end }}
```

在Helm3中专门为此提供了[解决方案](https://helm.sh/docs/topics/chart_best_practices/custom_resource_definitions/#install-a-crd-declaration-before-using-the-resource). 通过在目录下增加一个叫做```crds```的文件夹.

1. CRD不存在, 则Helm3安装所有的CRD
2. CRD存在, 则Helm3跳过安装CRD并打印一个warning
3. CRD存在, 指定 ```--skip-crds```, 则跳过安装crd

但默认情况下, 这仍然存在问题. 假设我们的程序从1.0升级到2.0, CRD定义必须更新才能正确匹配新版的Controller, 那么这时, Helm3并不能帮助我们更新CRD定义. 它所做的, 只是打印一个warning. 因此Helm2时代的钩子仍然有用武之地. 当然, 更简单的做法是在 NOTES.txt 中告诉用户, 记得 ```kubectl apply -f crds/```

另一种做法, 就是将chart分离, 将这些集群资源分离成单独的chart包. 无疑这增加了安装步骤, 但也使得helm不必一次性使用集群级权限了.

# Helmignore: 忽略文件,不必打包

```.helmignore```文件与```.gitignore```写法完全一样,所有匹配的文件都会在打包时被忽略
