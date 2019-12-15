---
title: go-restful
layout: post
category: golang
author: 夏泽民
---
https://github.com/emicklei/go-restful
kubernetes源码的时候发现内部使用的是 go-restful实现的

REST
REST要求开发者显式地使用HTTP方法，并与HTTP协议的定义保持一致。REST的基本设计原则是再创建、读取、更新和删除（CRUD）操作和HTTP方法之间建立起一对一的映射。根据该映射：

GET = 检索资源的表示
POST = 创建，使用某种服务器端的算法将内容发送到服务器以创建指定资源或资源的集合
PUT = 创建，如果发送指定URI的完整内容
PUT = 更新，如果要更新指定URI的完整内容
DELETE = 删除，如果要请求服务器删除指定的资源
PATCH = 更新，更新某个资源的部分内容
OPTIONS = 获取请求URI相关的通信选项的信息
go-restful库的特征
路由请求：支持函数映射到路径参数（例如{id}）
路由可配置
（默认）快速路由算法，允许URL路径上出现静态元素、正则表达式、动态参数。比如 /meetings/{id} 或者 /static/{subpath:*}
JSR311规范实现的路由算法
提供了Request API用于从JSON/XML读取结构和访问各种参数（路径参数、查询参数、头部参数）
提供了Response API用于将结构（struct）写入到JSON/XML以及设置头部
支持使用EntityReaderWriter注册的自定义编码
支持在服务级或路由级对请求到响应流的过滤和拦截
支持使用属性来定义请求范围的变量
容器Container支持不同HTTP端点上的WebService
请求和响应的有效负载上的内容编码(gzip,deflate)
（使用过滤器）自动响应OPTIONS请求
（使用过滤器）自动处理CORS请求
支持Swagger UI编写的API文档声明（阅读go-restful-openapi和go-restful-swagger12）
针对HTTP 500状态码的恐慌恢复，使用RecoverHandler(…)自定义处理
路由错误产生HTTP 404/405/406/415等错误，使用ServiceErrorHandler(…)自定义处理
可配置的日志跟踪
使用CompressorProvider注册自定义gzip/deflate的读入器和输出器
<!-- more -->
package main

import (
    "github.com/emicklei/go-restful"
    "github.com/emicklei/go-restful-swagger12"
    "io"
    "log"
    "net/http"
)

func main(){
    wsContainer := restful.NewContainer()

    // 跨域过滤器
    cors := restful.CrossOriginResourceSharing{
        ExposeHeaders:  []string{"X-My-Header"},
        AllowedHeaders: []string{"Content-Type", "Accept"},
        AllowedMethods: []string{"GET", "POST"},
        CookiesAllowed: false,
        Container:      wsContainer}
    wsContainer.Filter(cors.Filter)

    // Add container filter to respond to OPTIONS
    wsContainer.Filter(wsContainer.OPTIONSFilter)



    config := swagger.Config{
        WebServices:    restful.DefaultContainer.RegisteredWebServices(), // you control what services are visible
        WebServicesUrl: "http://localhost:8080",
        ApiPath:        "/apidocs.json",
        ApiVersion:     "V1.0",
        // Optionally, specify where the UI is located
        SwaggerPath:     "/apidocs/",
        SwaggerFilePath: "D:/gowork/src/doublegao/experiment/restful/dist"}
    swagger.RegisterSwaggerService(config, wsContainer)
    //swagger.InstallSwaggerService(config)


    u := UserResource{}
    u.RegisterTo(wsContainer)

    log.Print("start listening on localhost:8080")
    server := &http.Server{Addr: ":8080", Handler: wsContainer}
    defer server.Close()
    log.Fatal(server.ListenAndServe())

}



type UserResource struct{}

func (u UserResource) RegisterTo(container *restful.Container) {
    ws := new(restful.WebService)
    //设置匹配的schema和路径
    ws.Path("/user").Consumes("*/*").Produces("*/*")

    //设置不同method对应的方法，参数以及参数描述和类型
    //参数:分为路径上的参数,query层面的参数,Header中的参数
    ws.Route(ws.GET("/{id}").
        To(u.result).
        Doc("方法描述：获取用户").
        Param(ws.PathParameter("id", "参数描述:用户ID").DataType("string")).
        Param(ws.QueryParameter("name", "用户名称").DataType("string")).
        Param(ws.HeaderParameter("token", "访问令牌").DataType("string")).
        Do(returns200, returns500))
    ws.Route(ws.POST("").To(u.result))
    ws.Route(ws.PUT("/{id}").To(u.result))
    ws.Route(ws.DELETE("/{id}").To(u.result))

    container.Add(ws)
}

func (UserResource) SwaggerDoc() map[string]string {
    return map[string]string{
        "":         "Address doc",//空表示结构本省的描述
        "country":  "Country doc",
        "postcode": "PostCode doc",
    }
}

func (u UserResource) result(request *restful.Request, response *restful.Response) {
    io.WriteString(response.ResponseWriter, "this would be a normal response")
}

func returns200(b *restful.RouteBuilder) {
    b.Returns(http.StatusOK, "OK", "success")
}

func returns500(b *restful.RouteBuilder) {
    b.Returns(http.StatusInternalServerError, "Bummer, something went wrong", nil)
}

二、详解
1、WebServices and Routes
WebService拥有一个Route对象的集合，这些Route对象负责分发即将到来的HTTP请求到相应的函数调用。一般来说，WebService有一个root根路径（例如：/users），还为路由定义了常见的MIME类型。WebService必须添加到某个容器以便能从服务器处理HTTP请求。

Route是根据HTTP分发、URL路径和MIME类型所定义的。

正则表达式匹配路由
路由参数可以使用格式“uri/{var[:regexp]}”或指定版本的“uri/{var:*}”匹配路径尾部进行指定。例如，/persons/{name:[A-Z][A-Z]} 可用来限制name参数的值只能包含大写字母。正则表达式必须使用regexp包中的标准Go语法进行描述。https://code.google.com/p/re2/wiki/Syntax，此功能需要使用CurlyRouter。

2、容器（Container）
容器可以hold住一套WebService的集合、过滤器、一个用于http请求多路复用的http.ServeMux。使用语句“restful.Add(…)”和“restful.Filter(…)”，前者可以在容器注册一个WebService，后者可以过滤。go-restful默认的容器使用http.DefaultServeMux。用户可以创建自己的容器，以及为指定的容器创建一个新的http.Server。

container := restful.NewContainer()
server := &http.Server{Addr: ":8081", Handler: container}
1
2
3、过滤器（Filter）
过滤器可以动态拦截请求和响应，以及转换或使用请求和响应中包含的信息。用户可以使用过滤器来执行常规的日志记录、测量、验证、重定向、设置响应头部Header等。restful包中有三个针对请求、响应流的钩子，还可以添加过滤器。每个过滤器必须定义一个FilterFunction：

func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain)
1
使用如下语句传递请求/响应对到下一个过滤器或RouteFunction：

chain.ProcessFilter(req, resp)
1
4、容器过滤器（Container Filter）
在注册WebService之前处理：

// 安装一个（全局的）过滤器到默认的容器
restful.Filter(globalLogging)
1
2
5、WebService过滤器（WebService Filter）
在路由WebService之前处理：

// 安装一个WebService过滤器
ws.Filter(webserviceLogging).Filter(measureTime)
1
2
6、路由过滤器（Route Filter）
在调用路由Route相关的函数之前处理：

// 安装2个链式的路由过滤器
ws.Route(ws.GET("/{user-id}").Filter(routeLogging).Filter(NewCountFilter().routeCounter))
1
2
例子可见：https://github.com/emicklei/go-restful/blob/master/examples/restful-filters.go

7、响应编码（Response Encoding）
支持两种响应编码：gzip和deflate。要为所有的响应启用它们：

restful.DefaultContainer.EnableContentEncoding(true)
1
如果某个Http请求包含了Accept-Encoding头部，那么响应内容必须使用指定的编码进行压缩。或者，可以创建一个过滤器执行编码并安装它到每一个WebService和Route。

见例子：https://github.com/emicklei/go-restful/blob/master/examples/restful-encoding-filter.go

8、OPTIONS支持
通过安装预定义的容器过滤器，你的WebService可以响应HTTP OPTIONS请求。

Filter(OPTIONSFilter())
1
9、CORS
通过安装CrossOriginResourceSharing过滤器，你的WebService可以处理CORS请求。

cors := CrossOriginResourceSharing{ExposeHeaders: []string{"X-My-Header"}, CookiesAllowed: false, Container: DefaultContainer}
Filter(cors.Filter)
1
2
10、异常处理
意想不到的事情发生。如果因为故障而不能处理请求，服务端需要通过响应告诉客户端发生了什么和为什么。因此使用HTTP状态码，更重要的是要正确的使用状态码。

400: Bad Request
如果路径或查询参数无效（内容或类型），那么使用http.StatusBadRequest。

404: Not Found
尽管URI有效，但请求的资源可能不可用。

500: Internal Server Error
如果应用程序逻辑无法处理请求（或编写响应），则使用http.StatusInternalServerError。

405: Method Not Allowed
请求的URL是有效的，但请求使用的HTTP方法（GET，PUT，POST，…）是不允许的。

406: Not Acceptable
请求的头部没有或设置了未知Accept Header。

415: Unsupported Media Type
请求的头部没有或设置了未知的Content-Type报头。

11、ServiceError
除了设置HTTP状态码，还应该为响应选择写适当的ServiceError消息。

12、性能选项（Performance Options）
这个包有几个选项，它们可能会影响服务的性能。重要的是要理解这些选项，正确地设置它们。

restful.DefaultContainer.DoNotRecover(false)
DoNotRecover控制是否因返回HTTP 500状态码而（恐慌）停止服务。如果设置为false，那么容器Container会恢复服务。默认值为true。

restful.SetCompressorProvider(NewBoundedCachedCompressors(20, 20))
如果启用了内容编码，那么获得新gzip/zlib输出器（writer）和读入器（reader）的默认策略是使用sync.Pool。由于输出器writer是昂贵的结构，当使用预加载缓存时性能提高非常明显。你也可以注入自己的实现。

13、故障排除（Trouble shooting）
这个包可以对完整的Http请求的匹配过程和过滤器调用产生详细的日志记录。启用此功能需要你设置restful.StdLogger的实现，例如log.Logger：

restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|logs.Lshortfile))
1
使用案例
https://github.com/emicklei/mora
https://github.com/emicklei/landskape

