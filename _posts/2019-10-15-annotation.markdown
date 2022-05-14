---
title: thrift annotation
layout: post
category: golang
author: 夏泽民
---
1,作为注释(desc="errno")
2，在 thrift idl 语法的基础上, 加入一些扩展的 annotation 字段, 用于指导生成 http 以及 thrift 下游服务的 sdk
3，根据不同语言加上不同的前缀来做区分, 比如 go.type, go.filed_name 等.
4，Go语言特有语法
不支持无符号标量类型，如uint64
字段类型codegen和字段标识关系：
require：对于标量类型，codegen代码为标量类型json注解没有omitempty
optional：对于标量类型，codegen代码为指针类型，json注解有omitemtpy
结构体
require和optional结构体都是指针类型、包含omitempty
slice、map
require和optional结构体都是非指针类型、包含omitempty。如果元素是struct，则struct不能是指针。
<!-- more -->
Golang内置了对RPC支持，但只能适用于go语言程序之间调用，且貌似序列化、反序列化性能不高。如果go语言能使用Thrift开发，那么就如虎添翼了。可惜，thrift虽然很早就包含了golang的代码，但一直都存在各种问题无法正确执行，以至于GitHub上有许多大牛小牛自行实现的Thrift代码

gRPC
一个高性能、通用的开源RPC框架，其由Google主要面向移动应用开发并基于HTTP/2协议标准而设计，基于ProtoBuf(Protocol Buffers)序列化协议开发，且支持众多开发语言。
gRPC基于HTTP/2标准设计，带来诸如双向流控、头部压缩、单TCP连接上的多复用请求等特性。这些特性使得其在移动设备上表现更好，更省电和节省空间占用。


annotation 在service级别提供注解，用于生成服务的配置信息，方法级别的注解用于生成方法的配置信息
/**
 * MyService 服务
 */
service MyService {
    MyResp MyFunc(1: MyReq req) (
            /* 会覆盖 service 级别的超时配置 */
            /* url 路径*/
            path="/home／index"
            httpMethod="post"
            contentType="json"
            go.requestProto="JsonProtocol"
            go.responseProto="JsonProtocol"
        )
} (
    /* 此服务对应的 url 前缀.  */
    prefix="／home"
    /* 取值可以是 post/get, 不区分大小写 */
    httpMethod="post"
    /* 连接超时, 重试相关. 这几个参数可以在每个接口中覆盖 */
    /* 故障摘除, 服务健康度相关 */
)

.thrift 官方支持 annotation 语法, annotation 语法提出的目的就是为了满足各种语言 binding 的定制化需求. 默认情况下 annotation 本身是不会作为产出存在于最终生成的代码中 

Annotations are used to associate metadata with types defined in the Thrift definition (".thrift") file. The AnnotationThrift.test file in the source distribution has examples.

Here, for instance, is a struct with annotations (in parentheses):

struct foo {
  1: i32 bar ( presence = "required" );
  2: i32 baz ( presence = "manual", cpp.use_pointer = "", );
  3: i32 qux;
  4: i32 bop;
} (
  cpp.type = "DenseFoo",
  python.type = "DenseFoo",
  java.final = "",
  annotation.without.value,
)
Looking at the code, it seems annotations are only ever used to provide directives to the compiler—for instance the C++ compiler uses the cpp.type annotation, if it's present, to override a type's name in the generated code.

I see nothing that suggests the annotations themselves are ever reproduced in or accessible to the generated code, though if such code does exist it'd be located in compiler/cpp/src/generate/.

annotations can be imagined as extension points to the basic IDL language to control certain aspects of the code generation. They are quite often language specific, but not limited to that. In some cases annotations have become part of the core IDl later. Another thing is, that annotations are still quite poorly documented (I'm being kind here).

can the information that they contain be accessed in the generated client and/or server code?

It depends on the annotation and on what you mean by "accessed". If you change the base class by means of an annotation, or if you make a Java class final, the information is of course present in the generated code.

https://github.com/samuel/go-thrift

annotation 中可以添加任何我们想要的字段, 只要不和某个语言 binding 的 annotation 字段冲突就好


google发布了protobuf v3，为了pb更好用，更跨语言，他对protobuf v2做了以下change：

      1. Removal of field presence logic for primitive value fields（匪夷所思，留存以待以后翻译出来）, 删除required（大大地赞同，即保留repeated，required和optional都不要了，默认就是optional），删除默认值（不明白）。谷歌生成这些改变视为了更好的兼容Android Java、Objective C和Go语言；

      2. 删除对unknown field的支持；

      3. 不再支持继承，以Any type代之；

      4 修正了enum中的unknown类型；

      5 支持map；

         protobuf v2和v3都支持map了，其声明形式如下： 

         message Foo {

                map<string, string> values = 1;

         }

         注意，此处的map是unordered_map。

      6 添加了一些类型集，以支持表述时间、动态数据等；

      7 默认以json形式代替二进制进行编码。

目前v3 alpha版仅仅实现了1-5这五个feature，6和7还未支持。新添加了syntax关键字，以指明proto文件的protobuf协议版本，不指明则是v2。如：

 // foo.proto

      syntax = "proto3";

      message Bar {...}

如果你目前使用了v2，那么暂时不支持你切换到v3，我们还会对v2提供支持。如果你是新手，那就大胆使用v3吧。

