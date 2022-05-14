---
title: FreeSwitch
layout: post
category: web
author: 夏泽民
---
Voice over Internet Protocol，缩写为VoIP）是一种语音通话技术，经由网际协议（IP）来达成语音通话与多媒体会议，也就是经由互联网来进行通信。其他非正式的名称有IP电话（IP telephony）、互联网电话（Internet telephony）、宽带电话（broadband telephony）以及宽带电话服务（broadband phone service）。

VoIP可用于包括VoIP电话、智能手机、个人计算机在内的诸多互联网接入设备，通过蜂窝网络、Wi-Fi进行通话及发送短信。
<!-- more -->
https://baike.baidu.com/item/voip/110300?fr=aladdin
https://blog.csdn.net/vevenlcf/article/details/50887119


http://www.freeswitch.org.cn/
FreeSWITCH是一个开源的电话软交换平台，主要开发语言是C，某些模块中使用了C++，以 MPL1.1发布。更多的说明请参考 什么是FreeSWITCH?和 FreeSWITCH新手指南。您也可以阅读这本 《FreeSWITCH权威指南》。

https://baike.baidu.com/item/freeswitch/1074133?fr=aladdin

https://freeswitch.com/

http://wiki.freeswitch.org.cn/

http://wiki.freeswitch.org.cn/wiki/New_User_s_guide.html

https://freeswitch.com/

https://github.com/freeswitch-cn/freeswitch-cn.github.com/commits/master

https://github.com/signalwire/freeswitch
1.FreeSwitch的概念 
FreeSwitch是一个开源的电环交换平台，是一个跨平台的/伸缩性极好的/免费的/多协议的电话软交换平台。 
1.1.FreeSwitch的特性 
FreeSwitch是跨平台的。他能原生地运行于Windows、Max OS X、Linux、BSD及Solaris等诸多32/64位平台。 
FreeSwitch具有很强的可伸缩性。FreeSwitch从一个简单的软电话客户端到运营商用级软交换设备几乎无所不能。 
FreeSwitch是免费的。 
FreeSwitch支持SIP、H323、Skype、Google Talk等多种通信协议，并能很容易的与各种开源的PBX系统通信，他也可以与商用的交换系统（如华为、中兴的交换机或思科、Avaya的交换机等）互通。 
FreeSwitch可以用作一个简单的交换引擎、一个PBX、一个媒体网关或媒体支持IVR的服务器，或在运营商的IMS网络中担当CSCF或Application Server等。 
FreeSwitch遵循相关RFC并支持很多高级的SIP特性，如Presence、BLF、SLA以及TCP、TLS和sRTP等，它也可以在用作一个SBC进行透明的SIP代理以支持其他媒体。 
FreeSwitch支持宽带及窄带语音编码，电话会议桥接可同时支持8、12、16、24、32及48kHz的语音。 
1.2.FreeSwitch的典型功能 
在线计费、预付费功能。 
电话路由服务器。 
语音转码服务器。 
支持资源优先权和QoS的服务器。 
多点会议服务器。 
IVR、语音通知服务器。 
VoiceMail服务器。 
PBX应用和软交换。 
应用层网关。 
防火墙/NAT穿越应用。 
私有服务器。 
第三方呼叫控制应用。 
业务生成环境运行时引擎。 
会话边界控制器。 
IMS中的S-CSCF/P-CSCF/I-CSCF。 
SIP网间互联网关。 
SBC及安全网关。 
传真服务器、T.30到T.38网关。 
2.在windows上安装FreeSwitch 
（1）使用安装包安装。 
（2）从源代码安装。

https://blog.csdn.net/acliyu/article/details/81556893
https://blog.csdn.net/yetyongjin/article/details/96838256
  实时通信（Real-time communication (RTC)）始于19世纪下半页电话的发明，几乎立刻演变成一个由大公司和基础设施组成的 全球互联网络。

        直到几年前，电话还是一座由大公司严格把守的带围墙的堡垒，几乎没有人能完全理解它是如何运作的。通过参加内部技术研讨会和内部学校，你有机会获得这些深奥的知识，甚至这些知识也仅限于你要工作的系统部分（中央办公室、最后一英里、PBX等）。承载和路由呼叫的基础设施以及应答和管理呼叫的应用程序都是临时的、僵化的、彼此不兼容的，并且需要大量的投资。

        两次革命彻底摧毁了旧世界。第一次是VoIP电话革命，它带来了一个开放协议（SIP），首先革新的是应用程序/PBX领域，然后是基础设施领域。原有体系中，有一个专用的PBX，它只能通过更换同一家公司生产的内部硬件卡来扩展，只提供自己的类型和型号的硬件电话。变革后，它被标准PC架构的标准服务器取代，使用现成的扩展卡，能够支持任何类型的标准兼容电话。紧接着，SIP从下到上爬到了大型电信基础设施的核心。今天，包括大型电信公司和运营商在内的所有电信基础设施都在运行SIP协议的某些版本。

        第二次革命，现在正在进行，还需要几年才能完成并取得成果，它就是WebRTC。这是一个误导性很强的名字，WebRTC根本不需要web页面和浏览器。WebRTC是一系列用于通信终端间加密互联的标准。只不过这些WebRTC标准正好首先在浏览器中实现而已。与此同时，它已成为物联网的主流通信方式，从智能手机应用到汽车、电梯、商店出纳和销售点。

        如今，我们有可能建立一个比传统的语音、视频和会议服务更好的通信系统，并以相对低廉的成本提供现高级的功能。Freeswitch的设计就是为了让所有这些变得更简单，我们将仔细研究它的架构，以便更好地理解它是如何工作的。

        如果你不能一下子掌握所有内容，不必担心。学习需要时间，尤其是RTC和VoIP。事实上，我们建议您反复阅读这一章的内容。第一次通读时尽可能多地吸收，然后在读完第6章XML拨号方案后再回头。你会惊讶于你对FreeSWITCH的理解有多大的提升。然后，在读过第10章（拨号方案、目录、通过XML_CURL和脚本实现一切）之后，回来第三次浏览它，您将牢固掌握FreeSWITCH的概念。给自己一些时间去消化这些概念，您将成为熟练的FreeSWITCH管理员。

        今天，我们生活在多种实时通信技术的生态系统中，它们同时共存：电话系统（如电话交换机和PBX）、传统的模拟电话（POTS线路或普通的老式电话服务）、运营商运营的传统电话网（PSTN或公共交换电话网）、移动电话（CDMA、GSM、LTE等）、传真、WebRTC、智能手机应用程序、VoIP电话和企业系统。

        FreeSWITCH就位于这些生态系统的中心：它连接并接受来自所有这些技术的连接，将它们桥接并混合在一起，它为用户提供交互应用和服务，无论他们的终端是什么类型的。FreeSwitch能够连接到外部数据服务、遗留的内部系统、计算机程序和业务程序等等。FreeSWITCH是跨平台的，可以运行在Linux、Windows、Mac OS X、BSD和Solaris上。我们有很多硬件可以选择，从大型服务器到Raspberry Pi都支持。因此，您可以在您的笔记本上开发，然后把它部署到数据中心或嵌入式设备上。第2章详细描述了FreeSWITCH的安装部署过程。

 

FreeSWITCH设计-模块化、可扩展和稳定性
 

        FreeSWITCH的设计目标是提供一个模块化的可扩展的系统，围绕一个稳定的交换内核，为开发人员添加和控制系统提供一个强大的接口。Freeswitch中的各种元素彼此独立，除了“Freeswitch API”中提供的内容外，对其他模块的工作方式没有太多的了解。Freeswitch的功能可以通过可加载模块进行扩展，这些模块将特定功能或外部技术绑定到核心中。

        在FreeSWITCH中，有许多不同类型的模块围绕着内核，就像一个机器人的大脑连接着许多传感器和接口一样。这些模块的类型列表如下：

 

模块类型

用途

Endpoint终端

电话协议，比如SIP，PSTN。

Application应用

执行某种任务，比如播放音频，发送数据。

Automated Speech Recognition ( ASR ) 自动语音识别

语音识别系统的处理接口

Chat聊天

桥接交换各种聊天协议。

Codec编解码

在音频编码格式间相互转换。

Dialplan拨号方案

解析话务详情，并决定话务路由。

Directory目录

连接内核查询与目录信息服务（比如LDAP）的API

Event handlers事件处理

允许外部程序控制FreeSWITCH。

File文件

提供一个用于提取和播放各种格式音频文件的接口。

Formats格式

播放各种格式的音频文件

Languages语言

用于控制话务的编程语言接口。

Loggers

日志控制接口，输出到控制台、系统日志或日志文件。

Say

将各种语言的音频文件串在一起，以提供反馈，以便说出电话号码、一天中的时间、单词拼写等内容。

Text-To-Speech (TTS)

TTS引擎的接口。

Timers

应用程序中的POSIX或Linux内核时钟。

XML Interfaces

XML接口，用于CDR、CURL、LDAP、RPC等。

 

 

        下图体现了FreeSWITCH的架构和FreeSWITCH的内核与外围模块的关系：




        通过组合各种功能模块，可以把FreeSwitch配置为IP电话、POTS线路、WebRTC和基于IP的电话服务。它还可以转换音频格式，自定义交互式语音菜单（IVR）系统。还可以从另一台计算机控制FreeSwitch服务器。接下来，我们介绍两个重要的模块。

 

重要模块—终端和拨号方案
        终端模块是非常重要 ，正因为它实现了一些关键功能，才使FreeSWITCH成为强大的平台。终端模块的主要功能是采用某些公共的通信技术，并将它们规范化为一个公共抽象实体，我们称之为会话(session)。一个会话描述了FreeSWITCH和一种特定协议间的连接。FreeSWITCH附带了几个终端模块，它们实现了几种不同的协议，如SIP、H.323、Jingle、Verto、WebRTC和其他一些协议。我们将花点时间来解释其中一个最为流行的模块，它的名字叫mod_sofia。

        Sofia-SIP (http://sofia-sip.sourceforge.net)是一个开源项目，最早由诺基亚开发，它提供会议初始化协议（SIP）的接口实现。Freeswitch开发人员对它进行了大量的开发和修复，进一步增强了其健壮性和特性。在FreeSWITCH中我们使用自己的定制版本库(/usr/src/freeswitch/libs/sofia-sip)，这是通过一个名为mod_sofia的模块实现的。这个模块注册FreeSwitch所需的所有必要钩子，以生成终端模块，它将FreeSwitch的原生结构转换为Sofia SIP结构，反之亦然。配置信息取自freeswitch配置文件，这使mod_sofia可以加载用户定义的首选项和连接详细信息。它允许FreeSwitch接受来自SIP电话和设备的注册，可以向其他SIP服务器（如服务提供商）注册，发送通知，并为话机设备提供服务（如亮灯和语音邮件服务）。

当FreeSWITCH与其它SIP设备建立一个SIP音/视频通话时，这个通话在FreeSWITCH里表现为一个活跃的会话。如果通话是呼入的，它可以被转移或桥接到IVR菜单、保持音乐（音频或视频的），或其它分机。或者，它可以被桥接到一个新的外呼话务，连接PSTN线，或WebRTC浏览器。让我们研究一个典型的场景，一个内部注册的分机2000，拨打2001分机，希望建立通话。

        首先，SIP话机通过网络向FreeSWITCH发送一条建立通话的消息（mod_sofia时刻监听这类消息）。收到消息后，mod_sofia反向解析相关细节，建立一个内核能够理解的抽象呼叫会话数据结构实例，并将它传递给FreeSWITCH的内核状态机。然后状态机迁移到ROUTING状态（意味着正在寻找目的地）。

        内核的下一步是根据呼叫终端的配置数据确定拨号方案模块。XML拨号方案模块是缺省的和使用最广泛的模块。这个模块的设计逻辑是从FreeSWITCH内存中的XML树查找一系列指令集。XML拨号方案模块会用正则表达式匹配来解析一系列的XML对象。

        当我们正尝试呼叫2001时，会在XML节点中寻找一个destination_number字段与2001相匹配的分机。此外，我们要记住拨号方案并没有限定只能匹配一个分机。一个呼叫在拨号方案中允许有多个匹配，在第6章，XML 拨号方案主题里，您将会得到分机这个术语的扩展定义。XML拨号方案模块会为呼叫建立一个TODO列表。每个匹配的分机将它的操作要求添加到呼叫的TODO列表中。

        假设FreeSWITCH至少找到一个条件匹配的分机，XML拨号方案模块将会向会话对象中插入指令，其中包含尝试呼叫2001所需要的信息（建立这个呼叫的TODO列表）。一旦指令就绪，呼叫会议的状态会从ROUTING迁移到EXECUTE，内核将钻取列表，并逐个执行ROUTING状态下所堆积的指令。这就是API发挥作用的地方。

        每条指令都以一个命名application形式添加到会话中，可以向application传递一个data参数数据。在本例中，我们将会用到bridge app。这个app的作用是创建一个新的会话，建立外呼连接，然后把两个会话连接在一起进行音频数据交换。我们将给bridge提供一个参数：user/2001，这是生成内部注册分机呼叫电话对象的最简单方式。针对2001的拨号方案条目，看起来是这样的：

<extension name="example">
 
<condition field="destination_number"
 
expression="^2001$">
 
<action application="bridge" data="user/2001"/>
 
</condition>
 
</extension>
 

        这个分机被命名为example，它有一个单一的匹配条件。如果条件匹配成功，它需要执行一条单一的app。可以理解为：如果主叫方拨打2001，那么将在主叫方与终端间建立一条连接

        一旦指令在会话的TODO列表中插入完毕，会话的状态立刻迁移到EXECUTE，FreeSWITCH内核将开始使用所搜集到的数据执行所需的操作。它首先解析列表，发现它必须对user/2001执行bridge；接下来，查找bridge app，并将参数user/2001传递给它。这将导致FreeSWITCH内核按所需类型创建一个新的外呼会话（也就是一路呼叫）。我们假定2001分机的用户用一个SIP电话注册到FreeSWITCH ，那么user/2001将被解析为一个SIP拨号串，这个拨号串将被传递给mod_sofia，请求建立一个外呼会话（指向2001的SIP话机）。

        如果新的会话建立成功，那么FreeSWITCH内核中将会存在两个会话：新建会话和原始会话（主叫侧建立的会话）。Bridge app将会接管这两个会话，将对它们调用bridge函数。在被叫接听后，bridge调用将使音频和/或视频流在双向间流动。如果用户不能应答或正忙，会触发超时（也就是呼叫失败），那么会给主叫话机回传一条失败消息。如果一个电话无人应答或分机占线，拨号方案中可以有许多种不同的响应，其中包括呼叫转接或语音信箱。

        FreeSWITCH继承了SIP的所有复杂性，并将其简化为一个公共的(内部的)接口。然后通过允许我们使用拨号方案中的一条指令将2000分机的电话连接到2001分机的电话，进一步降低了复杂性。如果我们还想让2001分机能够呼叫2000分机，我们可以在拨号方案中添加另一个条目:

<extension name="example 2">
 
<condition field="destination_number" expression="^2000$">
 
<action application="bridge" data="user/2000"/>
 
</condition>
 
</extension>
        在这个场景中，终端模块（mod_sofia）把呼入的SIP呼叫转换为FreeSWITCH会话，而拨号方案模块（mod_dialplan_xml）把XML转换为分机。然后，bridge app（由mod_dptools模块提供）把复杂的外呼和连接媒体流的过程，变成了一个简单的应用/数据对。拨号方案模块和app模块接口都是围绕FreeSWITCH会话设计的。抽象不仅让用户层的工作变得更容易，而且还简化了app和拨号方案的设计，因为它可以使呼叫过程中涉及的终端技术细节透明化。正是由于这种抽象技术，当明天我们为三维全息通信网络编写一个新的终端模块时，我们将能够重用所有的app模块和拨号方案模块。对于FreeSWITCH内核和其它所有FreeSWITCH模块来说，仅是多了一种全息三维呼叫方式，它也仅仅是一个标准的会话抽象而已。您可能希望使用终端的协议提供的一些特定数据。比如说在SIP协议中，包含了几个特定的头域或其它一些有趣的数据，您可能希望获取这些特别的信息。我们通过向通道（channel，会话结构中与终端交互的部分）添加变量来解决这个问题。使用通道变量，mod_sofia可以创建在SIP数据中遇到的任意值，您可以在拨号方案或app中通过变量名检索这些值。这些最初是底层SIP消息的特定（任意）部分，然而，FreeSWITCH内核只是将它们视为内核可以忽略的普通通道变量。此外，还有一些特殊的保留通道变量，它们可以以许多有用的方式影响FreeSWITCH的行为。如果您曾经使用过某种脚本语言或配置引擎，而且它们有变量的概念（有时叫属性值或AVP），那么通道变量的概念和它们几乎是一样的。它们仅是一个简单的名字与值的配对，可以传递给通道，可以设置数据值。有一个用于设置通道变量的app叫set，在拨号方案中，我们可以用它设置自己的通道变量：

<extension name="example 3">
 
<condition field="destination_number" expression="^2000$">
 
<action application="set" data="foo=bar"/>
 
<action application="bridge" data="user/2000"/>
 
</condition>
 
</extension>
 

        这个例子与之前的几乎相同，只是在外呼之前，我们将通道变量foo的值设置为bar。这个变量会一直保留到呼叫结束，甚至可以在呼叫结束后的CDR中引用它。

        用小块构建的东西越多，可以重用的相同底层资源就越多，从而使整个系统的使用变得更简单。例如：codec接口，除了它自己世界中的音视频包编解码之外，对内核的其它内容一无所知。一旦写了某个特定的codec模块，任何可以携带相应媒体流的终端接口都可以直接使用它。这意味着，如果我们有一个文语转换（TTS）模块可以工作，那么我们就可以在所有的FreeSWITCH支持的终端上生成合成的语音。如果系统支持更多的codec，TTS也会变得更有用；反过来，如果我们添加一可以使用某个codec的功能，也会使对应的codec模块变得更有价值。如果我们写了一个新的app模块，现有的终端立刻就能运行和使用这个app。
        
        https://blog.csdn.net/yetyongjin/article/details/96838256
        
 https://www.oschina.net/p/freeswitch?hmsr=aladdin1e1
 
 https://www.zhihu.com/topic/19892529/hot
 
 https://book.dujinfang.com/
 
 https://github.com/signalwire/freeswitch
 
 https://www.cnblogs.com/Braveliu/p/10943511.html
 
 https://www.cnblogs.com/Braveliu/p/10943511.html
 
  软交换技术是NGN网络的核心技术，为下一代网络(NGN)具有实时性要求的业务提供呼叫控制和连接控制功能。软交换技术独立于传送网络，主要完成呼叫控制、资源分配、协议处理、路由、认证、计费等主要功能，同时可以向用户提供现有电路交换机所能提供的所有业务，并向第三方提供可编程能力。
  
  https://baike.baidu.com/item/%E8%BD%AF%E4%BA%A4%E6%8D%A2/1180831?fr=aladdin
  
  SIP（Session Initiation Protocol，会话初始协议）是由IETF（Internet Engineering Task Force，因特网工程任务组）制定的多媒体通信协议。
它是一个基于文本的应用层控制协议，用于创建、修改和释放一个或多个参与者的会话。SIP 是一种源于互联网的IP 语音会话控制协议，具有灵活、易于实现、便于扩展等特点。

https://baike.baidu.com/item/SIP/33921?fr=aladdin

https://www.cs.columbia.edu/sip/

SIP中继或者VOIP网关、GOIP网关

GOIP的网关，这是一种能插SIM卡的将GSM通话转换成VOIP协议的设备

https://blog.csdn.net/weixin_44081297/article/details/88388578

http://www.yunzongji.cn/news/hujiaozhongxin-254.html

https://blog.csdn.net/hnzwx888/article/details/83414609


FreeSWITCH拨号计划dialplan详解
https://blog.csdn.net/dujiajiyu/article/details/93676154
https://blog.csdn.net/smileyan9/article/details/88542307

https://blog.csdn.net/mao834099514/article/details/73174229

https://www.jianshu.com/p/8dd30f06c974


https://blog.csdn.net/u012377333/article/details/44567615?locationNum=6

命令行的形式来控制拨打电话

https://blog.csdn.net/cyq129445/article/details/81865823

https://blog.csdn.net/hnzwx888/article/details/83188087

https://blog.csdn.net/weixin_34351321/article/details/93809948

使用FreeSWITCH做电话自动回访设置
https://blog.csdn.net/irizhao/article/details/88635237

http://www.360doc.com/content/13/1213/15/12747488_336873552.shtml

https://blog.csdn.net/dittychen/article/details/79280046

ESL（Event Socket Library

https://blog.csdn.net/huoyin/article/details/39394189?utm_source=blogxgwz9

https://blog.csdn.net/mao834099514/article/details/73565205

http://www.freeswitch.net.cn/174.html

https://blog.csdn.net/huoyin/article/details/39394189?locationNum=8

https://www.cnblogs.com/MikeZhang/archive/2016/09/27/freeswitch_mod_vent_socket_20160926.html

https://github.com/0x19/goesl