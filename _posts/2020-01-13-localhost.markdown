---
title: Charles不能捕获localhost请求
layout: post
category: web
author: 夏泽民
---
官网解决方案
Some systems are hard coded to not use proxies for localhost traffic, so when you connect to http://localhost/ it doesn't show up in Charles.
The workaround is to connect to localhost.charlesproxy.com/ instead. This points to the IP address 127.0.0.1, so it should work identically to localhost, but with the advantage that it will go through Charles. This will work whether or not Charles is running or you're using Charles. If you use a different port, such as 8080, just add that as you usually would, e.g. localhost.charlesproxy.com:8080.
You can also put anything in front of that domain, e.g. myapp.localhost.charlesproxy.com, which will also always resolve to 127.0.0.1.
Alternatively you can try adding a '.' after localhost, or replace localhost with the name of your machine, or use your local link IP address (eg. 192.168.1.2).
If Charles is running and you're using Charles as your proxy, you can also use local.charles as an alternative for localhost. Note that this only works when you're using Charles as your proxy, so the above approaches are preferred, unless you specifically want requests to fail if not using Charles

解决方案就是配置 在host文件中添加一行

127.0.0.1 localhost.charlesproxy.com

或者直接用http://localhost.charlesproxy.com:端口号   发请求
<!-- more -->
连接到http://localhost.charlesproxy.com/。这指向IP地址127.0.0.1，因此它应该与localhost完全相同，但它的优势在于它将通过Charles。无论Charles是在跑，还是在使用Charles，这都会有效。如果您使用其他端口，例如8080，只需像往常一样添加它，例如localhost.charlesproxy.com:8080。

您还可以在该域前放置任何内容，例如myapp.localhost.charlesproxy.com，它也将始终解析为127.0.0.1。

或者，您可以尝试添加'。' 在localhost之后，或用本机名称替换localhost，或使用本地链接IP地址（例如192.168.1.2）。

如果Charles正在运行并且您使用Charles作为代理，那么您也可以使用local.charles作为localhost的替代方案。请注意，这仅在您使用Charles作为代理时才有效，因此上述方法是首选方法，除非您特别希望请求在不使用Charles时失败

第一种：Charles官方对不能捕获localhost本地流量的说明，以及解决方法。全文大致意思如下：

Localhost流量不会出现在Charles中
某些系统被硬编码为不使用代理进行本地主机流量，因此当您连接到http：// localhost /时，它不会显示在Charles中。

解决方法是连接到http://localhost.charlesproxy.com/。这指向IP地址127.0.0.1，因此它应该与localhost完全相同，但它的优势在于它将通过Charles。无论Charles是在跑还是你在使用Charles，这都会有效。如果您使用其他端口，例如8080，只需像往常一样添加它，例如localhost.charlesproxy.com:8080。

您还可以在该域前放置任何内容，例如myapp.localhost.charlesproxy.com，它也将始终解析为127.0.0.1。

或者，您可以尝试添加'。' 在localhost之后，或用本机名称替换localhost，或使用本地链接IP地址（例如192.168.1.2）。

如果Charles正在运行并且您使用Charles作为代理，那么您也可以使用local.charles作为localhost的替代方案。请注意，这仅在您使用Charles作为代理时才有效，因此上述方法是首选方法，除非您特别希望请求在不使用Charles时失败。

原文链接：https://www.charlesproxy.com/documentation/faqs/localhost-traffic-doesnt-appear-in-charles/

这个声明说明的缘由，charles是不再支持localhost流量的捕捉，并说了localhost.charlesproxy.com:8080可以代替localhost本地服务使用。

痛点：可惜的是，这个说的很清楚，可是小白对这些还不太熟，使用localhost.charlesproxy.com:8080切实可以打开一个本地服务（目测是charles软件提供的），但是我们项目的服务请求地址是localhost:8080呀使用他这个并不会转到我们想要的服务目录下---折腾了一会，放弃了（不知道会不是我没有改下服务启动的服务默认地址，比如localhost该成localhost.charlesproxy.com，没有尝试，总之放弃了）

第二种，是朋友提供的hosts对127.0.0.1 进行映射（以实测成功，可以参照：Charles不能捕获localhost请求解决方法-修改host方法）

host文件在C:\Windows\System32\drivers\etc

在这里面我新增了 127.0.0.1 映射成localhost.charlesproxy.com，保存重启电脑，启动服务，输入127.0.01 或localhost.charlesproxy.com 还是不管用啊（Tips:可能我太菜了，设置或理解映射有问题，知道的朋友欢迎告知一下呀，共享一下知识）



于似乎，上面两种方法，都不能解决我的charles捕捉localhost本地流量问题，更要命的是，博友们都是简洁的给了一点提示，但是又给具体的步骤，菜菜的自己无力回天实现不成功。

该方法已经设置成功，可以参考：Charles不能捕获localhost请求解决方法-修改host方法

3.曲线救国-使用Fiddler完成localhost抓取和对其重定向

Fiddler是跟Charles类似的软件，可以自动捕获localhost请求，最重要的开源免费，只是操作上有点不同，想要仔细了解它的特性，建议大家去看看文档，本文只会介绍简单的和完成localhost捕获或重定向请求步骤。

下载地址:https://www.telerik.com/download/fiddler

安装默认即可

使用Fidder,页面如下，可以看到本地流量localhost是自动被捕捉到的。后面就是介绍我怎么完成对本地请求的重定向配置。
3.1 重定向配置

和charles一样如果我们要对某个url进行重定向配置，第一步就是选中要配置的地址，Fidder也是一样，一localhost:8088为例。

1.选中要配置的url（这是前端本地开启的服务，请求接口在线上，所以报404-->任务就是将这个配成线上地址）
