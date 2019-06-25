---
title: docker for mac k8s
layout: post
category: web
author: 夏泽民
---
https://github.com/maguowei/k8s-docker-desktop-for-mac
下载最新的Docker for Mac Edge 版本，跟普通mac软件一样安装，然后运行它，会在右上角菜单栏看到多了一个鲸鱼图标，这个图标表明了 Docker 的运行状态。
配置镜像加速地址
鉴于国内网络问题，国内从 Docker Hub 拉取镜像有时会遇到困难，此时可以配置镜像加速器。Docker 官方和国内很多云服务商都提供了国内加速器服务。

点击设置菜单
preferences->Daemon->basic->reposeriory mirrors	
设置镜像加速地址
<!-- more -->
启用k8s
点击设置菜单
preferences->Kubernetes->enable Kubernetes
点击启动k8s的checkbox，这里会拉取比较多的镜像，可能要等好一会儿。
检查k8s环境
可执行以下命令检查k8s环境
$ kubectl get nodes
NAME                 STATUS    ROLES     AGE       VERSION
docker-for-desktop   Ready     master    3h        v1.9.6
$ kubectl cluster-info
Kubernetes master is running at https://localhost:6443
KubeDNS is running at https://localhost:6443/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy
To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.

部署kubernetes-dashboard服务
按以下步骤部署k8s-dashboard服务
$ kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml
# 开发环境推荐用NodePort的方式访问dashboard，因此编辑一下该部署
$ kubectl -n kube-system edit service kubernetes-dashboard
# 这里将type: ClusterIP修改为type: NodePort
# 获取dashboard服务暴露的访问端口
$ kubectl -n kube-system get service kubernetes-dashboard
NAME                   TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)         AGE
kubernetes-dashboard   NodePort   10.98.82.248   <none>        443:31241/TCP   2h

按上述输出，dashboard服务暴露的访问端口是31241，因此可以用浏览器访问https://localhost:31241/，我们可以看到登录界面
此时可暂时直接跳过，进入到控制面板中
使用k8s本地开发环境
这里尝试用Skaffold往本地开发环境部署微服务应用。

安装Skaffold
1
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-darwin-amd64 && chmod +x skaffold && sudo mv skaffold /usr/local/bin
获取微服务示例代码
git clone https://github.com/GoogleContainerTools/skaffold
cd skaffold/examples/microservices
部署到本地k8s环境
skaffold run
# 获取leeroy-web服务暴露的访问端口
$ kubectl get service leeroy-web
NAME         TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
leeroy-web   NodePort   10.98.162.88   <none>        8080:30789/TCP   56m
按上述输出，dashboard服务暴露的访问端口是30789，因此可以用浏览器访问http://localhost:30789/

k8s的dashboard中检查部署

