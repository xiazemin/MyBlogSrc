---
title: BitTorrent协议分析
layout: post
category: web
author: 夏泽民
---
BitTorrent（简称BT）是一个文件分发协议，每个下载者在下载的同时不断向其他下载者上传已下载的数据。而在FTP、HTTP协议中，每个下载者从FTP或HTTP服务器处下载自己所需要的文件，各个下载者之间没有交互。当非常多的用户同时访问和下载服务器上的文件时，由于FTP服务器的处理能力和带宽的限制，下载速度会急剧下降，有的用户根本访问不了服务器。BT协议与FTP协议不同，它的特点是下载的人越多下载的速度越快，其原因在于每个下载者将已下载的数据提供给其他下载者下载，它充分利用了用户的上载带宽。BT协议通过一定的策略保证上传的速度越快，下载的速度也越快。

2.2  基于BT协议的文件分发系统的构成
基于BT协议的文件分发系统由以下几个实体构成。

（1）一个Web服务器。

（2）一个种子文件。

（3）一个Tracker服务器。

（4）一个原始文件提供者。

（5）一个网络浏览器。

（6）一个或多个下载者。

Web服务器上保存着种子文件，下载者使用网络浏览器（如IE浏览器）从Web服务器上下载种子文件。种子文件，又称为元原文件或metafile，它保存了共享文件的一些信息，如共享文件的文件名、文件大小、Tracker服务器的地址。种子文件通常很小，一般大小为1GB的共享文件，其种子文件不足100KB，种子文件以.torrent为后缀。Tracker服务器保存着当前下载某共享文件的所有下载者的IP和端口。原始文件提供者提供完整的共享文件供其他下载者下载，它也被称为种子，种子文件就是提供者使用BT客户端生成的。每个下载者通过运行BT客户端软件下载共享文件。我们把某个下载者本身称为客户端，把其他下载者称为peer。

BT客户端下载一个共享文件的过程是：客户端首先解析种子文件，获取待下载的共享文件的一些信息，其中包括Tracker服务器的地址。然后客户端连接Tracker获取当前下载该文件的所有下载者的IP和端口。之后客户端根据IP和端口连接其他下载者，从它们那里下载文件，同时把自己已下载的部分提供给其他下载者下载。

共享文件在逻辑上被划分为大小相同的块，称为piece，每个piece的大小通常为256KB。对于共享文件，文件的第1字节到第256K（即262144）字节为第一个piece，第256K＋1字节到第512K字节为第二个piece，依此类推。种子文件中包含有每个piece的hash值。BT协议规定使用Sha1算法对每个piece生成20字节的hash值，作为每个piece的指纹。每当客户端下载完一个piece时，即对该peice使用Sha1算法计算其hash值，并与种子文件中保存的该peice的hash值进行比较，如果一致即表明下载了一个完整而正确的piece。一旦某个piece被下载，该piece即提供给其他peer下载。在实际上传和下载中，每个piece又被划分为大小相同的slice，每个slice的大小固定为16KB（16384字节）。peer之间每次传输以slice为单位。

从以上描述可以得知，待开发的BT软件（即BT客户端）主要包含以下几个功能：解析种子文件获取待下载的文件的一些信息，连接Tracker获取peer的IP和端口，连接peer进行数据上传和下载、对要发布的提供共享文件制作和生成种子文件。种子文件和Tracker的返回信息都以一种简单而高效的编码方式进行编码，称为B编码。客户端与Tracker交换信息基于HTTP协议，Tracker本身作为一个Web服务器存在。客户端与其他peer采用面向连接的可靠传输协议TCP进行通信。下面将进一步作详细的介绍。

2.3  B编码
种子文件和Tracker的返回信息都是经过B编码的。要解析和处理种子文件以及Tracker的返回信息，首先要熟悉B编码的规则。B编码中有4种类型：字符串、整型、列表、字典。

字符串的编码格式为：<字符串的长度>:<字符串>，其中<>括号中的内容为必需。例如，有一个字符串spam，则经过B编码后为4:spam。

整型的编码格式为：i<十进制的整型数>e，即B编码中的整数以i作为起始符，以e作为终结符，i为integer的第一个字母，e为end的第一个字母。例如，整数3，经过B编码后为i3e，整数−3的B编码为i−3e，整数0的B编码为i0e。

注意i03e不是合法的B编码，因为03不是十进制整数，而是八进制整数。

列表的编码格式为：l<任何合法的类型>e，列表以l为起始符，以e为终结符，中间可以为任何合法的经过B编码的类型，l为list的第一个字母。例如，列表l4:spam4:eggse表示两个字符串，一个是spam，一个是eggs。

字典的编码格式为：d<关键字><值>e，字典以d为起始符，以e为终结符，关键字是一个经过B编码的字符串，值可以是任何合法的B编码类型，在d和e之间可以出现多个关键字和值对，d是dictionary的第一个字母。例如，d4:spaml3:aaa3:bbbee，它是一个字典，该字典的关键字是spam，值是一个列表（以l开始，以e结束），列表中有两个字符串aaa和bbb。

又如：d9:publisher3:bob17:publisher-webpage15:www.example.come，它也是一个字典，第一个关键字是publisher，对应的值为bob，第二个关键字是publisher-webpage，对应的值是www.example.com。

2.4  种子文件的结构
种子文件包含了提供共享的文件的一些信息，它以.torrent为后缀名，种子文件也被称为元信息文件或metafile，它是经过B编码的。种子文件事实上就是一个B编码的字典，它含有以下关键字如表13-1所示。

表13-1                                                       种子文件的关键字

关  键  字

含    义

info

该关键字对应的值是一个字典，它有两种模式，“singel file”和“multiple file”，文件模式和多文件模式。单文件模式是指待共享的文件只有一个，多文件模式是指提供共享的不止一个文件，而是两个或两个以上。如使用BT软件下载一部影片时，影片的上下部可能分别放在不同的文件里

announce

该关键字的值为Tracker的URL

announce－list

可选，它的值存放的是备用Tracker的URL

creation－date

可选，该关键字对应的值存放的是创建种子文件的时间

comment

可选，它的值存放的是种子文件制作者的备注信息，对于下载来说，该关键字基本没有用处，因此不必理会

created by

可选，该关键字对应的值存放的是生成种子文件的BT客户端软件的信息，如客户端名、版本号等，一般不必理会

 

info是最重要的一个关键字，它的值是一个字典，下面对它再作进一步的介绍。无论是单文件模式还是多文件模式，该字典都包含关键字如表13-2所示。

表13-2                                                       info包含的关键字

关  键  字

含    义

 

l       piece length

每个piece的长度，它的值是一个B编码的整型，该值通常为i262144e，即256K，也有可能为512K或128K

 

l       pieces

对应的值为一个字符串，它存放的是各个piece的hash值，这个字符串的长度一定是20的倍数，因为每个piece的hash值的长度为20字节

 

l       private

该值如果为1，则表明客户端必须通过连接Tracker来获取其他下载者，即peer的IP地址和端口号；如果为0，则表明客户端还可以通过其他方式来获取peer的IP地址和端口号，如DHT方式。DHT即分布式哈希表（Distribute Hash Tabel），它是一种以分布式的方式来获取peer的方法，现在许多BT客户端既支持通过连接Tracker来获取peer，也支持通过DHT来获取peer。如果种子文件中没有private这个关键字，则表明不限制一定要通过连接Tracker来获取peer

 	 	 
 

对于单文件模式的种子文件，info的值还含有的关键字如表13-3所示。

表13-3                                                 单模式种子文件的关键字

关  键  字

含    义

name

共享文件的文件名，也就是要下载的文件的文件名

length

共享文件的长度，以字节为单位

md5sum

可选，它是共享文件的md5值，这个值在BT协议中根本没有使用，所以不必理会

对于多文件模式的种子文件，info的值还含有的关键字如表13-4所示

表13-4                                             多文件模式种子文件的关键字

关  键  字

含    义

name

存放所有共享文件的文件夹名

files

它的值是一个列表，列表中含有多个字典，每个共享文件为一个字典。该字典中含有三个关键词

 

Files的每个共享文件为一个字典字典的关键词如表13-5所示

表13-5                                                        Files字的关键诃

关  键  词

含    义

length

共享文件的长度，以字节为单位

md5sum

可选，同上

path

存放的是共享文件的路径和文件名

 

建议读者到一些提供BT种子文件下载的网站，如bt.greedland.net、www.btchina.net，下载几个种子文件并在Windows操作系统下使用记事本打开进行分析，就可以清楚的了解上述概念。

2.5  与Tracker交互
完成解析种子文件并从中获取Tracker服务器的URL后，即可开始与Tracker进行交互。与Tracker进行交互主要有两个目的：一是将自己的下载进度告知给Tracker以便Tracker进行一些相关的统计；二是获取当前下载同一个共享文件的peer的IP地址和端口号。

客户端使用HTTP协议与Tracker进行通信。Tracker通过HTTP GET方法获取请求，请求的构成为Tracker的URL后面跟一个？以及参数和值对，如http://tk.greedland.net/ announce? param1= value1&param2= value2。

在客户端发往Tracker的GET请求中，通常包含参数如表13-6所示。

表13-6                                                        GET请求的参数

参    数

含    义

info_hash

与种子文件中info关键字对应的值，通过Sha1算法计算其hash值，该hash值就是info_hash参数对应的值，该hash值的长度固定为20字节

peer_id

每个客户端在下载文件前以随机的方式生成的20字节的标识符，用于标识自己，它的长度也是固定不变的

port

监听端口号，用于接收其他peer的连接请求

uploaded

当前总的上传量，以字节为单位

downloaded

当前总的下载量，以字节为单位

left

还剩余多少字节需要下载，以字节为单位

compact

该参数的值一般为1。

event

它的值为started、completed、stopped其中之一。客户端第一次与Tracker进行通信时，该值为started；下载完成时，该值为completed；客户端即将关闭时，该值为stopped

ip

可选，将客户端的IP地址告知给Tracker，Tracker可以通过分析客户端发给Tracker的IP数据包来获取客户端的IP地址，因此该参数是可选的，一般不用指明客户端的IP

numwant

可选，希望Tracker返回多少个peer的IP地址和端口号。如果该参数缺省，则默认返回50个peer的IP地址和端口号

key

可选，它的值为一个随机数，用于进一步标识客户端。因为已经由peer_id来标识客户端，因此该参数一般不使用

trackerid

可选，一般不使用

Tracker服务器的返回信息是一个经过B编码的字典。它含有关键字如表13-7所示。

表13-7                                           Tracker服务器返回信息关键字

关  键  字

含    义

failure reason

该关键字对应的值是一个可以读懂的字符串，指明GET请求失败的原因，如果返回信息中含有这个关键字，就不会再包含其他任何关键字

warnging message

该关键字对应的值是一个可以读懂的警告字符串

interval

指明客户端在下一次连接Tracker前所需等待的时间，以秒为单位

min interval

指明客户端在下一次连接Tracker前所需等待的最少时间，以秒为单位

tracker id

指明Tracker的ID

complete

一个整数，指明当前有多少个peer已经完成了整个共享文件的下载

incomplete

一个整数，指明当前有多少个peer还没有完成共享文件的下载

peers

返回各个peer的IP和端口号，它的值是一个字符串。首先是第一个peer的IP地址，然后是其端口号；接着是第二个peer的IP地址，然后是其端口号；依此类推

 

以下是一个发往Tracker服务器的HTTP GET请求的示例：

http://tk.greedland.net/announce?info_hash=01234567890123456789&

peer_id=01234567890123456789&port=3210&compact=1&uploaded=0&downloaded=0&left=8000000&event=started

以下是一个Tracker服务器回应的示例：

d8:completei100e10:incompletei200e8:intervali1800e5:peers300:......e

其中，“......”是一个长度为300的字符串，含有50个peer的IP地址和端口号。IP地址占4字节，端口号占2字节，即一个peer占6字节。

 


 

发往Tracker服务器的HTTP GET请求中，info_hash和peer_id可能含有非数字、非字母的字符，即含有除0～9、a～z、A～Z之外的字符，此时要对字符进行编码转换。例如，空格应该转换为 。否则Tracker无法正确处理GET请求。

 

2.6  peer之间的通信协议
peer之间的通信协议又称为peer wire protocal，即peer连线协议，它是一个基于TCP协议的应用层协议。

为了防止有的peer只下载不上传，BitTorrent协议建议，客户端只给那些向它提供最快下载速度的4个peer上传数据。简单地说就是谁向我提供下载，我也提供数据供它下载；谁不提供数据给我下载，我的数据也不会上传给它。客户端每隔一定时间，比如10秒，重新计算从各个peer处下载数据的速度，将下载速度最快的4个peer解除阻塞，允许这4个peer从客户端下载数据，同时将其他peer阻塞。

一个例外情况是，为了发现下载速度更快的peer，协议还建议，在任一时刻，客户端保持一个优化非阻塞peer，即无论该peer是否提供数据给客户端下载，客户端都允许该peer从客户端这里下载数据。由于客户端向peer上传数据，peer接着也允许客户端从peer处下载数据，并且下载速度超过4个非阻塞peer中的一个。客户端每隔一定的时间，如30秒，重新选择优化非阻塞peer。

当客户端与peer建立TCP连接后，客户端必须维持的几个状态变量如表13-8所示。

表13-8                                               客户端必须维持的状态变量

状 态 变 量

含    义

am_chocking

该值若为1，表明客户端将远程peer阻塞。此时如果peer发送数据请求给客户端，客户端将不会理会。也就是说，一旦将peer阻塞，peer就无法从客户端下载到数据；该值若为0，则刚好相反，即表明peer未被阻塞，允许peer从客户端下载数据

am_interested

该值若为1，表明客户端对远程的peer感兴趣。当peer拥有某个piece，而客户端没有，则客户端对peer感兴趣。该值若为0，则刚好相反，即表明客户端对peer不感兴趣，peer拥有的所有piece，客户端都拥有

peer_chocking

该值若为1，表明peer将客户端阻塞。此时，客户端无法从peer处下载到数据。该值若为0，表明客户端可以向peer发送数据请求，客户端将进行响应

peer_interested

该值若为1，表明peer对客户端感兴趣。也即客户端拥有某个piece，而peer没有。该值若为0，表明peer对客户端不感兴趣

 

当客户端与peer建立TCP连接后，客户端将这几个变量的值设置为。

am_chocking     = 1。

am_interested  = 0。

peer_chocking  = 1。

peer_interested = 0。

当客户端对peer感兴趣且peer未将客户端阻塞时，客户端可以从peer处下载数据。当peer对客户端感兴趣，且客户端未将peer阻塞时，客户端向peer上传数据。

除非另有说明，所有的整数型在本协议中被编码为4字节值（高位在前低位在后），包括在握手之后所有信息的长度前缀。
<!-- more -->
就HTTP､FTP､PUB等下载方式而言,一般都是首先将文件放到服务器上,然后再由服务器传送到每位用户的机器上。因此如果同一时刻下载的用户数量太多,势必影响到所有用户的下载速度,如果某些用户使用了多线程下载,那对带宽的影响就更严重了,因此几乎所有的下载服务器都有用户数量和最高下载速度等方面的限制｡

目的

此规范的目的是详细介绍BitTorrent协议规范 v1.0 ｡Bram的协议规范网站 http://www.bittorrent.com/protocol.html 简要地叙述了此协议,在部分范围缺少详细行为阐述｡希望此文档能成为 一个正式的规范,明确的条款,将来能作为讨论和执行的基础｡

此文档规定由BitTorrent开发者维持和使用｡欢迎大家为它做贡献,其中的内容代表当前协议,它仍由许多客户使用｡

这里不是提出特性请求的地方｡如果有请求,请见邮箱列表｡

应用范围

本文档适用于BitTorrent协议规范的第一版(v1.0)｡目前,这份文档应用于 torrent 文件结构､用户线路协议和服务器(Tracker)HTTP/HTTPS 协议规范｡如果某个协议有了新的修订,请到对应页面查看,而不在这里｡

约定

在本文档中,使用了许多约定来简明和明确地表达信息｡

用户(peer)v/s 客户端(client):在本文档中,一个用户可以是任何参与下载的BitTorrent客户端｡客户端也是一个用户,尽管BitTorrent客户端运行在本地机器上｡本规范的读者可能会认为自己是连接了许多用户的客户端｡

片断(piece)v/s 块(block):在本文档中,片断是指在元信息文件中描述的一部分已下载的数据,它可通过 SHA-1 hash 来校验｡而块是指客户端向用户请求的一部分数据｡两块或更多块组成一个完整的片断,它能被校验｡

实际标准:大的斜体字文本指出普通的准则在不同客户端BitTorrent协议的执行,它被当作为实际标准｡(对照英文原文,common应该翻译成通用或者常见,这句话的大概意思是一个规范由于被许多不同的BitTorrent客户端实现所通用,以至于被当做是实际标准)
BitTorrent 使用"分布式哈希表"(DHT)来为无 tracker 的种子(torrents)存储 peer 之间的联系信息。这样每个 peer 都成了 tracker。这个协议基于 Kademila[1] 网络并且在 UDP 上实现。

请注意本文档中使用的术语，以免混乱。

"peer" 是在一个 TCP 端口上监听的客户端/服务器，它实现了 BitTorrent 协议。
"节点" 是在一个 UDP 端口上监听的客户端/服务器，它实现了 DHT(分布式哈希表) 协议。
DHT 由节点组成，它存储了 peer 的位置。BitTorrent 客户端包含一个 DHT 节点，这个节点用来联系 DHT 中其他节点，从而得到 peer 的位置，进而通过 BitTorrent 协议下载。

概述 Overview
每个节点有一个全局唯一的标识符，作为 "node ID"。节点 ID 是一个随机选择的 160bit 空间，BitTorrent infohash[2] 也使用这样的 160bit 空间。 "距离"用来比较两个节点 ID 之间或者节点 ID 和 infohash 之间的"远近"。节点必须维护一个路由表，路由表中含有一部分其它节点的联系信息。其它节点距离自己越近时，路由表信息越详细。因此每个节点都知道 DHT 中离自己很"近"的节点的联系信息，而离自己非常远的 ID 的联系信息却知道的很少。

在 Kademlia 网络中，距离是通过异或(XOR)计算的，结果为无符号整数。distance(A, B) = |A xor B|，值越小表示越近。

当节点要为 torrent 寻找 peer 时，它将自己路由表中的节点 ID 和 torrent 的 infohash 进行"距离对比"。然后向路由表中离 infohash 最近的节点发送请求，问它们正在下载这个 torrent 的 peer 的联系信息。如果一个被联系的节点知道下载这个 torrent 的 peer 信息，那个 peer 的联系信息将被回复给当前节点。否则，那个被联系的节点则必须回复在它的路由表中离该 torrent 的 infohash 最近的节点的联系信息。最初的节点重复地请求比目标 infohash 更近的节点，直到不能再找到更近的节点为止。查询完了之后，客户端把自己作为一个 peer 插入到所有回复节点中离种子最近的那个节点中。

请求 peer 的返回值包含一个不透明的值，称之为"令牌(token)"。如果一个节点宣布它所控制的 peer 正在下载一个种子，它必须在回复请求节点的同时，附加上对方向我们发送的最近的"令牌(token)"。这样当一个节点试图"宣布"正在下载一个种子时，被请求的节点核对令牌和发出请求的节点的 IP 地址。这是为了防止恶意的主机登记其它主机的种子。由于令牌仅仅由请求节点返回给收到令牌的同一个节点，所以没有规定他的具体实现。但是令牌必须在一个规定的时间内被接受，超时后令牌则失效。在 BitTorrent 的实现中，token 是在 IP 地址后面连接一个 secret(通常是一个随机数)，这个 secret 每五分钟改变一次，其中 token 在十分钟以内是可接受的。

路由表 Routing Table
每个节点维护一个路由表保存已知的好节点。路由表中的节点是用来作为在 DHT 中请求的起始点。路由表中的节点是在不断的向其他节点请求过程中，对方节点回复的。

并不是我们在请求过程中收到得节点都是平等的，有的节点是好的，而另一些则不是。许多使用 DHT 协议的节点都可以发送请求并接收回复，但是不能主动回复其他节点的请求。节点的路由表只包含已知的好节点，这很重要。好节点是指在过去的 15 分钟以内，曾经对我们的某一个请求给出过回复的节点，或者曾经对我们的请求给出过一个回复(不用在15分钟以内)，并且在过去的 15 分钟给我们发送过请求。上述两种情况都可将节点视为好节点。在 15 分钟之后，对方没有上述 2 种情况发生，这个节点将变为可疑的。当节点不能给我们的一系列请求给出回复时，这个节点将变为坏的。相比那些未知状态的节点，已知的好节点会被给于更高的优先级。

路由表覆盖从 0 到 2^160 全部的节点 ID 空间。路由表又被划分为桶(bucket)，每个桶包含一部分的 ID 空间。空的路由表只有一个桶，它的 ID 范围从 min=0 到 max=2^160。当 ID 为 N 的节点插入到表中时，它将被放到 ID 范围在 min <= N < max 的 桶 中。空的路由表只有一个桶，所以所有的节点都将被放到这个桶中。每个桶最多只能保存 K 个节点，当前 K=8。当一个桶放满了好节点之后，将不再允许新的节点加入，除非我们自身的节点 ID 在这个桶的范围内。在这样的情况下，这个桶将被分裂为 2 个新的桶，每个新桶的范围都是原来旧桶的一半。原来旧桶中的节点将被重新分配到这两个新的桶中。如果一个新表只有一个桶，这个包含整个范围的桶将总被分裂为 2 个新的桶，每个桶的覆盖范围从 0..2^159 和 2^159..2^160。

当桶装满了好节点，新的节点会被丢弃。一旦桶中的某个节点变为了坏的节点，那么我们就用新的节点来替换这个坏的节点。如果桶中有在 15 分钟内都没有活跃过的节点，我们将这样的节点视为可疑的节点，这时我们向最久没有联系的节点发送 ping。如果被 ping 的节点给出了回复，那么我们向下一个可疑的节点发送 ping，不断这样循环下去，直到有某一个节点没有给出 ping 的回复，或者当前桶中的所有节点都是好的(也就是所有节点都不是可疑节点，他们在过去 15 分钟内都有活动)。如果桶中的某个节点没有对我们的 ping 给出回复，我们最好再试一次(再发送一次 ping，因为这个节点也许仍然是活跃的，但由于网络拥塞，所以发生了丢包现象，注意 DHT 的包都是 UDP 的)，而不是立即丢弃这个节点或者直接用新节点来替代它。这样，我们得路由表将充满稳定的长时间在线的节点。

每个桶都应该维持一个 lastchange 字段来表明桶中节点的"新鲜"度。当桶中的节点被 ping 并给出了回复，或者一个节点被加入到了桶，或者一个节点被新的节点所替代，桶的 lastchange 字段都应当被更新。如果一个桶的lastchange 在过去的 15 分钟内都没有变化，那么我们将更新它。这个更新桶操作是这样完成的：从这个桶所覆盖的范围中随机选择一个 ID，并对这个 ID 执行 find_nodes 查找操作。常常收到请求的节点通常不需要常常更新自己的桶，反之，不常常收到请求的节点常常需要周期性的执行更新所有桶的操作，这样才能保证当我们用到 DHT 的时候，里面有足够多的好的节点。

在插入第一个节点到路由表并启动服务后，这个节点应试着查找 DHT 中离自己更近的节点，这个查找工作是通过不断的发出 find_node 消息给越来越近的节点来完成的，当不能找到更近的节点时，这个扩散工作就结束了。路由表应当被启动工作和客户端软件保存（也就是启动的时候从客户端中读取路由表信息，结束的时候客户端软件记录到文件中）。

BitTorrent 协议扩展 BitTorrent Protocol Extension
BitTorrent 协议已经被扩展为可以在通过 tracker 得到的 peer 之间互相交换节点的 UDP 端口号(也就是告诉对方我们的 DHT 服务端口号)，在这样的方式下，客户端可以通过下载普通的种子文件来自动扩展 DHT 路由表。新安装的客户端第一次试着下载一个无 tracker 的种子时，它的路由表中将没有任何节点，这是它需要在 torrent 文件中找到联系信息。

peers 如果支持 DHT 协议就将 BitTorrent 协议握手消息的保留位的第 8 字节的最后一位置为 1。这时如果 peer 收到一个 handshake 表明对方支持 DHT 协议，就应该发送 PORT 消息。它由字节 0x09 开始，payload的长度是 2 个字节，包含了这个 peer 的 DHT 服务使用的网络字节序的 UDP 端口号。当 peer 收到这样的消息是应当向对方的 IP 和消息中指定的端口号的节点发送 ping。如果收到了 ping 的回复，那么应当使用上述的方法将新节点的联系信息加入到路由表中。

Torrent 文件扩展 Torrent File Extensions
一个无 tracker 的 torrent 文件字典不包含 announce 关键字，而使用 nodes 关键字来替代。这个关键字对应的内容应该设置为 torrent 创建者的路由表中 K 个最接近的节点。可供选择的，这个关键字也可以设置为一个已知的可用节点，比如这个 torrent 文件的创建者。请不要自动加入 router.bittorrent.com 到 torrent 文件中或者自动加入这个节点到客户端路由表中。

nodes = [["<host>", <port>], ["<host>", <port>], ...]
nodes = [["127.0.0.1", 6881], ["your.router.node", 4804]]
KRPC 协议 KRPC Protocol
KRPC 协议是由 bencode 编码组成的一个简单的 RPC 结构，他使用 UDP 报文发送。一个独立的请求包被发出去然后一个独立的包被回复。这个协议没有重发。它包含 3 种消息：请求，回复和错误。对DHT协议而言，这里有 4 种请求：ping，find_node，get_peers 和 announce_peer。

一条 KRPC 消息由一个独立的字典组成，其中有 2 个关键字是所有的消息都包含的，其余的附加关键字取决于消息类型。每条消息都包含 t 关键字，它是一个代表了 transaction ID 的字符串类型。transaction ID 由请求节点产生，并且回复中要包含回显该字段，所以回复可能对应一个节点的多个请求。transaction ID 应当被编码为一个短的二进制字符串，比如 2 个字节，这样就可以对应 2^16 个请求。另外每个 KRPC 消息还应该包含的关键字是 y，它由一个字节组成，表明这个消息的类型。y 对应的值有三种情况：q 表示请求，r 表示回复，e 表示错误。

联系信息编码 Contact Encoding
Peers 的联系信息被编码为 6 字节的字符串。又被称为 "CompactIP-address/port info"，其中前 4 个字节是网络字节序的 IP 地址，后 2 个字节是网络字节序的端口。

节点的联系信息被编码为 26 字节的字符串。又被称为 "Compactnode info"，其中前 20 字节是网络字节序的节点 ID，后面 6 个字节是 peers 的 "CompactIP-address/port info"。

请求 Queries
请求，对应于 KPRC 消息字典中的 y 关键字的值是 q，它包含 2 个附加的关键字 q 和 a。关键字 q是字符串类型，包含了请求的方法名字。关键字 a 一个字典类型包含了请求所附加的参数。

回复 Responses
回复，对应于 KPRC 消息字典中的 y 关键字的值是 r，包含了一个附加的关键字 r。关键字 r 是字典类型，包含了返回的值。发送回复消息是在正确解析了请求消息的基础上完成的。

错误 Errors
错误，对应于 KPRC 消息字典中的 y 关键字的值是 e，包含一个附加的关键字 e。关键字 e 是列表类型。第一个元素是数字类型，表明了错误码。第二个元素是字符串类型，表明了错误信息。当一个请求不能解析或出错时，错误包将被发送。下表描述了可能出现的错误码：

错误码 描述
201	一般错误
202	服务错误
203	协议错误，比如不规范的包，无效的参数，或者错误的 token
204	未知方法
错误包例子 Example Error Packets:

generic error = {"t":"aa", "y":"e", "e":[201, "A Generic Error Ocurred"]}
bencoded = d1:eli201e23:A Generic Error Ocurrede1:t2:aa1:y1:ee
DHT 请求 DHT Queries
所有的请求都包含一个关键字 id，它包含了请求节点的节点 ID。所有的回复也包含关键字id，它包含了回复节点的节点 ID。

ping
最基础的请求就是 ping。这时 KPRC 协议中的 "q" = "ping"。Ping 请求包含一个参数 id，它是一个 20 字节的字符串包含了发送者网络字节序的节点 ID。对应的 ping 回复也包含一个参数 id，包含了回复者的节点 ID。

参数: {"id" : "<querying nodes id>"}
回复: {"id" : "<queried nodes id>"}
报文包例子 Example Packets

ping Query = {"t":"aa", "y":"q", "q":"ping", "a":{"id":"abcdefghij0123456789"}}
bencoded = d1:ad2:id20:abcdefghij0123456789e1:q4:ping1:t2:aa1:y1:qe
Response = {"t":"aa", "y":"r", "r": {"id":"mnopqrstuvwxyz123456"}}
bencoded = d1:rd2:id20:mnopqrstuvwxyz123456e1:t2:aa1:y1:re
find_node
find_node 被用来查找给定 ID 的节点的联系信息。这时 KPRC 协议中的 "q" == "find_node"。find_node 请求包含 2 个参数，第一个参数是 id，包含了请求节点的ID。第二个参数是target，包含了请求者正在查找的节点的 ID。当一个节点接收到了 find_node 的请求，他应该给出对应的回复，回复中包含 2 个关键字 id 和 nodes，nodes 是字符串类型，包含了被请求节点的路由表中最接近目标节点的 K(8) 个最接近的节点的联系信息。

参数: {"id" : "<querying nodes id>", "target" : "<id of target node>"}
回复: {"id" : "<queried nodes id>", "nodes" : "<compact node info>"}
报文包例子 Example Packets

find_node Query = {"t":"aa", "y":"q", "q":"find_node", "a": {"id":"abcdefghij0123456789", "target":"mnopqrstuvwxyz123456"}}
bencoded =d1:ad2:id20:abcdefghij01234567896:target20:mnopqrstuvwxyz123456e1:q9:find_node1:t2:aa1:y1:qe
Response = {"t":"aa", "y":"r", "r": {"id":"0123456789abcdefghij", "nodes": "def456..."}}
bencoded = d1:rd2:id20:0123456789abcdefghij5:nodes9:def456...e1:t2:aa1:y1:re
get_peers
get_peers 与 torrent 文件的 infohash 有关。这时 KPRC 协议中的 "q" = "get_peers"。get_peers 请求包含 2 个参数。第一个参数是 id，包含了请求节点的 ID。第二个参数是 info_hash，它代表 torrent 文件的 infohash。如果被请求的节点有对应 info_hash 的 peers，他将返回一个关键字 values，这是一个列表类型的字符串。每一个字符串包含了 "CompactIP-address/portinfo" 格式的 peers 信息。如果被请求的节点没有这个 infohash 的 peers，那么他将返回关键字 nodes，这个关键字包含了被请求节点的路由表中离 info_hash最近的 K 个节点，使用 "Compactnodeinfo" 格式回复。在这两种情况下，关键字 token 都将被返回。token 关键字在今后的 annouce_peer 请求中必须要携带。token 是一个短的二进制字符串。

参数: {"id" : "<querying nodes id>", "info_hash" : "<20-byte infohash of target torrent>"}
回复: {"id" : "<queried nodes id>", "token" :"<opaque write token>", "values" : ["<peer 1 info string>", "<peer 2 info string>"]}
或: {"id" : "<queried nodes id>", "token" :"<opaque write token>", "nodes" : "<compact node info>"}
报文包例子 Example Packets:

get_peers Query = {"t":"aa", "y":"q", "q":"get_peers", "a": {"id":"abcdefghij0123456789", "info_hash":"mnopqrstuvwxyz123456"}}
bencoded =d1:ad2:id20:abcdefghij01234567899:info_hash20:mnopqrstuvwxyz123456e1:q9:get_peers1:t2:aa1:y1:qe
Response with peers = {"t":"aa", "y":"r", "r": {"id":"abcdefghij0123456789", "token":"aoeusnth", "values": ["axje.u", "idhtnm"]}}
bencoded =d1:rd2:id20:abcdefghij01234567895:token8:aoeusnth6:valuesl6:axje.u6:idhtnmee1:t2:aa1:y1:re
Response with closest nodes = {"t":"aa", "y":"r", "r": {"id":"abcdefghij0123456789", "token":"aoeusnth", "nodes": "def456..."}}
bencoded =d1:rd2:id20:abcdefghij01234567895:nodes9:def456...5:token8:aoeusnthe1:t2:aa1:y1:re
announce_peer
这个请求用来表明发出 announce_peer 请求的节点，正在某个端口下载 torrent 文件。announce_peer 包含 4 个参数。第一个参数是 id，包含了请求节点的 ID；第二个参数是 info_hash，包含了 torrent 文件的 infohash；第三个参数是 port 包含了整型的端口号，表明 peer 在哪个端口下载；第四个参数数是 token，这是在之前的 get_peers 请求中收到的回复中包含的。收到 announce_peer 请求的节点必须检查这个token 与之前我们回复给这个节点 get_peers 的 token 是否相同。如果相同，那么被请求的节点将记录发送 announce_peer 节点的 IP 和请求中包含的 port 端口号在 peer 联系信息中对应的 infohash 下。

参数: {"id" : "<querying nodes id>", "implied_port": <0 or 1>, "info_hash" : "<20-byte infohash of target torrent>", "port" : <port number>, "token" : "<opaque token>"}
回复: {"id" : "<queried nodes id>"}
报文包例子 Example Packets:

announce_peers Query = {"t":"aa", "y":"q", "q":"announce_peer", "a": {"id":"abcdefghij0123456789", "implied_port": 1, "info_hash":"mnopqrstuvwxyz123456", "port": 6881, "token": "aoeusnth"}}
bencoded = d1:ad2:id20:abcdefghij01234567899:info_hash20:<br /> mnopqrstuvwxyz1234564:porti6881e5:token8:aoeusnthe1:q13:announce_peer1:t2:aa1:y1:qe
Response = {"t":"aa", "y":"r", "r": {"id":"mnopqrstuvwxyz123456"}}
bencoded = d1:rd2:id20:mnopqrstuvwxyz123456e1:t2:aa1:y1:re

