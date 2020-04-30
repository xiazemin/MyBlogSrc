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