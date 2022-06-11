---
title: cache Control
layout: post
category: web
author: 夏泽民
---
Cache-Control
计算一个response的Freshness Lifetime（新鲜生命周期、有效期）值的规则为：

如果cache是共享的、且response中含有s-maxage，则值为s-maxage；
如果response中含有max-age，则值为max-age；
如果response的header中含有Expires，则值为Expires减去response header中的Date；
否则的话，表示response中没有显式的过期时间，可能会采用探测型的策略（Last-Modified）；同时，Cache会在任何age超过24小时的response上附加113告警（如果还没有添加过的话）。
Request Cache-Control Directives：
max-age: <delta-seconds>

表示client愿意接收age在delta-seconds秒内的response。在Request和Response的Header中都可以设置。
在Response中设置max-age的时间信息，可以在APP端生成缓存文件，在缓存不过期的情况下，APP不会直接向服务器请求数据，在缓存过期的情况下，APP客户端会向服务器直接请求生成新的缓存。
如果是仅仅设置Request中的max-age时间，是不会生成缓存文件，并且没有缓存是否过期的情况，都是直接向后台服务器直接请求数据。
如果同时设置了Response和Request中的max-age 缓存时间，如果Request中的max-age时间小于Response中的max-age时间，APP会根据Request中max-age时间周期去直接进行网络请求，如果碰到断网或者网络请求不通的情况，即使缓存还在有效期内（Response中设置的max-age时间足够大），在Request设置的max-age过期之后，APP也会直接去进行网络请求。

max-stale: <delta-seconds>
表示client希望接收一个超过其新鲜生命周期（freshness lifetime)、且超出值不大于特定值（秒）的response。如果max-stale未被赋值，则表示client愿意接收任何年龄的response。仅仅在Request的Header中有效。

min-fresh
表示client希望接收一个其新鲜生命周期（freshness lifetime）在现有的基础上再增加不少于特定值（秒）的response。

no-cache

表示在没有成功验证server时，

no-store

no-transform

only-if-cached

表示client仅仅愿意接受一个stored response。server如果有缓存则返回缓存，否则返回504错误。

表示client愿意接受超过其freshness lifetime一定时间的response。如果max-stale没有被赋值，则表示client愿意接受任何陈旧（stale）的值。

Response Cache-Control Directives
must-revalidate

no-cache

在使用任何cached response之前，必须先重新验证服务器（revalidate with the server）。首先连接服务器、并比较服务器上资源的ETag和缓存中的ETag是否一样，如果一样，则返回缓存的资源；

否则，就意味着资源已经发生更新，client需要下载最新的资源并返回。

no-store

no-transform

public

private

proxy-revalidate

max-age=<delta-seconds>

用于计算该response的有效期，该字段的优先级高于Header中的Expires字段。

s-maxage=<delta-seconds>

在共享的缓存中，s-maxage会覆盖max-age指令值或者Expires header的值。

Response中的Age字段
表示该response从产生到现在经过的时间。

If-Modified-Since
If-Modified-Since 是一个条件式请求首部，服务器只在所请求的资源在给定的日期时间之后对内容进行过修改的情况下才会将资源返回，状态码为 200 。如果请求的资源从那时起未经修改，那么返回一个不带有消息主体的 304 响应，而在 Last-Modified 首部中会带有上次修改时间。 If-Modified-Since 只可以用在 GET 或 HEAD 请求中。

If-None-Match
If-None-Match 是一个条件式请求首部。对于 GETGET 和 HEAD 请求方法来说，当且仅当服务器上没有任何资源的 ETag 属性值与这个首部中列出的相匹配的时候，服务器端会才返回所请求的资源，响应码为 200 。对于其他方法来说，当且仅当最终确认没有已存在的资源的 ETag 属性值与这个首部中所列出的相匹配的时候，才会对请求进行相应的处理。

对于 GET 和 HEAD 方法来说，当验证失败的时候，服务器端必须返回响应码 304 （Not Modified，未改变）。对于能够引发服务器状态改变的方法，则返回 412 （Precondition Failed，前置条件失败）。需要注意的是，服务器端在生成状态码为 304 的响应的时候，必须同时生成以下会存在于对应的 200 响应中的首部：Cache-Control、Content-Location、Date、ETag、Expires 和 Vary 。

ETag 属性之间的比较采用的是弱比较算法，即两个文件除了每个比特都相同外，内容一致也可以认为是相同的。例如，如果两个页面仅仅在页脚的生成时间有所不同，就可以认为二者是相同的。

当与 If-Modified-Since 一同使用的时候，If-None-Match 优先级更高（假如服务器支持的话）。
以下是两个常见的应用场景：

采用 GET 或 HEAD 方法，来更新拥有特定的ETag 属性值的缓存。
采用其他方法，尤其是 PUT，将 If-None-Match used 的值设置为 * ，用来生成事先并不知道是否存在的文件，可以确保先前并没有进行过类似的上传操作，防止之前操作数据的丢失。这个问题属于更新丢失问题的一种。

示例：

If-None-Match: "bfc13a64729c4290ef5b2c2730249c88ca92d82d"
If-None-Match: W/"67ab43", "54ed21", "7892dd"
If-None-Match: *
<!-- more -->
https://www.cnblogs.com/xifengcoder/p/15085556.html
