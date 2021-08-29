---
title: annotations
layout: post
category: k8s
author: 夏泽民
---
将这些Kubernetes批注添加到特定的Ingress对象以自定义其行为

可以使用 --annotations-prefix command line argument, 命令行参数更改注释前缀，但是默认值为nginx.ingress.kubernetes.io
https://rocdu.gitbook.io/ingress-nginx-docs-cn/docs/user-guide/nginx-configuration/annotations

Lua Resty WAF
使用lua-resty-waf- *注释，我们可以启用和控制 lua-resty-waf 每个location的Web应用防火墙.
https://docs.openshift.com/container-platform/4.6/rest_api/network_apis/ingress-networking-k8s-io-v1.html


<!-- more -->
The Kubernetes API
The core of Kubernetes' control plane is the API server. The API server exposes an HTTP API that lets end users, different parts of your cluster, and external components communicate with one another.

https://kubernetes.io/docs/concepts/overview/kubernetes-api/

Kubernetes API是系统描述性配置的基础。 Kubectl 命令行工具被用于创建，更新，删除，获取API对象。

Kubernetes 通过API资源存储自己序列化状态(现在存储在etcd)。

Kubernetes 被分成多个组件，各部分通过API相互交互。

http://kubernetes.kansea.com/docs/api/

Nginx Ingress注解Annotations
Nginx Ingress 注解使用在 Ingress 资源实例中，用以设置当前 Ingress 资源实例中 Nginx 虚拟主机的相关配置，对应配置的是 Nginx 当前虚拟主机的 server 指令域内容。在与 Nginx Ingress 配置映射具有相同功能配置时，将按照所在指令域层级遵循 Nginx 配置规则覆盖。

Nginx Ingress注解按照配置功能有如下分类。
1、Nginx原生配置指令
支持在注解中添加 Nginx 原生配置指令。配置说明如下表所示。

注解 类型 功能描述
nginx.ingress.kubernetes.io/server-snippet string 在 server 指令域添加 Nginx 配置指令
nginx.ingress.kubernetes.io/configuration-snippet string 在 location 指令域添加 Nginx 配置指令

https://putianhui.cn/index.php/K8S/145.html

https://github.com/kubernetes/ingress-nginx/issues/3570

https://putianhui.cn/index.php/K8S/145.html
Nginx Ingress 注解使用在 Ingress 资源实例中，用以设置当前 Ingress 资源实例中 Nginx 虚拟主机的相关配置，对应配置的是 Nginx 当前虚拟主机的 server 指令域内容。在与 Nginx Ingress 配置映射具有相同功能配置时，将按照所在指令域层级遵循 Nginx 配置规则覆盖。

Nginx Ingress注解按照配置功能有如下分类。
1、Nginx原生配置指令
支持在注解中添加 Nginx 原生配置指令。配置说明如下表所示。

注解 类型 功能描述
nginx.ingress.kubernetes.io/server-snippet string 在 server 指令域添加 Nginx 配置指令
nginx.ingress.kubernetes.io/configuration-snippet string 在 location 指令域添加 Nginx 配置指令

https://kubernetes.io/zh/docs/tasks/manage-kubernetes-objects/declarative-config/

https://github.com/kubernetes/ingress-nginx/issues/4940

https://cloud.ibm.com/docs/containers?topic=containers-comm-ingress-annotations&locale=ko

https://zivgitlab.uni-muenster.de/c_wied05/wordpress-kube/-/blob/master/wordpress-helm/templates/ingress.yml

kubectl.kubernetes.io/last-applied-configuration

