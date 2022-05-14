---
title: authentication cookie token 
layout: post
category: node
author: 夏泽民
---
权限认证基础:区分Authentication,Authorization以及Cookie、Session、Token

1. 认证 (Authentication) 和授权 (Authorization)的区别是什么？
这是一个绝大多数人都会混淆的问题。首先先从读音上来认识这两个名词，很多人都会把它俩的读音搞混，所以我建议你先先去查一查这两个单词到底该怎么读，他们的具体含义是什么。

说简单点就是：

认证 (Authentication)： 你是谁。
授权 (Authorization)： 你有权限干什么。
稍微正式点（啰嗦点）的说法就是：

Authentication（认证） 是验证您的身份的凭据（例如用户名/用户ID和密码），通过这个凭据，系统得以知道你就是你，也就是说系统存在你这个用户。所以，Authentication 被称为身份/用户验证。
Authorization（授权） 发生在 Authentication（认证） 之后。授权嘛，光看意思大家应该就明白，它主要掌管我们访问系统的权限。比如有些特定资源只能具有特定权限的人才能访问比如admin，有些对系统资源操作比如删除、添加、更新只能特定人才具有。
这两个一般在我们的系统中被结合在一起使用，目的就是为了保护我们系统的安全性。

2. 什么是Cookie ? Cookie的作用是什么?如何在服务端使用 Cookie ?
2.1 什么是Cookie ? Cookie的作用是什么?
Cookie 和 Session都是用来跟踪浏览器用户身份的会话方式，但是两者的应用场景不太一样。

维基百科是这样定义 Cookie 的：Cookies是某些网站为了辨别用户身份而储存在用户本地终端上的数据（通常经过加密）。简单来说： Cookie 存放在客户端，一般用来保存用户信息。

下面是 Cookie 的一些应用案例：

我们在 Cookie 中保存已经登录过的用户信息，下次访问网站的时候页面可以自动帮你登录的一些基本信息给填了。除此之外，Cookie 还能保存用户首选项，主题和其他设置信息。
使用Cookie 保存 session 或者 token ，向后端发送请求的时候带上 Cookie，这样后端就能取到session或者token了。这样就能记录用户当前的状态了，因为 HTTP 协议是无状态的。
Cookie 还可以用来记录和分析用户行为。举个简单的例子你在网上购物的时候，因为HTTP协议是没有状态的，如果服务器想要获取你在某个页面的停留状态或者看了哪些商品，一种常用的实现方式就是将这些信息存放在Cookie
2.2 如何能在 服务端使用 Cookie 呢？
这部分内容参考：https://attacomsian.com/blog/cookies-spring-boot，更多如何在Spring Boot中使用Cookie 的内容可以查看这篇文章。

1)设置cookie返回给客户端

@GetMapping("/change-username")
public String setCookie(HttpServletResponse response) {
    // 创建一个 cookie
    Cookie cookie = new Cookie("username", "Jovan");
    //设置 cookie过期时间
    cookie.setMaxAge(7 * 24 * 60 * 60); // expires in 7 days
    //添加到 response 中
    response.addCookie(cookie);
 
    return "Username is changed!";
}
2) 使用Spring框架提供的@CookieValue注解获取特定的 cookie的值

@GetMapping("/")
public String readCookie(@CookieValue(value = "username", defaultValue = "Atta") String username) {
    return "Hey! My username is " + username;
}
3) 读取所有的 Cookie 值

@GetMapping("/all-cookies")
public String readAllCookies(HttpServletRequest request) {
 
    Cookie[] cookies = request.getCookies();
    if (cookies != null) {
        return Arrays.stream(cookies)
                .map(c -> c.getName() + "=" + c.getValue()).collect(Collectors.joining(", "));
    }
 
    return "No cookies";
}
3. Cookie 和 Session 有什么区别？如何使用Session进行身份验证？
Session 的主要作用就是通过服务端记录用户的状态。 典型的场景是购物车，当你要添加商品到购物车的时候，系统不知道是哪个用户操作的，因为 HTTP 协议是无状态的。服务端给特定的用户创建特定的 Session 之后就可以标识这个用户并且跟踪这个用户了。

Cookie 数据保存在客户端(浏览器端)，Session 数据保存在服务器端。相对来说 Session 安全性更高。如果使用 Cookie 的一些敏感信息不要写入 Cookie 中，最好能将 Cookie 信息加密然后使用到的时候再去服务器端解密。

那么，如何使用Session进行身份验证？

很多时候我们都是通过 SessionID 来实现特定的用户，SessionID 一般会选择存放在 Redis 中。举个例子：用户成功登陆系统，然后返回给客户端具有 SessionID 的 Cookie，当用户向后端发起请求的时候会把 SessionID 带上，这样后端就知道你的身份状态了。关于这种认证方式更详细的过程如下：

Session Based Authentication flow

用户向服务器发送用户名和密码用于登陆系统。
服务器验证通过后，服务器为用户创建一个 Session，并将 Session信息存储 起来。
服务器向用户返回一个 SessionID，写入用户的 Cookie。
当用户保持登录状态时，Cookie 将与每个后续请求一起被发送出去。
服务器可以将存储在 Cookie 上的 Session ID 与存储在内存中或者数据库中的 Session 信息进行比较，以验证用户的身份，返回给用户客户端响应信息的时候会附带用户当前的状态。
另外，Spring Session提供了一种跨多个应用程序或实例管理用户会话信息的机制。如果想详细了解可以查看下面几篇很不错的文章：

Getting Started with Spring Session
Guide to Spring Session
4. 什么是 Token?什么是 JWT?如何基于Token进行身份验证？
我们在上一个问题中探讨了使用 Session 来鉴别用户的身份，并且给出了几个 Spring Session 的案例分享。 我们知道 Session 信息需要保存一份在服务器端。这种方式会带来一些麻烦，比如需要我们保证保存 Session 信息服务器的可用性、不适合移动端（依赖Cookie）等等。

有没有一种不需要自己存放 Session 信息就能实现身份验证的方式呢？使用 Token 即可！JWT （JSON Web Token） 就是这种方式的实现，通过这种方式服务器端就不需要保存 Session 数据了，只用在客户端保存服务端返回给客户的 Token 就可以了，扩展性得到提升。

JWT 本质上就一段签名的 JSON 格式的数据。由于它是带有签名的，因此接收者便可以验证它的真实性。

下面是 RFC 7519 对 JWT 做的较为正式的定义。

JSON Web Token (JWT) is a compact, URL-safe means of representing claims to be transferred between two parties. The claims in a JWT are encoded as a JSON object that is used as the payload of a JSON Web Signature (JWS) structure or as the plaintext of a JSON Web Encryption (JWE) structure, enabling the claims to be digitally signed or integrity protected with a Message Authentication Code (MAC) and/or encrypted. ——JSON Web Token (JWT)

JWT 由 3 部分构成:

Header :描述 JWT 的元数据。定义了生成签名的算法以及 Token 的类型。
Payload（负载）:用来存放实际需要传递的数据
Signature（签名）：服务器通过Payload、Header和一个密钥(secret)使用 Header 里面指定的签名算法（默认是 HMAC SHA256）生成。
在基于 Token 进行身份验证的的应用程序中，服务器通过Payload、Header和一个密钥(secret)创建令牌（Token）并将 Token 发送给客户端，客户端将 Token 保存在 Cookie 或者 localStorage 里面，以后客户端发出的所有请求都会携带这个令牌。你可以把它放在 Cookie 里面自动发送，但是这样不能跨域，所以更好的做法是放在 HTTP Header 的 Authorization字段中： Authorization: Bearer Token。

Token Based Authentication flow

用户向服务器发送用户名和密码用于登陆系统。
身份验证服务响应并返回了签名的 JWT，上面包含了用户是谁的内容。
用户以后每次向后端发请求都在Header中带上 JWT。
服务端检查 JWT 并从中获取用户相关信息。
 

5 什么是OAuth 2.0？
OAuth 是一个行业的标准授权协议，主要用来授权第三方应用获取有限的权限。而 OAuth 2.0是对 OAuth 1.0 的完全重新设计，OAuth 2.0更快，更容易实现，OAuth 1.0 已经被废弃。详情请见：rfc6749。

实际上它就是一种授权机制，它的最终目的是为第三方应用颁发一个有时效性的令牌 token，使得第三方应用能够通过该令牌获取相关的资源。

OAuth 2.0 比较常用的场景就是第三方登录，当你的网站接入了第三方登录的时候一般就是使用的 OAuth 2.0 协议。
<!-- more -->
https://blog.csdn.net/a967333/article/details/105404998

https://stackoverflow.com/questions/17000835/token-authentication-vs-cookies

https://www.cnblogs.com/chucklu/p/13166131.html

https://www.cnblogs.com/wuwuyong/p/12210842.html
发展史
1、很久很久以前，Web 基本上就是文档的浏览而已， 既然是浏览，作为服务器， 不需要记录谁在某一段时间里都浏览了什么文档，每次请求都是一个新的HTTP协议， 就是请求加响应，  尤其是我不用记住是谁刚刚发了HTTP请求，   每个请求对我来说都是全新的。这段时间很嗨皮

2、但是随着交互式Web应用的兴起，像在线购物网站，需要登录的网站等等，马上就面临一个问题，那就是要管理会话，必须记住哪些人登录系统，  哪些人往自己的购物车中放商品，  也就是说我必须把每个人区分开，这就是一个不小的挑战，因为HTTP请求是无状态的，所以想出的办法就是给大家发一个会话标识(session id), 说白了就是一个随机的字串，每个人收到的都不一样，  每次大家向我发起HTTP请求的时候，把这个字符串给一并捎过来， 这样我就能区分开谁是谁了

3、这样大家很嗨皮了，可是服务器就不嗨皮了，每个人只需要保存自己的session id，而服务器要保存所有人的session id ！  如果访问服务器多了， 就得由成千上万，甚至几十万个。

这对服务器说是一个巨大的开销 ， 严重的限制了服务器扩展能力， 比如说我用两个机器组成了一个集群， 小F通过机器A登录了系统，  那session id会保存在机器A上，  假设小F的下一次请求被转发到机器B怎么办？  机器B可没有小F的 session id啊。

有时候会采用一点小伎俩： session sticky ， 就是让小F的请求一直粘连在机器A上， 但是这也不管用， 要是机器A挂掉了， 还得转到机器B去。

那只好做session 的复制了， 把session id  在两个机器之间搬来搬去， 快累死了。

　　　　　　

后来有个叫Memcached的支了招： 把session id 集中存储到一个地方， 所有的机器都来访问这个地方的数据， 这样一来，就不用复制了， 但是增加了单点失败的可能性， 要是那个负责session 的机器挂了，  所有人都得重新登录一遍， 估计得被人骂死。

　　　　　　  

也尝试把这个单点的机器也搞出集群，增加可靠性， 但不管如何， 这小小的session 对我来说是一个沉重的负担

 

4 于是有人就一直在思考， 我为什么要保存这可恶的session呢， 只让每个客户端去保存该多好？

 

可是如果不保存这些session id ,  怎么验证客户端发给我的session id 的确是我生成的呢？  如果不去验证，我们都不知道他们是不是合法登录的用户， 那些不怀好意的家伙们就可以伪造session id , 为所欲为了。

 

嗯，对了，关键点就是验证 ！

 

比如说， 小F已经登录了系统， 我给他发一个令牌(token)， 里边包含了小F的 user id， 下一次小F 再次通过Http 请求访问我的时候， 把这个token 通过Http header 带过来不就可以了。

 

不过这和session id没有本质区别啊， 任何人都可以可以伪造，  所以我得想点儿办法， 让别人伪造不了。

 

那就对数据做一个签名吧， 比如说我用HMAC-SHA256 算法，加上一个只有我才知道的密钥，  对数据做一个签名， 把这个签名和数据一起作为token ，   由于密钥别人不知道， 就无法伪造token了。



这个token 我不保存，  当小F把这个token 给我发过来的时候，我再用同样的HMAC-SHA256 算法和同样的密钥，对数据再计算一次签名， 和token 中的签名做个比较， 如果相同， 我就知道小F已经登录过了，并且可以直接取到小F的user id ,  如果不相同， 数据部分肯定被人篡改过， 我就告诉发送者： 对不起，没有认证。



Token 中的数据是明文保存的（虽然我会用Base64做下编码， 但那不是加密）， 还是可以被别人看到的， 所以我不能在其中保存像密码这样的敏感信息。

 

当然， 如果一个人的token 被别人偷走了， 那我也没办法， 我也会认为小偷就是合法用户， 这其实和一个人的session id 被别人偷走是一样的。

 

这样一来， 我就不保存session id 了， 我只是生成token , 然后验证token ，  我用我的CPU计算时间获取了我的session 存储空间 ！

 

解除了session id这个负担，  可以说是无事一身轻， 我的机器集群现在可以轻松地做水平扩展， 用户访问量增大， 直接加机器就行。   这种无状态的感觉实在是太好了！

Cookie
cookie 是一个非常具体的东西，指的就是浏览器里面能永久存储的一种数据，仅仅是浏览器实现的一种数据存储功能。

cookie由服务器生成，发送给浏览器，浏览器把cookie以kv形式保存到某个目录下的文本文件内，下一次请求同一网站时会把该cookie发送给服务器。由于cookie是存在客户端上的，所以浏览器加入了一些限制确保cookie不会被恶意使用，同时不会占据太多磁盘空间，所以每个域的cookie数量是有限的。

Session
session 从字面上讲，就是会话。这个就类似于你和一个人交谈，你怎么知道当前和你交谈的是张三而不是李四呢？对方肯定有某种特征（长相等）表明他就是张三。

session 也是类似的道理，服务器要知道当前发请求给自己的是谁。为了做这种区分，服务器就要给每个客户端分配不同的“身份标识”，然后客户端每次向服务器发请求的时候，都带上这个“身份标识”，服务器就知道这个请求来自于谁了。至于客户端怎么保存这个“身份标识”，可以有很多种方式，对于浏览器客户端，大家都默认采用 cookie 的方式。

服务器使用session把用户的信息临时保存在了服务器上，用户离开网站后session会被销毁。这种用户信息存储方式相对cookie来说更安全，可是session有一个缺陷：如果web服务器做了负载均衡，那么下一个操作请求到了另一台服务器的时候session会丢失。

Token
在Web领域基于Token的身份验证随处可见。在大多数使用Web API的互联网公司中，tokens 是多用户下处理认证的最佳方式。

以下几点特性会让你在程序中使用基于Token的身份验证

1.无状态、可扩展

 2.支持移动设备

 3.跨程序调用

 4.安全

 

那些使用基于Token的身份验证的大佬们

大部分你见到过的API和Web应用都使用tokens。例如Facebook, Twitter, Google+, GitHub等。

 

Token的起源

在介绍基于Token的身份验证的原理与优势之前，不妨先看看之前的认证都是怎么做的。

　　基于服务器的验证

　  我们都是知道HTTP协议是无状态的，这种无状态意味着程序需要验证每一次请求，从而辨别客户端的身份。

在这之前，程序都是通过在服务端存储的登录信息来辨别请求的。这种方式一般都是通过存储Session来完成。

下图展示了基于服务器验证的原理

 

随着Web，应用程序，已经移动端的兴起，这种验证的方式逐渐暴露出了问题。尤其是在可扩展性方面。

 

基于服务器验证方式暴露的一些问题

1.Seesion：每次认证用户发起请求时，服务器需要去创建一个记录来存储信息。当越来越多的用户发请求时，内存的开销也会不断增加。

2.可扩展性：在服务端的内存中使用Seesion存储登录信息，伴随而来的是可扩展性问题。

3.CORS(跨域资源共享)：当我们需要让数据跨多台移动设备上使用时，跨域资源的共享会是一个让人头疼的问题。在使用Ajax抓取另一个域的资源，就可以会出现禁止请求的情况。

4.CSRF(跨站请求伪造)：用户在访问银行网站时，他们很容易受到跨站请求伪造的攻击，并且能够被利用其访问其他的网站。

在这些问题中，可扩展行是最突出的。因此我们有必要去寻求一种更有行之有效的方法。

 

基于Token的验证原理

基于Token的身份验证是无状态的，我们不将用户信息存在服务器或Session中。

这种概念解决了在服务端存储信息时的许多问题

　　NoSession意味着你的程序可以根据需要去增减机器，而不用去担心用户是否登录。

基于Token的身份验证的过程如下:

1.用户通过用户名和密码发送请求。

2.程序验证。

3.程序返回一个签名的token 给客户端。

4.客户端储存token,并且每次用于每次发送请求。

5.服务端验证token并返回数据。

 每一次请求都需要token。token应该在HTTP的头部发送从而保证了Http请求无状态。我们同样通过设置服务器属性Access-Control-Allow-Origin:* ，让服务器能接受到来自所有域的请求。需要主要的是，在ACAO头部标明(designating)*时，不得带有像HTTP认证，客户端SSL证书和cookies的证书。

  实现思路：



1.用户登录校验，校验成功后就返回Token给客户端。

2.客户端收到数据后保存在客户端

3.客户端每次访问API是携带Token到服务器端。

4.服务器端采用filter过滤器校验。校验成功则返回请求数据，校验失败则返回错误码

 

 

当我们在程序中认证了信息并取得token之后，我们便能通过这个Token做许多的事情。

我们甚至能基于创建一个基于权限的token传给第三方应用程序，这些第三方程序能够获取到我们的数据（当然只有在我们允许的特定的token）

 

Tokens的优势

无状态、可扩展

在客户端存储的Tokens是无状态的，并且能够被扩展。基于这种无状态和不存储Session信息，负载负载均衡器能够将用户信息从一个服务传到其他服务器上。

如果我们将已验证的用户的信息保存在Session中，则每次请求都需要用户向已验证的服务器发送验证信息(称为Session亲和性)。用户量大时，可能会造成

 一些拥堵。

但是不要着急。使用tokens之后这些问题都迎刃而解，因为tokens自己hold住了用户的验证信息。

安全性

请求中发送token而不再是发送cookie能够防止CSRF(跨站请求伪造)。即使在客户端使用cookie存储token，cookie也仅仅是一个存储机制而不是用于认证。不将信息存储在Session中，让我们少了对session操作。 

token是有时效的，一段时间之后用户需要重新验证。我们也不一定需要等到token自动失效，token有撤回的操作，通过token revocataion可以使一个特定的token或是一组有相同认证的token无效。

可扩展性（）

Tokens能够创建与其它程序共享权限的程序。例如，能将一个随便的社交帐号和自己的大号(Fackbook或是Twitter)联系起来。当通过服务登录Twitter(我们将这个过程Buffer)时，我们可以将这些Buffer附到Twitter的数据流上(we are allowing Buffer to post to our Twitter stream)。

使用tokens时，可以提供可选的权限给第三方应用程序。当用户想让另一个应用程序访问它们的数据，我们可以通过建立自己的API，得出特殊权限的tokens。

多平台跨域

我们提前先来谈论一下CORS(跨域资源共享)，对应用程序和服务进行扩展的时候，需要介入各种各种的设备和应用程序。

Having our API just serve data, we can also make the design choice to serve assets from a CDN. This eliminates the issues that CORS brings up after we set a quick header configuration for our application.

只要用户有一个通过了验证的token，数据和资源就能够在任何域上被请求到。

          Access-Control-Allow-Origin: *       
基于标准

创建token的时候，你可以设定一些选项。我们在后续的文章中会进行更加详尽的描述，但是标准的用法会在JSON Web Tokens体现。

最近的程序和文档是供给JSON Web Tokens的。它支持众多的语言。这意味在未来的使用中你可以真正的转换你的认证机制。

https://blog.csdn.net/augnita/article/details/96966561

HTTP Cookie（也叫Web Cookie或浏览器Cookie）是服务器发送到用户浏览器并保存在本地的一小块数据，它会在浏览器下次向同一服务器再发起请求时被携带并发送到服务器上。通常，它用于告知服务端两个请求是否来自同一浏览器，如保持用户的登录状态。Cookie使基于无状态的HTTP协议记录稳定的状态信息成为了可能。

Cookie主要用于以下三个方面：

会话状态管理（如用户登录状态、购物车、游戏分数或其它需要记录的信息）
个性化设置（如用户自定义设置、主题等）
浏览器行为跟踪（如跟踪分析用户行为等）

Cookie曾一度用于客户端数据的存储，因当时并没有其它合适的存储办法而作为唯一的存储手段，但现在随着现代浏览器开始支持各种各样的存储方式，Cookie渐渐被淘汰。由于服务器指定Cookie后，浏览器的每次请求都会携带Cookie数据，会带来额外的性能开销（尤其是在移动环境下）。新的浏览器API已经允许开发者直接将数据存储到本地，如使用 Web storage API （本地存储和会话存储）或 IndexedDB 。

1. cookie身份验证
用户输入登陆凭据；
服务器验证凭据是否正确，并创建会话，然后把会话数据存储在数据库中；
具有会话id的cookie被放置在用户浏览器中；
服务器验证凭据是否正确，并创建会话；
在后续请求中，服务器会根据数据库验证会话id，如果验证通过，则继续处理；
一旦用户登出，服务端和客户端同时销毁该会话在后续请求中，服务器会根据数据库验证会话id，如果验证通过，则继续处理；
在这里插入图片描述

2. token身份验证
用户输入登陆凭据；
服务器验证凭据是否正确，然后返回一个经过签名的token；
客户端负责存储token，可以存在localstorage，或者cookie中
对服务器的请求带上这个token；
服务器对JWT进行解码，如果token有效，则处理该请求；
一旦用户登出，客户端销毁token。
在这里插入图片描述

3. 二者特性对比
3.1 cookie
用户登录成功后，会在服务器存一个session，同时发送给客户端一个cookie

数据需要客户端和服务器同时存储

用户进行操作时，需要带上cookie，在服务器进行验证

cookie是有状态的

3.2 token
用户进行任何操作时，都需要带上一个token

token的存在形式有很多种，header/requestbody/url 都可以

这个token只需要存在客户端，服务器在收到数据后，进行解析

token是无状态的

3.3总结
Token 完全由应用管理，所以它可以避开同源策略。

Token 可以避免 CSRF 攻击。

Token 可以是无状态的，可以在多个服务间共享。

https://blog.csdn.net/haoaiqian/article/details/85061385

http://fairysoftware.com/token_session_cookie.html
	

cookie、session与token之间的关系

token
令牌，是用户身份的验证方式。
最简单的token组成:uid(用户唯一的身份标识)、time（当前时间的时间戳）、sign（签名）。
对Token认证的五点认识

一个Token就是一些信息的集合；
在Token中包含足够多的信息，以便在后续请求中减少查询数据库的几率；
服务端需要对cookie和HTTP Authrorization Header进行Token信息的检查；
基于上一点，你可以用一套token认证代码来面对浏览器类客户端和非浏览器类客户端；
因为token是被签名的，所以我们可以认为一个可以解码认证通过的token是由我们系统发放的，其中带的信息是合法有效的；
 

session
会话，代表服务器与浏览器的一次会话过程，这个过程是连续的，也可以时断时续。
cookie中存放着一个sessionID，请求时会发送这个ID；
session因为请求（request对象）而产生；
session是一个容器，可以存放会话过程中的任何对象；
session的创建与使用总是在服务端，浏览器从来都没有得到过session对象；
session是一种http存储机制，目的是为武装的http提供持久机制。

cookie
储存在用户本地终端上的数据，服务器生成，发送给浏览器，下次请求统一网站给服务器。

cookie与session区别
cookie数据存放在客户端上，session数据放在服务器上；
cookie不是很安全，且保存数据有限；
session一定时间内保存在服务器上,当访问增多，占用服务器性能。

session与token
作为身份认证，token安全行比session好；
Session 认证只是简单的把User 信息存储到Session 里，因为SID 的不可预测性，暂且认为是安全的。这是一种认证手段。 而Token ，如果指的是OAuth Token 或类似的机制的话，提供的是 认证 和 授权 ，认证是针对用户，授权是针对App 。其目的是让 某App有权利访问 某用户 的信息。

token与cookie
Cookie是不允许垮域访问的，但是token是支持的， 前提是传输的用户认证信息通过HTTP头传输；

token就是令牌，比如你授权（登录）一个程序时，他就是个依据，判断你是否已经授权该软件；cookie就是写在客户端的一个txt文件，里面包括你登录信息之类的，这样你下次在登录某个网站，就会自动调用cookie自动登录用户名；session和cookie差不多，只是session是写在服务器端的文件，也需要在客户端写入cookie文件，但是文件里是你的浏览器编号.Session的状态是存储在服务器端，客户端只有session id；而Token的状态是存储在客户端。

HTTP协议与状态保持：Http是一个无状态协议

 

1. 实现状态保持的方案：

1)修改Http协议，使得它支持状态保持(难做到)

2)Cookies：通过客户端来保持状态信息

Cookie是服务器发给客户端的特殊信息

cookie是以文本的方式保存在客户端，每次请求时都带上它

3)Session：通过服务器端来保持状态信息

Session是服务器和客户端之间的一系列的交互动作

服务器为每个客户端开辟内存空间，从而保持状态信息

由于需要客户端也要持有一个标识(id)，因此，也要求服务器端和客户端传输该标识，

标识(id)可以借助Cookie机制或者其他的途径来保存

2. COOKIE机制

1)Cookie的基本特点

Cookie保存在客户端

只能保存字符串对象，不能保存对象类型

需要客户端浏览器的支持：客户端可以不支持，浏览器用户可能会禁用Cookie

2)采用Cookie需要解决的问题

Cookie的创建

通常是在服务器端创建的(当然也可以通过javascript来创建)

服务器通过在http的响应头加上特殊的指示，那么浏览器在读取这个指示后就会生成相应的cookie了

Cookie存放的内容

业务信息("key","value")

过期时间

域和路径

浏览器是如何通过Cookie和服务器通信？

通过请求与响应，cookie在服务器和客户端之间传递

每次请求和响应都把cookie信息加载到响应头中；依靠cookie的key传递。

3. COOKIE编程

1)Cookie类

Servlet API封装了一个类：javax.servlet.http.Cookie，封装了对Cookie的操作，包括：

public Cookie(String name, String value) //构造方法，用来创建一个Cookie
 
HttpServletRequest.getCookies() //从Http请求中可以获取Cookies
 
HttpServletResponse.addCookie(Cookie) //往Http响应添加Cookie
 
public int getMaxAge() //获取Cookie的过期时间值
 
public void setMaxAge(int expiry) //设置Cookie的过期时间值
2)Cookie的创建

Cookie是一个名值对(key=value)，而且不管是key还是value都是字符串

如： Cookie visit = new Cookie("visit", "1");
3)Cookie的类型——过期时间

会话Cookie

Cookie.setMaxAge(-1);//负整数

保存在浏览器的内存中，也就是说关闭了浏览器，cookie就会丢失

普通cookie

Cookie.setMaxAge(60);//正整数，单位是秒

表示浏览器在1分钟内不继续访问服务器，Cookie就会被过时失效并销毁(通常保存在文件中)

注意：

cookie.setMaxAge(0);//等价于不支持Cookie；
 

4. SESSION机制

每次客户端发送请求，服务断都检查是否含有sessionId。

如果有，则根据sessionId检索出session并处理；如果没有，则创建一个session，并绑定一个不重复的sessionId。

1)基本特点

状态信息保存在服务器端。这意味着安全性更高

通过类似与Hashtable的数据结构来保存

能支持任何类型的对象(session中可含有多个对象)

2)保存会话id的技术(1)

Cookie

这是默认的方式，在客户端与服务器端传递JSeesionId

缺点：客户端可能禁用Cookie

表单隐藏字段

在被传递回客户端之前，在 form 里面加入一个hidden域，设置JSeesionId：

<input type=hidden name=jsessionid value="3948E432F90932A549D34532EE2394" />
URL重写

直接在URL后附加上session id的信息

HttpServletResponse对象中，提供了如下的方法：

encodeURL(url); //url为相对路径
5. SESSION编程

1)HttpSession接口

Servlet API定义了接口：javax.servlet.http.HttpSession， Servlet容器必须实现它，用以跟踪状态。

当浏览器与Servlet容器建立一个http会话时，容器就会通过此接口自动产生一个HttpSession对象

2)获取Session

HttpServletRequest对象获取session，返回HttpSession：
 
request.getSession(); //表示如果session对象不存在，就创建一个新的会话
 
request.getSession(true); //等价于上面这句；如果session对象不存在，就创建一个新的会话
 
request.getSession(false); //表示如果session对象不存在就返回 null，不会创建新的会话对象
3)Session存取信息

session.setAttribute(String name,Object o) //往session中保存信息
 
Object session.getAttribute(String name) //从session对象中根据名字获取信息
4)设置Session的有效时间

public void setMaxInactiveInterval(int interval)

设置最大非活动时间间隔，单位秒；

如果参数interval是负值，表示永不过时。零则是不支持session。

通过配置web.xml来设置会话超时，单位是分钟

<seesion-config>
 
<session-timeout>1</session-timeout>
 
</session-config>
允许两种方式并存，但前者优先级更高

5)其他常用的API

HttpSession.invalidate() //手工销毁Session
 
boolean HttpSession.isNew() //判断Session是否新建
 
如果是true，表示服务器已经创建了该session，但客户端还没有加入(还没有建立会话的握手)
 
HttpSession.getId() //获取session的id
6). 两种状态跟踪机制的比较

Cookie Session

保持在客户端 保存在服务器端

只能保持字符串对象 支持各种类型对象

通过过期时间值区分Cookie的类型 需要sessionid来维护与客户端的通信

会话Cookie——负数 Cookie(默认)

普通Cookie——正数 表单隐藏字段

不支持Cookie——0 url重写

 
 

JSON Web Token（缩写 JWT）是目前最流行的跨域认证解决方案



一、跨域认证的问题
互联网服务离不开用户认证。一般流程是下面这样。

1、用户向服务器发送用户名和密码。

2、服务器验证通过后，在当前对话（session）里面保存相关数据，比如用户角色、登录时间等等。

3、服务器向用户返回一个 session_id，写入用户的 Cookie。

4、用户随后的每一次请求，都会通过 Cookie，将 session_id 传回服务器。

5、服务器收到 session_id，找到前期保存的数据，由此得知用户的身份。

这种模式的问题在于，扩展性（scaling）不好。单机当然没有问题，如果是服务器集群，或者是跨域的服务导向架构，就要求 session 数据共享，每台服务器都能够读取 session。

举例来说，A 网站和 B 网站是同一家公司的关联服务。现在要求，用户只要在其中一个网站登录，再访问另一个网站就会自动登录，请问怎么实现？

一种解决方案是 session 数据持久化，写入数据库或别的持久层。各种服务收到请求后，都向持久层请求数据。这种方案的优点是架构清晰，缺点是工程量比较大。另外，持久层万一挂了，就会单点失败。

另一种方案是服务器索性不保存 session 数据了，所有数据都保存在客户端，每次请求都发回服务器。JWT 就是这种方案的一个代表。

二、JWT 的原理
JWT 的原理是，服务器认证以后，生成一个 JSON 对象，发回给用户，就像下面这样。

 
{
  "姓名": "张三",
  "角色": "管理员",
  "到期时间": "2018年7月1日0点0分"
}
以后，用户与服务端通信的时候，都要发回这个 JSON 对象。服务器完全只靠这个对象认定用户身份。为了防止用户篡改数据，服务器在生成这个对象的时候，会加上签名（详见后文）。

服务器就不保存任何 session 数据了，也就是说，服务器变成无状态了，从而比较容易实现扩展。

三、JWT 的数据结构
实际的 JWT 大概就像下面这样。



它是一个很长的字符串，中间用点（.）分隔成三个部分。注意，JWT 内部是没有换行的，这里只是为了便于展示，将它写成了几行。

JWT 的三个部分依次如下。

Header（头部）
Payload（负载）
Signature（签名）
写成一行，就是下面的样子。

 
Header.Payload.Signature


下面依次介绍这三个部分。

3.1 Header
Header 部分是一个 JSON 对象，描述 JWT 的元数据，通常是下面的样子。

 
{
  "alg": "HS256",
  "typ": "JWT"
}
上面代码中，alg属性表示签名的算法（algorithm），默认是 HMAC SHA256（写成 HS256）；typ属性表示这个令牌（token）的类型（type），JWT 令牌统一写为JWT。

最后，将上面的 JSON 对象使用 Base64URL 算法（详见后文）转成字符串。

3.2 Payload
Payload 部分也是一个 JSON 对象，用来存放实际需要传递的数据。JWT 规定了7个官方字段，供选用。

iss (issuer)：签发人
exp (expiration time)：过期时间
sub (subject)：主题
aud (audience)：受众
nbf (Not Before)：生效时间
iat (Issued At)：签发时间
jti (JWT ID)：编号
除了官方字段，你还可以在这个部分定义私有字段，下面就是一个例子。

 
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
注意，JWT 默认是不加密的，任何人都可以读到，所以不要把秘密信息放在这个部分。

这个 JSON 对象也要使用 Base64URL 算法转成字符串。

3.3 Signature
Signature 部分是对前两部分的签名，防止数据篡改。

首先，需要指定一个密钥（secret）。这个密钥只有服务器才知道，不能泄露给用户。然后，使用 Header 里面指定的签名算法（默认是 HMAC SHA256），按照下面的公式产生签名。

 
HMACSHA256(
  base64UrlEncode(header) + "." +
  base64UrlEncode(payload),
  secret)
算出签名以后，把 Header、Payload、Signature 三个部分拼成一个字符串，每个部分之间用"点"（.）分隔，就可以返回给用户。

3.4 Base64URL
前面提到，Header 和 Payload 串型化的算法是 Base64URL。这个算法跟 Base64 算法基本类似，但有一些小的不同。

JWT 作为一个令牌（token），有些场合可能会放到 URL（比如 api.example.com/?token=xxx）。Base64 有三个字符+、/和=，在 URL 里面有特殊含义，所以要被替换掉：=被省略、+替换成-，/替换成_ 。这就是 Base64URL 算法。

四、JWT 的使用方式
客户端收到服务器返回的 JWT，可以储存在 Cookie 里面，也可以储存在 localStorage。

此后，客户端每次与服务器通信，都要带上这个 JWT。你可以把它放在 Cookie 里面自动发送，但是这样不能跨域，所以更好的做法是放在 HTTP 请求的头信息Authorization字段里面。

 
Authorization: Bearer <token>
另一种做法是，跨域的时候，JWT 就放在 POST 请求的数据体里面。

五、JWT 的几个特点
（1）JWT 默认是不加密，但也是可以加密的。生成原始 Token 以后，可以用密钥再加密一次。

（2）JWT 不加密的情况下，不能将秘密数据写入 JWT。

（3）JWT 不仅可以用于认证，也可以用于交换信息。有效使用 JWT，可以降低服务器查询数据库的次数。

（4）JWT 的最大缺点是，由于服务器不保存 session 状态，因此无法在使用过程中废止某个 token，或者更改 token 的权限。也就是说，一旦 JWT 签发了，在到期之前就会始终有效，除非服务器部署额外的逻辑。

（5）JWT 本身包含了认证信息，一旦泄露，任何人都可以获得该令牌的所有权限。为了减少盗用，JWT 的有效期应该设置得比较短。对于一些比较重要的权限，使用时应该再次对用户进行认证。

（6）为了减少盗用，JWT 不应该使用 HTTP 协议明码传输，要使用 HTTPS 协议传输。

六、参考链接
Introduction to JSON Web Tokens， by Auth0
Sessionless Authentication using JWTs (with Node + Express + Passport JS), by Bryan Manuele
Learn how to use JSON Web Tokens, by dwyl

https://blog.csdn.net/qq_37939251/article/details/83511451

https://segmentfault.com/a/1190000017831088

https://www.jianshu.com/p/bd1be47a16c1

https://www.jianshu.com/p/ce9802589143
https://www.cnblogs.com/moyand/p/9047978.html