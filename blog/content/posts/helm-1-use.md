---
title: "Helm(1):安装和使用"
date: 2020-02-21T11:40:26+08:00
tags: "Helm"
---
## 安装
可以从 https://github.com/helm/helm/releases/tag/v3.0.3 获取不同操作系统的版本
```
wget https://get.helm.sh/helm-v3.0.3-linux-amd64.tar.gz

tar xzvf helm-v3.0.3-linux-amd64.tar.gz

sudo cp helm /usr/local/bin/helm
```


与Helm2不同, helm不再需要在集群中维持一个tiller. Helm 3 默认使用与kubectl相同的配置(KUBECONFIG), 来对Kubernetes进行操作.

## 使用
类似Linux的rpm/deb, Helm能够安装和管理的软件包, 叫做Chart. Helm Chart能够从本地安装, 也可以从网络下载安装.

同样的, helm 3 可以通过命令,将网络地址作为自己的Repository.

而一个Chart安装到Kubernetes集群, 成为一组真实运行的资源, 就叫做一个Release. 同一个Chart在集群中的每次安装, 都会创建一个Release.



类似于Docker Hub, Helm 3 已经将 https://hub.helm.sh 作为自己的默认搜索源,  可以通过 helm search hub <keyword> 从 https://hub.helm.sh 寻找自己需要的Chart.

添加一个命名为azure的Repo:  
```
helm repo add azure http://mirror.azk8s.cn/kubernetes/charts
```
从Repo中搜索一个Chart:
```
helm search repo <keyword>
```


安装一个网络上的Chart:
```
helm install [NAME] [CHART] [flags]
```

例如: ```helm install redis azure/redis --namespace middleware```

将名为 azure/redis 的Chart安装到 middleware 命名空间中, 其release名为 redis



同样, 也可以从本地安装一个Chart.

将网络上的chart拉到本地: ```helm pull azure/redis```



卸载一个Release, 则 ```helm uninstall azure/redis```

如果有特定namespace, 可以使用 ```helm uninstall azure/redis --namespace xxx```



安装一个本地的Chart, 可以通过压缩包安装, 也可以通过未压缩的文件夹安装:
```
helm install pipeline -f config.yaml .
```
以 config.yaml 中的值作为配置, 安装当前文件夹下的Chart.

