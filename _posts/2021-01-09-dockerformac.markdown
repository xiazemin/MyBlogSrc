---
title: docker for mac k8s 无法启动	
layout: post
category: k8s
author: 夏泽民
---
安装k8s大致有2种方式，minikube和Docker Desktop

https://github.com/docker/for-mac/issues/2990

https://github.com/AliyunContainerService/k8s-for-docker-desktop/issues/116

https://github.com/kubernetes/kubernetes/issues/71647

https://github.com/gotok8s/k8s-docker-desktop-for-mac
https://github.com/AliyunContainerService/k8s-for-docker-desktop

但是由于众所周知的原因, 国内的网络下不能很方便的下载 Kubernetes 集群所需要的镜像, 导致集群启用失败. 这里提供了一个简单的方法, 利用 GitHub Actions 实现 k8s.gcr.io 上 kubernetes 依赖镜像自动同步到 Docker Hub 上指定的仓库中。 通过 load_images.sh 将所需镜像从 Docker Hub 的同步仓库中取回，并重新打上原始的tag.
<!-- more -->
第一步 克隆详细

git clone https://github.com/gotok8s/k8s-docker-desktop-for-mac.git

第二步 进入 k8s-docker-desktop-for-mac项目，拉取镜像

./load_images.sh

第三步 打开docker 配置页面，enable k8s。需要等k8s start一会

https://www.jianshu.com/p/a6abdc6f76e1

https://github.com/gotok8s/k8s-docker-desktop-for-mac

https://docs.docker.com/docker-for-mac/release-notes/

一、背景
最近想要下载 neo4j 的 docker 镜像，发现速度不是一般的慢，囧…

于是乎，类似于 maven 有国内镜像，docker 是不是也有呢？

搜了一下，的确有。

二、用法
1、打开 docker 选择 Preferences

图片描述

2、切换到 Daemon 选项卡，在 Registry mirrors 添加想要添加的国内镜像

图片描述

如：

{
  "features": {
    "buildkit": true
  },
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn"
  ],
  "experimental": false
}

https://docker.mirrors.ustc.edu.cn
https://hub-mirror.c.163.com


https://www.cnblogs.com/shengulong/p/10261707.html
https://hub.docker.com/editions/community/docker-ce-desktop-mac
https://github.com/wubiaowp/kubernetes-for-docker-desktop-mac/tree/master/v1.16.5
https://docs.docker.com/get-started/kube-deploy/

https://docs.docker.com/docker-for-mac/release-notes/

