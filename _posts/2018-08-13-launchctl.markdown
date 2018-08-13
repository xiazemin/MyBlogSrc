---
title: launchctl
layout: post
category: web
author: 夏泽民
---
launchctl是一个统一的服务管理框架，可以启动、停止和管理守护进程、应用程序、进程和脚本等。
launchctl是通过配置文件来指定执行周期和任务的。
编写plist文件
launchctl 将根据plist文件的信息来启动任务。
plist脚本一般存放在以下目录：

/Library/LaunchDaemons -->只要系统启动了，哪怕用户不登陆系统也会被执行

/Library/LaunchAgents -->当用户登陆系统后才会被执行

更多的plist存放目录：

~/Library/LaunchAgents 由用户自己定义的任务项
/Library/LaunchAgents 由管理员为用户定义的任务项
/Library/LaunchDaemons 由管理员定义的守护进程任务项
/System/Library/LaunchAgents 由Mac OS X为用户定义的任务项
/System/Library/LaunchDaemons 由Mac OS X定义的守护进程任务项

进入~/Library/LaunchAgents，创建一个plist文件com.demo.plist

<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <!-- Label唯一的标识 -->
  <key>Label</key>
  <string>com.demo.plist</string>
  <!-- 指定要运行的脚本 -->
  <key>ProgramArguments</key>
  <array>
    <string>/Users/demo/run.sh</string>
  </array>
  <!-- 指定要运行的时间 -->
  <key>StartCalendarInterval</key>
  <dict>
        <key>Minute</key>
        <integer>00</integer>
        <key>Hour</key>
        <integer>22</integer>
  </dict>
<!-- 标准输出文件 -->
<key>StandardOutPath</key>
<string>/Users/demo/run.log</string>
<!-- 标准错误输出文件，错误日志 -->
<key>StandardErrorPath</key>
<string>/Users/demo/run.err</string>
</dict>
</plist>
3. 加载命令
launchctl load -w com.demo.plist
这样任务就加载成功了。

更多的命令：

# 加载任务, -w选项会将plist文件中无效的key覆盖掉，建议加上
$ launchctl load -w com.demo.plist

# 删除任务
$ launchctl unload -w com.demo.plist

# 查看任务列表, 使用 grep '任务部分名字' 过滤
$ launchctl list | grep 'com.demo'

# 开始任务
$ launchctl start  com.demo.plist

# 结束任务
$ launchctl stop   com.demo.plist
如果任务呗修改了，那么必须先unload，然后重新load
start可以测试任务，这个是立即执行，不管时间到了没有
执行start和unload前，任务必须先load过，否则报错
stop可以停止任务

番外篇
plist支持两种方式配置执行时间：
StartInterval: 指定脚本每间隔多长时间（单位：秒）执行一次；
StartCalendarInterval: 可以指定脚本在多少分钟、小时、天、星期几、月时间上执行，类似如crontab的中的设置，包含下面的 key:
Minute <integer>
The minute on which this job will be run.
Hour <integer>
The hour on which this job will be run.
Day <integer>
The day on which this job will be run.
Weekday <integer>
The weekday on which this job will be run (0 and 7 are Sunday).
Month <integer>
The month on which this job will be run.
plist部分参数说明：
Label：对应的需要保证全局唯一性；
Program：要运行的程序；
ProgramArguments：命令语句
StartCalendarInterval：运行的时间，单个时间点使用dict，多个时间点使用 array <dict>
StartInterval：时间间隔，与StartCalendarInterval使用其一，单位为秒
StandardInPath、StandardOutPath、StandardErrorPath：标准的输入输出错误文件，这里建议不要使用 .log 作为后缀，会打不开里面的信息。
定时启动任务时，如果涉及到网络，但是电脑处于睡眠状态，是执行不了的，这个时候，可以定时的启动屏幕就好了。

<!-- more -->
