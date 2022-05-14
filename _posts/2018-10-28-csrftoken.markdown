---
title: csrf token
layout: post
category: web
author: 夏泽民
---
<!-- more -->
CSRF 攻击之所以能够成功，是因为黑客可以完全伪造用户的请求，该请求中所有的用户验证信息都是存在于 cookie 中，因此黑客可以在不知道这些验证信息的情况下直接利用用户自己的 cookie 来通过安全验证。要抵御 CSRF，关键在于在请求中放入黑客所不能伪造的信息，并且该信息不存在于 cookie 之中。可以在 HTTP 请求中以参数的形式加入一个随机产生的 token，并在服务器端建立一个拦截器来验证这个 token，如果请求中没有 token 或者 token 内容不正确，则认为可能是 CSRF 攻击而拒绝该请求。
这种方法要比检查 Referer 要安全一些，token 可以在用户登陆后产生并放于 session 之中，然后在每次请求时把 token 从 session 中拿出，与请求中的 token 进行比对，但这种方法的难点在于如何把 token 以参数的形式加入请求。对于 GET 请求，token 将附在请求地址之后，这样 URL 就变成 http://url?csrftoken=tokenvalue。 而对于 POST 请求来说，要在 form 的最后加上 <input type=”hidden” name=”csrftoken” value=”tokenvalue”/>，这样就把 token 以参数的形式加入请求了。
如果说这个Token是指的用户登录的凭据，并用以维持登录状态的话，也就是说一个用户必须要输入用户名密码并验证通过后，服务器才会分配一个Token，传回并储存在客户端作为凭证（同时储存在服务器上）。因此并不是每个人都可以获得这个Token，只有能提供正确用户密码的客户端才可以。
之后每一次操作，都需要客户端向服务器提供这个Token，以验证登录状态，如果考虑安全性的话，还可以增加对User-Agent、IP等信息的验证。
CSRF防范方法：
（1）验证码
（2）refer头
（3）Token
说明：理解token的作用，他是一个随机的值，是服务器端前一个请求给的，是一次性的，可以防止csrf这种恶意的携带自己站点的信息发请求或者提交数据（这个动作一般需要获取你的前一个请求的响应返回的token值，加大了难度，并不能完全杜绝）。
注意当然不能写到cookie中，因为浏览器在发出恶意csrf请求时，是自动带着你的cookie的。


