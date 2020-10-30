---
title: thrift_golang
layout: post
category: golang
author: 夏泽民
---
Thrift RPC 框架指南
认识Thrift框架
thrift是一个软件框架，用来进行可扩展且跨语言的服务的开发。它结合了功能强大的软件堆栈和代码生成引擎，以构建在 C++, Java, Python, PHP, Ruby, Erlang, Perl, Haskell, C#, Cocoa, JavaScript, Node.js, Smalltalk, and OCaml 这些编程语言间无缝结合的、高效的服务。

thrift最初由facebook开发，07年四月开放源码，08年5月进入apache孵化器。
thrift允许定义一个简单的定义文件中的数据类型和服务接口，以作为输入文件，编译器生成代码用来方便地生成RPC客户端和服务器通信的无缝跨编程语言。
类似Thrift的工具，还有Avro、protocol buffer,但相对于Thrift来讲，都没有Thrift支持全面和使用广泛。
Thrift自下到上可以分为4层
Server(single-threaded, event-driven etc)
服务器进程调度
Processor(compiler generated)
RPC接口处理函数分发，IDL定义接口的实现将挂接到这里面
Protocol (JSON, compact etc)
协议
Transport(raw TCP, HTTP etc)
网络传输
Thrift实际上是实现了C/S模式，通过代码生成工具将接口定义文件生成服务器端和客户端代码（可以为不同语言），从而实现服务端和客户端跨语言的支持。用户在Thirft描述文件中声明自己的服务，这些服务经过编译后会生成相应语言的代码文件，然后用户实现服务（客户端调用服务，服务器端提服务）便可以了。其中protocol（协议层, 定义数据传输格式，可以为二进制或者XML等）和transport（传输层，定义数据传输方式，可以为TCP/IP传输，内存共享或者文件共享等）被用作运行时库。

Thrift支持的传输及服务模型
支持的传输格式：
参数	描述
TBinaryProtocol	二进制格式
TCompactProtocol	压缩格式
TJSONProtocol	JSON格式
TSimpleJSONProtocol	提供JSON只写协议, 生成的文件很容易通过脚本语言解析。
TDebugProtocol	使用易懂的可读的文本格式，以便于debug
支持的数据传输方式：
参数	描述
TSocket	阻塞式socker
TFramedTransport	以frame为单位进行传输，非阻塞式服务中使用。
TFileTransport	以文件形式进行传输。
TMemoryTransport	将内存用于I/O. java实现时内部实际使用了简单的ByteArrayOutputStream。
TZlibTransport	使用zlib进行压缩， 与其他传输方式联合使用。当前无java实现。
支持的服务模型：
参数	描述
TSimpleServer	简单的单线程服务模型，常用于测试
TThreadPoolServer	多线程服务模型，使用标准的阻塞式IO。
TNonblockingServer	多线程服务模型，使用非阻塞式IO（需使用TFramedTransport数据传输方式）
Thrift 下载及安装
如何获取Thrift
官网：http://thrift.apache.org/
golang的Thrift包:
go get git.apache.org/thrift.git/lib/go/thrift
如何安装Thrift
mac下安装Thrift,参考上一篇介绍 
其他平台安装自行挖掘,呵呵。 
安装后通过

liuxinmingMacBook-Rro#:thrift -version
Thrift version 0.9.2 #看到这一行表示安装成功
Golang、PHP通过Thrift调用
先发个官方各种语言DEMO地址 https://git1-us-west.apache.org/repos/asf?p=thrift.git;a=tree;f=tutorial;h=d69498f9f249afaefd9e6257b338515c0ea06390;hb=HEAD

Thrift的协议库IDL文件
语法参考
参考资料
http://www.cnblogs.com/tianhuilove/archive/2011/09/05/2167669.html

http://my.oschina.net/helight/blog/195015
基本类型

bool: 布尔值 (true or false), one byte
byte: 有符号字节
i16: 16位有符号整型
i32: 32位有符号整型
i64: 64位有符号整型
double: 64位浮点型
string: Encoding agnostic text or binary string 
基本类型中基本都是有符号数，因为有些语言没有无符号数，所以Thrift不支持无符号整型。
特殊类型

binary: Blob (byte array) a sequence of unencoded bytes 
这是string类型的一种变形，主要是为java使用
struct结构体
thrift中struct是定义为一种对象，和面向对象语言的class差不多.,但是struct有以下一些约束： 
struct不能继承，但是可以嵌套，不能嵌套自己。 
1. 其成员都是有明确类型 
2. 成员是被正整数编号过的，其中的编号使不能重复的，这个是为了在传输过程中编码使用。 
3. 成员分割符可以是逗号（,）或是分号（;），而且可以混用，但是为了清晰期间，建议在定义中只使用一种，比如C++学习者可以就使用分号（;）。 
4. 字段会有optional和required之分和protobuf一样，但是如果不指定则为无类型–可以不填充该值，但是在序列化传输的时候也会序列化进去， 
optional是不填充则部序列化。 
required是必须填充也必须序列化。 
5. 每个字段可以设置默认值 
6. 同一文件可以定义多个struct，也可以定义在不同的文件，进行include引入。

struct Work {
  1: i32 num1 = 0,
  2: i32 num2,
  3: Operation op,
  4: optional string comment,
}
容器（Containers）
Thrift3种可用容器类型：


list(t): 元素类型为t的有序表，容许元素重复。
set(t):元素类型为t的无序表，不容许元素重复。对应c++中的set，java中的HashSet,python中的set，php中没有set，则转换为list类型。
map(t,t): 键类型为t，值类型为t的kv对，键不容许重复。对用c++中的map, Java的HashMap, PHP 对应 array, Python/Ruby 的dictionary。 
容器中元素类型可以是除了service外的任何合法Thrift类型（包括结构体和异常）。为了最大的兼容性，map的key最好是thrift的基本类型，有些语言不支持复杂类型的key，JSON协议只支持那些基本类型的key。 
容器都是同构容器，不失异构容器。
实现Thrift TDL文件
batu.thrift文件：

/**
 * BatuThrift TDL
 * @author liuxinming
 * @time 2015.5.13
 */

namespace go batu.demo
namespace php batu.demo

/**
 * 结构体定义
 */
struct Article{
 1: i32 id, 
 2: string title,
 3: string content,
 4: string author,
}

const map<string,string> MAPCONSTANT = {'hello':'world', 'goodnight':'moon'}

service batuThrift {        
        list<string> CallBack(1:i64 callTime, 2:string name, 3:map<string, string> paramMap),
        void put(1: Article newArticle),
}
编译IDL文件，生成相关代码
thrift -r --gen go batu.thrift
thrift -r --gen php batu.thrift  
thrift -r --gen php:server batu.thrift #生成PHP服务端接口代码有所不一样
Golang Service 实现
先按照golang的Thrift包

go get git.apache.org/thrift.git/lib/go/thrift

将Thrift生成的开发库复制到GOPATH中

cp -r /Users/liuxinming/wwwroot/testphp/gen-go/batu $GOPATH/src

开发Go server端代码（后面的代码，目录我们放在$GOPATH/src/thrift 中运行和演示） 
test.go文件:

package main

import (
    "batu/demo" #注意导入Thrift生成的接口包
    "fmt"
    "git.apache.org/thrift.git/lib/go/thrift"
    "os"
    "time"
)

const (
    NetworkAddr = "127.0.0.1:9090" #监听地址&端口
)

type batuThrift struct {
}

func (this *batuThrift) CallBack(callTime int64, name string, paramMap map[string]string) (r []string, err error) {
    fmt.Println("-->from client Call:", time.Unix(callTime, 0).Format("2006-01-02 15:04:05"), name, paramMap)
    r = append(r, "key:"+paramMap["a"]+"    value:"+paramMap["b"])
    return
}

func (this *batuThrift) Put(s *demo.Article) (err error) {
    fmt.Printf("Article--->id: %d\tTitle:%s\tContent:%t\tAuthor:%d\n", s.Id, s.Title, s.Content, s.Author)
    return nil
}

func main() {
    transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
    protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
    //protocolFactory := thrift.NewTCompactProtocolFactory()

    serverTransport, err := thrift.NewTServerSocket(NetworkAddr)
    if err != nil {
        fmt.Println("Error!", err)
        os.Exit(1)
    }

    handler := &batuThrift{}
    processor := demo.NewBatuThriftProcessor(handler)

    server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
    fmt.Println("thrift server in", NetworkAddr)
    server.Serve()
}
4.运行go服务端(监听9090端口)

liuxinmingdeMacBook-Pro:thrift liuxinming$ go run test.go 
thrift server in 127.0.0.1:9090

至此Go的Thrift服务端OK.

Golang Client 实现
goClient.go文件：

package main

import (
    "batu/demo"
    "fmt"
    "git.apache.org/thrift.git/lib/go/thrift"
    "net"
    "os"
    "strconv"
    "time"
)

const (
    HOST = "127.0.0.1"
    PORT = "9090"
)

func main() {
    startTime := currentTimeMillis()

    transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
    protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

    transport, err := thrift.NewTSocket(net.JoinHostPort(HOST, PORT))
    if err != nil {
        fmt.Fprintln(os.Stderr, "error resolving address:", err)
        os.Exit(1)
    }

    useTransport := transportFactory.GetTransport(transport)
    client := demo.NewBatuThriftClientFactory(useTransport, protocolFactory)
    if err := transport.Open(); err != nil {
        fmt.Fprintln(os.Stderr, "Error opening socket to "+HOST+":"+PORT, " ", err)
        os.Exit(1)
    }
    defer transport.Close()

    for i := 0; i < 10; i++ {
        paramMap := make(map[string]string)
        paramMap["a"] = "batu.demo"
        paramMap["b"] = "test" + strconv.Itoa(i+1)
        r1, _ := client.CallBack(time.Now().Unix(), "go client", paramMap)
        fmt.Println("GOClient Call->", r1)
    }

    model := demo.Article{1, "Go第一篇文章", "我在这里", "liuxinming"}
    client.Put(&model)
    endTime := currentTimeMillis()
    fmt.Printf("本次调用用时:%d-%d=%d毫秒\n", endTime, startTime, (endTime - startTime))

}

func currentTimeMillis() int64 {
    return time.Now().UnixNano() / 1000000
}
goClient运行后结果：

liuxinmingdeMacBook-Pro:thrift liuxinming$ go run goClient.go 
GOClient Call-> [key:batu.demo value:test1] 
GOClient Call-> [key:batu.demo value:test2] 
GOClient Call-> [key:batu.demo value:test3] 
GOClient Call-> [key:batu.demo value:test4] 
GOClient Call-> [key:batu.demo value:test5] 
GOClient Call-> [key:batu.demo value:test6] 
GOClient Call-> [key:batu.demo value:test7] 
GOClient Call-> [key:batu.demo value:test8] 
GOClient Call-> [key:batu.demo value:test9] 
GOClient Call-> [key:batu.demo value:test10] 
本次调用用时:1431583140857-1431583140855=2毫秒

PHP Client 实现
首先去下载Thrift，git库地址为：https://github.com/apache/thrift
新建项目目录testphp，然后把thrift/lib/php/lib复制到testphp目录下面
复制生成的gen-php到testphp目录下面
客户端代码
<?php
/**
 * Thrift RPC - PHPClient
 * @author liuxinming
 * @time 2015.5.13
 */
namespace  batu\testDemo;
header("Content-type: text/html; charset=utf-8");
$startTime = getMillisecond();//记录开始时间

$ROOT_DIR = realpath(dirname(__FILE__).'/');
$GEN_DIR = realpath(dirname(__FILE__).'/').'/gen-php';
require_once $ROOT_DIR . '/Thrift/ClassLoader/ThriftClassLoader.php';

use Thrift\ClassLoader\ThriftClassLoader;
use Thrift\Protocol\TBinaryProtocol;
use Thrift\Transport\TSocket;
use Thrift\Transport\TSocketPool;
use Thrift\Transport\TFramedTransport;
use Thrift\Transport\TBufferedTransport;

$loader = new ThriftClassLoader();
$loader->registerNamespace('Thrift',$ROOT_DIR);
$loader->registerDefinition('batu\demo', $GEN_DIR);
$loader->register();

$thriftHost = '127.0.0.1'; //UserServer接口服务器IP
$thriftPort = 9090;            //UserServer端口

$socket = new TSocket($thriftHost,$thriftPort);  
$socket->setSendTimeout(10000);#Sets the send timeout.
$socket->setRecvTimeout(20000);#Sets the receive timeout.
//$transport = new TBufferedTransport($socket); #传输方式：这个要和服务器使用的一致 [go提供后端服务,迭代10000次2.6 ~ 3s完成]
$transport = new TFramedTransport($socket); #传输方式：这个要和服务器使用的一致[go提供后端服务,迭代10000次1.9 ~ 2.1s完成，比TBuffer快了点]
$protocol = new TBinaryProtocol($transport);  #传输格式：二进制格式
$client = new \batu\demo\batuThriftClient($protocol);# 构造客户端

$transport->open();  
$socket->setDebug(TRUE);

for($i=1;$i<11;$i++){
    $item = array();
    $item["a"] = "batu.demo";
    $item["b"] = "test"+$i;
    $result = $client->CallBack(time(),"php client",$item); # 对服务器发起rpc调用
    echo "PHPClient Call->".implode('',$result)."<br>";
}

$s = new \batu\demo\Article();
$s->id = 1;
$s->title = '插入一篇测试文章';
$s->content = '我就是这篇文章内容';
$s->author = 'liuxinming';
$client->put($s);

$s->id = 2;
$s->title = '插入二篇测试文章';
$s->content = '我就是这篇文章内容';
$s->author = 'liuxinming';
$client->put($s);

$endTime = getMillisecond();

echo "本次调用用时: :".$endTime."-".$startTime."=".($endTime-$startTime)."毫秒<br>";

function getMillisecond() {
    list($t1, $t2) = explode(' ', microtime());
    return (float)sprintf('%.0f', (floatval($t1) + floatval($t2)) * 1000);
}

$transport->close();
PHP运行后结果：

PHPClient Call->key:batu.demo value:1 
PHPClient Call->key:batu.demo value:2 
PHPClient Call->key:batu.demo value:3 
PHPClient Call->key:batu.demo value:4 
PHPClient Call->key:batu.demo value:5 
PHPClient Call->key:batu.demo value:6 
PHPClient Call->key:batu.demo value:7 
PHPClient Call->key:batu.demo value:8 
PHPClient Call->key:batu.demo value:9 
PHPClient Call->key:batu.demo value:10 
本次调用用时: :1431582183296-1431582183290=6毫秒

Go服务端看到打印数据：

–>from client Call: 2015-05-13 22:43:03 php client map[a:batu.demo b:1] 
–>from client Call: 2015-05-13 22:43:03 php client map[a:batu.demo b:2] 
–>from client Call: 2015-05-13 22:43:03 php client map[a:batu.demo b:3] 
–>from client Call: 2015-05-13 22:43:03 php client map[a:batu.demo b:4] 
–>from client Call: 2015-05-13 22:43:03 php client map[a:batu.demo b:5] 
–>from client Call: 2015-05-13 22:43:03 php client map[a:batu.demo b:6] 
–>from client Call: 2015-05-13 22:43:03 php client map[a:batu.demo b:7] 
–>from client Call: 2015-05-13 22:43:03 php client map[b:8 a:batu.demo] 
–>from client Call: 2015-05-13 22:43:03 php client map[a:batu.demo b:9] 
–>from client Call: 2015-05-13 22:43:03 php client map[a:batu.demo b:10] 
Article—>id: 1 Title:插入一篇测试文章 Content:我就是这篇文章内容 Author:liuxinming 
Article—>id: 2 Title:插入二篇测试文章 Content:我就是这篇文章内容 Author:liuxinming

完结，至此一个Golang的Thrift服务端 和 PHP的Thrift客户端完成！
