---
title: golang后端graphql接口
layout: post
category: golang
author: 夏泽民
---
1 创建scripts/gqlgen.go
package main 
import "github.com/99designs/gqlgen/cmd" 
func main() 
{ cmd.Execute() }

2 初始化项目
go run scripts/gqlgen.go init
3 更改schema
type Query {
user(id:ID!): User
person(id:ID!): Person
users(name:String):[User!]!
}
type Mutation {
 signupUser(name: String!): User!
 signupPerson(name: String!): Person!
 createUser(name:String!):User!
 createPerson(name:String!):Person!
}
type User {
  id: ID!
  name: String!
}
type Person{
  id: ID!
  name: String!
}
删除schema.resolvers.go
go run scripts/gqlgen.go generate
4 启动服务
go run server.go
<!-- more -->
根据 GraphQL 中文官网代码 中找到graphql-go：一个 Go/Golang 的 GraphQL 实现。

这个库还封装 graphql-go-handler：通过HTTP请求处理GraphQL查询的中间件。
第一步
一般推荐SDL语法的.graphql文件，更强类型要求需要编写类似以下代码。

// schemaQuery 查询函数路由
var schemaQuery= graphql.NewObject(graphql.ObjectConfig{
    Name:        graphql.DirectiveLocationQuery,
    })
    
第二步
进行Schema文档组合


// Schema
var Schema graphql.Schema
Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
    Query:    schemaQuery, // 查询函数Schema
    Mutation: schemaMutation, // 如果有提交函数Schema
})

第三步
获取参数与Schema对应查询函数执行

// ExecuteQuery GraphQL查询器
func ExecuteQuery(params *graphql.Params) *graphql.Result {
    params.Schema = schema
    return graphql.Do(*params)
}
第四步
在路由入口解析参数，并使用查询
https://www.jianshu.com/p/16719baa1713

https://github.com/TsMask/graphql-server-go
https://zhuanlan.zhihu.com/p/35792985
https://blog.csdn.net/qq_41882147/article/details/82966783

super-graph基于golang编写的强大graphql 服务
super-graph 是基于golang 编写的一个graphql 服务（可作为library以及独立的服务）
super-graph 对于graphql 的支持是通过编译graphql查询为sql（hasura就是使用此方法）
https://www.cnblogs.com/rongfengliang/p/12941122.html
https://github.com/dosco/super-graph

类似的工具比较多，比如prisma 、qloo、golang 的gqlgen、apollo-codegen
graphql-code-generator 也是一个不错的工具（灵活、模版自定义。。。）
http://www.mamicode.com/info-detail-2416499.html
https://studygolang.com/articles/13825

https://www.infoq.cn/article/8CTAakhd*EsUtwqIcGNl
https://cloud.tencent.com/developer/article/1477870
