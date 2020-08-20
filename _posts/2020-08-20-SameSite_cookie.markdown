---
title: SameSite cookie
layout: post
category: web
author: 夏泽民
---
chrome 新策略导致跨域后会重定向
浏览器访问：chrome://flags/#same-site-by-default-cookies 设置由默认Default改为disabled
<!-- more -->
默认情况下，谷歌将开始为从Chrome 80开始的用户实施新的cookie策略，该策略计划于2020年初发布。 本文解释了您需要了解的有关新SameSite cookie策略的所有信息，如 Adobe Target 何支持这些策略，以及如何使用 Target 来符合Google Chrome的新SameSite cookie策略。
从Chrome 80开始，Web开发人员必须明确指定哪些cookies可以跨网站工作。 这是谷歌计划为改善网络隐私和安全性而发布的多项声明中的第一个。
鉴于Facebook在隐私和安全方面一直处于热门地位，苹果等其他主要企业也迅速利用这一机会创造新的隐私和安全捍卫者身份。 苹果率先宣布今年年初通过ITP 2.1和最近ITP 2.2对其cookie政策进行了修改。在ITP 2.1中，Apple完全阻止第三方Cookie，并将在浏览器上创建的Cookie保存仅七天。 在ITP 2.2中，cookies只保存一天。 谷歌的公告远没有苹果那么咄咄逼人，但这是朝着同一个最终目标迈出的第一步。 有关Apple策略的详细信息，请参 阅Apple Intelligent Tracking Prevention(ITP)2.x 。
什么是cookies，它们是如何使用的？
在深入了解Google对其cookies策略的更改之前，我们先了解Cookie是什么以及它们的使用方式。 简而言之，Cookie是存储在Web浏览器中的小文本文件，用于记住用户属性。
Cookies很重要，因为当用户浏览Web时，它们会增强用户的体验。 例如，如果您在电子商务网站上购物并向购物车中添加内容，但不登录或在该访问中购买，则Cookie会记住您的物品并将它们保存在购物车中，供您下次访问。 或者，想象一下，如果您每次访问您喜爱的社交媒体网站时都被迫重新输入用户名和密码。 Cookies也解决了这个问题，因为它们存储有助于网站识别您身份的信息。 这些类型的Cookie称为第一方Cookie，因为它们是由您访问的网站创建和使用的。
第三方Cookie也存在。 为了更好地理解这些问题，我们来考虑以下示例：
假设某家名为“朋友”的社交媒体公司提供了一个“共享”按钮，其他网站通过该按钮允许“朋友”用户在“朋友”源上共享该网站的内容。 现在，用户在使用“共享”按钮的新闻网站上阅读一篇新闻文章，并单击它以自动发布到其“朋友”帐户。
为了实现此目的，加载新闻文章后，浏览器会 platform.friends.com 从中获取“朋友共享”按钮。 在此过程中，浏览器将包含用户登录凭据的Friends cookie附加到好友服务器的请求中。 这样，好友就可以代表用户在其源中发布新闻文章，而无需用户登录。
这一切都可以通过使用第三方Cookie实现。 在这种情况下，第三方Cookie将保存在浏览器上 platform.friends.com ，以便 platform.friends.com 代表用户在“朋友”应用程序中发布。
如果您想象一下，如何在没有第三方Cookie的情况下实现此使用案例，用户将必须执行大量手动步骤。 首先，用户必须复制指向新闻文章的链接。 其次，用户必须单独登录“朋友”应用程序。 然后，用户单击“创建帖子”按钮。 然后，用户将复制并粘贴文本字段中的链接，最后单击“发布”。 正如您所看到的，第三方Cookie可以极大地减少手动步骤，从而帮助用户体验。
更一般地说，第三方Cookie使得数据存储在用户浏览器上成为可能，而无需该用户显式访问网站。
安全问题
虽然Cookies增强了用户体验和强大的广告功能，但它们也可能引入安全漏洞，如跨站点请求伪造(CSRF)攻击。 例如，如果用户登录银行站点以支付信用卡账单并离开该站点而不注销，然后浏览到同一会话中的恶意站点，则可能发生CSRF攻击。 恶意站点可能包含向银行站点发出请求的代码，该请求在页面加载时执行。 由于用户仍然通过银行站点身份验证，因此会话Cookie可用于启动CSRF攻击，以从用户的银行帐户发起资金转移事件。 这是因为，每次访问站点时，HTTP请求中都会附加所有Cookie。 由于这些安全问题，谷歌现在正在尝试缓解这些问题。
Target如何使用cookies?
尽管如此，让我们看看如何使 Target 用cookies。 为了首先使 Target 用，您需要在站点上安 Target 装JavaScript库。 这使您能够在访问您网站的用户的浏览器上放置第一方Cookie。 当用户与您的网站交互时，您可以通过JavaScript库将用户的行为和兴趣数 Target 据传递给您。 JavaScript Target 库使用第一方Cookie提取有关用户的标识信息以映射到用户的行为和兴趣数据。 然后，这些数据将被用 Target 于推动您的个性化活动。
Target还（有时）使用第三方Cookie。 如果您拥有多个位于不同域上的网站并且希望跟踪这些网站中的用户旅程，则可以通过利用跨域跟踪来使用第三方Cookie。 通过在 Target JavaScript库中启用跨域跟踪，您的帐户将开始使用第三方Cookie。 当用户从一个域跳到另一个域时，浏览器与的后端服务器通信 Target，在此过程中，会创建第三方Cookie并将其放置在用户的浏览器上。 通过用户浏览器上的第三方Cookie, Target 可以为单个用户跨不同域提供一致的体验。
Google的新Cookie菜谱
为了避免在跨站点发送cookies以保护用户时提供保护，Google计划添加对称为SameSite的IETF标准的支持，该标准要求Web开发人员使用Set-Cookie头中的SameSite属性组件管理cookies。
可以将三个不同的值传递到 SameSite 属性：Strict、Lax 或 None。
值
描述
Strict
只有在访问最初设置的域时，才可访问具有此设置的 Cookie。换言之，Strict 会完全阻止跨站点使用 Cookie。这一选择最适合需要高安全性的应用程序，如银行。
Lax
Cookies with this setting are sent only on same-site requests or top-level navigation with non-idempotent HTTP requests, like HTTP GET . 因此，如果第三方可以使用Cookie，但增加了安全优势，保护用户免受CSRF攻击的侵害，则使用此选项。
None
使用此设置的Cookie将像Cookies现在的工作方式一样工作。
请牢记以上几点，Chrome 80为用户引入了两个独立设置：“SameSite by default cookies”和“Cookies without SameSite必须是安全的。” 这些设置将在Chrome 80中默认启用。

SameSite（默认cookies） :设置后，所有未指定SameSite属性的Cookie将自动强制使用 SameSite = Lax 。
没有SameSite的Cookie必须是安全的 :设置后，没有SameSite属性或具有SameSite属性的Cookie SameSite = None 需要是安全的。 在此上下文中，安全是指所有浏览器请求都必须遵循 HTTPS 协议。不符合此要求的 Cookie 将被拒绝。所有网站都应使用HTTPS来满足此要求。
Target 遵循 Google 安全最佳实践
在Adobe，我们始终希望支持业界最新的安全和隐私最佳实践。 We are happy to announce that Target supports the new security and privacy settings introduced by Google.
对于“SameSite by default cookies”设置，将继续提 Target 供个性化，而不会受到任何影响和干预。 Target 会使用第一方 Cookie，并在 Google Chrome 应用标记 SameSite = Lax 时继续正常运行。
对于“不带SameSite的Cookies必须是安全的”选项，如果您不选择加入中的跨域跟踪功能，则中的第一方Cookies将继续 Target 工作。
However, when you opt-in to use cross-domain tracking to leverage Target across multiple domains, Chrome requires SameSite = None and Secure flags to be used for third-party cookies. 这意味着您必须确保您的站点使用HTTPS协议。 中的客户端库 Target 将自动使用HTTPS协议，并将 SameSite = None “和安全”标记附加到第三方Cookie中，以 Target 确保所有活动继续交付。
您需要执行哪些操作？
要了解继续为Google Chrome 80 Target +用户工作需要做什么，请查阅下表，其中将显示以下列：
目标JavaScript库 :如果您使用mbox.js,at.js 1。 x 还是 at.js 2. x 。
默认情况下，SameSite cookies =已启用 :如果您的用户启用了“SameSite by default cookies”，它会对您产生什么影响，您还需要做什么才能继 Target 续工作。
没有SameSite的Cookies必须是安全的=已启用 :如果您的用户启用了“无SameSite的Cookies必须是安全的”，它会对您产生什么影响，您还需要做什么才能继 Target 续工作。
目标JavaScript库
默认为 Cookie 设置 SameSite = 已启用
不具有 SameSite 的 Cookie 必须是安全的 = 已启用
mbox.js（仅包含第一方Cookie）。
没有影响。
如果您不使用跨域跟踪，则不会受到影响。
启用了跨域跟踪的mbox.js。
没有影响。
必须为站点启用HTTPS协议。
Target 使用第三方Cookie跟踪用户，而Google要求第三方Cookie具有和安 SameSite = None 全标记。 安全标志要求您的站点必须使用HTTPS协议。
at.js 1. x 使用第一方cookie。
没有影响。
如果您不使用跨域跟踪，则不会受到影响。
at.js 1. x 并启用跨域跟踪。
没有影响。
必须为站点启用HTTPS协议。
Target 使用第三方Cookie跟踪用户，而Google要求第三方Cookie具有和安 SameSite = None 全标记。 安全标志要求您的站点必须使用HTTPS协议。
at.js 2. x
没有影响。
没有影响。
Target需要做什么？
那么，我们需要在我们的平台中做什么来帮助您遵守新的Google Chrome 80+ SameSite cookie策略？
目标JavaScript库
默认为 Cookie 设置 SameSite = 已启用
不具有 SameSite 的 Cookie 必须是安全的 = 已启用
mbox.js（仅包含第一方Cookie）。
没有影响。
如果您不使用跨域跟踪，则不会受到影响。
启用了跨域跟踪的mbox.js。
没有影响。
Target 在调 SameSite = None 用服务器时向第三方Cookie添加安全 Target 标志。
at.js 1. x 使用第一方cookie。
没有影响。
如果您不使用跨域跟踪，则不会受到影响。
at.js 1. x 并启用跨域跟踪。
没有影响。
at.js 1. x 并启用跨域跟踪。
at.js 2. x
没有影响。
没有影响。
如果不转而使用HTTPS协议，会产生什么影响？
影响您的唯一用例是，如果您使用的跨域跟踪功能 Target 是通过mbox.js或at.js 1。 x 中不再对跨域跟踪提供开箱即用支持。如果不改用Google要求的HTTPS，您会发现您域内的独特访客数量会激增，因为Google会丢弃我们使用的第三方Cookie。 由于第三方Cookie将被删除，当用户从一个域导航到另一个域时， Target 将无法为该用户提供一致的个性化体验。 第三方Cookie主要用于标识在您拥有的域中导航的单个用户。
结论
随着行业在为消费者创建更安全的Web方面取得长足进步， Adobe 我们始终致力于帮助客户以确保最终用户安全和隐私的方式提供个性化体验。 您只需遵循上述最佳实践并充分利用 Target Google Chrome的新SameSite cookie策略。

https://docs.adobe.com/content/help/zh-Hans/target/using/implement-target/before-implement/privacy/google-chrome-samesite-cookie-policies.translate.html

随着疫情的减弱，2020年3月份公司开始复工，复工的第n天我打开Chrome浏览器发现可以升级到80版本了，于是升级，看着关于Chrome前面那个对勾真的开心。
在这里插入图片描述
完了后开始一天的日常任务，打开vscode，启动项目，然后登录系统，WTF，为什么登录后又跳到登录页了？？？控制台也多了下面这个警告。
在这里插入图片描述

找到问题
换了Firefox、360极速都是可以正常登录的，所以这个情况只在Chrome浏览器中有问题。然后看控制台报的警告，发现是SameSite这个属性造成的。经过百度80版本更新日志和查找一些博客发现80版本SameSite这个属性默认值由None变为Lax了。

Cookie 的SameSite属性用来限制第三方 Cookie，防止CRSF（跨站请求攻击）攻击，从而减少安全风险。
SameSite可以设置以下三个值：

Strict
Strict最为严格，完全禁止第三方 Cookie，跨站点时，任何情况下都不会发送 Cookie。换言之，只有当前网页的 URL 与请求目标一致，才会带上 Cookie。
Set-Cookie: CookieName=CookieValue; SameSite=Strict;
1
Lax
Lax规则稍稍放宽，大多数情况也是不发送第三方 Cookie，但是导航到目标网址的 Get 请求除外。
Set-Cookie: CookieName=CookieValue; SameSite=Lax;
1
None
接受第三方Cookie
Set-Cookie: widget_session=abc123; SameSite=None
1
由于我项目在本地，连接的是线上的服务器，请求地址不同，后台接口响应头Set-Cookie中也没设置SameSite这个属性，所以用的默认值，80版本以前是SameSite=“None”,可以接收到线上地址的Cookie，所以就可以登录成功。但是80版本后变为了SameSite=“Lax”，获取不到线上地址的Cookie，所以登录后后台服务器没接收到Cookie，就重新跳转到登录页。

解决方法
方法一：找后台小姐姐帮忙
找后台在接口响应头中的Set-Cookie中加SameSite=None; Secure

80版本后，网站可以选择显式关闭SameSite属性，将其设为None。不过，前提是必须同时设置Secure属性（Cookie 只能通过 HTTPS 协议发送），否则无效。感觉Chrome强迫你支持HTTPS啊。

下面的设置无效。

Set-Cookie: widget_session=abc123; SameSite=None
1
下面的设置有效。

Set-Cookie: widget_session=abc123; SameSite=None; Secure
1
方法二：修改浏览器配置
这个方法只适用于测试环境。线上环境的话就不行了，你不可能强迫别人浏览器都改配置吧。

chrome://flags/#same-site-by-default-cookies

Treat cookies that don’t specify a SameSite attribute as if they were SameSite=Lax. Sites must specify SameSite=None in order to enable third-party usage.

SameSite by default cookies这个配置主要是说如果不设置SameSite就默认SameSite=“Lax”，站点必须指定SameSite=None以启用第三方使用。

chrome://flags/#cookies-without-same-site-must-be-secure

If enabled, cookies without SameSite restrictions must also be Secure. If a cookie without SameSite restrictions is set without the Secure attribute, it will be rejected. This flag only has an effect if “SameSite by default cookies” is also enabled.

Cookies without SameSite must be secure这个配置是说如果启用，没有SameSite限制的Cookie也必须是安全的。如果设置了SameSite=“None”，但是没有设置Secure，它将被拒绝。此配置只有在“SameSite by default cookies”配置开启的情况下才有效。

所以只需要将SameSite by default cookies设置为Disabled就好了。

方法三：回退到80版本前或者换个浏览器（不推荐）
如果你不习惯换个浏览器的话就重新安装个80版本之前的。如果你习惯换个浏览器，那也只能解决一时之需，后面其他浏览器应该也会更新这个功能。这个方法也只适用于测试环境。

https://blog.csdn.net/qwe435541908/article/details/105048484

