---
title: build.gradle
layout: post
category: java
author: 夏泽民
---
<!-- more -->
在语法上是基于Groovy语言的（Groovy 是一种基于JVM的敏捷开发语言，可以简单的理解为强类型语言java的弱类型版本），在项目管理上是基于Ant和Maven概念的项目自动化建构工具。

基础知识准备
Java基础，命令行使用基础
官方文档：https://docs.gradle.org/current/dsl/
** Gradle使用指南：** https://gradle.org/docs/current/userguide/userguide
Android插件文档：https://github.com/google/android-gradle...
AndroidGradle使用文档：http://tools.android.com/tech-docs/new-build-system/user-guide
Groovy基础： http://attis-wong-163-com.iteye.com/blog/1239819
Groovy闭包的Delegate机制：http://www.cnblogs.com/davenkin/p/gradle-learning-3.html

搭建Gradle运行环境
Gradle 运行依赖JVM，也就是java运行的环境。所以要安装jdk和jre，好像目前的Gradle的运行环境要求jdk的版本在1.6以上，应该的，现在jdk都到1.8了。
然后到Gradle官网现在Gradle的压缩包。地址，这个页面里面又两种方式，一种手动安装，一种通过脚本安装。我一般喜欢自己动手，这样将来清理起来比较方便。
下载压缩包后，解压，然后配置环境变量，手动安装过jdk的人应该都配置环境变量很熟了吧。每个平台下配置环境变量的方式不一样
MacOS 下配置。在~/.bash_profile中添加如下代码

#gradle  注意gradle-2.14.1是自己解压的路径
export GRADLE_HOME=${HOME}/gradle-2.14.1
PATH=${PATH}:${GRADLE_HOME}/bin
export PATH
保存后在终端输入source ~/.bash_profile回车执行让刚刚的配置生效。然后命令行输入gradle -v查看是否安装成功。

$ gradle -v

------------------------------------------------------------
Gradle 2.14.1
------------------------------------------------------------

Build time:   2016-07-18 06:38:37 UTC
Revision:     d9e2113d9fb05a5caabba61798bdb8dfdca83719

Groovy:       2.4.4
Ant:          Apache Ant(TM) version 1.9.6 compiled on June 29 2015
JVM:          1.8.0_111 (Oracle Corporation 25.111-b14)
OS:           Mac OS X 10.12.2 x86_64
弄一个HelloWorld看看
创建一个test_gralde文件夹。然后在文件夹里面创建一个build.gradle文件。注意文件名不要乱起。在build.gradle中添加如下代码：

task helloworld{
    doLast{
        println'Hello World!'
    }
}
#后者等同于下面的代码,
task helloworld2 <<{
    println "Hello World!"
}
然后来运行一下：

liuqiangs-MacBook-Pro:test_gralde liuqiang$ gradle helloworld
:helloworld
Hello World!

BUILD SUCCESSFUL

Total time: 1.52 secs

This build could be faster, please consider using the Gradle Daemon: https://docs.gradle.org/2.14.1/userguide/gradle_daemon.html
我们分析一下执行步骤。build.gradle是Gradle默认的构建脚本文件，执行Gradle命令的时候，会默认加载当前目录下的build.gradle脚本文件，当然你也可以通过 -b 参数指定想要加载执行的文件。这只是个最简单的task例子，后面详细介绍task的常见定义。

这个构建脚本定义一个任务（Task），任务名字叫helloworld，并且给任务helloworld添加了一个动作，官方名字是Action，阅读Gradle源代码你会到处见到它，其实他就是一段Groovy语言实现的闭包，doLast就意味着在Task执行完毕之后要回调doLast的这部分闭包的代码实现。第二个方法中的“<<”表示向helloworld中加入执行代码。至于语法部分，基本是Groovy语法（包括一些语法糖，也就是写简写方式，如果写个JavaScript或者Python会好理解一些，但是还是建议去读一下groovy的基本语法），加上一些DSL（domain specific language）的约定。

执行流程和基本术语
和Maven一样，Gradle只是提供了构建项目的一个框架，真正起作用的是Plugin。Gradle在默认情况下为我们提供了许多常用的Plugin，其中包括有构建Java项目的Plugin，还有Android等。与Maven不同的是，Gradle不提供内建的项目生命周期管理，只是java Plugin向Project中添加了许多Task，这些Task依次执行，为我们营造了一种如同Maven般项目构建周期。

Gradle是一种声明式的构建工具。在执行时，Gradle并不会一开始便顺序执行build.gradle文件中的内容，而是分为两个阶段，第一个阶段是配置阶段，然后才是实际的执行阶段。
配置阶段，Gradle将读取所有build.gradle文件的所有内容来配置Project和Task等，比如设置Project和Task的Property，处理Task之间的依赖关系等。

看一个基本结构的Android多Moudule（也就是gradle中的多Project Multi-Projects Build）的基本项目结构。

├── app #Android App目录
│   ├── app.iml
│   ├── build #构建输出目录
│   ├── build.gradle #构建脚本
│   ├── libs #so相关库
│   ├── proguard-rules.pro #proguard混淆配置
│   └── src #源代码，资源等
├── module #Android 另外一个module目录
│   ├── module.iml
│   ├── build #构建输出目录
│   ├── build.gradle #构建脚本
│   ├── libs #so相关库
│   ├── proguard-rules.pro #proguard混淆配置
│   └── src #源代码，资源等
├── build
│   └── intermediates
├── build.gradle #工程构建文件
├── gradle
│   └── wrapper
├── gradle.properties #gradle的配置
├── gradlew #gradle wrapper linux shell脚本
├── gradlew.bat
├── LibSqlite.iml
├── local.properties #配置Androod SDK位置文件
└── settings.gradle #工程配置
上面的是完整的AndroidStudio中的项目结构，我们抽象成Gradle多个Project的样子

├── app 
│   ├── build.gradle #构建脚本
├── module 
│   ├── build.gradle #构建脚本
├── build.gradle #工程构建文件
├── gradle
│   └── wrapper    #先不去管它
├── gradle.properties #gradle的配置
├── gradlew #gradle wrapper linux shell脚本
├── gradlew.bat
└── settings.gradle #工程配置
Gradle为每个build.gradle都会创建一个相应的Project领域对象，在编写Gradle脚本时，我们实际上是在操作诸如Project这样的Gradle领域对象。在多Project的项目中，我们会操作多个Project领域对象。Gradle提供了强大的多Project构建支持。
要创建多Project的Gradle项目，我们首先需要在根（Root）Project中加入名为settings.gradle的配置文件，该文件应该包含各个子Project的名称。Gradle中的Project可以简单的映射为AndroidStudio中的Module。
在最外层的build.gradle。一般干得活是：配置其他子Project的。比如为子Project添加一些属性。
在项目根目录下有个一个名为settings.gradle。这个文件很重要，名字必须是settings.gradle。它里边用来告诉Gradle，这个multiprojects包含多少个子Project（可以理解为AndroidStudio中Module）。
读懂Gradle配置语法
Gradle向我们提供了一整套DSL，所以在很多时候我们写的代码似乎已经脱离了groovy，但是在底层依然是执行的groovy所以很多语法还是Groovy的语法规则。
看一个AndroidStudio中app下的build.gradle的配置

apply plugin: 'com.android.application'

android {
    compileSdkVersion 25
    buildToolsVersion "25.0.0"
    defaultConfig {
        applicationId "me.febsky.demo"
        minSdkVersion 15
        targetSdkVersion 25
        versionCode 1
        versionName "1.0"
        testInstrumentationRunner "android.support.test.runner.AndroidJUnitRunner"
    }
    buildTypes {
        release {
            minifyEnabled false
            proguardFiles getDefaultProguardFile('proguard-android.txt'), 'proguard-rules.pro'
        }
    }
}

dependencies {
    compile fileTree(dir: 'libs', include: ['*.jar'])
    compile 'com.android.support:appcompat-v7:25.1.0'
}
分析第一行apply plugin: 'com.android.application'
这句其实是Groovy语法糖，像Ruby和Js都有这种语法糖，apply实际上是个方法，补上括号后的脚本：apply (plugin: 'com.android.application'),看起来还是有点别扭是不？还有个语法糖，如果方法参数是个map类型，那么方括号可以省略，进一步还原apply([ plugin: 'com.android.application']),不理解的可以去看下Groovy的map的写法，和js一样。所以这行的意思是：apply其实是个方法，接收一个Map类型的参数。

总结两点：1. 方法调用，圆括号可以省略 2. 如果方法参数是个Map，方括号可以省略。

Groovy语言的闭包语法
看上面的dependencies 这其实是个方法调用。调用了Project的dependencies方法。只不过参数是个闭包，闭包的用法在文章开始给出了链接。我们对其进行还原一下：

#方法调用省略了（）我们加上
dependencies ({
    compile fileTree(dir: 'libs', include: ['*.jar'])
    compile 'com.android.support:appcompat-v7:25.1.0'
})
提示一点：如果闭包是方法的最后一个参数，那么闭包可以放在圆括号外面

#所以代码还能写成这样
dependencies (){
    compile fileTree(dir: 'libs', include: ['*.jar'])
    compile 'com.android.support:appcompat-v7:25.1.0'
}
Getter和Setter
Groovy语言中的两个概念，一个是Groovy中的Bean概念，一个是Groovy闭包的Delegate机制。
Java程序员对JavaBeans和Getter/Setter方法肯定不陌生，被设计用来获取/设置类的属性。但在Groovy中就不用那些没用的方法了。即Groovy动态的为每一个字段都会自动生成getter和setter，并且我们可以通过像访问字段本身一样调用getter和setter。比如Gradle的Project对象有个version属性（Property）下面这两行代码执行结果是一样的:

println project.version // Groovy  
println(project.getVersion()) // Java  
Project，Task ，Action
Gradle的Project之间的依赖关系是基于Task的，而不是整个Project的。

Project:是Gradle最重要的一个领域对象，我们写的build.gradle脚本的全部作用，其实就是配置一个Project实例。在build.gradle脚本里，我们可以隐式的操纵Project实例，比如，apply插件、声明依赖、定义Task等，如上面build.gradle所示。apply、dependencies、task等实际上是Project的方法，参数是一个代码块。如果需要，也可以显示的操纵Project实例，比如：project.ext.myProp = 'myValue'

Task:被组织成了一个有向无环图（DAG）。Gradle中的Task要么是由不同的Plugin引入的，要么是我们自己在build.gradle文件中直接创建的。Gradle保证Task按照依赖顺序执行，并且每个Task最多只被执行一次。

Gradle在默认情况下为我们提供了几个常用的Task，比如查看Project的Properties、显示当前Project中定义的所有Task等。可以通过一下命令行查看Project中所有的Task：$ gradle tasks （具体log不再贴出来）。可以看到，Gradle默认为我们提供了dependencies、projects和properties等Task。dependencies用于显示Project的依赖信息，projects用于显示所有Project，包括根Project和子Project，而properties则用于显示一个Project所包含的所有Property。

**Tips: **查看Project中所有的Task：$ gradle tasks 
查看Project中所有的properties：$ gradle properties

在上面的build.gradle中加入如下代码：

task myTask {  
    doFirst {  
        println 'hello'  
    }  
    doLast {  
        println 'world'  
    }  
}  
这段代码的含义：给Project添加一个名为“myTask”的任务
用一个闭包来配置这个任务,Task提供了doFirst和doLast方法来给自己添加Action。

其实build.gradle脚本的真正作用，就是配置一个Project实例。在执行build脚本之前，Gradle会为我们准备好一个Project实例，执行完脚本之后，Gradle会按照DAG依次执行任务。

自定义Task的写法
看下面代码文件路径~/Test/build.gradle：

#1
task helloWorld << {
    println "Hello World"
}
#2 Test文件夹下建一个src目录，建一个dst目录，src目录下建立一个文件，命名为test.txt
task copyFile(type: Copy){
    from "src"
    into "dst"
}
第一个这里的helloWorld是一个DefaultTask类型的对象，这也是定义一个Task时的默认类型，当然我们也可以显式地声明Task的类型，甚至可以自定义一个Task类型。
第二个代码中（type：Copy）就是“显式地声明Task的类型”，执行gradle copyFile test.txt也跑到dst中去了。

如果task声明在根Project的build.gradle中的allprojects()方法中，那么这个Task会应用于所有的Project。

task的依赖关系
Gradle不提供内建的项目生命周期管理，只是java Plugin向Project中添加了许多Task，这些Task依次执行，为我们营造了一种如同Maven般项目构建周期。那么这些task是如何依次执行的这就用到声明的依赖关系taskA.dependsOn taskB看下面代码：

task taskA << {
   println 'this is taskA from project 1'
}

task taskB << {
   println 'this is taskB from project 1'
}

taskA.dependsOn taskB
然后我们在命令行运行：
$ gradle taskA
运行结果会先执行taskB的打印，然后执行taskA的打印

如果是Muliti-Project的模式，依赖关系要带着所属的Project，如taskA.dependsOn ':other-project:taskC' 其中taskC位于和taskA不同的Project中，相对于AndroidStudio来说，就是位于不同的Module下的build.gradle中，而other-project为Module名字。

Task 的type可以自定义（没有深入研究）
自定义Plugin的写法
没有深入研究，给出一个网上的例子：

apply plugin: DateAndTimePlugin

dateAndTime {
    timeFormat = 'HH:mm:ss.SSS'
    dateFormat = 'MM/dd/yyyy'
}

class DateAndTimePlugin implements Plugin<Project> {
    //该接口定义了一个apply()方法，在该方法中，我们可以操作Project，
    //比如向其中加入Task，定义额外的Property等。
    void apply(Project project) {
        project.extensions.create("dateAndTime", DateAndTimePluginExtension)

        project.task('showTime') << {
            println "Current time is " + new Date().format(project.dateAndTime.timeFormat)
        }

        project.tasks.create('showDate') << {
            println "Current date is " + new Date().format(project.dateAndTime.dateFormat)
        }
    }
}
//每个Gradle的Project都维护了一个ExtenionContainer，
//我们可以通过project.extentions进行访问
//比如读取额外的Property和定义额外的Property等。
//向Project中定义了一个名为dateAndTime的extension
//并向其中加入了2个Property，分别为timeFormat和dateFormat
class DateAndTimePluginExtension {
    String timeFormat = "MM/dd/yyyyHH:mm:ss.SSS"
    String dateFormat = "yyyy-MM-dd"
}
每一个自定义的Plugin都需要实现Plugin接口，除了给Project编写Plugin之外，我们还可以为其他Gradle类编写Plugin。该接口定义了一个apply()方法，在该方法中，我们可以操作Project，比如向其中加入Task，定义额外的Property等。

原文地址

Gradle Wrapper
Wrapper，顾名思义，其实就是对Gradle的一层包装，便于在团队开发过程中统一Gradle构建的版本，然后提交到git上，然后别人可以下载下来，这样大家都可以使用统一的Gradle版本进行构建，避免因为Gradle版本不统一带来的不必要的问题。（所以要明白这个东西可以没有，有了只是为了统一管理，更加方便）

生成wrapper
gradle 内置了生成wrapper的task，我们可以命令行下执行：
$ gradle wrapper

生成后的目录结构如下(用过AndroidStudio的很熟悉了)：

├── gradle
│   └── wrapper
│       ├── gradle-wrapper.jar
│       └── gradle-wrapper.properties
├── gradlew
└── gradlew.bat
gradlew和gradlew.bat分别是Linux和Window下的可执行脚本，他们的用法和gradle原生命令是一样的，gradle怎么用，他们也就可以怎么用。在MacOS下运行$ ./gradlew myTask
gradle-wrapper.jar是具体业务逻辑实现的jar包，gradlew最终还是使用java执行的这个jar包来执行相关gradle操作。
gradle-wrapper.properties是配置文件，用于配置使用哪个版本的gradle等
详细的看下gradle-wrapper.properties内容
#Sat Jan 21 14:02:40 CST 2017
distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
distributionUrl=https\://services.gradle.org/distributions/gradle-2.14.1-bin.zip
从上面内容和文件的名称都可以看出，这就是个java的配置文件,上面看到的是自动生成的，我们也可以手动修改。然后看下各个字段的含义：

distributionBase 下载的gradle压缩包解压后存储的主目录
distributionPath 相对于distributionBase的解压后的gradle压缩包的路径
zipStoreBase 同distributionBase，只不过是存放zip压缩包的
zipStorePath 同distributionPath，只不过是存放zip压缩包的
distributionUrl gradle发行版压缩包的下载地址，也就是你现在这个项目将要依赖的gradle的版本。
生成wrapper可以指定参数
生成wrapper可以通过指定参数的方式来指定gradle-wrapper.properties内容。
使用方法如gradle wrapper –gradle-version 2.14这样，这样就意味着我们配置wrapper使用2.14版本的gradle，它会影响gradle-wrapper.properties中的distributionUrl的值，该值的规则是http://services.gradle.org/distributions/gradle-${gradleVersion}-bin.zip
如果我们在调用gradle wrapper的时候不添加任何参数呢，那么就会使用你当前Gradle的版本作为生成的wrapper的gradle version。例如你当前安装的gradle是2.10版本的，那么生成的wrapper也是2.10版本的。注：当前版本指的是环境变量中配置的那个版本。
