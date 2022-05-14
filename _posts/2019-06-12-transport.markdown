---
title: thrift transport
layout: post
category: web
author: 夏泽民
---
transport类体系架构与TProtocol类体系架构一样，所以这里就不重复叙述了，想了解可转去TProtocol类体系架构分析那篇。
<!-- more -->
　下面将对transport层的几种transport类进行介绍：

1、TSocket　阻塞型socket,　用于客户端，采用系统函数read和write进行读写数据；


2、TServerSocket　非阻塞型socket, 用于服务器端,　accecpt到的socket类型都是TSocket（即阻塞型socket）；


3、TBufferedTransport和TFramedTransport都是有缓存的，均继承TBufferBase，调用下一层TTransport类进行读写操作,结构极为相似。只是TFramedTransport以帧为传输单位，帧结构为：4个字节（int32_t）+传输字节串，头4个字节是存储后面字节串的长度，该字节串才是正确需要传输的数据，因此TFramedTransport每传一帧要比TBufferedTransport和TSocket多传4个字节；



4、TMemoryBuffer继承TBufferBase，用于程序内部通信用，不涉及任何网络I/O，可用于三种模式：（1）OBSERVE模式，不可写数据到缓存；（2）TAKE_OWNERSHIP模式，需负责释放缓存；（3）COPY模式，拷贝外面的内存块到TMemoryBuffer。

5、TFileTransport直接继承TTransport，用于写数据到文件。对事件的形式写数据，主线程负责将事件入列，写线程将事件入列，并将事件里的数据写入磁盘。这里面用到了两个队列，类型为TFileTransportBuffer，一个用于主线程写事件，另一个用于写线程读事件，这就避免了线程竞争，在读完队列事件后，就会进行队列交换，由于由两个指针指向这两个队列，交换只要交换指针即可。它还支持以chunk（块）的形式写数据到文件。


6、TFDTransport是非常简单地写数据到文件和从文件读数据，它的write和read函数都是直接调用系统函数write和read进行写和读文件。


7、TSimpleFileTransport直接继承TFDTransport，没有添加任何成员函数和成员变量，不同的是构造函数的参数和在TSimpleFileTransport构造函数里对父类进行了初始化（打开指定文件并将fd传给父类和设置父类的close_policy为CLOSE_ON_DESTROY）。


8、TZlibTransport跟TBufferedTransport和TFramedTransport一样，调用下一层TTransport类进行读写操作。它采用<zlib.h>提供的zlib压缩和解压缩库函数来进行压解缩，写时先压缩再调用底层TTransport类发送数据，读时先调用TTransport类接收数据再进行解压，最后供上层处理。

9、TSSLSocket继承TSocket，阻塞型socket,　用于客户端；采用openssl的接口进行读写数据。checkHandshake(）函数调用SSL_set_fd将fd和ssl绑定在一起，之后就可以通过ssl的SSL_read和SSL_write接口进行读写网络数据。

10、TSSLServerSocket继承TServerSocket，非阻塞型socket, 用于服务器端；accecpt到的socket类型都是TSSLSocket类型。

11、THttpClient和THttpServer是基于http1.1协议的继承Transport类型，均继承THttpTransport，其中
THttpClient用于客户端，THttpServer用于服务器端。两者都调用下一层TTransport类进行读写操作，均用到TMemoryBuffer作为读写缓存，只有调用flush（）函数才会将真正调用网络I/O接口发送数据。

TTransport是所有Transport类的父类，为上层提供了统一的接口而且通过TTransport即可访问各个子类不同实现，类似多态。
