---
title: github不再支持密码方式登录需要用token
layout: post
category: web
author: 夏泽民
---
2021年8月13日起，github不再支持密码方式push，需要把密码换成token的方式，我们推送的时候报错如下：
```
git push
remote: Support for password authentication was removed on August 13, 2021. Please use a personal access token instead.
remote: Please see https://github.blog/2020-12-15-token-authentication-requirements-for-git-operations/ for more information.
fatal: 无法访问 'https://github.com/xiazemin/MyBlogSrc/'：The requested URL returned error: 403
```
<!-- more -->
如何解决呢，参考官方的wiki
https://docs.github.com/en/github/authenticating-to-github/keeping-your-account-and-data-secure/creating-a-personal-access-token
先到个人设置里设置token，按照上面这个wiki一步步操作就行了，很详细

接着在钥匙串里修改github 的密码为刚刚获得的token，具体参考下面这个wiki：
https://docs.github.com/en/get-started/getting-started-with-git/updating-credentials-from-the-macos-keychain

然后保存，push 就成功了

https://stackoverflow.com/questions/68779331/use-token-to-push-some-codes-to-github



https://github.blog/changelog/2021-08-12-git-password-authentication-is-shutting-down/

https://djc8.cn/archives/github-started-from-august-13-2021-and-does-not-accept-the-user-password-for-git-operation-verification.html
