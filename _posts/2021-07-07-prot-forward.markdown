---
title: kubectrl port-forward超时
layout: post
category: k8s
author: 夏泽民
---
https://github.com/kubernetes/kubernetes/issues/19231

如果想将超过5分钟（或无限制）的内容传递到你的kubelet中，可以指定streaming-connection-idle-timeout。例如    --streaming-connection-idle-timeout=4h，将其设置为4小时。或者：    --streaming-connection-idle-timeout=0使其无限制。

<!-- more -->
https://cloud.tencent.com/developer/ask/171689


