---
title: Serverless
layout: post
category: k8s
author: 夏泽民
---
https://zhuanlan.zhihu.com/p/65914436
Serverless 简介
根据 CNCF 的定义，Serverless 是指构建和运行不需要服务器管理的应用程序的概念。(serverless-overview)

Serverless computing refers to the concept of building and running applications that do not require server management. --- CNCF
其实 Serverless 早已和前端产生了联系，只是我们可能没有感知。比如 CDN，我们把静态资源发布到 CDN 之后，就不需要关心 CDN 有多少个节点、节点如何分布，也不需要关心它如何做负载均衡、如何实现网络加速，所以 CDN 对前端来说是 Serverless。再比如对象存储，和 CDN 一样，我们只需要将文件上传到对象存储，就可以直接使用了，不需要关心它如何存取文件、如何进行权限控制，所以对象存储对前端工程师来说是 Serverless。甚至一些第三方的 API 服务，也是 Serverless，因为我们使用的时候，不需要去关心服务器。

当然，有了体感还不够，我们还是需要一个更精确的定义。从技术角度来说，Serverless 就是 FaaS 和 BaaS 的结合。

Serverless = FaaS + BaaS。






简单来讲，FaaS（Function as a Service） 就是一些运行函数的平台，比如阿里云的函数计算、AWS 的 Lambda 等。

BaaS（Backend as a Service）则是一些后端云服务，比如云数据库、对象存储、消息队列等。利用 BaaS，可以极大简化我们的应用开发难度。

Serverless 则可以理解为运行在 FaaS 中的，使用了 BaaS 的函数。

Serverless 的主要特点有：

事件驱动
函数在 FaaS 平台中，需要通过一系列的事件来驱动函数执行。
无状态
因为每次函数执行，可能使用的都是不同的容器，无法进行内存或数据共享。如果要共享数据，则只能通过第三方服务，比如 Redis 等。
无运维
使用 Serverless 我们不需要关心服务器，不需要关心运维。这也是 Serverless 思想的核心。
低成本
使用 Serverless 成本很低，因为我们只需要为每次函数的运行付费。函数不运行，则不花钱，也不会浪费服务器资源
<!-- more -->
https://blog.csdn.net/broadview2006/article/details/80132302
BaaS
　　BaaS（Backend as a Service，后端即服务）是指我们不再编写和/或管理所有服务端组件。与虚拟实例和容器相比，在概念上它更接近SaaS（软件即服务）。SaaS主要是业务流程的外包——HR、销售工具，或者从技术端来讲，像Github这样的产品，而BaaS则是要把应用拆分为更小的颗粒，其中一部分完全使用外部产品实现。

　　BaaS 服务都是领域通用的远程组件（而不是进程内的库），可以以 API 的形式使用，深受移动 App 或者单页Web app开发团队的欢迎。因为这些团队通常会使用大量的第三方服务，否则他们就要自己花很多精力做这些事情。我们来看一些例子。

　　Google Firebase是完全由云厂商（Google）管理的数据库，可以直接在移动或者Web应用中使用，而无须经过我们自己的中间层应用服务器。这解释了BaaS的一个方面：用服务替我们管理数据组件。

　　BaaS服务还允许我们倚赖其他人已经实现的应用逻辑。对于这点，认证就是一个很好的例子。很多应用都要自己编写实现注册、登录、密码管理等逻辑的代码，而其实对于不同的应用这些代码往往大同小异。完全可以把这些重复性的工作提取出来，再做成外部服务，而这正是Auth0和Amazon Cognito等产品的目标。它们能实现全面的认证和用户管理，开发团队再也不用自己编写或者管理实现这些功能的代码。

　　BaaS这个词是随着移动应用开发火起来的。事实上，它有时指的是MBaaS（Mobile Backend as a Service）。然而，使用完全托管在外部的产品来开发应用，这个理念并不是移动开发，或者更一般地说，前端开发，所独有的。比如，我们可以不再管理EC2机器上的MySQL数据库服务器，转而使用Amazon的RDS服务，或者我们可以用Kinesis取代我们自己的Kafka消息总线。其他数据基础设施服务还有：文件系统/对象存储（如Amazon S3）、数据仓库（如Amazon Redshift），而更面向逻辑的服务，比如语音分析（如Amazon Lex）以及前面提到的认证，也可以直接在服务端组件中使用。这其中有很多服务都可以认为是Serverless，但并非全部都是。
FaaS/Serverless计算
　　事实上，Serverless 还有一半是 FaaS（Function asa Service，也即函数即服务）。FaaS 是Compute as a Service（计算即服务）的一种形式。事实上，有些人（特别是AWS）说FaaS就是Serverless计算。当然，不可否认，AWS的Lambda是如今被采用得最广泛的FaaS实现。

　　FaaS是一种构建和部署服务端软件的新方式，面向部署单个的函数或者操作。关于Serverless许多时髦的词儿都来自FaaS。很多人认为Serverless就是FaaS，其实他们是只知其一不知其二。

　　我们部署服务端软件的传统方式都是这样开始的：先要有一个主机实例，一般是一个虚拟机（VM）或者容器（见图1），然后把应用部署在主机上。如果主机是VM或者容器，那么我们的应用就是一个操作系统进程。通常我们的应用里包含各种相关操作的代码——比如一个Web服务可能要收回和更新资源。

图1 传统服务端软件的部署

　　FaaS改变了这种部署模式（见图2）。我们去掉主机实例和应用进程，仅关注表达应用逻辑的那些操作或者函数。我们把这些函数上传至由云厂商提供的FaaS平台。

图2 FaaS软件部署

　　但是在一个服务器进程中，函数不是一直处于运行状态的，它们只会在需要的时候才运行，其他时间都是空闲状态（见图3）。我们可以对FaaS平台进行配置，让它为每一个操作监听特定事件。一旦该事件发生，平台就会实例化Lambda函数，然后再用这个触发事件来调用该函数。

图3 FaaS函数生命周期

　　一旦这个函数执行完毕，FaaS平台就可以随意销毁它。或者，平台将其保留一会儿，直到有另一个事件需要处理。

　　FaaS本质上是事件驱动的途径。除了提供一个平台保存和执行代码，FaaS供应商还会将各种同步和异步事件源集成起来。比如 HTTP API Gateway 就是一个同步源；而托管的消息总线、对象存储，或者协调的事件就是异步源。

　　2014年秋 Amazon 发布了 AWS Lambda，经过3年时间，该产品已经逐渐成熟，开始被一些企业采纳。有些Lambda函数的使用量非常少，一天就几次，而也有些公司使用Lambda每天处理数十亿事件。截至本文写作之时，Lambda已经集成了15种以上的不同事件源，可以满足各种不同应用的需求。

　　除了大家所熟识的 AWS Lambda 之外，微软、IBM及Google等大公司，以及一些更小的厂商比如Auth0，也提供商业FaaS。正如各种Computer-as-a-Service平台（IaaS、PaaS、Container-as-a-Service）一样，现在也有一些开源FaaS项目，你可以在自己的硬件或者公有云平台上运行。目前这类私有FaaS还处于混战时代，并没有明显的冒尖者。Galactic Fog、IronFunctions及Fission（使用的是Kubernetes），以及IBM公司自己的OpenWhisk均属于开源系的FaaS。

Serverless的关键
　　从表面上看，BaaS和FaaS是两码事——前者是把应用中的各个部分完全外包出去，后者是一种新的运行代码的托管环境。那么，为什么要把它们都划归为Serverless呢？关键在于，它们都不需要你管理自己的服务器主机或者服务器进程。一个完全Serverless的app不需要你考虑架构中的任何东西。你的应用逻辑——不管是自己编程实现，还是使用第三方服务集成——运行在一个完全弹性的操作环境里。你的状态也是以同样弹性的形式存储的。

　　Serverless并不意味着没有服务器，而是你不需要操心服务器相关的事情。

跨越式变革
　　Serverless是变革。过去十年来，我们已经把应用和环境中很多通用的部分变成了服务。Serverless也有这样的趋势——如把主机管理、操作系统管理、资源分配、扩容，甚至是应用逻辑的全部组件都外包出去，把它们都看作某种形式的商品——厂商提供服务，我们掏钱购买。这是云计算向纵深发展的一种自然而然的过程。

　　但是，Serverless给应用架构带来巨大的变化。直到现在，大多数云服务并没有从根本上改变我们设计应用的方式。比如，使用Docker时，把一个小“箱子”放到应用边上，但是它仍然是一个箱子，而我们的应用也没有显著改变。当我们把自己的MySQL实例托管到云上时，还是要思考需要怎样的虚拟机来处理负载，考虑故障转移问题。

　　Serverless带来跃进式的变化。Serverless FaaS开启的是一种全新的应用架构，完全由事件驱动。更细粒度的部署，需要在 FaaS 组件外面持久化状态。Serverless BaaS把我们从编写逻辑组件中解放出来，但是我们必须将应用与云厂商提供的特定接口与模式集成。

http://www.broadview.com.cn/book/5084

SaaS鼻祖SalesForce喊出的口号『No Software』吗？SalesForce在这个口号声中开创了SaaS行业，并成为当今市值460亿美元的SaaS之王。今天谈谈『No Server』有关的事。继OpenStack、Docker 、MiscroService、Unikernel、Kubernetes和Mesos之后，ServerLess正成为Google亚马逊乃至创业公司暗战的新战场，它们能否成为云计算领域的颠覆性趋势？
 
1、开始于Eucalyptus终结于OpenStack，不仅是从众而且想取巧并弯道超车
在2008、2009、2010那三年，虽然过来还处于云计算的蛮荒年代，但国外敏锐的开发者和公司已经看到的云计算的星星之火。像国人一样，他们 也想在这火成燎原之势前，以最小的风险、最小的投入、最快的时间，分得较大一杯羹。

当然，最先行、最敏锐的是开发者们。

彼时，从虚拟化管理到公有云API，热闹异常。

虚拟化管理曾领过风骚的就有Virtualiron、3tera、qlayer、OpenNebula、Abiquo、virt-manager、oVirt、XenServer、ConVirt、Ganeti、Proxmox、Enomalism。相信我，2009年4月的时候，最牛的就是这几个：Enomalism、ConVirt、oVirt、Virtual Manager。那会连Eucalyptus都感觉很难用。

当然，最后归结于Eucalyptus 、CloudStack和OpenStack。

关于他们的优劣和成败原因的分析，已有很多。三者中，Eucalyptus出身最学术，CloudStack出身商业味最浓，OpenStack介于两者之间。或许，这是OpenStack成功的重要因素？我认为采用Python语言也是重要因素之一。

Eucalyptus出现最早，2008年就出来了，是由加州大学圣塔芭芭拉分校的Rich Wolski和他的博士弟子们开始的。NASA最开始也是使用的Eucalyptus啊。但是，学术机构，还是从事HPC的嘛，虽然一开始就对标和兼容EC2 API，但可用性还是差了那么些，特别是对商业运作敏感性查了一些。及至后来引进了MySQL的创始人Martin，加快了商业化步伐，怎奈OpenStack势头已成，且OpenStack没有商业公司控制，这一点很重要，大咖们都喜欢玩不受商业公司控制的，由基金会管理的项目。

OpenStack出现的晚多了，2010年才出现。NASA最开始也是使用Eucalyptus，据说NASA想给Eucalyptus开源版本贡献patch，但Eucalyptus没接受，Eucalyptus不会玩社区嘛。NASA 的六个开发人员，用了一个星期时间用Python 做出来一套原型，结果，能跑。这就是Nova的起源。NASACIO 跟 Raskspace一个副总走得近，于是NASA 贡献 Nova，Raskspace贡献Swift ，在2010年的7月发起了 OpenStack 项目。

CloudStack也在2008年成立。CloudStack一开始就用了cloud.com这个域名，让我觉得背后的团队太有商业眼光了。这个域名肯定花了不少钱，但将来一定能挣回来。因为那会用全部力量扑在虚拟化上的团队，不多。那会，OpenStack刚成立的时候，CloudStack还是OpenStack的成员呢。为嘛？我也没想通。

但是，CloudStack和Eucalyptus一样，最开始并不开源，开源后还留了点尾巴，而且自己控制着商业版本。等发现这种模式玩不转了，2011年7月卖给了Citrix，全部开源后，发现大家已经都在玩OpenStack了。

OpenStack太好了？为什么？有人已经贡献了很多代码了，其实OpenStack发布后直到CloudStack被Citrix再转卖出去为止，OpenStack的易用性和稳定性一直和CloudStack有差距。但是架不住OpenStack免费、完全开源、没有商业公司控制。

这多好，天上掉下的馅饼。拿来就能用，改改界面就是自己的版本啦。就有软件和产品了，有解决方案了，甚至可以做公有云和亚马逊AWS一决雌雄了。

嘿，这事RackSpace最有发言权。虽然Rackspace 2015年才明显放弃公有云的全面竞争，但在2010年RackSpace决定发起和开源OpenStack项目是，不说明确，至少已经隐隐觉得肯定搞不过亚马逊AWS了。那时，他们与亚马逊AWS的竞争已有三年。

OpenStack本来想另起炉灶，搞一套和AWS EC2不同的API，利用代码和API的开放性，纠结开发者和业内公司，形成一个生态系统，对抗AWS。但是后来，从API到架构，越来越像AWS。

RackSpace当然是这么想的，但是后来的发展却不受自己控制了。来的小公司当然很多，大腕也是不少。当然，RackSpace的投入也是不少。股票得靠云计算支撑啊。公有云发展慢了，OpenStack也是个招牌。RackSpace的云主机不是收购的slicehost嘛，后来有没有用OpenStack不知道，但云存储基本用的是Swift，基于哪个版本就不好说了。

我猜的是，OpenStack对RackSpace的公有云没有明显帮助。因为OpenStack这样的软件能在一个公有云的运营里面占据的角色太基础了，而且OpenStack需要考虑众多厂家和参与者的需求。

接着上面说，天上掉下的馅饼，谁遇到这好事不嫌弃。其实OpenStack就功能和稳定性来说，那几个大厂家复制一个是没有问题的。但是，声势不够，没有名声。IT圈也是个圈啊，没有名头也没人关注啊。纯商业版的，VMware和微软已经够了，再搞一个没人要的，还不如在当红炸子鸡OpenStack上投入和榨点油水。当然，也有不少，把大把银子和大堆的人力投入在OpenStack上的。

不投不行啊，自己搞一个没人关注。还不如找个有名头的继续包装。大公司无奈如此。小公司反而更好办了，反正一无所有，拿个现成的起步，有是在东西，还有响亮的名头，必须上啊。

国外不知道，中国想以OpenStack为生的公司么有100家，也有50家。

当然，后来的结果大家都知道了。OpenStack搞AWS没戏，投入较大的HP Helion都要关闭了，其他拿OpenStack搞公有云的就更不用说了。

从RackSpace开始，大家拿OpenStack开始搞私有云。私有云？从开头说的那些开源的，到VMware VShpere、微软ystem center，那都是相当成熟的。恩，就是贵了点。

现在OpenStack开始往下走了，私有云了嘛。以前是管理和集成，现在深入到更底层的了，从虚拟机的大页内存、CPU绑定、IO调度，到网络的SDN、NFV，这都得有啊。私有云软件，这个都是可以控制的。OpenStack现在回过神来搞私有云，那这些都得搞。

但是，有多少人和公司想过，自己需要的是一个新玩意还是一个虚拟化管理工具？是OpenStack的复杂性可扩展性还是顺手和可靠？当然没有多少人在用之前，对虚拟化管理软件领域有充分的了解，也不一定有资源做充分的调研，流行的就是好的。

后面的事，大家都知道了。CloudStack2015年底被Citrix转卖给Accelerite，而Eucalyptus早在2014年9月就已经委身于HP。

竞争对手一个个倒下，看似势头无敌的时候，也就是最危险的时候。这不，还没等干掉对手的时候，Docker就来了。Docker的渊源虽然很古老，但本身诞生于2013年，是OpenStack红得发紫，各公司蜂拥而上的那年。

谁也想不到，Docker这个老古董能掀起这么大的波澜，差点让OpenStack翻船。OpenStack最牛的是势，而Docker也是来势汹汹。看看下面的图，IT圈曝光率也是基础啊，面对Docker，OpenStack不心慌才怪。


但其实Docker是半个老古董。

2、Docker这半个老古董，掀起了第二波从众而且想取巧并弯道超车的浪潮
没错，第二波终于来了。

因为不但搞OpenStack的没能搞定公有云，不搞OpenStack的也没能搞掂公有云。得有点新东西出来。

我们先从看看Docker有多老开始。现在看到的Docker的起源都不能算起源，只能说出生。其实Docker也有老祖先七大姑八大姨的。

任何东西都能追溯到老祖先，虚拟化还能追溯到70年代的大型机呢。那个是有点牵强了，但是Docker的直系亲属那也是够老的。

远房的亲戚就不说了，新生代的KVM都6、7年的历史了，老一代的Xen和QEMU从2003年算起都十二三年的历史了。在IT行业，10年已经很久了。

但Docker的近亲，历史也不短，甚至有的更久。Docker它爸LXC从2006年开始开发，到2008年发布0.1.0版本，Docker直接间接使用的其他组件Chroot、cgroup、namespaces、libcontainer的历史当然也不会比LXC短。它叔叔Linux-VServer在2003年就已经发布了1.0版本，这是基于内核态隔离的容器。它还有众多表兄弟Cloudfoundry Warden、BSD Jails、AIX workload partitions、iCore Virtual Accounts、Sandboxie、Virtuozzo等等，其中Virtuozzo、OpenVZ和Solaris Container历史都在十年以上。

关注虚拟化和IDC的，有些Docker的亲戚应该是很熟悉的。这回又提到IDC了，云计算真实上辈子就跟IDC纠缠在一起。收费的Virtuozzo、免费的OpenVZ，那是云计算和云主机出现之前，Xen和KVM出现之前，搞VPS的利器。10年前VPS卖的多火，被视为虚拟主机的升级版。

OpenVZ是Virtuozzo容器技术的linux版，LXC是OpenVZ为了融入内核开发的对应版，开发者大部分是同一批人。LXC及周边工具就是现在Docker的重要组成部分。

Docker出现也就5，6年，但它的大部分身体都出现差不多10年了，你说它是不是半个老古董呢？

Docker本身大家都很熟很熟了，不用赘言。不过经常有人拿Docker和Xen、KVM支持的虚拟机对比，说占用体积小、启动速度快，不是一个层级的东西好嘛，一个在操作系统之上，一个在操作系统之下，不赘言了。Docker当然也不能和LXC等Linux容器对比，都是用的同一系列核心工具，只是管理层不同而已。

Docker在2013年底，2014年初，突然吸引了众人的目光，并在2014年2015年集万千宠爱于一身，就如前两年的OpenStack一样。

回到Docker诞生的年代，2010年，几个年轻人在旧金山成立了一家做 PaaS 平台的公司dotCloud。大家现在都知道Heroku，也是PaaS型，而且，也用到了容器，可能是LXC吧。当然不是新堆栈PaaS，而是传统堆栈PaaS。这和Heroku一样，为开发人员提供操作系统、通用库、特定语言的运行环境，应用的部署、管理、监控等。

dotCloud把需要花费大量时间的重复性工作，抽象成组件和服务，如规范容器的格式、便利容器的生命周期管理等，方便开发者管理和监控自己的应用。

正如我在《云计算时代》一书中所言，新堆栈PaaS离开发者现有技能太远，传统堆栈PaaS离现有堆栈太近。不管哪种，都挡住了开发者掌控一切的意愿。所以，PaaS多年来虽然独立作为一类云服务，并没有足够大的市场。

虽然dotCloud 在2011年初拿到了1000万美元的融资，但这个市场本没有那么大，也没有快速的增长，容纳不下太多的企业。也没有干过Heroku、Engine Yard等公司，自然前景不妙。

dotCloud 的创始人 Solomon Hykes 把大伙召集到一起，大家琢磨了一下，商业是没戏了，那也要搞一把非商业的事情，把我们在容器上的工作起名Docker开源出来。能不能把竞争对手干掉不好说，希望是挺渺茫，但是起码能在社区留个名！万一，开放开源能够搞成更大的势头，公司还能起来呢。

是不是和RackSpace当初搞OpenStack有几分相像？狗急了还要跳墙呢，绝处逢生不是不可能。所谓坚持，耐心，就是这个意思。只是这是一条高风险高收益的路罢了。

LXC，包括OpenVZ，在2013年之前，可不只是Heroku等PaaS公司在用，有些公司内部都在使用，甚至包括我在内的一些公有云从业者也慎重考虑过用容器作为公有云的基础。当然，还真有这么干的。Joyent大概是除了AWS之外，干公有云最早的了吧，可能比RackSpace还早点，就是基于Solaris Zone卖云主机的。Joyent的技术骨干是从SUN出来的，精通Solaris，他们整了一个基于Solaris的精简内核，专门用来跑Zone，类似于CoreOS，叫做SmartOS，基于这个做出了私有云软件SmartDataCenter。说这些可能没几个人知道，但是鼎鼎大名的Node.js很多人熟悉，就是Joyent开发维护的。

没错，LXC和OpenVZ用在企业内部是很好的，非常好。但是限于LXC当时的管理工具欠缺，并没能大规模流行起来。它需要一个契机。这个契机就是Docker，Docker解决了LXC的管理问题。

电视剧总是那么相似，相遇、离别、重逢，受苦受难、遇到高人、报仇雪恨，IT圈的故事也逃脱不了这样的情节。Docker的故事，真的和OpenStack，包括此前的Linux等其他开源软件，有几分相似。

dotCloud把自己在容器上的工作成果Docker开放开源了，开发者和小公司雀跃了：又来一个馅饼啊。这个馅饼解决了一些问题，让其他人和开发者免除了管理和开发工作量，这是次要的。更重要的是，后来参与的开发者、小公司、大公司对Docker的期待，远不止解决容器本身的管理问题，也不只是因为有一批人喜欢而从众，还有一个很多人很多公司没有说出的理由：容器是未来，干掉现在的私有云和公有云。

Docker 如此受欢迎，2013年10月 dotCloud 公司更名为 Docker 公司。2014年8月 Docker把平台即服务的业务dotCloud出售给位于德国柏林的平台即服务提供商cloudControl，2016年初 dotCloud被cloudControl关闭。而Docker公司开始全身心运营Docker社区，并开展Docker商业化。

Docker 迅速成长为云计算相关领域更受欢迎的开源项目，Amazon、Google、IBM、Microsoft、Red Hat 和 VMware 都表示支持 Docker 技术或准备支持。据说，有 Linux 的地方，就可以运行 Docker。Mac上也可以，Windows 上都有直接运行 Docker 的 Windows Docker 和 Nano Server，当然，这已经是 Windows 版的了。

截止2016年初，Docker 获得 5 轮 1.8 亿美元融资，推出了Docker Hub、Docker Trusted Registry、Docker Tutum等产品，试图控制Docker容器的存储、管理。在2015年上半年与OpenStack的论战之后，双方握手言和，以OpenStack支持Docker管理告终。

Docker当然不甘心只是一个工具软件，也是要做产品、平台的，拿投资人的钱可不能做公益做开源啊。凡是有威胁的就要干掉，或者收掉。于是Docker收购Unikernels。

3、Docker为何收购Unikernels？既是看好更是感到威胁
容器作为虚拟化技术的一种，在云计算出现之前就出现了。之所以在2014年前后流行，是因为大家需要一种与硬件虚拟化更轻量级的技术，来有效运转私有的基础设施。这个私有的，既可以在自家机房里也可以在公有云上。

在私有的基础设施上，至少某些应用场景，Docker因为其轻量级，应用启动更为迅速，资源利用更为高效。但是，循着这个思路，Unikernels更进了一步。

我们先来看看怎么回事。

从操作系统诞生以来，它就扮演了一个应用和硬件之间的平台的角色，提供对硬件的控制。除了操作系统内核和基本的控制台，还有软件开发接口、语言运行时环境、语言库、输入输出设备控制，也需要支持各种古老的和新兴的硬件标准。它需要支持多用户、多进程、多设备并发。

虽然操作系统的的用途各异，有桌面的、内部IT系统的、有面向网络的，但操作系统的架构和模块基本相同，一致没有大的改动。但在上世纪90年代后期，Hypervisor被引入了主流的操作系统。Hypervisor运行于硬件之上，能支持多个虚拟机运行多个不同的操作系统。但这一变化，并未对操作系统的设计产生大的影响。每一个虚拟机仍然运行着一个传统的操作系统。

但是当Hypervisor推动了云时代的到来后，通用操作系统的局限就更明显了。在云环境中，由于规模更大，负载被明显分成了不同的类型：Web服务器负责处理网络请求，数据库服务器负责数据库的运行和数据库访问，等等。这些服务器可能永远都用不上显示器、多用户、多进程。这些场景下的VM和OS的任务很明确，就是提供较好的存储、计算、网络性能。

开发者们开始质疑操作系统的的设计和架构：为什么一个Web Server上要安装它从来用不到的软件？其实在此之前，Windows服务器就遇到过类似的问题。我们只需要能快速扩展和收缩VM的规模，VM越精简越轻量级越有利于这种弹性。但由传统操作系统构成的VM，只能勉强完成这个任务。

Docker所代表的容器恰逢其时。因为基础技术早已就绪，流行很快。因为能在同一个操作系统上快速运行多个隔离的轻量级的，容器基本解决了上述疑问。容器封装了应用程序所需要的一切，除了共用的操作系统内核，它封装了运行时环境、框架和库、代码、配置和相关的依赖。容器大大削减了操作系统作为一个全能平台所承担的角色。容器之下的操作系统这时只需要承担一个非常轻量级的角色：启动和控制容器。这时的操作系统更为专业化，而容器承担了运行各种不同场景所需资源的角色。

这种趋势催生了容器的Hypervisor的产生：CoreOS,RancherOS,Redhat Atomic Hosts,Snappy Ubuntu Core,Microsoft Nano。他们是为了支持在其上运行容器而存在，没有其他的任务，已经非常精简。甚至Hypervisor之上的容器，有人又将其区分为操作系统容器和应用容器，应用容器就是只为单一应用为目标的容器。到这里，微服务（MicroService）终于可以见天日了，终于能够使用了，而所谓的SOA、Mashable是不是能够换发新的光彩呢？

但这种精简、轻量级还没有到此为止，Unikernals出现了。在2014年初就出现了，那时Docker刚刚开始崛起。

Unikernals将这种最小化操作系统的理念往前进一步推进。它是一个定制的操作系统，专为特定的应用程序的运行而编译。因此，开发者能够创建一个极度精简的OS，只包含应用和它所需的操作系统组件。Unikernals是单用户、单进程、单任务的定制操作系统，它在编译时去除了所有不需要的功能，但包含了一个软件运行所需的全部堆栈：OS内核、系统库、语言运行时环境、应用，这些被编译成一个可启动的VM镜像，直接运行在标准的Hypervisor上。

Unikernals让操作系统之上的容器又变成了一个操作系统，不过这是一个重新吧便宜的极为紧凑的，直接运行在Hypervisor而不是精简的通用操作系统上的操作系统。Unikernel有着更小的尺寸，更小的可攻击面，启动时间也以毫秒计。Unikernals的论文在这里：https://queue.acm.org/detail.cfm?id=2566628。

不过，和Docker一样，灵活带来的伴生问题，就是管理、监控、回溯、审计，有运维工程师说，Docker就是运维的噩梦，那Unikernals可能是运维和开发者的噩梦？为啥，因为对应用改一行代码要重新编译整个镜像并部署，无法对一个Unikernals实例打补丁，也就是说Unikernals的实例是静态的。

Unikernals的例子包括MirageOS、Clive、OSv，目前都跑在Xen Hypervisor上。它有多小呢？一个MirageOS DNS镜像是200KB，而一个目前全球90%DNS使用的开源DNS服务器BIND镜像是400MB。而MirageOS DNS镜像的性能据称比BIND好45%。

咦？这不是嵌入式系统吗？Unikernals当然可以编译出镜像，运转在低功耗的设备上。如果Unikernals被移植到ARM平台上，开发者就可以编译出运行在ARM设备的hypervisor上的镜像。这将让嵌入式应用的开发更为容易。

那么，看起来，Unikernals虽然现在更像一个极客玩具，但是，不可否认，Unikernals有代替容器和虚拟机的可能，至少在某些场景下。既然容器比虚拟机在某些场景下更好，为什么更精简的Unikernals就不能代替容器和虚拟机呢？

有可能。而且这个理念和Docker代表的容器理念相符。于是，Docker收购了它。大家一起玩，一起干掉虚拟机。

Docker看起来无敌，前景光明。但是，道路还是曲折的。

大佬们是想干掉私有云和公有云的啊，你Docker公司守着这个热门技术不放，控制的紧紧的， 我们怎么玩？不光是大佬，创业公司也不干啊。

首先发难的是CoreOS和谷歌。

4、CoreOS反水，Kubernetes和Mesos把docker打回工具原型

CoreOS首先不干了。

CoreOS原本是Docker初期的铁杆盟友，CoreOS可以说就是为Docker而生的Linux，它的任务就是管理好Docker。但是随着Docker开始商业化，Docker对镜像格式和代码收紧了控制，甚至建立开放平台存储Docker镜像和认证，当然，还发布了Docker管理工具。

那CoreOS的位置在哪里？当然，理由还是要像样一点的：Docker的镜像格式不够灵活，工具链太庞大，我们要灵活而精简的东西。

于是CoreOS自己搞了一个容器：Rocket。本来嘛，大家都是基于LXC等工具的，这些工具都是开源开放，而且CoreOS也搞容器管理的，新搞个格式和管理工具还不是手到擒来。

当然，双方都承认，Docker和Rocket不是直接竞争关系。Docker是一个产品，并正在成为一个平台。Rocket只想做一个容器管理的组件。但是，双方还是分道扬镳了。

CoreOS一个人可没这么大的气势，背后还有谷歌撑腰。Rocket很快被Kubernetes支持。

Kubernetes的灵感来源于Google的内部borg系统，吸收了包括Omega在内的容器管理器的经验和教训。2014年6月谷歌宣布Kubernetes 开源。谷歌想靠容器翻身呢？怎么能让另一个商业云公司控制最流行的容器。

Kubernetes算是一个与Docker平台竞争的容器管理工具，自然首先支持Docker，但是也支持竞争对手Rocket。

但是Kubernetes也有一个潜在对手：Mesos。Mesos比Kubernetes出现得早，而且两者都深受谷歌的数据中心管理你项目Borg和Omega的影响。问题是，Mesos不是谷歌自家的，虽然属于Apache项目，但仍被商业公司Mesosphere，也是Mesos重要维护者主导。Mesos被称为数据中心操作系统，软件定义数据中心的典范。

Mesos既是计算框架也是一个集群管理器，是Apache下的开源分布式资源管理框架，它被称为是分布式系统的内核，可以运行Hadoop、MPI、Hypertable、Spark。使用ZooKeeper实现容错复制，使用Linux Containers来隔离任务，支持多种资源计划分配。Mesos使用了与Linux内核相似的规则来构造，仅仅是不同抽象层级的差别。Mesos从设备（物理机或虚拟机）抽取 CPU，内存，存储和其他计算资源，让容错和弹性分布式系统更容易使用。Mesos内核运行在每个机器上，在整个数据中心和云环境内向应用程序（Hadoop、Spark、Kafka、Elastic Serarch等等）提供资源管理和资源负载的API接口。

Mesos也不是也不是没有隐忧，Apache yarn似乎有一统分布式计算之势，MapReduce，Spark，Storm，MPI，HBase都在向yarn上迁移。当然，好在Mesos不仅仅是做分布式计算的框架。

Mesos也起源于Google的数据中心资源管理系统Borg。Twitter从Google的Borg系统中得到启发，然后就开发一个类似的资源管理系统来帮助他们摆脱可怕的“失败之鲸”， 后来他们注意到加州大学伯克利分校AMPLab正在开发的名为Mesos的项目，这个项目的负责人是Ben Hindman，Ben是加州大学伯克利分校的博士研究生。后来Ben Hindman加入了Twitter，负责开发和部署Mesos。现在Mesos管理着Twitter超过30,0000台服务器上的应用部署。

Borg的论文2015年四月才发布：http://research.google.com/pubs/pub43438.html
Mesos的论文：https://www.cs.berkeley.edu/~alig/papers/mesos.pdf
Omega的论文：http://research.google.com/pubs/pub41684.html。

这一回，谷歌论文发晚了，虽然也很有价值，但可能没有三大论文那么有影响力。

2014年7月，Mircrosoft、RedHat、IBM、Docker、 CoreOS、Mesosphere 和Saltstack 加入Kubernetes。

2015年1月，Google和Mirantis及伙伴将Kubernetes引入OpenStack，开发者可以在OpenStack上部署运行Kubernetes 应用。

2015年7月，Google正式加入OpenStack基金会，Kubernetes的产品经理Craig McLuckie宣布Google将成为OpenStack基金会的发起人之一，Google将把它容器计算的专家技术带入OpenStack,成一体提高公有云和私有云的互用性。

同时，谷歌联合Linux基金会及其他合作伙伴共同成立了CNCF基金会(Cloud Native Computing Foundation)，并将Kubernetes作为较早的编入CNCF管理体系的开源项目。来，我们来看一下发起人：AT&T, Box, Cisco, Cloud Foundry Foundation, CoreOS, Cycle Computing, Docker, eBay, Goldman Sachs, Google, Huawei, IBM, Intel, Joyent, Kismatic, Mesosphere, Red Hat, Switch SUPERNAP, Twitter, Univa, VMware and Weaveworks。

到此是不是大团圆了？所有跟容器有点关系的都来了。谷歌加入了OpenStack基金会，OpenStack上可以部署运行Kubernetes，OpenStack支持Docker，Mesos支持Docker和Kubernetes。大家互相都支持，谁能发展好，各自努力吧。这关系够乱的。

但是，容器的另外两个大玩家-亚马逊和微软，没有到场。没错，容器界就是要掀翻现有的云计算格局，当然不能让云计算老大和老二进来了。

谷歌纠集了一帮小兄弟，誓要利用容器浪潮，干翻亚马逊AWS和微软Azure。当然，谷歌也没有奔到准备靠一招制胜，暗战还有另一个战场：Serverless。

5、Serverless是云计算的决胜负战场？
2014年11月14日，亚马逊AWS发布了新产品Lambda。当时Lambda被描述为：一种计算服务，根据时间运行用户的代码，无需关心底层的计算资源。从某种意义上来说，Lambda姗姗来迟，它更像S3，更像云计算的PaaS理念：客户只管业务，无需担心存储和计算资源。

比如你要架构一个视频服务，你需要用一堆服务器，设计出一套上传、解码、转码的架构。但是，可能这套系统99%的时间都是空闲的。而使用一个Lambda function，你就不需要操心服务器和这套架构，当AWS探测到用户定义的时间，比如用户上传了一个视频文件，Lambda自动运行响应的程序，结束后关闭程序。为客户生了时间和金钱。

Lambda识别Event的速度非常快，以毫秒来计算。它会在图片上传、应用内活动、点击网站或联网设备的输出等事件发生后的几毫秒内，开始运行代码，分配合适的计算资源来执行这个行动。它可以自动扩展到数百万个请求，如需要可跨越多个可用区。根据AWS Lambda是按计算时间收费，计费单位为100毫秒，可以轻松地把应用从每天几次请求扩展到所需要的任何规模的请求

而在此之前不久，2014年10月22日，谷歌今天收购了实时后端数据库创业公司Firebase。Firebase声称开发者只需引用一个API库文件就可以使用标准REST API的各种接口对数据进行读写操作，只需编写 HTML＋CSS＋JavaScrip前端代码，不需要服务器端代码（如需整合，也极其简单）。Firebase是一个实时应用平台，它可以为几乎所有应用的通用需求提供支持，包括数据库、API和认证处理。数据库的特点是基于NoSQL的实时处理能力，并且允许直接存储JSON文档。Firebase具有双向同步的能力，在客户端侧，开发者通过Firebase的客户端库来支持典型场景的需求，比如屏幕刷新、离线时数据访问或者设备重新连接后的再次同步。Firebase对数据存储容量没有限制，较高能处理百万级的并发和TB级的数据传输，数据发生更改，同步敏感颗粒度基本达到10毫秒级别。

相对于上两者，Facebook 在2014年二月收购的 Parse，则侧重于提供一个通用的后台服务。不过这些服务被称为Serverless或no sever。想到PaaS了是吗？很像，用户不需要关心基础设施，只需要关心业务，这是迟到的PaaS，也是更实用的PaaS。这很有可能将会变革整个开发过程和传统的应用生命周期，一旦开发者们习惯了这种全自动的云上资源的创建和分配，或许就再也回不到那些需要微应用配置资源的时代里去了。

AWS的Lambda既不是最早，也不是较好，但它仍然是serverless最有影响力的产品，谁让AWS的用户多，产品全呢？

Serverless是未来吗？也许是。

Serverless能决定云计算的胜负吗？不能。

什么能决定云计算的胜负？

不仅Serverless不能，其他的产品、技术、项目、工具，都不能单独决定云计算的胜负。

从云计算初期的OpenStack和PaaS，到后来的Docker、Kubernets、Mesos、Unikernals，以及最近沸沸扬扬的人工智能，还有大数据分析，IBM甚至宣称Watson是它的云计算秘密武器，甚至可能即将光大的Serverless，都不能单独决定云计算的胜负。

它们都是优秀的产品、技术、项目、工具，但只是一项产品、技术、项目、工具，它们只能用来更好的服务客户，开辟新产品和加强已有产品，可以用来建立新业务或新公司。

IBM或谷歌也许能成为人工智能的王者，Docker也许能成为容器的领袖，Cloudera也许能成为大数据的领军企业，即使如此，这都是另一个领域的事情，是一种应用场景的事情，它们也许能也许不能成就新的行业，但都与云计算基础服务无关，与IaaS和PaaS无关，与公有云胜负无关。

公司管理者和技术控们：指望拿热门开源项目，个人赚点钱可以，要让一个公司鲤鱼跳龙门或翻身，没门，那就不仅仅是你抓的项目名字火不火的问题。

这个世界从来没有独门秘籍。改变云计算格局，颠覆云计算市场的，不会是一个独门技术，也不会是一项秘密产品。

能决定云计算胜负的，仍然是远见、魄力、耐心。如果已经有了早行者，那就需要持续的创新，来蚕食领先者的优势。云计算是一个庞大的市场，从来不是一项技术、一个项目、一个产品的事情。
http://www.dataguru.cn/article-9011-1.html
最近半年来，随着AWS的各线服务都开始支持lambda，serverless architecture便渐渐成为一个火热的话题。lambda是amzon推出的一个受控的运行环境，起初仅仅支持nodejs（之后添加了java/python的支持）。你可以写一段nodejs的代码，为其创建一个lambda资源，这样，当指定的事件来临的时候，aws的运行时会创建你的运行环境，执行你的代码，执行完毕（或者timeout）后，回收一切资源。这看起来并不稀罕，整个运行环境还受到很多限制，比如目前aws为lambda提供了哪些事件支持，你就能用哪些事件，同时你的代码无法超过timeout指定的时间执行（目前最大是5min），内存使用最多也就是1.5G。那么问题来了，这样一个看起来似乎有那么点鸡肋的服务，为什么还受到如此热捧？原因就在于无比低廉的价格（每百万次请求0.2美元 + 每百万GB秒运行时间16.67美元），毋须操心infrastructure，以及近乎无限扩容的能力。


使用lambda处理事件触发

在服务器端，我们所写的大部分代码是事件触发的：


处理用户对某个URI的请求（打开某个页面，点击某个按钮）

用户注册时发邮件验证邮箱地址

用户上传图片时将图片切割成合适的尺寸

当log持续出现若干次500 internal error时将错误日志聚合发给指定的邮箱

半夜12点，分析一天来收集的数据（比如clickstream）并生成报告

当数据库的某个字段修改时做些事后处理

同时，在处理一个事件的过程中，往往会触发新的事件。基本上我们做一个系统，如果能厘清内部的数据流和事件流，以及对应的行为，那么这个系统的架构也就八九不离十了。如果要让我们自己来设计一个分布式的事件处理系统，一般会使用Message Qeueue，比如RabbitMQ或者Kafka作为事件激发和事件处理的中枢。这往往意味着在现有的infrastructure之外至少添置事件处理的broker（MQ）和worker（读取并处理事件的例程）。如果你用aws的服务，SQS（或者SNS+SQS）可以作为broker，然后配置若干台EC2做worker。如果某个事件流的产生速度大大超过这个事件流的处理速度，那么我们还得考虑使用auto scaling group在queue的长度超过一定阈值或者低于一定阈值时scale up / down。这不仅麻烦，也无法满足某些要求一定访问延迟保障的场景，因为，新的EC2的启动直至在auto scaling group里被标记为可用是数十秒级的动作。


lambda就很好地弥补了这个问题。lambda的执行是置于container之中的，所以启动速度可以低至几十到数百毫秒之间，而且它可以被已知的事件或者某段代码触发，所以基本上你可以在不同的上下文中直接调用或者触发lambda函数，当然也可以使用SNS（kenisis）+lambda取代原本用MQ+worker完成的工作。


我们看上述事件的处理：


处理用户对某个URI的请求（打开某个页面，点击某个按钮）：使用API gateway + lambda

用户注册时发邮件验证邮箱地址：


可以在用户注册的流程里直接调用lambda函数发送邮件

如果使用dynamodb，可以配置lambda函数使其使用dynamodb stream在用户数据写入数据时调用lambda

用户上传图片时将图片切割成合适的尺寸


可以配置lambda函数被S3的Object Create Event触发，在lambda函数里使用libMagic的衍生库处理图片。

当log持续出现若干次500 internal error时将错误日志聚合发给指定的邮箱


如果用kenisis来收集log，那么可以配置lambda函数使其使用kenisis stream

半夜12点，分析一天来收集的数据（比如clickstream）并生成报告


使用aws最新支持lambda cronjob

当数据库的某个字段修改时做些事后处理


如果使用dynamodb，同上（配置lambda函数使其使用dynamodb）

如果使用RDBMS，可以使用database trigger + lambda cronjob

想进一步深入代码的童鞋，可以看我的这个repo: tyrchen/aws-lambda-thumbnail · GitHub 。它接收S3的Object Create Event，并对event中所述的图片做resize，代码使用es6完成。为了简便起见，没有使用cloudformation创建/更新lambda function，而是使用了aws CLI（见makefile）。如果想要运行此代码，你需要定义自己的 $(LAMBDA_ROLE)，手工创建S3 bucket并将其和lambda函数关联（目前aws cli不支持S3 event）。


使用lambda处理大数据

lambda近乎无限扩容的能力使得我们可以很轻松地进行大容量数据的map/reduce。你可以使用一个lambda函数分派数据给另一个lambda函数，使其执行成千上万个相同的实例。假设在你的S3里存放着过去一年间每小时一份的日志文件，为做security audit，你需要从中找出非正常访问的日志并聚合。如果使用lambda，你可以把访问高峰期（7am-11pm）每两小时的日志，或者访问低谷期每四小时的日志交给一个lambda函数处理，处理结果存入dynamodb。这样会同时运行近千个lambda函数（24 x 365 / 10），在不到一分钟的时间内完成整个工作。同样的事情交给EC2去做的话，单单为这些instance配置网络就让人头疼：instance的数量可能已经超出了子网中剩余IP地址的数量（比如，你的VPC使用了24位掩码）。


同时，这样一个量级的处理所需的花费几乎可以忽略不计。而EC2不足一小时按一小时计费，上千台t2.small运行一小时的花费约等于26美金，相当可观。


使用lambda带来的架构优势

如果说lambda为事件处理和某些大容量数据的快速处理带来了新的思路，并实实在在省下了在基础设施和管理上的真金白银，那么，其在架构上也带来了新的思路和优势。


web系统是天然离散的系统，里面涵盖了众多大大小小，或并联，或串联的子系统。因为基础设施的成本问题，很多时候我们做了逻辑上的分层和解耦，却在物理上将其部署在一起，这为scalability和management都带来了一些隐患。scalability上的隐患好理解，management上的隐患是指这无形把dev和ops分成不同的team：一个个dev team可以和逻辑上的子系统一一对应，但ops却要集中起来处理部署的问题，免得一个逻辑上「解耦」的功能更新，在物理上却影响了整个系统的正常运行。这种混搭的管理架构势必会影响部署的速度和效率，和「一个team负责从功能开发到上线所有的事情」的思路是相悖的。


举个例子：「用户上传图片时将图片切割成合适的尺寸」这一需求可能在不断变化和优化。对于任何失焦的照片我们还希望做一些焦距上的优化，此外，如果上传的是头像，那么我们希望切割的位置是最合适的头像的位置，如果上传的是照片，除了切割外，我们可能还要生成黑白/灰度等等不同主题的图片。这个功能在不改变已有接口的前提下，并不会影响其他团队的工作，但因为和其他功能放在一起部署，所以部署的工作并不能自己说了算。因为部署交由专门的ops团队完成（可能一天部署一次，也可能一周部署一次），这个团队无法很快地把一些有意思的点子拿出来在生产环境试验，拖累了试错和创新。


而lambda解决了基础设施上的问题，每个子系统甚至子功能（小到函数级的粒度）都可以独立部署，这就让功能开发无比轻松。只要界定好事件流的输入输出，任何事件处理的功能本身可以按照自己团队的节奏更新。


部署和管理上的改变反过来会影响架构，促成以micro-service为主体的系统架构。micro-service孰好孰坏目前尚有争论，但micros-ervice不仅拥抱软件设计上的解耦，同时拥抱软件部署上的解耦是不争的事实。一个web系统的成败和其部署方案有着密切的关系，耦合度越低的部署方案，其局部部署更新的能力也就越强，而一个系统越大越复杂，就越不容易整体部署，所以对局部部署的要求也越来越高。这如同一个有机体，其自我更新从来不靠「死亡-重生」，而是通过新陈代谢。


此外，lambda还是一个充分受限的环境，给代码的撰写带来很多约束条件。我之前在谈架构的时候曾经提到，约束条件是好事，设计软件首先要搞明白约束条件。lambda最强的几个约束是：


lambda函数必须设计成无状态的，因为其所有状态（内存，磁盘）都会在其短短的生命周期结束后消失

lambda函数有最大内存限制

lambda函数有最大运行时间限制

这些限制要求你把每个lambda函数设计得尽可能简单，一次只做一件事，但把它做到最好。很符合unix的哲学。反过来，这些限制强迫你接受极简主义之外，为你带来了无限扩容的好处。


JAWS和server-less architecture

两三个月前，我介绍了JAWS，当时它是一个利用aws刚刚推出的API gateway和lambda配合，来提供REST API的工具，如果辅以架设在S3上的静态资源，可以打造一个完全不依赖EC2的网站。这个项目从另一个角度诠释了lambda的巨大威力，所以demo一出炉，就获得了一两千的github star。如今JAWS羽翼臻至丰满，推出了尚处在beta的jaws fraemwork v1版本：jaws-framework/JAWS · GitHub，并且在re:invent 2015上做了相当精彩的主题演讲（见github）。JAWS framework大量使用API gateway，cloudformation和lambda来提供serverless architecture，值得关注。


一个完整的serverless website可以这么考虑：


用户注册使用：API gateway，lambda，dynamodb，SES（发邮件）

用户登录使用：API gateway，lambda，或者（cognito和IAM，如果要集成第三方登录）

用户UGC各种内容：API gateway，lambda，dynamodb

其他REST API：API gateway + lambda

各种事件处理使用lambda

所有的静态资源放在S3上，使能static website hosting，然后通过javascript访问cognito或者REST API

日志存放在cloudWatch，并在需要的时候触发lambda

clickstream存在在kenisis，并触发lambda

如此这般，一个具备基本功能的serverless website就搭起来了。


如果你对JAWS感兴趣，可以尝试我生成的 https://github.com/tyrchen/jaws-test。


避免失控

lambda带来的部署上的解耦同时是把双刃剑。成千上万个功能各异的lambda函数（再加上各自不同的版本），很容易把系统推向失控的边缘。所以，最好通过以下手段来避免失控：


为lambda函数合理命名：使用一定规格的，定义良好的前缀（可类比ARN）

使用cloudformation处理资源的分配和部署（可以考虑JAWS）

可视化系统的实时数据流/事件流（类似下图）



（图片来自youtube视频截图：A Day in the Life of a Netflix Engineer，图片和本文关系不大，但思想类似）


由于基于lambda的诸多应用场景还处在刚刚起步的阶段，所以很多orchestration的事情还需要自己做，相信等lambda的使用日趋成熟时，就像docker生态圈一样，会产生众多的orchestration的工具，解决或者缓解系统失控的问题。

https://zhuanlan.zhihu.com/p/20297696
https://github.com/tyrchen/aws-lambda-thumbnail
https://github.com/serverless/serverless
https://github.com/ZeroSharp/serverless-php
https://github.com/serverless/examples
https://github.com/serverless/components
https://github.com/tyrchen/jaws-test
https://developers.weixin.qq.com/miniprogram/dev/api/
https://developers.weixin.qq.com/miniprogram/dev/wxcloud/guide/openapi/openapi.html#%E4%BA%91%E8%B0%83%E7%94%A8

https://developers.weixin.qq.com/miniprogram/dev/wxcloud/basis/capabilities.html#%E4%BA%91%E5%87%BD%E6%95%B0

