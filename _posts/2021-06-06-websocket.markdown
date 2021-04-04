---
title: websocket 抓包
layout: post
category: web
author: 夏泽民
---
Chrome控制台
(1)F12进入控制台，点击Network，选中ws栏，注意选中Filter。
(2)刷新页面会得到一个ws链接。
(3)点击链接可以查看链接详情

https://www.cnblogs.com/songwenjie/p/8575579.html
<!-- more -->

https://stackoverflow.com/questions/13364243/websocketserver-node-js-how-to-differentiate-clients

WebSocket 构造函数接受一个必要参数和一个可选参数：

WebSocket WebSocket(
  in DOMString url,
  in optional DOMString protocols
);
url
要连接的URL；这应当是 WebSocket  服务器会响应的URL。
protocols 可选
一个协议字符串或一个协议字符串数组。这些字符串用来指定子协议，这样一个服务器就可以实现多个WebSocket子协议（比如你可能希望一个服务器可以根据指定的 protocol 来应对不同的互动情况）。如果不指定协议字符串则认为是空字符串。
构造函数可能抛出以下异常：

SECURITY_ERR
尝试连接的端口被阻塞。

https://developer.mozilla.org/zh-CN/docs/Web/API/WebSockets_API/Writing_WebSocket_client_applications

https://javascript.info/websocket

