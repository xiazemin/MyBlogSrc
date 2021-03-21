---
title: nat64_v4_requires_synthesis docker for mac k8s 无法使用
layout: post
category: k8s
author: 夏泽民
---
https://github.com/AliyunContainerService/k8s-for-docker-desktop

https://github.com/daniel-hutao/k8s-source-code-analysis

Docker启用K8S后新增的两个目录
~/.kube
~/Library/Group Containers/group.com.docker/pki
注意部分网友对解决Docker长时间无法启动K8S建议的解决方案是删除这两个目录然后重启

https://www.jianshu.com/p/f09f7421e841
https://blog.csdn.net/weixin_31558841/article/details/112072166
https://github.com/AliyunContainerService/k8s-for-docker-desktop/issues/106
https://github.com/maguowei/k8s-docker-desktop-for-mac/issues/16

https://github.com/AliyunContainerService/k8s-for-docker-desktop/issues/123

如果在Kubernetes部署的过程中出现问题，可以通过docker desktop应用日志获得实时日志信息：

pred='process matches ".*(ocker|vpnkit).*"
  || (process in {"taskgated-helper", "launchservicesd", "kernel"} && eventMessage contains[c] "docker")'
/usr/bin/log stream --style syslog --level=debug --color=always --predicate "$pred"

http://hutao.tech/k8s-source-code-analysis/prepare/debug-environment-3node.html

https://docs.docker.com/docker-for-mac/release-notes/
https://github.com/docker/awesome-compose
<!-- more -->
https://www.jianshu.com/p/22c497ffe191
https://github.com/wubiaowp/kubernetes-for-docker-desktop-mac


https://blog.csdn.net/weixin_39954682/article/details/111554523


rm -rf ~/.kuberm -rf ~/.minikuberm -rf /usr/local/bin/minikube

rm -rf ~/Library/Group\ Containers/group.com.docker/pki

rm -rf ~/.kube

 pred='process matches ".*(ocker|vpnkit).*"\n  || (process in {"taskgated-helper", "launchservicesd", "kernel"} && eventMessage contains[c] "docker")'\n/usr/bin/log stream --style syslog --level=debug --color=always --predicate "$pred"
  
  
sh load_images.sh


功能的方法


vi ~/Library/Group\ Containers/group.com.docker/settings.json

  "kubernetesEnabled": false,
  "showKubernetesSystemContainers": false,
  "kubernetesInitialInstallPerformed": false,
  
  
  
Docker 3.2.2版本

修改k8s-for-docker-desktop
的images.properties
v1.19.3 为v1.19.7