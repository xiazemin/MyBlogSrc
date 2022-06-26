---
title: timeStamp
layout: post
category: golang
author: 夏泽民
---
fmt.Println(time.Now())
fmt.Println(time.Now().Local())
fmt.Println(time.Now().UTC())
fmt.Println(time.Now().Location())
<!-- more -->
2018-11-21 11:50:39.540473 +0800 CST m=+0.000311562
2018-11-21 11:50:39.540628 +0800 CST
2018-11-21 03:50:39.540632 +0000 UTC
Local


获取时间的函数为time.now()，这里加不加.Local()，获取的都是当地时间。
加.UTC得到的是0时区（也就是伦敦）的时间。
func Now() Time这个函数的返回值是Time，也就是时间类型。

时间戳
时间戳函数的返回值都是int64，是一个大整数。

获取时间戳
fmt.Println(time.Now().Unix())
fmt.Println(time.Now().Local().Unix())
fmt.Println(time.Now().UTC().Unix())
fmt.Println(time.Now().UnixNano())
运行结果

1542772752
1542772752
1542772752
1542772752846107000
这次，加不加.Local()、.UTC()结果都是一样的。
那什么是时间戳呢，时间戳就是

格林威治时间1970年01月01日00时00分00秒(北京时间1970年01月01日08时00分00秒)起到此时此刻的【总秒数】

那么，在go语言中，time.Now().Unix()或者time.Now().Local().Unix()就是【北京时间1970年01月01日08时00分00秒】到【北京时间此时此刻】的总秒数。

相应的time.Now().UTC().Unix()就是【格林威治时间1970年01月01日00时00分00秒】到【格林威治时间此时此刻】的总秒数。

因此上面得到的几个时间戳是一样的。

时间戳是一个【总秒数】，所以时间戳函数的返回值都是int64。所以go语言中有时间类型，但并没有一个单独的【时间戳类型】。

将时间类型格式化，得到一个表示时间的字符串
t := time.Now()
str := t.Format("2006-01-02 15:04:05")
str1 := t.Format("2006年1月2日 15:04:05")
fmt.Println(t)
fmt.Println(str)
fmt.Println(str1)

运行结果

2018-11-21 12:48:19.870047 +0800 CST m=+0.000503740
2018-11-21 12:48:19
2018年11月21日 12:48:19
第一行是time.Now()的结果，是时间类型【Time】
下面两行是t.Format()的结果，是字符串。

将表示时间类型的字符串转换为时间类型Time
t := time.Now()
str := t.Format("2006-01-02 15:04:05")
str1 := t.Format("2006年1月2日 15:04:05")
timestamp, _ := time.Parse("2006-01-02 15:04:05", str)
timestamp1, _ := time.Parse("2006年1月2日 15:04:05", str1)
fmt.Println(timestamp)
fmt.Println(timestamp1)
运行结果

2018-11-21 12:48:19 +0000 UTC
2018-11-21 12:48:19 +0000 UTC
函数func Parse(layout, value string) (Time, error)的第一个参数是需要转换的字符串的格式，第二个参数是需要转换的字符串。返回值是时间类型和一个err。

【注意】

在将字符串转为时间类型的时候，是直接转为了【伦敦时间】，go语言并不会去判断这个字符串表示的是北京时间，还是伦敦时间，因为没法判断，只有你知道它表示的是哪里的时间。比如16:08:05在中国那当然表示的是北京时间，但是如果把这个字符串转为时间类型，就直接变成伦敦时间的16:08:05了。

将时间类型转换为时间戳
直接调用方法func (t Time) Unix() int64即可。
将上面的两个时间变量timestamp和timestamp1转为时间戳

fmt.Println(timestamp.Unix())
fmt.Println(timestamp1.Unix())
运行结果

1542804499
1542804499
将时间戳转换为时间类型
用函数func Unix(sec int64, nsec int64) Time进行转换，第一个参数是秒，第二个参数是纳秒，会被加到结果的小数点后面。

tmsp := time.Now().Unix()
fmt.Println(tmsp)
t1 := time.Unix(tmsp, 0).UTC()
t2 := time.Unix(tmsp, tmsp).Local()
fmt.Println(t1)
fmt.Println(t2)
运行结果

1542779708
2018-11-21 05:55:08 +0000 UTC
2018-11-21 13:55:09.542779708 +0800 CST
这里的转换可以选择是转换为当地时间还是伦敦时间。
https://blog.csdn.net/Charliewolf/article/details/84323574
