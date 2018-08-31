---
title: Jenkins和SonarQube集成
layout: post
category: web
author: 夏泽民
---
Jenkins与SonarQube 集成插件的安装与配置
Jenkins 是一个支持自动化框架的服务器，我们这里不做详细介绍。Jenkins 提供了相关的插件，使得 SonarQube 可以很容易地集成 ，登陆 jenkins，点击"Manage Jenkins"，选择“Mange Plugins”点击“Avzilable”，搜索“Sonar”选中“SonarQube Scanner for Jenkins”点击安装插件，安装后好如下图：
<img src="{{site.url}}{{site.baseurl}}/img/JenkinsSonar1.png"/>
点击"Manage Jenkins"，选择“Configure System”将SonarQube server的信息填入，点击保存。如图：
<img src="{{site.url}}{{site.baseurl}}/img/JenkinsSonar2.png"/>
在jenkinse服务器上下载sonar-scanner，下载地址：https://sonarsource.bintray.com/Distribution/sonar-scanner-cli/sonar-scanner-cli-3.0.3.778-linux.zip

将下载文件解压至/usr/local/目录下
<img src="{{site.url}}{{site.baseurl}}/img/JenkinsSonar3.png"/>
点击"Manage Jenkins"，选择“Global Tool Configuration”，填入jenkins服务器上的SonarQube 客户端路径，点击保存。如图：
<img src="{{site.url}}{{site.baseurl}}/img/JenkinsSonar4.png"/>
在 Jenkins项目构建过程中加入 SonarScanner 进行代码分析
首先需要在新建的 Jenkins 项目的构建环境标签页中勾选"Prepare SonarQube Scanner evironment"，增加 Execute SonarQube Scanner 构建步骤。如图：
<img src="{{site.url}}{{site.baseurl}}/img/JenkinsSonar5.png"/>
配置 Execute SonarQube Scanner 构建步骤
<img src="{{site.url}}{{site.baseurl}}/img/JenkinsSonar6.png"/>

sonar.projectKey=testSonar 

sonar.projectName=cms

sonar.projectVersion=1.0 

sonar.language=java 

sonar.java.binaries=/var/lib/jenkins/workspace/cms/7-brand-web-cms/target/classes/
 sonar.sources=/var/lib/jenkins/workspace/cms/7-brand-web-cms/src
查看分析结果
在新建的 Jenkins 项目的构建的 Console Output 中可以得到 SonarQube 分析结果的链接，如图：
分析结果报告
<img src="{{site.url}}{{site.baseurl}}/img/JenkinsSonar7.png"/>
具体问题展示：
<img src="{{site.url}}{{site.baseurl}}/img/JenkinsSonar8.png"/>
<!-- more -->
