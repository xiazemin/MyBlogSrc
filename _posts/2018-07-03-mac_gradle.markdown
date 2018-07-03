---
title: mac安装gradle
layout: post
category: java
author: 夏泽民
---
地址为：https://gradle.org/releases/ 
• 选择最新版本，选择complete
v4.8.1
Jun 21, 2018
Download: binary-only or complete
<!-- more -->
• 下载后把压缩包解压到本地任意安装地址 
• 打开终端，把gradle的本地存放地址配置到path地址中
• 输入 export GRADLE_HOME=（你本地解压存放的地址） 
• 输入 export path=（你本地解压存放的地址）\bin 
• 保存后并启动 
• 在终端输入检测配置是否正确，如果正确显示具体的版本信息 $gradle -version  
$gradle -v

Welcome to Gradle 4.8.1!

Here are the highlights of this release:
 - Dependency locking
 - Maven Publish and Ivy Publish plugins improved and marked stable
 - Incremental annotation processing enhancements
 - APIs to configure tasks at creation time

For more details see https://docs.gradle.org/4.8.1/release-notes.html

$gradle
Starting a Gradle Daemon (subsequent builds will be faster)

> Task :help

Welcome to Gradle 4.8.1.

To run a build, run gradle <task> ...
For troubleshooting, visit https://help.gradle.org

BUILD SUCCESSFUL in 6s
1 actionable task: 1 executed