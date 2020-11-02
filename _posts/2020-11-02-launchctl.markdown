---
title: Mac使用Launchd命令行lauchctl
layout: post
category: linux
author: 夏泽民
---
mac的守护进程目录有以下几处：

~/library/launchagents 	# 用户的进程
/library/launchagents 	# 管理员设置的用户进程
/library/launchdaemons 	# 管理员提供的系统守护进程
/system/library/launchagents 	# mac操作系统提供的用户进程
/system/library/launchdaemons 	# mac操作系统提供的系统守护进程

https://sspai.com/post/37258

% sudo launchctl unload /System/Library/LaunchDaemons/com.apple.securityd.plist
/System/Library/LaunchDaemons/com.apple.securityd.plist: Operation not permitted while System Integrity Protection is engaged

https://www.cnblogs.com/EasonJim/p/7173859.html
<!-- more -->
https://www.sunzhongwei.com/process-of-launchd-is-that-it-produces-a-large-number-of-disk-io-on-mac-system

https://blog.csdn.net/hk_5788/article/details/51213072

https://blog.csdn.net/kevinlou2008/article/details/49155497

https://www.cnblogs.com/findumars/p/6891408.html

https://blog.csdn.net/skymingst/article/details/39371935
