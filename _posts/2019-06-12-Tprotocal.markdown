---
title: TProtocol 协议和编解码
layout: post
category: web
author: 夏泽民
---
客户端和服务器通过约定的协议来传输消息(数据)，通过特定的格式来编解码字节流，并转化成业务消息，提供给上层框架调用。

Thrift的协议比较简单，它把协议和编解码整合在了一起。抽象类TProtocol定义了协议和编解码的顶层接口。个人感觉采用抽象类而不是接口的方式来定义顶层接口并不好，TProtocol关联了一个TTransport传输对象，而不是提供一个类似getTransport()的接口，导致抽象类的扩展性比接口差。

TProtocol主要做了两个事情:

1. 关联TTransport对象

2.定义一系列读写消息的编解码接口，包括两类，一类是复杂数据结构比如readMessageBegin, readMessageEnd,  writeMessageBegin, writMessageEnd.还有一类是基本数据结构，比如readI32, writeI32, readString, writeString
<!-- more -->
所谓协议就是客户端和服务器端约定传输什么数据，如何解析传输的数据。对于一个RPC调用的协议来说，要传输的数据主要有:
调用方

1. 方法的名称，包括类的名称和方法的名称

2. 方法的参数，包括类型和参数值

3.一些附加的数据，比如附件，超时事件，自定义的控制信息等等

返回方

1. 调用的返回码

2. 返回值

3.异常信息



从TProtocol的定义我们可以看出Thrift的协议约定如下事情:

1. 先writeMessageBegin表示开始传输消息了，写消息头。Message里面定义了方法名，调用的类型，版本号，消息seqId

2. 接下来是写方法的参数，实际就是写消息体。如果参数是一个类，就writeStructBegin

3. 接下来写字段，writeFieldBegin, 这个方法会写接下来的字段的数据类型和顺序号。这个顺序号是Thrfit对要传输的字段的一个编码，从１开始

4. 如果是一个集合就writeListBegin/writeMapBegin，如果是一个基本数据类型，比如int, 就直接writeI32

5. 每个复杂数据类型写完都调用writeXXXEnd，直到writeMessageEnd结束

6. 读消息时根据数据类型读取相应的长度



每个writeXXX都是采用消息头+消息体的方式。我们来看TBinaryProtocol的实现。

1. writeMessgeBegin方法写了消息头，包括4字节的版本号和类型信息，字符串类型的方法名，４字节的序列号seqId

2. writeFieldBegin，写了１个字节的字段数据类型，和2个字节字段的顺序号

3. writeI32，写了４个字节的字节数组

4. writeString,先写４字节消息头表示字符串长度，再写字符串字节

5. writeBinary,先写４字节消息头表示字节数组长度，再写字节数组内容

6.readMessageBegin时，先读４字节版本和类型信息，再读字符串，再读４字节序列号

7.readFieldBegin，先读1个字节的字段数据类型，再读2个字节的字段顺序号

8. readString时，先读４字节字符串长度，再读字符串内容。字符串统一采用UTF-8编码
TProtocol定义了基本的协议信息，包括传输什么数据，如何解析传输的数据的基本方法。

还存在一个问题，就是服务器端如何知道客户端发送过来的数据是怎么组合的，比如第一个字段是字符串类型，第二个字段是int。这个信息是在IDL生成客户端时生成的代码时提供了。Thrift生成的客户端代码提供了读写参数的方法，这两个方式是一一对应的，包括字段的序号，类型等等。客户端使用写参数的方法，服务器端使用读参数的方法。

thrift可分为以下几个组件，其中传输层被细分为低级传输层和复写传输层。    代码生成器：根据thrift idl文件生成各个语言代码，位于compiler目录内。    低级传输层：靠近网络层、作为rpc框架接收报文的入口，提供各种底层实现如socket创建、读写、接收连接等。    复写传输层：基于低级传输层，实现各种复写传输层包括http、framed、buffered、压缩传输层等，复写传输层可以被协议层直接使用，用户也可以通过重写低级传输层和复写传输层实现自己的传输层。    协议层：协议层主要负责解析请求、应答报文为具体的结构体、类实例，供处理层直接使用，目前的协议包括Binary(最为常用)、json、多路混合协议等。    处理层：由代码生成器生成，根据获取到的具体信息如method name，进行具体的接口处理，处理层构造函数的入口包含一个handler，handler由业务方进行具体的实现，然后在处理层内被调用，并应答处理结果。    服务层：融合低级传输层、复写传输层、协议层、处理层，自身包含各种不同类型的服务模型，如非阻塞单进程服务、one request per fork、one request per thread、thread pool等模型。



