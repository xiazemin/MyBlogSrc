---
title: k8s权限控制
layout: post
category: k8s
author: 夏泽民
---
正常运行一段时间的POD，突然有一天运行报错了，错误是没有操作目录的权限，查其原因，原来是镜像被更新了，镜像添加了操作用户，而被操作的目录（NFS目录）并不具备普通用户操作的权限。


DOCKER容器如何控制用户#
　　容器的操作用户其实是镜像控制的，打造镜像的时候有一个USER选项，运行USER [uid:gid] 就可以指定用什么用户来运行该镜像，这方面知识可以查看一下DOCKFILE编写语法。如果没有添加该语句，镜像默认使用root权限运行，而添加了该语句，DOCKFILE后面的命令就使用该用户权限运行。那么问题就来了，镜像里面的用户跟操作系统有什么关系，或者说，跟NFS目录权限有什么关系？
<!-- more -->
1.linux主机通过uid和gid来控制用户对目录的操作权限，docker容器中也是如此。

　　2.当docker容器中的操作用户为root时，他相当于宿主机上的root

　　3.当docker容器中的操作用户为非root时，根据其uid在宿主机上的权限限制获取对应权限

https://www.cnblogs.com/garfieldcgf/p/12055384.html

为 Pod 或容器配置安全性上下文
安全上下文（Security Context）定义 Pod 或 Container 的特权与访问控制设置。 安全上下文包括但不限于：

自主访问控制（Discretionary Access Control）：基于 用户 ID（UID）和组 ID（GID）. 来判定对对象（例如文件）的访问权限。
安全性增强的 Linux（SELinux）： 为对象赋予安全性标签。
以特权模式或者非特权模式运行。
Linux 权能: 为进程赋予 root 用户的部分特权而非全部特权。
AppArmor：使用程序框架来限制个别程序的权能。
Seccomp：过滤进程的系统调用。
AllowPrivilegeEscalation：控制进程是否可以获得超出其父进程的特权。 此布尔值直接控制是否为容器进程设置 no_new_privs标志。 当容器以特权模式运行或者具有 CAP_SYS_ADMIN 权能时，AllowPrivilegeEscalation 总是为 true。
readOnlyRootFilesystem：以只读方式加载容器的根文件系统。

https://kubernetes.io/zh/docs/tasks/configure-pod-container/security-context/

USER 指定当前用户
格式：USER <用户名>[:<用户组>]
USER 指令和 WORKDIR 相似，都是改变环境状态并影响以后的层。WORKDIR 是改变工作目录，USER 则是改变之后层的执行 RUN, CMD 以及 ENTRYPOINT 这类命令的身份

https://yeasy.gitbook.io/docker_practice/image/dockerfile/user
