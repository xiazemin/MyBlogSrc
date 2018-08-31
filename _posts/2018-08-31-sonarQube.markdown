---
title: sonarQube之平台搭建及sonar-scanner扫描
layout: post
category: web
author: 夏泽民
---
SonarQube为静态代码检查工具，采用B/S架构，帮助检查代码缺陷，改善代码质量，提高开发速度，通过插件形式，可以支持Java、C、C++、JavaScripe等等二十几种编程语言的代码质量管理与检测。
通过客户端插件分析源代码，sonar客户端可以采用IDE插件、Sonar-Scanner插件、Ant插件和Maven插件方式，并通过各种不同的分析机制对项目源代码进行分析和扫描，并把分析扫描后的结果上传到sonar的数据库，通过sonar web界面对分析结果进行管理
架构图
<img src="{{site.url}}{{site.baseurl}}/img/sonar.png"/>
可以从七个维度检测代码质量:
(1)复杂度分布(complexity):代码复杂度过高将难以理解
(2) 重复代码(duplications):程序中包含大量复制、粘贴的代码而导致代码臃肿，sonar可以展示源码中重复严重的地方
(3) 单元测试统计(unit tests):统计并展示单元测试覆盖率，开发或测试可以清楚测试代码的覆盖情况
(4) 代码规则检查(coding rules):通过Findbugs,PMD,CheckStyle等检查代码是否符合规范
(5) 注释率(comments):若代码注释过少，特别是人员变动后，其他人接手比较难接手；若过多，又不利于阅读
(6) 潜在的Bug(potential bugs):通过Findbugs,PMD,CheckStyle等检测潜在的bug
(7) 结构与设计(architecture & design):找出循环，展示包与包、类与类之间的依赖、检查程序之间耦合度
SonarQube安装
搭建分两大步：服务端跟客户端
服务端
◆进入条件：
    1、准备Java环境，这里略去配置
    2、需要安装MySQL (支持数据库种类见sonar.properties)，这里略去配置
    3、sonar https://docs.sonarqube.org
◆数据库配置：
     1、创建sonar数据库
     2、选择conf/sonar.properties文件，配置数据库设置，默认已经提供了各类数据库的支持，这里选择MySQL数据库，默认已经准备了支持各种数据库，只需将MySQL注释部分去掉，顺便改了sonarQube的端口sonar.web.port=1011
sonar.jdbc.url=jdbc:mysql://localhost:1010/sonar?useUnicode=true&characterEncoding=utf8&rewriteBatchedStatements=true&useConfigs=maxPerformance
sonar.jdbc.driver=com.mysql.jdbc.Driver
sonar.jdbc.username=root
sonar.jdbc.password=root
◆sonar
将下载的soar安装包后，解压，随意放置一个地方
   注：JDK的环境和系统环境，要对应，我是windows系统，JDK位64位，选windows-x86-64
进入bin目录后，点击SonarStart.bat，页面输入http://localhost:1011/，进入页面，配置成功
客户端
前面已经说了客户端可以通过IDE插件、Sonar-Scanner插件、Ant插件和Maven插件方式进行扫描分析，这一节先记录Sonar-Scanner扫描
◆下载sonar-scanner解压，将bin文件加入环境变量path中如我的路径E:\sonar\sonar-scanner\bin将此路径加入path中
◆修改sonar scanner配置文件， conf/sonar-scanner.properties。根据数据库使用情况进行取消相关的注释即可,同时需要添加数据库用户名和密码信息，即配置要访问的sonar服务和mysql服务器地址
◆创建sonar-project.properties文件，以java工程为例在工程根目录下新建立一个sonar-project.properties配置文件

◆开始scanner，只需三步，即可完成
  1、打开CMD命令行，
  2、cd进入你的工作空间，某个工程的代码路径，
  3、敲入sonar-scanner，即可进行分析
◆结果展示，分析完后进入http://localhost:1011/，projectKey点击你分析的工程，查看分析结果
	<img src="{{site.url}}{{site.baseurl}}/img/sonar1.png"/>
<!-- more -->
安装sonar:
下载地址:https://www.sonarqube.org/downloads/

wget https://sonarsource.bintray.com/Distribution/sonarqube/sonarqube-5.6.zip

unzip sonarqube-5.6.zip

mv sonarqube-5.6 /usr/local/

ln -s /usr/local/sonarqube-5.6/ /usr/local/sonarqube

准备数据库：

CREATE DATABASE sonar CHARACTER SET utf8 COLLATE utf8_general_ci;

GRANT ALL ON sonar.* TO 'sonar'@'localhost' IDENTIFIED BY 'sonar@pw';

GRANT ALL ON sonar.* TO 'sonar'@'%' IDENTIFIED BY 'sonar@pw';

FLUSH PRIVILEGES;

启动sonar,如果报错可以看看web.log等日志

/usr/local/sonarqube/bin/linux-x86-64/sonar.sh start
安装sonar插件-中文包
藏的还是比较深的,费劲才找到.参考这里找到的

http://www.jianshu.com/p/a8d4825146a6
