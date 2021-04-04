---
title: goimports 分组导入
layout: post
category: golang
author: 夏泽民
---
https://github.com/golang/tools/blob/master/cmd/goimports/goimports.go

分组格式化 goimports -local pkg_prefix_a,pkg_prefix_b -w test/a_test.go

https://blog.jetbrains.com/go/2021/02/26/goland-2021-1-eap-5/

import (
    stdlib

    current_project

    company

    all others
)

有没有可以自动执行此操作的工具？
goimports的最新版本支持-local标志。 引用此提交消息：
例如，运行goimports -local example.com
<!-- more -->
golangci-lint
https://github.com/alecthomas/gometalinter

错误

pkg/skaffold/kubernetes/wait.go:23: File is not `goimports`-ed with -local github.com/GoogleContainerTools/skaffold (goimports)
        "github.com/GoogleContainerTools/skaffold/pkg/skaffold/kubectl"

仅执行常规的Sort imports可能无法正常工作。我认为您已启用goimports的local-prefixes linting，这就是为什么出现File is not 'goimports'-ed with -local ...错误的原因
通常，goimports以某种方式对导入的库进行排序，以使标准pkg和其他库位于单独的组中。但是，当启用了本地前缀时，linting会期望标准pkg，第三方pkg和具有指定本地前缀的pkg(在您的情况下为github.com/GoogleContainerTools/skaffold，又名您自己的项目pkg)，这3种类型在单独的组中

https://www.coder.work/article/7185002

https://www.cnblogs.com/davygeek/p/6387385.html
https://www.jb51.cc/go/187164.html
https://github.com/Masterminds/sprig

