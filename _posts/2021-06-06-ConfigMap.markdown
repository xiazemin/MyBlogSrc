---
title: ConfigMap
layout: post
category: k8s
author: 夏泽民
---
ConfigMap 是一种 API 对象，用来将非机密性的数据保存到键值对中。使用时， Pods 可以将其用作环境变量、命令行参数或者存储卷中的配置文件。

ConfigMap 将您的环境配置信息和 容器镜像 解耦，便于应用配置的修改

ConfigMap 并不提供保密或者加密功能。 如果你想存储的数据是机密的，请使用 Secret， 或者使用其他第三方工具来保证你的数据的私密性，而不是用 ConfigMap。
<!-- more -->
使用 ConfigMap 来将你的配置数据和应用程序代码分开。

比如，假设你正在开发一个应用，它可以在你自己的电脑上（用于开发）和在云上 （用于实际流量）运行。 你的代码里有一段是用于查看环境变量 DATABASE_HOST，在本地运行时， 你将这个变量设置为 localhost，在云上，你将其设置为引用 Kubernetes 集群中的 公开数据库组件的 服务。

这让你可以获取在云中运行的容器镜像，并且如果有需要的话，在本地调试完全相同的代码。

ConfigMap 在设计上不是用来保存大量数据的。在 ConfigMap 中保存的数据不可超过 1 MiB。如果你需要保存超出此尺寸限制的数据，你可能希望考虑挂载存储卷 或者使用独立的数据库或者文件服务。

ConfigMap 对象
ConfigMap 是一个 API 对象， 让你可以存储其他对象所需要使用的配置。 和其他 Kubernetes 对象都有一个 spec 不同的是，ConfigMap 使用 data 和 binaryData 字段。这些字段能够接收键-值对作为其取值。data 和 binaryData 字段都是可选的。data 字段设计用来保存 UTF-8 字节序列，而 binaryData 则 被设计用来保存二进制数据作为 base64 编码的字串。

ConfigMap 的名字必须是一个合法的 DNS 子域名。

data 或 binaryData 字段下面的每个键的名称都必须由字母数字字符或者 -、_ 或 . 组成。在 data 下保存的键名不可以与在 binaryData 下 出现的键名有重叠。

从 v1.19 开始，你可以添加一个 immutable 字段到 ConfigMap 定义中，创建 不可变更的 ConfigMap。



apiVersion: v1
kind: ConfigMap
metadata:
  name: game-demo
data:
  # 类属性键；每一个键都映射到一个简单的值
  player_initial_lives: "3"
  ui_properties_file_name: "user-interface.properties"

  # 类文件键
  game.properties: |
    enemy.types=aliens,monsters
    player.maximum-lives=5    
  user-interface.properties: |
    color.good=purple
    color.bad=yellow
    allow.textmode=true    
    

你可以使用四种方式来使用 ConfigMap 配置 Pod 中的容器：

在容器命令和参数内
容器的环境变量
在只读卷里面添加一个文件，让应用来读取
编写代码在 Pod 中运行，使用 Kubernetes API 来读取 ConfigMap

        - name: PLAYER_INITIAL_LIVES # 请注意这里和 ConfigMap 中的键名是不一样的
          valueFrom:
            configMapKeyRef:
              name: game-demo           # 这个值来自 ConfigMap
              key: player_initial_lives # 需要取值的键
            
            https://kubernetes.io/zh/docs/concepts/configuration/configmap/
            
            https://kubernetes.io/zh/docs/tasks/configure-pod-container/configure-pod-configmap/
            
            https://cloud.google.com/kubernetes-engine/docs/concepts/configmap
            
            https://stackoverflow.com/questions/42101808/ingress-gives-502-error
            
            https://rancher.com/docs/rancher/v1.6/en/cattle/adding-load-balancers/
            
            