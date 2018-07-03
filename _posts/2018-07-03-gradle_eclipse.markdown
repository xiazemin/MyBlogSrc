---
title: gradle_eclipse
layout: post
category: java
author: 夏泽民
---
配置系统环境：GRADLE_HOME设置为解压缩之后的地址，PATH属性追加%GRADLE_HOME%\bin;注意前后的分号
elipse插件下载
Help->Eclipse Marketplace 搜索buildShip并安装 
<!-- more -->
Eclipse导入Gradle项目有两种方式： 
1.直接导入gradle项目 
如果Eclipse中没有安装Gradle插件，需要参考：Eclipse安装Gradle插件这篇文章，安装Eclipse的Gradle插件。 
我们以SpringBoot的初始demo项目为例：

1>在spring官网中http://start.spring.io/下载Gradle项目。
下载要导入的Gradle项目 
2>把下载的项目解压，放到Eclipse工作空间中，也可以根据自己习惯放到一个固定的文件夹。在Eclipse中右键，选择Gradle： 
Existing Gradle Project
3>选择需要导入的Gradle项目，然后点击finish完成。 
这里写图片描述 
此时导入的项目会自动下载项目的jar包，速度会很慢，需要耐心等待。 
2.导入gradle编译后的项目 
1>安装Gradle安装包（网上教程很多，百度下就有） 
2>运行dos命令 
3>进入到项目路径下 
4>执行：gradlew eclipse（注意在build.gradle中加入apply plugin: ‘eclipse’） 
5>编译完成后直接在eclipse中导入。导入方式和正常导入Eclipse java项目一致。 

参考：https://github.com/davenkin/gradle-learning
https://www.w3cschool.cn/gradle/34zp1huk.html


点击eclipse Package Explorer 右上角的倒三角
点击 Fliters 
Deselect All  就显示所有文件了，否则只显示没有被选中的文件