---
title: Eclipse+maven+scala+spark环境搭建
layout: post
category: spark
author: 夏泽民
---
1.安装Scala-IDE
在Eclipse中开发Scala程序需要有scala插件，我们现在安装scala插件 
2.安装m2e-scala插件
m2e-scala用来支持scala开发中对maven的一些定制功能。通过eclipse的Install New Software安装。 
安装过程 
   1.Help->Install New Software 
   2.输入m2e-scala下载的url 
具体URL为http://alchim31.free.fr/m2e-scala/update-site/
	<img src="{{site.url}}{{site.baseurl}}/img/m2e_scala.png"/>
    3.安装完成后，可在Help->Installation Details中查看 
    4.添加远程的原型或模板目录
    	<img src="{{site.url}}{{site.baseurl}}/img/archetypes.png"/>
    	Catalog file:http://repo1.maven.org/maven2/archetype-catalog.xml
Description:Remote Catalog Scala
    5、出现过mvn连不上公共库的问题;
     解决方法：vi eclipse.ini
      add : -vmargs -Djava.net.preferIPv4Stack=true


    	
3.新建Eclipse+scala+maven工程
新建maven工程
此时的maven的Archetype需要设置为 org.scala-tools.archetypes 
如果没有安装Scala-IDE的话，会找不到org.scala-tools.archetypes这个类别 

新建Archetype，因为maven默认没有Group Id: net.alchim31.maven Artifact Id: scala-archetype-simple Version:1.6

　　Select New -> Project -> Other and then select Maven Project. On the next window, search forscala-archetype. Make sure you select the one in group net.alchim31.maven, and click Next。
　　<img src="{{site.url}}{{site.baseurl}}/img/configureScala.png"/>configure
　　
<!-- more -->
<plugin>
  <groupId>net.alchim31.maven</groupId>
  <artifactId>scala-maven-plugin</artifactId>
  <version>3.1.3</version>
  <executions>
    <execution>
      <goals>
        <goal>compile</goal>
        <goal>testCompile</goal>
      </goals>
    </execution>
  </executions>
  
  
  scala的新版本对老版本的兼容似乎并不好。这里可以自己修正pom.xml文件，不过估计代码可能也要修改。从git上下载了一个现成的基于scala2.11.5的maven工程。
git网址：https://github.com/scala/scala-module-dependency-sample
使用git clone下来之后，在eclipse中导入maven工程（maven-sample

或者直接编译scala-maven-plugin
https://github.com/davidB/scala-maven-plugin

 运行Maven是报错：No goals have been specified for this build
 pom.xml文件<build>标签后面加上<defaultGoal>compile</defaultGoal>即可  
 
 一个错误示例，子项目引用了父项目，子项目parent标签处报错如下：
Multiple annotations found at this line:
- maven-enforcer-plugin (goal "enforce") is ignored by m2e.
- Plugin execution not covered by lifecycle configuration: org.codehaus.mojo:aspectj-maven-plugin:1.3.1:compile (execution: 
 default, phase: compile)
 
解决办法
官网给出解释及解决办法：http://wiki.eclipse.org/M2E_plugin_execution_not_covered

这里有人说下面这样也可以解决， 即 <plugins> 标签外再套一个 <pluginManagement> 标签，我试验是成功的：
http://stackoverflow.com/questions/6352208/how-to-solve-plugin-execution-not-covered-by-lifecycle-configuration-for-sprin
<build>
    <pluginManagement>
        <plugins>
            <plugin> ... </plugin>
            <plugin> ... </plugin>
                  ....
        </plugins>
    </pluginManagement>
</build>

 
 
 scala配置
很多时候我们希望可以使用java+scala混合开发模式，此时只需要在maven进行如下配置即可：

<dependencies>
    <dependency>
      <groupId>org.scala-lang</groupId>
      <artifactId>scala-library</artifactId>
      <version>${scala.version}</version>
      <scope>compile</scope>
    </dependency>
</dependencies>

<build>
    <plugins>
      <plugin>
        <groupId>org.scala-tools</groupId>
        <artifactId>maven-scala-plugin</artifactId>
        <version>2.15.2</version>
        <executions>
          <execution>
            <id>scala-compile-first</id>
            <goals>
              <goal>compile</goal>
            </goals>
            <configuration>
              <includes>
                <include>**/*.scala</include>
              </includes>
            </configuration>
          </execution>
          <execution>
            <id>scala-test-compile</id>
            <goals>
              <goal>testCompile</goal>
            </goals>
          </execution>
        </executions>
      </plugin>
    </plugins>   
</build> 

可运行jar打包
<plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-shade-plugin</artifactId>
                <executions>
                    <execution>
                        <phase>package</phase>
                        <goals>
                            <goal>shade</goal>
                        </goals>
                        <configuration>
                            <transformers>
                                <transformer implementation="org.apache.maven.plugins.shade.resource.ManifestResourceTransformer">
                                    <mainClass>{此处填写main主类}</mainClass>
                                </transformer>
                            </transformers>
                            <filters> 
                                <filter>
                                    <artifact>*:*</artifact>
                                    <excludes>
                                        <exclude>META-INF/*.SF</exclude>
                                        <exclude>META-INF/*.DSA</exclude>
                                        <exclude>META-INF/*.RSA</exclude>
                                    </excludes>
                                </filter>
                            </filters>
                        </configuration>
                    </execution>
                </executions>
            </plugin>
            
            
     
