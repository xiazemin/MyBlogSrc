---
title: google authenticator 工作原理
layout: post
category: web
author: 夏泽民
---
Google authenticator是一个基于TOTP原理实现的一个生成一次性密码的工具，用来做双因素登录，市面上已经有很多这些比较成熟的东西存在，像是一些经常用到的U盾，以及数字密码等
实现源码 Google authenticator版本
https://github.com/google/google-authenticator-android

实现原理：
一、用户需要开启Google Authenticator服务时，
1.服务器随机生成一个类似于『DPI45HKISEXU6HG7』的密钥，并且把这个密钥保存在数据库中。
2.在页面上显示一个二维码，内容是一个URI地址（otpauth://totp/账号?secret=密钥），如『otpauth://totp/kisexu@gmail.com?secret=DPI45HCEBCJK6HG7』
3.客户端扫描二维码，把密钥『DPI45HKISEXU6HG7』保存在客户端。
二、用户需要登陆时
1.客户端每30秒使用密钥『DPI45HKISEXU6HG7』和时间戳通过一种『算法』生成一个6位数字的一次性密码，如『684060』
2.用户登陆时输入一次性密码『684060』。
3.服务器端使用保存在数据库中的密钥『DPI45HKISEXU6HG7』和时间戳通过同一种『算法』生成一个6位数字的一次性密码。大家都懂控制变量法，如果算法相同、密钥相同，又是同一个时间（时间戳相同），那么客户端和服务器计算出的一次性密码是一样的。服务器验证时如果一样，就登录成功了。

Tips：
1.这种『算法』是公开的，所以服务器端也有很多开源的实现，比如php版的：https://github.com/PHPGangsta/GoogleAuthenticator 。上github搜索『Google Authenticator』可以找到更多语言版的Google Authenticator。
2.所以，你在自己的项目可以轻松加入对Google Authenticator的支持，在一个客户端上显示多个账户的效果可以看上面android版界面的截图。目前dropbox、lastpass、wordpress，甚至vps等第三方应用都支持Google Authenticator登陆，请自行搜索。
3.现实生活中，网银、网络游戏的实体动态口令牌其实原理也差不多
<!-- more -->
名词解释
OTP 是 One-Time Password的简写，表示一次性密码。
HOTP 是HMAC-based One-Time Password的简写，表示基于HMAC算法加密的一次性密码。
TOTP 是Time-based One-Time Password的简写，表示基于时间戳算法的一次性密码。

TOTP 是时间同步，基于客户端的动态口令和动态口令验证服务器的时间比对，一般每60秒产生一个新口令，要求客户端和服务器能够十分精确的保持正确的时钟，客户端和服务端基于时间计算的动态口令才能一致。
　　HOTP 是事件同步，通过某一特定的事件次序及相同的种子值作为输入，通过HASH算法运算出一致的密码。
　　
OTP基本原理
计算OTP串的公式

OTP(K,C) = Truncate(HMAC-SHA-1(K,C))
其中，

K表示秘钥串；

C是一个数字，表示随机数；

HMAC-SHA-1表示使用SHA-1做HMAC；

Truncate是一个函数，就是怎么截取加密后的串，并取加密后串的哪些字段组成一个数字。

对HMAC-SHA-1方式加密来说，Truncate实现如下。

HMAC-SHA-1加密后的长度得到一个20字节的密串；
取这个20字节的密串的最后一个字节，取这字节的低4位，作为截取加密串的下标偏移量；
按照下标偏移量开始，获取4个字节，按照大端方式组成一个整数；
截取这个整数的后6位或者8位转成字符串返回。

HOTP基本原理
知道了OTP的基本原理，HOTP只是将其中的参数C变成了随机数

公式修改一下

HOTP(K,C) = Truncate(HMAC-SHA-1(K,C))
HOTP： Generates the OTP for the given count

即：C作为一个参数，获取动态密码。

一般规定HOTP的散列函数使用SHA2，即：基于SHA-256 or SHA-512 [SHA2] 的散列函数做事件同步验证；

TOTP基本原理
TOTP只是将其中的参数C变成了由时间戳产生的数字。

TOTP(K,C) = HOTP(K,C) = Truncate(HMAC-SHA-1(K,C))
不同点是TOTP中的C是时间戳计算得出。

C = (T - T0) / X;
T 表示当前Unix时间戳

技术分享

T0一般取值为 0.

X 表示时间步数，也就是说多长时间产生一个动态密码，这个时间间隔就是时间步数X，系统默认是30秒；

例如:

T0 = 0;

X = 30;

T = 30 ~ 59, C = 1; 表示30 ~ 59 这30秒内的动态密码一致。

T = 60 ~ 89, C = 2; 表示30 ~ 59 这30秒内的动态密码一致。

不同厂家使用的时间步数不同；

阿里巴巴的身份宝使用的时间步数是60秒；
宁盾令牌使用的时间步数是60秒；
Google的 身份验证器的时间步数是30秒；
腾讯的Token时间步数是60秒；

python的otp实现
https://pypi.python.org/pypi/pyotp
https://github.com/pyotp/pyotp

结合pyotp 和expect 可以很方便用python实现自动登陆

Google基于TOTP的开源实现
https://github.com/google/google-authenticator
RFC6238中TOTP基于java代码的实现。
golang的一个otp做的不错的实现
https://github.com/gitchs/gootp
RFC参考
RFC 4226 One-Time Password and HMAC-based One-Time Password.
RFC 6238 Time-based One-Time Password.
RFC 2104 HMAC Keyed-Hashing for Message Authentication.


（1）HOTP算法分析。Google Authenticator机制主要基于TOTP算法来实现的，而TOTP算法
本身是基于HOTP（HMAC-based One-Time Password，一种基于HMAC的一次性口令算法）算法
改进的。算法核心内容包括三个参数：一个双方共享的密钥（一个比特序列）K，用于一次
性密码的生成；双方各持有一个计数器C，并且实现将计数值同步；一个签署函数即如下公
式：HOTP(K,C) = Truncate(HMAC-SHA-1(K,C))，上面使用了HMAC-SHA-1，也可以使用HMAC-MD5
等。简单步骤如下：
1、客户端利用持有的密钥K和计数器数值C，通过上述算法公式生成HOTP值，同时计数器值
加1，然后要求用户输入该HOTP值，然后客户端将用户名、密码，连同生成的HOTP值发送给
服务器；
2、服务器获取到客户端发送的信息，解析并验证用户名、密码，然后服务器利用持有的密
钥及自身的计数器数值C，通过相同的算法生成HOTP值，再与客户端的HOTP值对比，如果成
功，则计数器值加1，并允许该客户端登录用户账户，否则拒绝登录。

算法存在的问题：客户端每次请求生成一次性密码操作都会使得计数器值加1，而同时如果
验证失败或者客户端不小心多进行了一次生成密码操作，那么服务器和客户端之间的计数器
C将不再同步，因此需要有一个重新同步（Resynchronization）的机制。由于该重同步机制
对于TOTP算法分析没有太多帮助，因此在此不再赘述，需要了解重同步机制详细的，请参看
RFC4226。

（2）TOTP算法分析。利用HOTP算法，将其中的计数器C用当前时间T来替代，同样可以得到
随着时间变化的一次性密码，并且减小计数器的代价，毕竟更多的使用场景中获取系统时间
是方便的。TOTP算法的三个核心内容也容易理解了：共享密钥K；客户端与服务器的时间同
步；签署函数TOTP = Truncate(HMAC-SHA-1(K, (T - T0) / X))，其中T0是Unix epoch（1970
年1月1日 00:00:00），X为时间分片长度。算法执行流程大致如HOTP算法，不再赘述，不同
的是，本算法中用系统时间T代替了计数器数值C作为HMAC算法的输入。

 需要注意的有几点：1、HMAC算法得出的值位数比较多，不方便用户输入，因此需要截断
 （Truncate）成一组不太长的是进制数（至少6位）2、由于时间是一直动态变化的，这就导
致服务器接收到客户端的消息，再进行TOTP值的计算时，时间上会有延时。因此，需要将时
间划片，当然，时间划片要合理，过短导致用户来不及输入并传输给服务器验证，过长导致
攻击者有足够的时间对用户账户进行攻击。Google默认的采用30秒时间分片，即每过三十秒，
系统时间T的值就会发生变化，得到的动态口令也不相同。3、同样利用系统时间的TOTP算法
也是需要重同步机制。由于网络延时、用户输入延迟等因素，可能服务器端接收到一次性密
码时，T数值已经发生了变化，这样就会导致验证失败。解决方法是，服务器计算当前时间
片以及前面的n个时间片内的TOTP值，只要其中有一个与用户输入的TOTP值相同，则验证通
过。同时也容易理解，n不能设置过大，否则将会降低安全性。该方法还有另外的功能，有
时候客户端与服务器的时钟会有偏差，这样也会造成上面类似的问题。但是如果服务器通过
计算前n个时间片的密码并且成功验证之后，服务器就知道了客户端的时钟偏差。因此，下
一次验证时，服务器就可以直接将偏差考虑在内进行计算，而不需要进行n次计算。

关于文中HOTP算法公式和TOTP算法公式的详细实现，也可以参看RFC4226，其中对于算法通
用性的考虑而做的改进值得学习。

 
值得一提的是，由于利用了系统时间，在一般情况下，生成该动态口令内嵌与程序中，本身
是无需联网的，并且系统时间普遍存在于各个设备中，算法通用性良好。

 一些总结：
1、无论是HOTP还是TOTP算法，都存在重同步问题。参考RFC4226可知，在前者的计数器C
重同步问题中，客户端计数器的值可以预料到必然大于服务器计数器的值，在验证过程中，
服务器向后计算N个HOTP值用来匹配一个客户端HOTP值，验证并通过。至于后者，TOTP
算法要求客户端与服务器时间同步，并且服务器计算TOTP值时会向前计算N个TOTP值用来
匹配一个客户端的TOTP值，从而保证了重同步。这是二者重同步问题中的区别。

2、由于上述的重同步方式，相对也会造成一定的系统安全性问题。同样参考RFC4226，
两种方式的系统服务器端，都会给服务器检测HOTP/TOTP值设置一个阀值S，用来保证服务
器不会不停的检测数值，从而限制了试图制造HOTP/TOTP值的攻击者的可能空间。


