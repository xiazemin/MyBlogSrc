---
title: validator
layout: post
category: golang
author: 夏泽民
---
awesome-go 要求项目测试覆盖率达到 80% 以上才符合入选标准。有一些公司也会要求项目有相对合理的测试覆盖率(如 70% 以上才符合代码准入条件等等)。

但有时，我们的逻辑代码却挺难做到这么高的覆盖率，主要还是因为目前 Go 的错误处理逻辑：

func Register(req RegisterReq) error{
 if len(req.Username) == 0 {
  return errors.New("length of username cannot be 0")
 }

 if len(req.PasswordNew) == 0 || len(req.PasswordRepeat) == 0 {
  return errors.New("password and password reinput must be longer than 0")
 }

 if req.PasswordNew != req.PasswordRepeat {
  return errors.New("password and reinput must be the same")
 }

 if emailFormatValid(req.Email) {
  return errors.New("invalid email")
 }

 createUser(req.Username, req.PasswordNew, req.Email)
 return nil
}
上面是一个做接口的请求校验的例子(当然请求校验应该直接使用 validator 库，这里只是举个例子)，整个 register 代码中和业务逻辑相关的代码只有最后 createUser(req.Username, req.PasswordNew, req.Email) 这一行，在这个函数里，业务逻辑：非业务逻辑基本都快 1:5 了。

如果我们想要让 Register 函数的测试覆盖率达到 100%，那么我们需要把每种可能的错误都构造出来，这要求我们得写大量的 test case，来处理和业务逻辑没什么关系的用户输入，稍微有点舍本逐末。上面这样的例子，可能的直接结果就是，我们的 test case 数，业务相关:业务无关也是 1:5。如果项目的入口比较多，会有很多很多重复的 test case 散落在各种地方，并不是很好维护。

错误处理代码占比太高的情况下，碰上偷鸡程序员甚至可以在不写任何逻辑代码测试的情况下，只构造错误输入就让该文件的局部覆盖率上升到 80%。

当然，我们用 validator 的话就没这么麻烦了，因为 validate 本来也只有一行：

err := v.validate(req)
if err != nil {
    return err
}
这样我们只要对 validator 本身进行大量严格的测试，甚至可能是各种 fuzz test。就不用在业务项目里去写太多业务无关的 case 了。这里暴露出来的问题，主要还是 Go 本身错误处理的问题，
<!-- more -->
https://mp.weixin.qq.com/s/5KMqKgHC7demT1WeqrAt6A
