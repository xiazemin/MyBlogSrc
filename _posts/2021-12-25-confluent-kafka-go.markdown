---
title: confluent-kafka-go mac m1
layout: post
category: golang
author: 夏泽民
---
./configure --prefix=/opt/homebrew --arch=arm64 --enable-static

https://github.com/xiazemin/confluent-kafka-go

https://github.com/edenhill/librdkafka

STATIC_LIB_libzstd=$(brew ls -v zstd | grep libzstd.a$) ./configure --enable-static --prefix=/opt/homebrew --arch=arm64

make

Using ctags to generate TAGS /Library/Developer/CommandLineTools/usr/bin/ctags: illegal option -- e usage: ctags [-BFadtuwvx] [-f tagsfile] file ... cmp: TAGS: No such file or directory

brew install universal-ctags

cmp: TAGS: No such file or directory

make clean
<!-- more -->
