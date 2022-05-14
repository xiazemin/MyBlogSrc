---
title: x-request-id
layout: post
category: node
author: 夏泽民
---
当您操作由客户端访问的Web服务时，可能难以将请求(客户端可以看到)与服务器日志(服务器可以看到)相关联。
X-Request-ID的想法是，客户端可以创建一些随机ID并将其传递给服务器。然后，服务器在其创建的每个日志语句中包含该ID。如果客户端收到错误，它可以将ID包含在错误报告中，允许服务器运算符查找相应的日志语句(而不必依赖于时间戳，IP等)。

由于该ID由客户端生成(随机)，它不包含任何敏感信息，因此不会违反用户的隐私。由于每个请求都创建了唯一的ID，因此跟踪用户也无助于此。

nginx 生成
server {
...

  set_by_lua $uuid '
    if ngx.var.http_x_request_id == nil then
        return uuid4.getUUID()
    else
        return ngx.var.http_x_request_id
    end
  ';

...
}

https://www.cnblogs.com/junneyang/p/5512369.html
<!-- more -->
https://segmentfault.com/a/1190000019682573?utm_source=tag-newest

https://github.com/shfshanyue/blog

https://www.136.la/tech/show-86887.html
