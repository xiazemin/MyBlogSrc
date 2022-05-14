---
title: monkeypatching
layout: post
category: golang
author: 夏泽民
---
https://github.com/bouk/monkey

https://bou.ke/blog/monkey-patching-in-go/

在进行单元测试的时候，通过 testify框架 对测试函数的数据和所依赖的方法做 mock，但是单测出现 panic。 根据错误提示，被测试函数调用了 time.Now()， 因为会对比这个函数返回值， 所以本次单测没有跑通过。下面介绍通过 monkey patch 来解决这个问题。
<!-- more -->

Monkey Patch
Monkey Patch 是程序在本地扩展、或修改程序属性的一种方式。是指在运行时对类或模块的动态修改，其目的是给现有的第三方代码打上补丁，以解决没有达到预期效果的问题或功能。 一般用于动态语言，比如 Python 和 Ruby。有以下应用场景：

在运行时替换掉 classes/methods/attributes/functions
修改/扩展第三方 Lib 的行为，而不依赖源代码
在运行时将 Patch 的结果应用到内存中的状态
修复原来代码存在的安全问题或行为修正
简单来说就是 Monkey Patch 可以修改当前运行的实例的变量状态和行为。以上面说到的问题，就是修改 time.Now()来返回我们约定好的时间值。

虽然 Go 是静态编译语言，Mockey Patch 的作用域在 Runtime，但是通过 Go 的 unsafe 包，能够将内存中函数的地址替换为运行时函数的地址。具体的原理和实现方式参考 => Monkey Ptching in Go。

解决方案
Monkey 库是 Monkey Patch 的一个 Go 版本实现。通过这个依赖包，修改 time.Now() 返回的时间：

func TestService_HandleEvent_OK(t *testing.T) {
    createdTime = time.Now()
  
      ...
  
      // resolve current time inconsistencies
    monkey.Patch(time.Now, func() time.Time {
        return createdTime
    })
  
      ...
  
}
Patch 后，当主代码执行到 time.Now()时，将指向到这个给定的函数，返回自定义的 Mock 值。

注意： 因为 unsafe操作是不安全的，绕过了 Go 的内存安全原则，所以应该在测试环境中使用 Monkey Patch，并且只在需要的时候使用，确保真正需要 Mocking 的 testing 函数只使用这种方式。

https://segmentfault.com/a/1190000024537099

https://blog.csdn.net/scun_cg/article/details/88395041
