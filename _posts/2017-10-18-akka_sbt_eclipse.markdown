---
title: akka_sbt_eclipse
layout: post
category: spark
author: 夏泽民
---
<!-- more -->
akka项目编译上有两种方法：
第一、 用sbt ，这个是akka 官方推荐的方法，可以用sbt生成Eclipsse项目，然后导入到Eclipse里面，可以运行。 但是我的编译还是通过sbt 命令行的方式来做的。 如果大家有好的方法，请指正。 
sbt的安装方法参考：http://www.scala-sbt.org/download.html
how to write a scala file , build and run  with sbt.
http://www.scala-sbt.org/0.13/tutorial/Hello.html
akka：http://akka.io/
scala for eclipse IDE bundle ： http://scala-ide.org/  （也推荐这个集成开发工具，内含编译所需要的akka actor 和 scala library）

第二、是用maven， 这是akka essentials 这本书所采用的。 我亲自实验过，可以编译akka 2.0.5， 2.1.2的旧有项目。 而且这本书的大多数例子，我都验证过。 所以，如果是想从头到尾、编译调试程序，不妨试试maven 的elipse 插件。
安装mvn  http://maven.apache.org/download.cgi 然后在readme安装步骤做。
Akka-Essentials 源代码 https://github.com/write2munish/Akka-Essentials
参考： 
Maven实战（三）Eclipse构建Maven项目
http://tangyanbo.iteye.com/blog/1503782



1,Eclipse安装Scala的开发插件

2,安装sbt 
  Mac 系统安装很简单：$ brew install sbt
  更多可参考：http://www.scala-sbt.org/0.13/docs/zh-cn/Setup.html
3,创建 akka 项目，并导入Eclipse
 通过sbt创建akka 项目。
   $ mkdir hw
   $ vim build.sbt
   录入(更多详细参考：http://www.scala-sbt.org/1.x/docs/zh-cn/index.html）：
name:="hw"   --项目名
version:= "1.0" –项目版本
scalaVersion:= "2.11.8"  --scala 版本，注意与Eclipse里的Scala版本一致
lazyval akkaVersion = "2.5.4" –akka 的版本
libraryDependencies++= Seq(
  "com.typesafe.akka" %%"akka-actor" % akkaVersion
   )
创建测试例子
$ mkdir -p  scr/main/scala
$ vim src/main/scala/hw.scala
录入并保存：
import akka.actor.Actor
import akka.actor.ActorSystem
import akka.actor.Props

class hw9actor extends Actor
{
  def receive={
    case "hello" => println("hello akka")
    case _=>println("hi")
  }
}
object hw9 extends App {
  val sy=ActorSystem("hw9")
  val helloActor=sy.actorOf(Props[hw9actor], "hello")
  helloActor ! "hello"
  helloActor ! "ww"
  
}
执行、编译、更新
$ sbt        --第一次执行时，会下载相关jar库引用。会保存在当前用户的目录下.sbt .ivy /cache中
$ run       --执行，会搜索src所有可执行的文件，如果有多个会，提供选择执行的列表。如当前只会执行，hw9.scala文件。
$ compile  --如果有修改，重新编译，会使用，$ ~compile 当文件有修改时，自动编译
$ reload   --如果配置文件有修改时，重新加载。
将项目导入Eclipse
进入项目要目录，创建project/plugins.sbt
$ mkdir project
$ vim project/plugins.sbt
录入保存：
addSbtPlugin("com.typesafe.sbteclipse"% "sbteclipse-plugin" % "5.1.0")
$ sbt 
$> reload.   -- 重新加载
$> eclipse.  ---将下载相关引用库
进入Eclipse工具，导入该项目：file -> Import -> Gerneral -> existing projects into workspace

注意：sbt 里引用的scala的版本与Eclipse里scala的版本一致。不然会报版本错误。
参考：https://github.com/typesafehub/sbteclipse

版本不正确，参考：
https://github.com/typesafehub/sbteclipse
设置com.typesafe.sbteclipse

sbt参考文档：
http://www.scala-sbt.org/0.13/docs/zh-cn/Hello.html