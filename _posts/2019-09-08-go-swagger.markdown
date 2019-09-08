---
title: go-swagger
layout: post
category: golang
author: 夏泽民
---
https://github.com/go-swagger/go-swagger
https://legacy.gitbook.com/book/huangwenchao/swagger/details

"Swagger UI 允许任何人（无论是你的开发团队还是最终用户）在没有任何实现逻辑的情况下对 API 资源进行可视化和交互。它（API文档）通过 Swagger 定义自动生成，可视化文档使得后端实现和客户端消费变得更加容易。"

简而言之，通过提供 Swagger（OpenAPI）定义，您可以获得与 API 进行交互的界面，而不必关心编程语言本身。你可以将 Swagger（OpenAPI） 视为 REST 的 WSDL 。

作为参考，Swagger Codegen 可以从这个定义中，用几十种编程语言来生成客户端和服务器代码。

回到那个时候，我使用的是 Java 和 SpringBoot ，觉得 Swagger 简单易用。你仅需创建一次 bean ，并添加一两个注解到端点上，再添加一个标题和一个项目描述。此外，我习惯将所有请求从 “/” 重定向到 “/swagger-ui” 以便在我打开 host:port 时自动跳转到 SwaggerUI 。在运行应用程序的时候， SwaggerUI 在同一个端口依然可用。（例如，您的应用程序运行在[host]:[port]， SwaggerUI 将在[host]:[port]/swagger-ui上访问到）。

如何为项目加上swagger注释，然后一键生成API文档
开始之前需要安装两个工具：

swagger-editor:用于编写swagger文档，UI展示，生成代码等...
go-swagger:用于一键生成API文档

安装swagger-editor,我这里使用docker运行，其他安装方式，请查看官方文档：
docker pull swaggerapi/swagger-editor
docker run --rm -p 80:8080 swaggerapi/swagger-editor
安装go-swagger,我这边使用brew安装


brew tap go-swagger/go-swagger
brew install go-swagger
<!-- more -->
开始编写注释
1.假设有一个user.server，提供一些REST API，用于对用户数据的增删改查。
比如这里有一个getOneUser接口，是查询用户信息的：
package service

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "user.server/models"
    "github.com/Sirupsen/logrus"
)

type GetUserParam struct {
    Id int `json:"id"`
}

func GetOneUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    decoder := json.NewDecoder(r.Body)
    var param GetUserParam
    err := decoder.Decode(&param)
    if err != nil {
        WriteResponse(w, ErrorResponseCode, "request param is invalid, please check!", nil)
        return
    }

    // get user from db
    user, err := models.GetOne(strconv.Itoa(param.Id))
    if err != nil {
        logrus.Warn(err)
        WriteResponse(w, ErrorResponseCode, "failed", nil)
        return
    }
    WriteResponse(w, SuccessResponseCode, "success", user)
}
复制代码根据swagger文档规范，一个swagger文档首先要有swagger的版本和info信息。利用go-swagger只需要在声明package之前加上如下注释即可：
// Package classification User API.
//
// The purpose of this service is to provide an application
// that is using plain go code to define an API
//
//      Host: localhost
//      Version: 0.0.1
//
// swagger:meta
package service
复制代码然后在项目根目录下使用swagger generate spec -o ./swagger.json命令生成swagger.json文件：

此命令会找到main.go入口文件，然后遍历所有源码文件，解析然后生成swagger.json文件

{
  "swagger": "2.0",
  "info": {
    "description": "The purpose of this service is to provide an application\nthat is using plain go code to define an API",
    "title": "User API.",
    "version": "0.0.1"
  },
  "host": "localhost",
  "paths": {}
}
复制代码2.基本信息有了，然后就要有路由，请求，响应等，下面针对getOneUser接口编写swagger注释：
// swagger:parameters getSingleUser
type GetUserParam struct {
    // an id of user info
    //
    // Required: true
    // in: path
    Id int `json:"id"`
}

func GetOneUser(w http.ResponseWriter, r *http.Request) {
    // swagger:route GET /users/{id} users getSingleUser
    //
    // get a user by userID
    //
    // This will show a user info
    //
    //     Responses:
    //       200: UserResponse
    decoder := json.NewDecoder(r.Body)
    var param GetUserParam
    err := decoder.Decode(&param)
    if err != nil {
        WriteResponse(w, ErrorResponseCode, "request param is invalid, please check!", nil)
        return
    }

    // get user from db
    user, err := models.GetOne(strconv.Itoa(param.Id))
    if err != nil {
        logrus.Warn(err)
        WriteResponse(w, ErrorResponseCode, "failed", nil)
        return
    }
    WriteResponse(w, SuccessResponseCode, "success", user)
}
复制代码可以看到在GetUserParam结构体上面加了一行swagger:parameters getSingleUser的注释信息，这是声明接口的入参注释，结构体内部的几行注释指明了id这个参数必填，并且查询参数id是在url path中。详细用法，参考: swagger:params
在GetOneUser函数中：

swagger:route指明使用的http method，路由，以及标签和operation id,详细用法，参考： swagger:route
Responses指明了返回值的code以及类型

然后再声明响应:
// User Info
//
// swagger:response UserResponse
type UserWapper struct {
    // in: body
    Body ResponseMessage
}

type ResponseMessage struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}
复制代码使用swagger:response语法声明返回值，其上两行是返回值的描述（我也不清楚，为啥描述信息要写在上面，欢迎解惑）,详细用法，参考； swagger:response
然后浏览器访问localhost,查看swagger-editor界面,点击工具栏中的File->Impoprt File上传刚才生成的 swagger.json文件，就可以看到界面：

这样一个简单的api文档就生成了
3.怎么样？是不是很简单？可是又感觉那里不对，嗯，注释都写在代码里了，很不美观，而且不易维护。想一下go-swagger的原理是扫描目录下的所有go文件，解析注释信息。那么是不是可以把api注释都集中写在单个文件内，统一管理，免得分散在各个源码文件内。
新建一个doc.go文件，这里还有一个接口是UpdateUser,那么我们在doc.go文件中声明此接口的api注释。先看一下UpdateUser接口的代码：
func UpdateUser(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    // decode body data into user struct
    decoder := json.NewDecoder(r.Body)
    user := models.User{}
    err := decoder.Decode(&user)
    if err != nil {
        WriteResponse(w, ErrorResponseCode, "user data is invalid, please check!", nil)
        return
    }

    // check if user exists
    data, err := models.GetUserById(user.Id)
    if err != nil {
        logrus.Warn(err)
        WriteResponse(w, ErrorResponseCode, "query user failed", nil)
        return
    }
    if data.Id == 0 {
        WriteResponse(w, ErrorResponseCode, "user not exists, no need to update", nil)
        return
    }

    // update
    _, err = models.Update(user)
    if err != nil {
        WriteResponse(w, ErrorResponseCode, "update user data failed, please try again!", nil)
        return
    }
    WriteResponse(w, SuccessResponseCode, "update user data success!", nil)
}
复制代码然后再doc.go文件中编写如下声明：
package service

import "user.server/models"

// swagger:parameters UpdateUserResponseWrapper
type UpdateUserRequest struct {
    // in: body
    Body models.User
}

// Update User Info
//
// swagger:response UpdateUserResponseWrapper
type UpdateUserResponseWrapper struct {
    // in: body
    Body ResponseMessage
}

// swagger:route POST /users users UpdateUserResponseWrapper
//
// Update User
//
// This will update user info
//
//     Responses:
//       200: UpdateUserResponseWrapper
复制代码这样就把api声明注释给抽离出来了，然后使用命令swagger generate spec -o ./swagger.json生成json文件,就可以看到结果：
很简单吧，参照文档编写几行注释，然后一个命令生成API文档。



<div class="container">
	<div class="row">
	<img src="{{site.url}}{{site.baseurl}}/img/jupyterSlider.png"/>
	</div>
	<div class="row">
	如果你的 API 仅提供在 HTTP 或 HTTPS 上，且只生成 JSON ，您应在此处添加它 - 允许你从每个路由中删除该注释。

安全也被添加在 swagger:meta 中，在 SwaggerUI 上添加一个授权按钮。为了实现JWT，我使用安全类型承载进行命名并将其定义为：

//     Security:
//     - bearer
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
Swagger:route [docs]
有两种方式两个注释你的路由，swagger:operation 和swagger:route 。两者看起来都很相似，那么主要区别是什么？

把 swagger:route 看作简单 API 的短注释,它适用于没有输入参数（路径/查询参数）的 API 。那些（带有参数）的例子是 /repos/{owner} ， /user/{id} 或者 /users/search?name=ribice

如果你有一个那种类型，那么你就必须使用 swagger:operation ，除此之外，如 /user 或 /version 之类的 APIs 都可以用 swagger:route 来注释。

swagger:route注释包含以下内容：

// swagger:route POST /repo repos users createRepoReq
// Creates a new repository for the currently authenticated user.
// If repository name is "exists", error conflict (409) will be returned.
// responses:
//  200: repoResp
//  400: badReq
//  409: conflict
//  500: internal
swagger:route - 注解
POST - HTTP方法
/repo - 匹配路径，端点
repos - 路由所在的空间分割标签，例如，“repos users”
createRepoReq - 用于此端点的请求（详细的稍后会解释）
Creates a new repository … - 摘要（标题）。对于swager:route注释，在第一个句号（.）前面的是标题。如果没有句号，就会没有标题并且这些文字会被用于描述。
If repository name exists … - 描述。对于swager:route类型注释，在第一个句号（.）后面的是描述。
responses: - 这个端点的响应
200: repoResp - 一个（成功的）响应HTTP状态 200，包含 repoResp（用 swagger:response 注释的模型）
400: badReq, 409: conflict, 500: internal - 此端点的错误响应（错误请求，冲突和内部错误， 定义在 cmd/api/swagger/model.go 下）


请记住，您还可能需要使用其他注释，具体取决于您的 API 。由于我将我的项目定义为仅使用单一模式（ https ），并且我的所有 API 都使用 https ，所以我不需要单独注释方案。如果您为端点使用多个模式，则需要以下注释：

// Schemes: http, https, ws, wss
同样适用于 消费者/生产者 媒体类型。我所有的 API 都只消费/生成 application/json 。如果您的 API 正在 消费/生成 其他类型，则需要使用该媒体类型对其进行注释。例如：

// consumes:
// - application/json
// - application/x-protobuf
//
// produces:
// - application/json
// - application/x-protobuf
安全性：

// security:
//   api_key:
//   oauth: read, write
//   basicAuth:
//      type: basic
//   token:
//      type: apiKey
//      name: token
//      in: query
//   accessToken:
//      type: apiKey
//      name: access_token
//      in: query
另一方面，swagger:operation 用于更复杂的端点。三个破折号（-）下的部分被解析为 YAML ，允许更复杂的注释。确保您的缩进是一致的和正确的，否则将无法正确解析。

Swagger:operation docs
使用 Swagger:operation 可以让你使用所有OpenAPI规范，你可以描述你的复杂的端点。如果你对细节感兴趣，你可以阅读规范文档。

简单来说 - swagger:operation 包含如下内容：

// swagger:operation GET /repo/{author} repos repoList
// ---
// summary: List the repositories owned by the given author.
// description: If author length is between 6 and 8, Error Not Found (404) will be returned.
// parameters:
// - name: author
//   in: path
//   description: username of author
//   type: string
//   required: true
// responses:
//   "200":
//     "$ref": "#/responses/reposResp"
//   "404":
//     "$ref": "#/responses/notFound"
swagger:operation - 注释
GET - HTTP 方法
/repo/{author} - 匹配路径，端点
repos - 路由所在的空间分割标签，例如，“repos users”
repoList - 用于此端点的请求。这个不存在（没有定义），但参数是强制性的，所以你可以用任何东西来替换repoList（noReq，emptyReq等）
--- - 这个部分下面是YAML格式的swagger规范。确保您的缩进是一致的和正确的，否则将无法正确解析。注意，如果你在YAML中定义了标签，摘要，描述或操作标签，将覆盖上述常规swagger语法中的摘要，描述，标记或操作标签。
summary: - 标题
description: - 描述
parameters: - URL参数（在这个例子中是{author}）。字符串格式，强制性的（Swagger不会让你调用端点而不输入），位于路径（/{author}）中。另一种选择是参数内嵌的请求 (?name="")
定义你的路由后，你需要定义你的请求和响应。从示例中，你可以看到，我创建了一个新的包，命名为 swagger 。这不是强制性的，它把所有样板代码放在一个名为 swagger 的包中。但缺点是你必须导出你的所有 HTTP 请求和响应。

如果你创建了一个单独的 Swagger 包，确保将它导入到你的主/服务器文件中（你可以通过在导入前加一个下划线来实现）：

_ "github.com/ribice/golang-swaggerui-example/cmd/swagger"
Swagger:parameters [docs]
根据您的应用程序模型，您的 HTTP 请求可能会有所不同（简单，复杂，封装等）。要生成 Swagger 规范，您需要为每个不同的请求创建一个结构，甚至包含仅包含数字（例如id）或字符串（名称）的简单请求。

一旦你有这样的结构（例如一个包含一个字符串和一个布尔值的结构），在你的Swagger包中定义如下：

// Request containing string
// swagger:parameters createRepoReq
type swaggerCreateRepoReq struct {
    // in:body
    api.CreateRepoReq
}
第 1 行包含一个在 SwaggerUI 上可见的注释
第 2 行包含 swagger:parameters 注释，以及请求的名称（operationID）。此名称用作路由注释的最后一个参数，以定义请求。
第 4 行包含这个参数的位置（in:body，in:query 等）
第 5 行是实际的内嵌结构。正如前面所提到的，你不需要一个独立的 swagger 批注包（你可以把swagger:parameters注释放在 api.CreateRepoReq 上），但是一旦你开始创建响应注释和验证，那么在 swagger 相关批注一个单独的包会更清晰。
swagger-parameters

如果你有大的请求，比如创建或更新，你应该创建一个新类型的变量,而不是内嵌结构。例如（注意第五行的区别）:

// Request containing string
// swagger:parameters createRepoReq
type swaggerCreateRepoReq struct {
    // in:body
    Body api.CreateRepoReq
}
这会产生以下 SwaggerUI 请求：

swagger-patameters-ui

Swagger 有很多验证注释提供给 swagger:parameters和 swagger:response ，在注释标题旁边的文档中有详细的描述和使用方法。

Swagger:response [docs]
响应注释与参数注释非常相似。主要的区别在于，经常将响应包裹到更复杂的结构中，所以你必须要在 swagger 中考虑到这点。

在我的示例中，我的成功响应如下所示：

{
     "code":200, // Code containing HTTP status CODE
     "data":{} // Data containing actual response data
}
虽然错误响应有点不同：

{
     "code":400, // Code containing HTTP status CODE
     "message":"" // String containing error message
}
要使用常规响应，像上面错误响应那样的，我通常在 swagger 包内部创建 model.go（或swagger.go）并在里面定义它们。在示例中，下面的响应用于 OK 响应（不返回任何数据）：

// Success response
// swagger:response ok
type swaggScsResp struct {
    // in:body
    Body struct {
        // HTTP status code 200 - OK
        Code int `json:"code"`
    }
}
对于错误响应，除了名称（和示例的情况下的 HTTP 代码注释）之外，它们中的大多数类似于彼此。尽管如此，你仍然应该为每一个错误的情况进行定义，以便把它们作为你的端点可能的响应：

// Error Forbidden
// swagger:response forbidden
type swaggErrForbidden struct {
    // in:body
    Body struct {
        // HTTP status code 403 -  Forbidden
        Code int `json:"code"`
        // Detailed error message
        Message string `json:"message"`
    }
}
data 中包含 model.Repository 的示例响应：

// HTTP status code 200 and repository model in data
// swagger:response repoResp
type swaggRepoResp struct {
    // in:body
    Body struct {
        // HTTP status code 200/201
        Code int `json:"code"`
        // Repository model
        Data model.Repository `json:"data"`
    }
}
data 中包含 model.Repository 切片的示例响应：

// HTTP status code 200 and an array of repository models in data
// swagger:response reposResp
type swaggReposResp struct {
    // in:body
    Body struct {
        // HTTP status code 200 - Status OK
        Code int `json:"code"`
        // Array of repository models
        Data []model.Repository `json:"data"`
    }
}
总之，这将足以生成您的 API 文档。您也应该向文档添加验证，但遵循本指南将帮助您开始。由于这主要是由我自己的经验组成，并且在某种程度上参考了 Gitea 的源代码，我将会听取关于如何改进这部分并相应更新的反馈。

如果您有一些问题或疑问，我建议您查看如何生成FAQ。

本地运行 SwaggerUI
一旦你的注释准备就绪，你很可能会在你的本地环境中测试它。要做到这一点，你需要运行两个命令：

Generate spec [docs]
Serve [docs]
这个命令我们用来生成 swagger.json 并使用 SwaggerUI：

swagger generate spec -o ./swagger.json --scan-models
swagger serve -F=swagger swagger.json
或者，如果你只想使它成为一个命令：

swagger generate spec -o ./swagger.json --scan-models && swagger serve -F=swagger swagger.json
执行该命令后，将使用 Petstore 托管的 SwaggerUI 打开一个新选项卡。服务器启用了 CORS，并将标准 JSON 的 URL 作为请求字符串附加到 petstore URL。

另外，如果使用 Redoc flavor（-F = redoc），则文档将托管在您自己的计算机上（localhost:port/docs）。

在服务器上部署
在服务器上部署生成的 SwaggerUI 有很多种方法。一旦你生成了 swagger.json，它应该相对容易地被运行。

例如，我们的应用程序正在 Google App Engine 上运行。Swagger Spec 由我们的 CI 工具生成，并在 /docs 路径上提供。

我们将 SwaggerUI 作为 Docker 服务部署在 GKE（Google Container/Kubernates Engine）上，它从 /docs 路径中获取swagger.json。

我们的 CI（Wercker）脚本的一部分：

build:
    steps:
        - script:
            name: workspace setup
            code: |
                mkdir -p $GOPATH/src/github.com/orga/repo
                cp -R * $GOPATH/src/github.com/orga/repo/
        - script:
            cwd: $GOPATH/src/bitbucket.org/orga/repo/cmd/api/
            name: build
            code: |
                go get -u github.com/go-swagger/go-swagger/cmd/swagger
                swagger generate spec -o ./swagger.json --scan-models
                CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o app .
                cp app *.template Dockerfile swagger.json "$WERCKER_OUTPUT_DIR"
路由：

func (d *Doc) docHandler(c context.Context, w http.ResponseWriter, r *http.Request) {
    r.Header.Add("Content-Type", "application/json")
    data, _ := ioutil.ReadFile("/swagger.json")
    w.Write(data)
}
Dockerfile：

FROM swaggerapi/swagger-ui
ENV API_URL "https://api.orga.com/swagger"
总结
SwaggerUI 是一个功能强大的 API 文档工具，可以让您轻松而漂亮地记录您的 API。在 go-swagger 项目的帮助下，您可以轻松地生成 SwaggerUI 所需的swagger规范文件（swagger.json）。


swagger 有一整套规范来定义一个接口文件，类似于 thrift 和 proto 文件，定义了服务的请求内容和返回内容，同样也有工具可以生成各种不同语言的框架代码，在 golang 里面我们使用 go-swagger 这个工具，这个工具还提供了额外的功能，可以可视化显示这个接口，方便阅读

下面通过一个例子来简单介绍一下这个框架的使用，还是之前的点赞评论系统：https://github.com/hatlonely/microservices

go-swagger 使用方法
api 定义文件
首先需要写一个 api 定义文件，这里我只展示其中一个接口 countlike，请求中带有某篇文章，返回点赞的次数

paths:
  /countlike:
    get:
      tags:
        - like
      summary: 有多少赞
      description: ''
      operationId: countLike
      consumes:
        - application/json
      produces:
        - application/json
      parameters:
        - name: title
          in: query
          description: 文章标题
          required: true
          type: string
      responses:
        '200':
          description: 成功
          schema:
            $ref: '#/definitions/CountLikeModel'
        '500':
          description: 内部错误
          schema:
            $ref: '#/definitions/ErrorModel'
definitions:
  CountLikeModel:
    type: object
    properties:
      count:
        type: integer
      title:
        type: string
        example: golang json 性能分析
  ErrorModel:
    type: object
    properties:
      message:
        type: string
        example: error message
      code:
        type: integer
        example: 400
这个是 yaml 语法，有点像去掉了括号的 json

这里完整地定义了请求方法、请求参数、正常返回接口、异常返回结果，有了这个文件只需要执行下面命令就能生成框架代码了

swagger generate server -f api/comment_like/comment_like.yaml
还可以下面这个命令可视化查看这个接口文件

swagger serve api/comment_like/comment_like.yaml
这个命令依赖 swagger 工具，可以通过下面命令获取

Mac

brew tap go-swagger/go-swagger
brew install go-swagger
Linux

go get -u github.com/go-swagger/go-swagger/cmd/swagger
export PATH=$GOPATH/bin:$PATH
执行完了之后，你发现多了几个文件夹，其中 cmd 目录里面包含 main 函数，是整个程序的入口，restapi 文件夹下面包含协议相关代码，其中 configure_xxx.go 是需要特别关注的，你需要在这个文件里面实现你具体的业务逻辑

现在你就其实已经可以运行程序了，go run cmd/comment-like-server/main.go，在浏览器里面访问一下你的 api，会返回一个错误信息，告诉你 api 还没有实现，下面就来实现一下吧

业务逻辑实现
api.LikeCountLikeHandler = like.CountLikeHandlerFunc(func(params like.CountLikeParams) middleware.Responder {
    count, err := comment_like.CountLike(params.Title)
    if err != nil {
        return like.NewCountLikeInternalServerError().WithPayload(&models.ErrorModel{
            Code: http.StatusInternalServerError,
            Message: err.Error(),
        })
    }
    return like.NewCountLikeOK().WithPayload(&models.CountLikeModel{
        Count: count,
        Title: params.Title,
    })
})
你只需要在这些 handler 里面实现自己的业务逻辑即可，这里对协议的封装非常好，除了业务逻辑以及打包返回，没有多余的逻辑

再次运行，现在返回已经正常了

统一处理
如果你对请求有一些操作需要统一处理，比如输出统一的日志之类的，可以重写这个函数，也在 configure_xxx.go 这个文件中

func setupGlobalMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        handler.ServeHTTP(w, r)
    })
}

swagger是一款绿色开源的后端工具，以yml或者json格式的说明文档为基点，包含了以此说明文档一站式自动生成后端(路由转发层)代码、api说明文档spec可视化，客户端与交互模型client and model自动生成等模块。本文主要从安装本地swagger editor 到 一键生成 server端，并自制client测试对象slot的增删改查。阅读本文的读者须有基本的go语法与http协议知识.

安装swagger editor 
https://swagger.io/ 官网安装下载swagger editor,解压后从cmd进入该解压文件路径,执行: 
npm install -g 
npm start 
使用npm的前提是要配置号nodeJS的环境，这里不再赘述 
执行完后，cmd窗口上，会显示出服务的url地址，在任意浏览器地址栏上输入即可进入本地swagger editor，其内容和在线版本基本完全一致，所以网络环境健康的读者可以跳过第一步 
http://editor.swagger.io/ 在线editor

编辑slot.yml文档 
yml的语法糖在这里有介绍https://legacy.gitbook.com/book/huangwenchao/swagger/details


该文档描述了slot的增删改查开放的api接口，新建文件夹slotSwagger,点击editor的’generate server’-‘go server’ 


既然goa框架自动生成啦swagger-json文件，那么如何用swagger－ui展示出来呢？

这里分三步：

1.下载swagger－ui的web代码

2.添加swagger.json 和 swagger－ui资源的导出

3.main.go里面mount这两个资源，然后编译启动程序，访问即可

 

为什么连swagger－ui一并导出？因为在swagger－ui中的test程序，需要请求api，如果时部署在不同端口，会有跨域请求问题（这个坑我踩了）。

跨域请求解决有很多方法：

1）把所有api设置为可接受跨域请求

2）把程序和swagger－ui部署到同一个域名下（或者设置代理访问）

3）其它


首先需要github.com/swaggo/gin-swagger和github.com/swaggo/gin-swagger/swaggerFiles（参见gin-swagger）。

然后根据 github.com/swaggo/swag/cmd/swag文档获取到swag工具；执行swag init在项目根目录下生成docs文件夹。然后在路由中import _ "/docs"。这时候编译程序，打开http://localhost:8080/swagger/index.html就可以看到API。有时候打开页面js报错，多刷新几次就有了（原因未知）。

github地址：https://github.com/swaggo/gin-swagger

1、下载swag

$ go get -u github.com/swaggo/swag/cmd/swag
2、在main.go所在目录执行

$ swag init
生成docs/doc.go以及docs/swagger.json,docs/swagger.yaml

3、下载gin-swagger

$ go get -u github.com/swaggo/gin-swagger
$ go get -u github.com/swaggo/files
然后在路由文件引入

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	
	_ "github.com/swaggo/gin-swagger/example/basic/docs" // docs is generated by Swag CLI, you have to import it.
)
并增加swagger访问路由

url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
 

3、一些注解，编写各API handler方法注释（注解格式传送门）

1）main.go主程序文件注释：

// @title Golang Esign API
// @version 1.0
// @description  Golang api of demo
// @termsOfService http://github.com

// @contact.name API Support
// @contact.url http://www.cnblogs.com
// @contact.email ×××@qq.com

//@host 127.0.0.1:8081
func main() {
}
2）handler方法注释：eg

//CreatScene createScene
// @Summary createScene
// @Description createScene
// @Accept multipart/form-data
// @Produce  json
// @Param app_key formData string true "AppKey"
// @Param nonce_str formData string true "NonceStr"
// @Param time_stamp formData string true "TimeStamp"
// @Success 200 {object} app.R
// @Failure 500 {object} app.R
// @Router /dictionaries/createScene [post]
