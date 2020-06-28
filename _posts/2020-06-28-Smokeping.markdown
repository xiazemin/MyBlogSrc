---
title: Smokeping
layout: post
category: linux
author: 夏泽民
---
https://oss.oetiker.ch/smokeping/
Smokeping is a latency measurement tool. It sends test packets out to the net and measures the amount of time they need to travel from one place to the other and back.
For every round of measurement smokeping sends several packets. It then sorts the different round trip times and selects the median, (ie. the middle one). This means when there are 10 time values, value number 5 is selected and drawn. The other values are drawn as successively lighter shades of gray in the background (smoke).

Sometimes a test packet is sent out but never returns. This is called packet-loss. The color of the median line changes according to the number of packets lost.

All this information together gives an indication of network health. For example, packet loss is something which should not happen out of the blue. It can mean that a device in the middle of the link is overloaded or a router configuration somewhere is wrong.

Heavy fluctuation of the RTT (round trip time) values also indicate that the network is overloaded. This shows on the graph as smoke; the more smoke, the more fluctuation.

Smokeping is not limited to testing just the roundtrip time of the packets. It can also perform some task at the remote end ("probe"), like download a webpage. This will give a combined 'picture' of webserver availability and network health.

https://oss.oetiker.ch/smokeping/doc/reading.en.html
<!-- more -->
https://github.com/oetiker/SmokePing

https://wiki.archlinux.org/index.php/Smokeping_(%E7%AE%80%E4%BD%93%E4%B8%AD%E6%96%87)

Smokeping允许你监测多台服务器。 Smokeping使用RRDtool来存储数据，另外，其可基于RRDtool输出生成相应的统计图表。 Smokeping由两个部分组成。一个运行在后台、定期收集数据的服务。一个以图表形式展示数据的Web界面。

https://wzfou.com/smokeping/
