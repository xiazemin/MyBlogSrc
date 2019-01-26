---
title: checkstyle
layout: post
category: web
author: 夏泽民
---
现在很多开源工具都可以对代码进行规范审核，比较流行的有以下几款，大致给个简单介绍。

     PMD：是一款采用BSD协议发布的Java程序代码检查工具，可以做到检查Java代码中是否含有未使用的变量、是否含有空的抓取块、是否含有不必要的对象等。
     FindBugs：是一个静态分析工具，它检查类或者 JAR 文件，将字节码与一组缺陷模式进行对比以发现可能的问题。
     Checkstyle：是一个静态分析工具，检查Java程序代码。
     Cppcheck是一种C/C++代码缺陷静态检查工具。
　  PC-Lint也是一种静态代码检测工具，检查C或C++。

      目前，中心使用的是Checkstyle工具。我个人觉得PMD和Checkstyle很类似，都可以以插件的形式集成到Eclipse或是MyEclipse开发环境中。下面就Checkstyle在Eclipse中的使用详细介绍下，同时，也欢迎大家继续补充。

      Checkstyle可以从其官网http://checkstyle.sourceforge.net/中下载。官网中还提供了Checkstyle的相关文档，如配置文件、代码检查项等，内容比较丰富，覆盖面也较齐全。可依据自身需要，参考官网上的相关资料。进入Checkstyle的官网后，进入Download页面，可以下载Checkstyle。目前，大多数开发项目使用Eclipse或是MyEclipse的集成开发环境，因此我推荐进入http://en.sourceforge.jp/projects/sfnet_eclipse-cs/releases/下载，该网站上有EclipseCheckstyle Plug-in的各种版本。下文的介绍中，我采用的是net.sf.eclipsecs-updatesite_5.5.0.201111092104-bin.zip版本的Checkstyle插件。


http://checkstyle.sourceforge.net/
Checkstyle是一种代码规约工具，可以帮助程序员编写符合编码标准的Java代码。使用

Checkstyle可以自定义一些代码规范用于在编译中强制执行。CheckStyle默认使用的是Sun的编码规范。本文以8.8版本进行编写。

参考资料：http://checkstyle.sourceforge.net/index.html

例如：

Sun：http://checkstyle.sourceforge.net/sun_style.html

Google：http://checkstyle.sourceforge.net/google_style.html，

对应文件下载地址（Github）：

https://github.com/checkstyle/checkstyle/tree/60f41e3c16e6c94b0bf8c2e5e4b4accf4ad394ab/src/main/resources

         我们可以在提供的文件基础上做修改，也可以自己创建新的文件完全自定义。

注意事项：

         CheckStyle7以后的版本需要JDK1.8编译要求。然而社区的成员已经创建了最新的Checkstyle版本的一个非官方的backport，以便可以在JDK1.6基础上运行，如果想在JDK1.8之前版本运行CheckStyle 7以后版本可以看看，具体的我就没去关注了。

         Checkstyle配置模块需要在checker根模块下。大部分模块是TreeWalker子模块。TreeWalker通过将每个Java源文件分别转换为抽象语法树，然后将结果交给每个子模块进行操作检查。例如典型的配置如下：

<module name="Checker">

   <module name="JavadocPackage"/>

   <module name="TreeWalker">

       <module name="AvoidStarImport"/>

       <module name="ConstantName"/>

       <module name="EmptyBlock"/>

   </module>

</module>
<!-- more -->
CheckStyle检验的主要内容
·Javadoc注释
·命名约定
·标题
·Import语句
·体积大小
·空白
·修饰符
·块
·代码问题
·类设计
·混合检查（包括一些有用的比如非必须的System.out和printstackTrace）
从上面可以看出，CheckStyle提供了大部分功能都是对于代码规范的检查，而没有提供像PMD和Jalopy那么多的增强代码质量和修改代码的功能。但是，对于团队开发，尤其是强调代码规范的公司来说，它的功能已经足够强大。

CheckStyle作为Eclipse插件在Eclipse-Luna-SR2中的安装和使用方法。

在Eclipse中点击Help->Install New Software...；

在弹出的窗口中Work with中填写“http://eclipse-cs.sourceforge.net/update”后点击右侧Add...按钮弹出对话框，在填写完Name栏（可以为空）后点击OK按钮；


enkins 配置checkstyle
首先，我们先在jenkins上新建一个item： 
然后，就给项目命名和选择项目类型： 
save完之后，项目就新建好了。 
接下来讲讲配置checkstyle，要支持checkstyle就要在pom文件里添加checkstyle的支持。
第一种就是刚刚在上面的pom.xml中提到的。 
<configLocation>fcm-cs-check.xml</configLocation>  
必须把文件的配置放在<build>元素里面。参考阅读：http://stackoverflow.com/questions/8975096/maven-checkstyle-configlocation-ignored 
并且fcm-cs-check.xml 必须要跟pom.xml是同一层目录的。 
第二种方法是： 
在mvn 命令中指定checkstyle.config.location,参考：https://dustplanet.de/howto-use-your-own-checkstyle-rules-in-your-jenkinsmaven-job/ 


golang：https://github.com/qiniu/checkstyle
