---
title: docker for mac k8s
layout: post
category: docker
author: 夏泽民
---
下载地址
https://docs.docker.com/docker-for-mac/release-notes/
由于我的版本是10.11.1，最后兼容的docker for mac 是
Docker Community Edition 18.06.1-ce-mac73 2018-08-29
其他版本不再支持，这也是支持k8s的最高版本的docker for mac
下载地址：https://download.docker.com/mac/stable/26764/Docker.dmg

<a href="{{site.url}}{{site.baseurl}}/img/Docker.dmg">本blog缓存</a>

<img src="{{site.url}}{{site.baseurl}}/img/k8s_docker_for_mac.png"/>
'''
$ docker container ls --format "table{\{.Names}\}\t{\{.Image }\}\t{\{.Command}\}"
NAMES                                                                                                                   IMAGE                                      COMMAND
k8s_sidecar_kube-dns-86f4d74b45-qtpg9_kube-system_ce960bd7-97c3-11e9-a0e2-025000000001_0                                k8s.gcr.io/k8s-dns-sidecar-amd64           "/sidecar --v=2 --lo…"
k8s_dnsmasq_kube-dns-86f4d74b45-qtpg9_kube-system_ce960bd7-97c3-11e9-a0e2-025000000001_0                                k8s.gcr.io/k8s-dns-dnsmasq-nanny-amd64     "/dnsmasq-nanny -v=2…"
'''

$kubectl get namespaces
NAME          STATUS    AGE
default       Active    12m
docker        Active    12m
kube-public   Active    12m
kube-system   Active    12m
$kubectl get posts --namespace kube-system
error: the server doesn't have a resource type "posts"
<!-- more -->
不同的版本对 kubernetes 支持的版本不一样，请在安装后查看具体的版本号。


Docker Version 18.06.1-ce-mac73 (26764) 》kubernetes 1.10.3
Docker 2.0.0.0 》kubernetes 1.10.3
Docker Version 2.0.0.3 (31259) 》kubernetes 1.10.11


从github上直接下载指定版本的 kubectl，根据你的操作系统，下载指定平台版本，比如 macOS 的：

https://github.com/maguowei/k8s-docker-desktop-for-mac

kubernetes 1.10.3  下载地址： kubernetes-client-darwin-amd64.tar.gz


kubernetes 1.10.11 下载地址：kubernetes-client-darwin-amd64.tar.gz


下载后解压到某个目录下。
打开终端命令行，cd 进入这个目录，执行以下脚本，将其变更为可执行命令，同时移动到系统特定目录下。
chmod +x kubectl && mv kubectl /usr/local/bin/kubectl

我们可以看下 kubectl 的版本号：
kubectl version

结果如下：

步骤2：手动拉取 Kubernetes 需要的镜像
根据目前我用的版本 Version 2.0.0.3 (31259) 集成 Kubernetes: v1.10.11。我要想把 Kubernetes 启动起来，需要手动下载 Kubernetes 组件的镜像。然后，进入 Docker 图标下拉菜单 “Preferences” > “Kubernetes”，启用它。
因为在阿里云上，有同步镜像的组件，我们就不需要翻到官网下载了。借鉴网上找到脚本 k8s-deploy，进行改良一下，加入了 Dashboard 组件进去。大家如果只使用 kubectl 来控制 Kubernetes 的话，可以自己将这部分去掉。对于新手来说，可能有个网页界面，看着舒服些。
不过需要注意的是，因为 Dashboard 的版本是单独演进的，要了解最新版本是多少，需要查看 kubernetes-dashboard.yaml 文件。
...
image: k8s.gcr.io/kubernetes-dashboard-amd64:v1.10.0
...

现在，创建一个脚本文件：docker-k8s-images.sh
#!/bin/bash

set -e 
# Check version in https://kubernetes.io/docs/reference/setup-tools/kubeadm/kubeadm-init/
# Search "Running kubeadm without an internet connection"
# For running kubeadm without an internet connection you have to pre-pull the required master images for the version of choice:
KUBE_VERSION=v1.10.11
KUBE_DASHBOARD_VERSION=v1.10.1
KUBE_PAUSE_VERSION=3.1
ETCD_VERSION=3.1.12
DNS_VERSION=1.14.8
GCR_URL=k8s.gcr.io
ALIYUN_URL=registry.cn-hangzhou.aliyuncs.com/google_containers

images=(kube-proxy-amd64:${KUBE_VERSION}
kube-scheduler-amd64:${KUBE_VERSION}
kube-controller-manager-amd64:${KUBE_VERSION}
kube-apiserver-amd64:${KUBE_VERSION}
pause-amd64:${KUBE_PAUSE_VERSION}
etcd-amd64:${ETCD_VERSION}
k8s-dns-sidecar-amd64:${DNS_VERSION}
k8s-dns-kube-dns-amd64:${DNS_VERSION}
k8s-dns-dnsmasq-nanny-amd64:${DNS_VERSION}
kubernetes-dashboard-amd64:${KUBE_DASHBOARD_VERSION}) 

for imageName in ${images[@]} ; do
docker pull $ALIYUN_URL/$imageName
docker tag $ALIYUN_URL/$imageName $GCR_URL/$imageName
docker rmi $ALIYUN_URL/$imageName
done

docker images
备注：如果你查看到的 Dashboard 有新版本了，修改一下脚本中的 KUBE_DASHBOARD_VERSION。
然后，我们运行一下这个脚本：
./docker-k8s-images.sh

看到最终的运行结果：
这时候我们再回到 Docker，就可以看到 Kubernetes 已正常启动了。
针对老版本 Kunbernetes 1.10.3 版本，请更改相关配置：
KUBE_VERSION=v1.10.3
KUBE_DASHBOARD_VERSION=v1.10.0
KUBE_PAUSE_VERSION=3.1
ETCD_VERSION=3.1.12
DNS_VERSION=1.14.8

步骤3：启动 Kubernetes Dashboard（可选）
接下来，我们要想启动 Kubernetes Dashboard，还得在集群中部署一下 kubernetes-dashboard.yaml。
kubectl create -f https://github.com/kubernetes/dashboard/tree/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml

部署成功后，我们进行启动 proxy。
kubectl proxy

Starting to serve on 127.0.0.1:8001
这时候，打开浏览器，访问 Kubernetes Dashboard
通过以下脚本，填写 kubeconfig 的 Token 信息（如果不操作这一步，就会提示 config 信息不全）。
#!/bin/bash
TOKEN=$(kubectl -n kube-system describe secret default| awk '$1=="token:"{print $2}')
kubectl config set-credentials docker-for-desktop --token="${TOKEN}"

选择 kubeconfig 文件，使用“shift + command + .”打开 $HOME 下隐藏目录文件 ./kube/config，点击“登录”，就可以认证成功，进入首页了。
