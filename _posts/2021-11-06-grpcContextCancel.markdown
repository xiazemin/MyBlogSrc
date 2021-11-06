---
title: grpc Context Cancel
layout: post
category: golang
author: 夏泽民
---
客户端取消后nginx 日志错误码是499，但是grpc 取消后，如果下游有其他不能认识的错误，会返回500，需要用 func FromContextError(err error) *Status  进行转换
<!-- more -->

https://github.com/grpc/grpc-go/blob/d590071c10a9ed4e4a453307a11b213827c1fb81/status/status.go#L123

https://github.com/grpc/grpc-go/issues/4696

