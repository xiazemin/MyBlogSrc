---
title: scala maven 版本冲突问题解决
layout: post
category: spark
author: 夏泽民
---

scalatest_2.10-1.9.1.jar of core build path is cross-compiled with an incompatible version of Scala (2.10.0)

Eclipse - Preferences - Scala - Compiler - Build manager
uncheck withVersionClasspathVariable


More than one scala library found in the build path (/home/hadoop/eclipse/plugins/org.scala-lang.scala-library_2.11.7.v20150622-112736-1fbce4612c.jar, /usr/local/spark/spark-1.5.1-bin-hadoop2.6/lib/spark-assembly-1.5.1-hadoop2.6.0.jar).At least one has an incompatible version. Please update the project build path so it contains only one compatible scala library. hello-test Unknown Scala Classpath Problem 




修改工程中的scala编译版本
右击 --> Scala --> set the Scala Installation

也可以

右击工程--> Properties --> Scala Compiler --> Use project Setting 中选择spark对应的scala版本，此处选择Lastest2.10 bundle


上述方法仍然没有解决
原因maven pom.xml 中的版本与eclipse里面设置的版本冲出
解决办法修改pom.xml 
<properties>
    <scala.compat.version>2.11</scala.compat.version>
    <scala.version>2.12.3</scala.version>
    
 问题解决
 
 Unsupported major.minor version 52.0
 You get this error because a Java 7 VM tries to load a class compiled for Java 8

Java 8 has the class file version 52.0 but a Java 7 VM can only load class files up to version 51.0

In your case the Java 7 VM is your gradle build and the class is com.android.build.gradle.AppPlugin
简单来说，就是java的编译环境版本太低，java 8 class file的版本是52，Java 7虚拟机只能支持到51。所以需要升级到java 8 vm才行


mvn -V
Apache Maven 3.5.2 (138edd61fd100ec658bfa2d307c43b76940a5d7d; 2017-10-18T15:58:13+08:00)
Maven home: /Users/didi/maven
Java version: 1.8.0_144, vendor: Oracle Corporation

Missing artifact org.scalatest:scalatest_2.12:jar:  2.2.4

http://mvnrepository.com/artifact/org.scalatest/scalatest_2.12/3.0.3


  <dependency>
      <groupId>org.specs2</groupId>
      <artifactId>specs2-core_${scala.compat.version}</artifactId>
      <version>${scala.compat.version}</version>
      <scope>test</scope>
    </dependency>
    <dependency>
      <groupId>org.scalatest</groupId>
      <artifactId>scalatest_${scala.compat.version}</artifactId>
      <version>3.0.3</version>
      <scope>test</scope>
    </dependency>
    
vi /Users/didi/maven/conf/settings.xml

 在maven的默认配置中，对于jdk的配置是1.4版本，那么创建/导入maven工程过程中，工程中未指定jdk版本。

对工程进行maven的update，就会出现工程依赖的JRE System Library会自动变成JavaSE-1.4。



解决方案1：修改maven的默认jdk配置

           maven的conf\setting.xml文件中找到jdk配置的地方，修改如下：


[html] view plaincopy在CODE上查看代码片派生到我的代码片

<profile>   
    <id>jdk1.6</id>    
    <activation>   
        <activeByDefault>true</activeByDefault>    
        <jdk>1.6</jdk>   
    </activation>    
    <properties>   
        <maven.compiler.source>1.6</maven.compiler.source>    
        <maven.compiler.target>1.6</maven.compiler.target>    
        <maven.compiler.compilerVersion>1.6</maven.compiler.compilerVersion>   
    </properties>   
</profile>  

解决方案2：修改项目中pom.xml文件，这样避免在导入项目时的jdk版本指定

 打开项目中pom.xml文件，修改如下：
<build>  
    <plugins>  
        <plugin>  
            <groupId>org.apache.maven.plugins</groupId>  
            <artifactId>maven-compiler-plugin</artifactId>  
            <configuration>  
                <source>1.6</source>  
                <target>1.6</target>  
            </configuration>  
        </plugin>  
    </plugins>  
</build>  

右键－》propertity  
   remove jre1.6  
      add jre1.8
      
      
<!-- more -->
运行成功

Could not resolve dependencies for project maven.scala:mavenScala:jar:0.0.1-SNAPSHOT: Failure to find org.specs2:specs2-core_2.12:jar:2.12 in https://repo.maven.apache.org/maven2 was cached in the local repository, resolution will not be reattempted until the update interval of central has elapsed or updates are forced

http://maven.outofmemory.cn/org.specs2/specs2-core_2.12.0-M4/3.8.4/


<dependency>
    <groupId>org.specs2</groupId>
    <artifactId>specs2-core_2.12.0-M4</artifactId>
    <version>3.8.4</version>
</dependency>


删除
<dependency>
    <groupId>org.specs2</groupId>
    <artifactId>specs2-core_2.12.0-M4</artifactId>
    <version>3.8.4</version>
</dependency>


scalac error: bad option: '-make:transitive'
解决方法：

（1）打开pom.xml，删除

       <parameter value="-make:transitive"/>
（2）添加dependance

        <dependency>
            <groupId>org.specs2</groupId>
            <artifactId>specs2_2.11</artifactId>
            <version>2.4.6</version>
            <scope>test</scope>
        </dependency>


测试报错  删除

mvn package

[INFO] Building jar: /Users/didi/PhpstormProjects/ProjGit/Spark/ScalaMaven/MavenScala/target/MavenScala-0.0.1-SNAPSHOT.jar
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ----------------


Description Resource  Path  Location  Type
Project configuration is not up-to-date with pom.xml. Select: Maven->Update Project... from the project context menu or use Quick Fix.  MavenScala    line 1  Maven Configuration Problem

右键  Maven->Update Project

至此没有错误了




