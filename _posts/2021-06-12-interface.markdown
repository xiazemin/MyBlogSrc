---
title: go mysql row 被scan 到 interface应该如何解析
layout: post
category: storage
author: 夏泽民
---
直接.(int64)会报错

正确解法raw, ok := v.([]uint8)

strconv.ParseInt(string(raw),10,64)
<!-- more -->
