---
title: artifactory
layout: post
category: web
author: 夏泽民
---
https://s0about0gitlab0com.icopy.site/devops-tools/jfrog-artifactory-vs-gitlab.html
JFrog Artifactory是一种工具，用于存储构建过程的二进制输出以用于分发和部署. Artifactory支持多种软件包格式，例如Maven，Debian，NPM，Helm，Ruby，Python和Docker. JFrog提供高可用性，复制，灾难恢复，可伸缩性，并且可以与许多本地和云存储产品一起使用.

GitLab还提供了高可用性，复制和可伸缩性，也可以使用本地或云存储来使用，但目前提供的包类型兼容性比Artifactory（Maven，Docker，NPM）要少. 但是，GitLab还提供了功能，可以自动完成从规划到创建，构建，验证，安全测试，部署和监视的整个DevOps生命周期. GitLab提供的内置二进制软件包存储库只是冰山一角.
<!-- more -->
Jfrog的Artifactory是一款Maven仓库服务端软件，可以用来在内网搭建maven仓库，供公司内部公共库的上传和发布，以提高公共代码使用的便利性。

1 Artifactory开源版本在Linux下的安装与启动

安装相对比较简单，从Jfrog网站下载当前最新版本的rpm包到本地，当前最新版是4.7.6，执行rpm -ivh命令进行安装。默认安装到/opt/jfrog目录下。

安装包里集成了tomcat，所以不需要再单独安装。但tomcat依赖于java1.8，所以还需要安装jre1.8。这一项就不多说了，去这个页面下载并安装。

成功安装后，切换到安装目录/opt/jfrog/artifactory/bin目录下，执行artifactoryctl start，默认会在8081端口开启服务。这时从浏览器里访问即可。初始用户为admin:password，对于管理员，在安装后需要先修改登录密码。


2 Artifactory的配置

详细的使用文档请见这里

首先介绍下仓库的分类，在Art中，repo有三种。本地Local型，远程Remote型，以及虚拟型。

本地私有仓库：用于内部使用，上传的组件不会向外部进行同步。

远程仓库：用于代理及缓存公共仓库，不能向此类型的仓库上传私有组件。

虚拟仓库：不是真实在存储上的仓库，它用于组织本地仓库和远程仓库。

了解了这点，就可以在admin -> repository 下相应的子类型中，创建新的仓库。


3 仓库的使用

Art安装好之后，预置了一些仓库。例如远程仓库默认配置了jcenter。因此在启动之后，就可以把原来引用公网仓库的组件，改为引用内网仓库了。

如果需要将一些私有代码打包到仓库，就需要使用仓库。创建好内部仓库后，可以将下面的gradle脚本加入到build.gradle中，然后执行uploadArchive的gradle任务，即可将代码编译打包并上传到仓库。

Artifactory  jfrog 家的用来做仓库管理和持续集成【配合 Jenkins 】的工具  免费版就够用了 【支持 maven gradle】
Maven Maven是Java开发者中流行的构建工具，Maven的好处之一是可以帮助减少构建应用程序时所依赖的软件构件的副本，Maven建议的方法是将所有软件构件存储于一个叫做repository的远程仓库中。
Gradle  是 Android Studio 中带的自动化构建工具 是 maven 的扩展
Nginx 是一个高性能的HTTP和反向代理服务器，也是一个IMAP/POP3/SMTP服务器，用来处理代理的。

注意三个文件:

artifactory.sh  用来直接运行 artifactory 的进程，运行之后就会打开 tomcat ，并且部署一个可视化的网页 http://<你的ip>/artifactory/webapp/#/home复制代码
installService.sh 这个是用来安装 artifactory 的服务，可以作为服务在后台自动运行，并会随服务器一起启动【我猜的。。我看他是移动到了init的目录下面】
artifactoryManage.sh 是用来做服务管理的提供几种方式    {start|stop|restart|redebug|status|check} 这个就不翻译了，可以看到当前的状态。
   使用方法这样  ./artifactoryManage.sh check  加命令复制代码

注意的地方：一般的教程都是让你，直接执行 artifactory.sh 就可以启动了。其实服务端更多的时候希望他是作为后台常驻的。所以这里我们要执行的  installService.sh 脚本执行完之后 会看到给我们的提示帮我们移动到了 /etc/init.d/artifactory 目录中通过这两个指令可以检查和启动后台的服务。/etc/init.d/artifactory check/etc/init.d/artifactory start
 这里要注意 artifactory.default 中的 user 的配置 ！默认的是设置为 artifactory 的，但是 artifactory 用户的权限不够【可能是我们服务器配置的原因】，会导致 /etc/init.d/artifactory start 由于权限不够而无法启动 tomcat 的。这是当时困扰了我很久问题。
按照上面的操作，你应该已经能看到 Artifactory 的界面了刚进去的时候会让你设置 admin 密码，同时设置仓库类型。都完成之后是这样的界面：
这里学到了如下几个 Linux 指令
ps -ef | grep artifactory复制代码ps -ef 查看所有的进程，通过 grep 进行过滤，可以看到和 artifactory 相关的进程，拿到 pid 之后通过
kill -9 目标id复制代码就可以停止目标进程
Lib 的上传
下面的内容就相对简单了。在你的 lib 的工程的 build.gradle 中增加如下插件的依赖
buildscript {
    repositories {
        jcenter()
    }
    dependencies {
        classpath 'com.android.tools.build:gradle:2.3.0'
        classpath 'org.jfrog.buildinfo:build-info-extractor-gradle:latest.release'
        // NOTE: Do not place your application dependencies here; they belong
        // in the individual module build.gradle files
    }
}复制代码接着在你需要上传的 lib 的 module 的 build.gradle 的文件中增加如下配置：这个是用来配置上传的路径和账号信息的
artifactory {
    contextUrl = MAVEN_LOCAL_PATH
    publish {
        repository {
            // 需要构建的路径
            repoKey = 'gradle-release-local'

            username = 'admin'
            password = '这里是密码'
        }
        defaults {
            // Tell the Artifactory Plugin which artifacts should be published to Artifactory.
            publications('aar')
            publishArtifacts = true

            // Properties to be attached to the published artifacts.
            properties = ['qa.level': 'basic', 'dev.team': 'core']
            // Publish generated POM files to Artifactory (true by default)
            // POM 文件
            publishPom = true
        }
    }
}复制代码还有是配置上传的版本信息
def MAVEN_LOCAL_PATH ='http://192.168.111.11:8081/artifactory'
def ARTIFACT_ID = 'testsdk'
def VERSION_NAME = '3.0.0'
def GROUP_ID = 'cn.test.test'

publishing {
    publications {
        aar(MavenPublication) {
            groupId GROUP_ID
            version = VERSION_NAME
            artifactId ARTIFACT_ID

            // Tell maven to prepare the generated "*.aar" file for publishing
            artifact("$buildDir/outputs/aar/${project.getName()}-release.aar")

        }
    }
}复制代码最主要的是配置这两个 task。然后开始执行上传！步骤如下

clean 初始化
assembleRelease 构建 aar 
artifactoryPublish 发布到 Artifactory 中


问题这里第一次 publish 的时候上传失败,遇到问题，说找不到 POM 文件。我去对应的路径里面找，确实没有生成。一开始我的操作是把上面的 artifactory 中的 defaults 的 publishPom 设置为 false 。这样能顺利 build 的，但是没有上传 POM 文件。

导致了后面在 demo 中通过 compile 'cn.test.test:testsdk:3.0.0' 这样的形式找不到包，必须通过明确的 aar 后缀的 complie 方式才能找到包，估计 POM 文件是起到类型配置的作用的。

正确的操作应该执行一下 artifactoryPublish 下面的 generatePomFileForAarPublication 就会生成了。🙂 呵呵 没想到吧反正我是翻了大量的资料，最后自己发现的。。。看的懂英文多重要！
上传上去需要配置的三个参数 

ARTIFACT_ID 你的库名字
GROUP_ID 库的包名【可以这么理解】
VERSION_NAME 库的版本号

Lib 的集成
当上一步的包上传完成之后，在你的本地通过下面两个配置就可以测试了。首先在 Demo 的项目 gradle 增加 maven 库的地址，记得和你上面的对应。大概是这样
allprojects {
    repositories {
        jcenter()
        maven { url "http://192.168.111.11:8081/artifactory/gradle-release-local" }
    }
}复制代码在 Demo 的 app 的 gradle dependencies 加上  
compile 'cn.test.test:testsdk:3.0.0'复制代码sync 一下你的 gradle 文件，需要的插件就下下来了~
注意这个 compile 的格式,是根据这个规则生成的。
GROUP_ID:ARTIFACT_ID:VERSION_NAME

进制文件的电脑语言，毫无疑问想要追踪是极其不易的。对任何一个开发者而言，这是一个常见的通病，他们在开发App途中写入的可读代码，比如Python，在实现其操作性之前必须先转化为二进制文件。

Santa Clara和来自以色列的JFrog——一家储存各种各样的二进制编码的公司，最近对外宣布称将发行他们最新的产品：一个可以处理所有形式的二进制构件的全球系统。这家公司称，他们是全球第一家发行这个系统的公司，这个系统可以支持所有类型的软件包和技术。

JFrog的服务初衷是为处于刚刚起步阶段的创业企业提供便利，后来又增加了一系列集成名单，其中有很多工具可以帮助开发者更快地工作，实现智能化工作。在这些集成工具中，有Docker的产品，有Black Duck、Maven、Bower、npm和Git LFS。

这款服务既可以供用户使用云服务，或是提供软件服务许可证。

Artifactory储存器有什么功能？

因为知道电子产品生命周期短，而artifactory一般只在产品使用过程中起作用，所以有必要在一开始就采取行动。当一个开发者编写一个新的App时，他们可以使用任意一种源代码，比如Java、Python或其他无数种编程语言。这些编码是可读的，这些语言和结构是可破译的。

简单来说，就目前的情况来看，像JFrog这样的服务商对于科技领域来说是不可或缺的，大量的编码都需要根据App和程序的更新而同步更新，而且这种更新通常都是大规模，所以几乎没有一种办法可以来管理这种混乱的情况。

然而在20年前，二进制编码都是按月更新或者是按季度更新的，再往前十年，更新速度也不是很快。但是现在，更新速度是按天计算的，随着编码的不断发展，出现了越来越多的新变化，想要管理这些编码，则需要跟踪管理世界上各个不同的地方，这几乎是不可能完成的的任务。

类似于Docker这样的游戏制造商对这个行业的影响是巨大的，随着游戏的不断升级，每天都会更新成百上千的新二进制编码。其中一个最艰巨的挑战就是跟踪管理这些编码，还要搞清楚这些编码究竟应该怎么区分。

这就是像JFrog这样的artifactory储存服务公司可以涉猎的范围，他们可以作为储存二进制代码的固定场所。而且，他们还可以添加独特的meta数据进去，这样就可以自动帮助用户过滤掉大量多余的二进制编码，还可以帮助他们找到正确的版本。在这种情况下，用户无需转码，只要用他们的Antifactory 查询语言，就可以解决一切烦恼。

彻底改革产业

在和JFrog的市场部执行副总监Adam Frankl交谈时，他解释说他们之所以能取得这样突破性的成就，主要是因为他们“创造了一个很灵活的系统，可以解决所有二进制的问题”，他接着说道，“我们扩大了数据库，所以它可以解决任何类型的开发工件，不管是像Docker这样的储存器还是从其他类型的包装发展而来的储存器。”

他还指出，在公司发展和研发新技术的过程中，他们的顾客起了至关重要的推动作用。他说他们的用户都在不断利用多种技来帮助他们研发新产品，这使得这个全球性的平台成为他们与用户交流的必需品。

通过单一的储存位置，用户现在可以更快地工作，而且可以节约时间，不用再花气力通过各种储存器将代码转化为二进制代码。

直面竞争

在看artifactory管理领域里的其他公司时，Frankl指出，Docker Trusted Registry是他们最大的直接竞争对手。然而，在他们自己的某些内容里面，他们却和Docker的技术展开了密切的合作关系，Trusted Registry提供一种让用户可以自行管理他们的代码的服务，这和JFrog的产品的功能很相似。

但是，两者之间最关键的不同点就是JFrog的全球储存系统可以让用户使用所有类型的技术，不会让他们只使用Docker生态系统里的工具，用户没有任何限制。

https://www.jfrogchina.com/
