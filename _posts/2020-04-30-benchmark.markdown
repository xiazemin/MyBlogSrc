---
title: benchmark
layout: post
category: golang
author: 夏泽民
---
Sometimes you have to solve a problem that comes in several flavours. Usually complicated problems do not offer a single solution, but there are several solutions that are optimal or terrible depending on which subset of that problem the program will have to solve at runtime.

One example I faced was to analyse some data flowing in some connections that I was proxying.

There are two main ways to extract some information from traffic: you can either record the entire traffic to analyse it as soon as it is done, or you can analyse it while it flows(with a buffer window) at the cost of slowing it down.

Memory is relatively cheap compared to processing power, so my first solution to the problem was the buffered one.
<!-- more -->
https://blogtitle.github.io/go-advanced-benchmarking/
https://dev.to/segflow/my-journey-optimizing-the-go-compiler-46jc

https://evrone.com/rob-pike-interview
https://github.com/facebookincubator/ent
https://dev.to/segflow/my-journey-optimizing-the-go-compiler-46jc
https://github.com/sunshinev/go-sword
https://segmentfault.com/a/1190000022523822
https://github.com/supanadit/jwt-go
https://github.com/golang/go/issues/38762
https://medium.com/a-journey-with-go/go-asynchronous-preemption-b5194227371c

https://www.infoq.cn/article/yDMrvVr1IJAAih3eh5fW
https://github.com/cespare/reflex

https://github.com/fatih/vim-go
https://medium.com/a-journey-with-go/go-how-does-defer-statement-work-1a9492689b6e
https://medium.com/a-journey-with-go/go-improve-the-usage-of-your-goroutines-with-godebug-4d1f33970c33

https://blog.min.io/accelerating-aggregate-md5-hashing-up-to-800-with-avx512-2/
https://github.com/minio/md5-simd

https://github.com/Cretezy/dSock
https://segflow.github.io/post/go-compiler-optimization/

https://benma.github.io/2020/05/05/golang-embeding-structs-breaks-modularity.html

https://github.com/explore/email

https://www.oreilly.com/

https://news.ycombinator/

https://thenewstack.io/

https://medium.com/

https://tools.ietf.org/html/rfc1180
https://oktop.tumblr.com/post/15352780846

https://medium.com/a-journey-with-go/go-samples-collection-with-pprof-2a63c3e8a142

https://medium.com/@bijeshos/building-command-line-interfaces-using-go-ce6a75d60bf5

https://github.com/mathetake/gasm

https://zhuanlan.zhihu.com/p/45492055

https://zhuanlan.zhihu.com/p/41251789

https://gocn.vip/topics/10359

https://blog.jetbrains.com/go/2020/05/06/debugging-a-go-application-inside-a-docker-container/

https://gocn.vip/topics/10358
https://medium.com/@arrafiv/basic-image-processing-with-go-combining-images-and-texts-8510d9214e55
https://medium.com/the-programming-hub/insanely-addictive-retro-looking-multiplayer-terminal-game-written-in-go-e820cfe8aa40
https://caddyserver.com/v2
https://github.com/caddyserver/caddy
https://github.com/nakabonne/golintui
https://github.com/codenotary/immudb