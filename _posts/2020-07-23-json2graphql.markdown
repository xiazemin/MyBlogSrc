---
title: json2graphql
layout: post
category: golang
author: 夏泽民
---
json2graphql
json2graphql 是一个根据 json 生成 GraphQL Schema 的工具。
可在 https://luojilab.github.io/json2graphql/ 在线体验其功能。

https://walmartlabs.github.io/json-to-simple-graphql-schema/
<!-- more -->
json2graphql
json2graphql 是一个根据 json 生成 GraphQL Schema 的工具。
可在 https://luojilab.github.io/json2graphql/ 在线体验其功能。

关于 GraphQL
GraphQL 是一个用于 API 的查询语言，是一个使用基于类型系统来执行查询的服务端运行时（类型系统由你的数据定义）。GraphQL 并没有和任何特定数据库或者存储引擎绑定，而是依靠你现有的代码和数据支撑。由于其强类型，返回结果可定制，自带聚合功能等特性，由 facebook 开源后，被 github 等各大厂广泛使用。

json protobuf 与 GraphQL
由于 protobuf 和 GraphQL 都是强类型的，所以可以直接从 protobuf 的 schema 生成 GraphQL Schema,因而才能有自动聚合 grpc 服务生成 GraphQL 接口的框架 rejoiner。但同样的方法不适用于 json，因为标准的 json 并不包含 schema,单纯根据 json 文件无法确定知道每个字段的类型（因为有空值，以及嵌套的情况)。因而目前无法实现类似 rejoiner for json 这样的全自动框架。
我们虽不能生成最终的 GraphQL Schema,但是基于对 json 的解析和一些约定，我们可以生成一个 GraphQL Schema 的草稿，生成 Schema 的绝大部分内容，并将有疑问的地方标记出来。
json2graphql 就是一个用 golang 实现的 json 生成 schema 的工具。如果你不熟悉 golang,可以使用其在线版本 https://luojilab.github.io/json2graphql/

在从 REST API 迁移到 GraphQL 的过程中，我们有很多接口会返回大量字段（几十个），如果完全手动编写这些 Schema，将是非常痛苦的，我们开发 json2graphql 的初衷就是解决这个问题，大大缩短开发时间。

https://www.jianshu.com/p/d473563c79ef
https://github.com/luojilab/json2graphql
https://www.jianshu.com/p/ea9b2fd7f647