---
title: "Helm(2): 创建和语法"
date: 2020-02-21T11:45:13+08:00
tags: "Helm"
---
首先, 创建一个Chart:  helm create mychart
接下来, 讨论的假设前提是你已经熟悉Go Template的基本用法.

# Flatten
```
configSvcName: mysrv
configSvcUrl: http://example.com
```
# Use Flatten
```
name: "{{ .Values.configSvcName }}
url: "{{ .Values.configSvcUrl }}
```
# Nested
```
config:
  svc:
    name: mysrv
    url: http://example.com
``` 
# Use Nested
name: "{{ .Values.config.svc.name }}"
url: "{{ .Values.config.svc.url }}"
移除左侧/右侧空格: -

通过在分隔符左侧或右侧增加-, 能够起到移除左侧或右侧所有空白的效果, 例如:
```
"{{23 -}} < {{- 45}}"
```
生成的字符串是:
```
23<45
```
此外, -也可以用来减少空格浪费的空间, 例如一些带缩进的控制语句if/else/with/for等, 本身不会嵌入到最终生成的YAML中, 但我们为了查看方便, 会对其增加缩进和换行, 这些空白没有必要带进最终渲染出的YAML中:
```
{{- with .Values.affinity }}
  affinity:
    {{- toYaml . | nindent 8 }}
{{- end }}
```
需要注意, 换行本身也是空白字符!



管道: |

与Linux下的命令管道类似, Go Template中同样也可以使用管道, 例如:
```
{{ "abc" | upper }}
```
生成:
```
ABC
```

转义字符串: ```quote```

yaml在处理数字,字符串和注释时存在微妙的变化, 为了安全的转义它们, 可以使用quote:
```
{{ .Values.text | quote }}
```
或
```
{{ quote .Values.text }}
```


左侧增加N个空格缩进: ```nindent N```

由于yaml通过缩进来区分层级, 在编写Kubernetes的yaml文件时, 可以结合 - 和 nindent, 来精确的制造偏移, 防止直接编辑空格出现失误, 而肉眼难以检查.
```
affinity:
  {{- toYaml . | nindent 8 }}
```

默认值: ```default XXX```
```
drink: {{ .Values.favorite.drink | default "tea" | quote }}
```