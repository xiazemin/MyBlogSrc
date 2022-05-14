---
title: Container
layout: post
category: golang
author: 夏泽民
---
https://github.com/mylxsw/container
Container 是一款为 Go 语言开发的运行时依赖注入库。Go 语言的语言特性决定了实现一款类型安全的依赖注入容器并不太容易，因此 Container 大量使用了 Go 的反射机制。如果你的使用场景对性能要求并不是那个苛刻，那 Container 非常适合你。

并不是说对性能要求苛刻的环境中就不能使用了，你可以把 Container 作为一个对象依赖管理工具，在你的业务初始化时获取依赖的对象。
使用方式

go get github.com/mylxsw/container
要创建一个 Container 实例，使用 containier.New 方法

cc := container.New()
此时就创建了一个空的容器。

你也可以使用 container.NewWithContext(ctx) 来创建容器，创建之后，可以自动的把已经存在的 context.Context 对象添加到容器中，由容器托管。
<!-- more -->
https://segmentfault.com/a/1190000022740651