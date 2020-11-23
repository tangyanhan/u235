---
title: "Helm(4) 让编排兼容不同版本K8S"
date: 2020-04-25T13:09:47+08:00
tags: "Helm"
---

在使用Helm过程中, 经常会遇到编排需要兼容不同K8S版本的问题. 考虑如下场景:

1. 以前编写的Deployment资源, 其apiVersion为 apps/v1beta1, 但后来新的版本中已经改为 apps/v1,希望能兼容

2. 在K8S 1.11以前, 默认CRD既不支持subResources, 也不能够通过 UpdateStatus 更新状态, 必须使用 Update. 而在这之后, 必须使用 UpdateStatus 才能更新状态.

Helm 内置了一系列[内部对象](https://helm.sh/docs/chart_template_guide/builtin_objects/),可以针对这些情况进行编排.

# Deployment兼容多版本K8S

针对以上第一个问题, 我们可以直接在 _helpers.tpl 中加入以下内容:

```
{{/*
Define apiVersion for Deployment
*/}}
{{- define "deployApiVersion" -}} 
{{- if .Capabilities.APIVersions.Has "apps/v1beta1/Deployment" -}}
apps/v1beta1
{{- else -}}
apps/v1
{{- end -}}
{{- end -}}
```

这里判断这套K8S是否具备 apps/v1beta/Deployment , 如果有, 使用 apps/v1 ,否则就是旧版本的 apps/v1beta1

然后, Deployment引用这一段即可:
```yaml
kind: Deployment
apiVersion: {{ include "deployApiVersion" . }}
```

# 检测配置 disableSubresources 是否开启

经过查询文档, disableSubresources 默认开启是在 1.11 版本, 因此我们可以简单一点, 只判断小版本是否大于等于 "11", 当然这里没有判断大版本.

```
disableSubresources: {{ if ge .Capabilities.KubeVersion.Minor "11" }}false{{ else }}true{{ end }}
```

# Capabilities 实现

我们可能在某些场景下, 希望自己的程序能够直接去判断我们对接的K8S集群有哪些能力. 这部分实现我们可以参考 helm 源码. 在 ```pkg/action/action.go```中揭示了其实现方式:

```go
// capabilities builds a Capabilities from discovery information.
func (c *Configuration) getCapabilities() (*chartutil.Capabilities, error) {
        if c.Capabilities != nil {
                return c.Capabilities, nil
        }
        dc, err := c.RESTClientGetter.ToDiscoveryClient()
        if err != nil {
                return nil, errors.Wrap(err, "could not get Kubernetes discovery client")
        }
        // force a discovery cache invalidation to always fetch the latest server version/capabilities.
        dc.Invalidate()
        kubeVersion, err := dc.ServerVersion()
        if err != nil {
                return nil, errors.Wrap(err, "could not get server version from Kubernetes")
        }
        // Issue #6361:
        // Client-Go emits an error when an API service is registered but unimplemented.
        // We trap that error here and print a warning. But since the discovery client continues
        // building the API object, it is correctly populated with all valid APIs.
        // See https://github.com/kubernetes/kubernetes/issues/72051#issuecomment-521157642
        apiVersions, err := GetVersionSet(dc)
        if err != nil {
                if discovery.IsGroupDiscoveryFailedError(err) {
                        c.Log("WARNING: The Kubernetes server has an orphaned API service. Server reports: %s", err)
                        c.Log("WARNING: To fix this, kubectl delete apiservice <service-name>")
                } else {
                        return nil, errors.Wrap(err, "could not get apiVersions from Kubernetes")
                }
        }

        c.Capabilities = &chartutil.Capabilities{
                APIVersions: apiVersions,
                KubeVersion: chartutil.KubeVersion{
                        Version: kubeVersion.GitVersion,
                        Major:   kubeVersion.Major,
                        Minor:   kubeVersion.Minor,
                },
        }
        return c.Capabilities, nil
}
```


# dry-run 没问题的chart一定能够在集群上部署吗?

在解决第一个问题时, 我事先通过 -dry-run 测试我写出的chart, 结果安装到集群却失败了. 其实, 在Helm Chart进行dryrun时, 是不会与集群进行交互的.

在 helm 源码中, 它的解释是编排chart的作者的期望可能是连接到集群,但用户却不一定如此期望:

```go
        // A `helm template` or `helm install --dry-run` should not talk to the remote cluster.
        // It will break in interesting and exotic ways because other data (e.g. discovery)
        // is mocked. It is not up to the template author to decide when the user wants to
        // connect to the cluster. So when the user says to dry run, respect the user's
        // wishes and do not connect to the cluster.
        if !dryRun && c.RESTClientGetter != nil {
                rest, err := c.RESTClientGetter.ToRESTConfig()
                if err != nil {
                        return hs, b, "", err
                }
                files, err2 = engine.RenderWithClient(ch, values, rest)
        } else {
                files, err2 = engine.Render(ch, values)
        }
```