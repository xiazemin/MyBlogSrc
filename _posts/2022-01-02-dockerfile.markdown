---
title: 一个dockerfile 编译多个镜像 默认镜像
layout: post
category: docker
author: 夏泽民
---
一个dockerfile 可以根据--target 编译成多个镜像

默认只会编译最后一个from的镜像

如果在第一个镜像的基础上编译了其他镜像，默认想编译出来的是第一个镜像，可以使用from final
From scratch as final

From final as xxx

FROM final


<!-- more -->
Docker 17.05版本以后，新增了Dockerfile多阶段构建。所谓多阶段构建，实际上是允许一个Dockerfile 中出现多个 FROM 指令
Docker镜像的每一层只记录文件变更，在容器启动时，Docker会将镜像的各个层进行计算，最后生成一个文件系统，这个被称为 联合挂载

能够将前置阶段中的文件拷贝到后边的阶段中，这就是多阶段构建的最大意义。
https://blog.csdn.net/weixin_43696020/article/details/107336940

https://blog.csdn.net/weixin_40046357/article/details/110790097
