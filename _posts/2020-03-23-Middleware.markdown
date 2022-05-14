---
title: 用Go编写Web中间件
layout: post
category: golang
author: 夏泽民
---
中间件（通常）是一小段代码，它们接收一个请求，对其进行处理，每个中间件只处理一件事情，完成后将其传递给另一个中间件或最终处理程序，这样就做到了程序的解耦。如果没有中间件那么我们必须在最终的处理程序中来完成这些处理操作，这无疑会造成处理程序的臃肿和代码复用率不高的问题。中间件的一些常见用例是请求日志记录， Header操纵、 HTTP请求认证和 ResponseWriter劫持等等。
<!-- more -->
https://mp.weixin.qq.com/s/3DwHTa-9Bjxei9woi1qeCw
https://mp.weixin.qq.com/s?__biz=MzUzNTY5MzU2MA==&mid=2247483692&idx=1&sn=fa20c127b08a8d35feb4420b09470adf&chksm=fa80d0bbcdf759ad9659f3b6e5ffc4e67217af4059ccadc4053f77ed1148bc783c365c29e4bf&token=140743472&lang=zh_CN&scene=21#wechat_redirect

创建中间件
接下来我们用 Go创建中间件，中间件只将 http.HandlerFunc作为其参数，在中间件里将其包装并返回新的 http.HandlerFunc供服务器服务复用器调用。这里我们创建一个新的类型 Middleware，这会让最后一起链式调用多个中间件变的更简单。

type 
Middleware
 func
(
http
.
HandlerFunc
)
 http
.
HandlerFunc

下面的中间件通用代码模板让我们平时编写中间件变得更容易。

中间件代码模板
中间件是使用装饰器模式实现的，下面的中间件通用代码模板让我们平时编写中间件变得更容易，我们在自己写中间件的时候只需要往样板里填充需要的代码逻辑即可。

func createNewMiddleware
()
 
Middleware
 
{

    
// 创建一个新的中间件

    middleware 
:=
 func
(
next
 http
.
HandlerFunc
)
 http
.
HandlerFunc
 
{

        
// 创建一个新的handler包裹next

        handler 
:=
 func
(
w http
.
ResponseWriter
,
 r 
*
http
.
Request
)
 
{



            
// 中间件的处理逻辑

                        
......

            
// 调用下一个中间件或者最终的handler处理程序

            
next
(
w
,
 r
)

        
}



        
// 返回新建的包装handler

        
return
 handler

    
}



    
// 返回新建的中间件

    
return
 middleware

}

使用中间件
我们创建两个中间件，一个用于记录程序执行的时长，另外一个用于验证请求用的是否是指定的 HTTPMethod，创建完后再用定义的 Chain函数把 http.HandlerFunc和应用在其上的中间件链起来，中间件会按添加顺序依次执行，最后执行到处理函数。完整的代码如下：

package
 main



import
 
(

    
"fmt"

    
"log"

    
"net/http"

    
"time"

)



type 
Middleware
 func
(
http
.
HandlerFunc
)
 http
.
HandlerFunc



// 记录每个URL请求的执行时长

func 
Logging
()
 
Middleware
 
{



    
// 创建中间件

    
return
 func
(
f http
.
HandlerFunc
)
 http
.
HandlerFunc
 
{



        
// 创建一个新的handler包装http.HandlerFunc

        
return
 func
(
w http
.
ResponseWriter
,
 r 
*
http
.
Request
)
 
{



            
// 中间件的处理逻辑

            start 
:=
 time
.
Now
()

            defer func
()
 
{
 log
.
Println
(
r
.
URL
.
Path
,
 time
.
Since
(
start
))
 
}()



            
// 调用下一个中间件或者最终的handler处理程序

            f
(
w
,
 r
)

        
}

    
}

}



// 验证请求用的是否是指定的HTTP Method，不是则返回 400 Bad Request

func 
Method
(
m 
string
)
 
Middleware
 
{



    
return
 func
(
f http
.
HandlerFunc
)
 http
.
HandlerFunc
 
{



        
return
 func
(
w http
.
ResponseWriter
,
 r 
*
http
.
Request
)
 
{



            
if
 r
.
Method
 
!=
 m 
{

                http
.
Error
(
w
,
 http
.
StatusText
(
http
.
StatusBadRequest
),
 http
.
StatusBadRequest
)

                
return

            
}



            f
(
w
,
 r
)

        
}

    
}

}



// 把应用到http.HandlerFunc处理器的中间件

// 按照先后顺序和处理器本身链起来供http.HandleFunc调用

func 
Chain
(
f http
.
HandlerFunc
,
 middlewares 
...
Middleware
)
 http
.
HandlerFunc
 
{

    
for
 _
,
 m 
:=
 range middlewares 
{

        f 
=
 m
(
f
)

    
}

    
return
 f

}



// 最终的处理请求的http.HandlerFunc 

func 
Hello
(
w http
.
ResponseWriter
,
 r 
*
http
.
Request
)
 
{

    fmt
.
Fprintln
(
w
,
 
"hello world"
)

}



func main
()
 
{

    http
.
HandleFunc
(
"/"
,
 
Chain
(
Hello
,
 
Method
(
"GET"
),
 
Logging
()))

    http
.
ListenAndServe
(
":8080"
,
 
nil
)

}

运行程序后会打开浏览器访问 http://localhost:8080会有如下输出：

2020
/
02
/
07
 
21
:
07
:
52
 
/
 
359.503
µ
s

2020
/
02
/
07
 
21
:
09
:
17
 
/
 
34.727
µ
s

到这里怎么用 Go编写和使用中间件就讲完，也就十分钟吧。不过这里更多的是探究实现原理，那么在生产环境怎么自己使用编写的这些中间件呢，我们接着往下看。

使用 gorilla/mux应用中间件
上面我们探讨了如何创建中间件，但是使用上每次用 Chain函数链接多个中间件和处理程序还是有些不方便，而且在上一篇文章中我们已经开始使用 gorilla/mux提供的 Router作为路由器了。好在 gorrila.mux支持向路由器添加中间件，如果发现匹配项，则按照添加中间件的顺序执行中间件，包括其子路由器也支持添加中间件。

gorrila.mux路由器使用 Use方法为路由器添加中间件， Use方法的定义如下：

func 
(
r 
*
Router
)
 
Use
(
mwf 
...
MiddlewareFunc
)
 
{

    
for
 _
,
 fn 
:=
 range mwf 
{

        r
.
middlewares 
=
 append
(
r
.
middlewares
,
 fn
)

    
}

}

它可以接受多个 mux.MiddlewareFunc类型的参数， mux.MiddlewareFunc的类型声明为：

type 
MiddlewareFunc
 func
(
http
.
Handler
)
 http
.
Handler

跟我们上面定义的 Middleware类型很像也是一个函数类型，不过函数的参数和返回值都是 http.Handler接口，在《深入学习用 Go 编写 HTTP 服务器》中我们详细讲过 http.Handler它 是 net/http中定义的接口用来表示处理 HTTP 请求的对象，其对象必须实现 ServeHTTP方法。我们把上面说的中间件模板稍微更改下就能创建符合 gorrila.mux要求的中间件：

func 
CreateMuxMiddleware
()
 mux
.
MiddlewareFunc
 
{



    
// 创建中间件

    
return
 func
(
f http
.
Handler
)
 http
.
Handler
 
{



        
// 创建一个新的handler包装http.HandlerFunc

        
return
 http
.
HandlerFunc
(
func
(
w http
.
ResponseWriter
,
 r 
*
http
.
Request
)
 
{



            
// 中间件的处理逻辑

            
......



            
// 调用下一个中间件或者最终的handler处理程序

            f
.
ServeHTTP
(
w
,
 r
)

        
})

    
}

}

接下来，我们把上面自定义的两个中间件进行改造，然后应用到我们一直在使用的 http_demo项目上，为了便于管理在项目中新建 middleware目录，两个中间件分别放在 log.go和 http_method.go中

//middleware/log.go

func 
Logging
()
 mux
.
MiddlewareFunc
 
{



    
// 创建中间件

    
return
 func
(
f http
.
Handler
)
 http
.
Handler
 
{



        
// 创建一个新的handler包装http.HandlerFunc

        
return
 http
.
HandlerFunc
(
func
(
w http
.
ResponseWriter
,
 r 
*
http
.
Request
)
 
{



            
// 中间件的处理逻辑

            start 
:=
 time
.
Now
()

            defer func
()
 
{
 log
.
Println
(
r
.
URL
.
Path
,
 time
.
Since
(
start
))
 
}()



            
// 调用下一个中间件或者最终的handler处理程序

            f
.
ServeHTTP
(
w
,
 r
)

        
})

    
}

}



// middleware/http_demo.go

func 
Method
(
m 
string
)
 mux
.
MiddlewareFunc
 
{



    
return
 func
(
f http
.
Handler
)
 http
.
Handler
 
{



        
return
 http
.
HandlerFunc
(
func
(
w http
.
ResponseWriter
,
 r 
*
http
.
Request
)
 
{



            
if
 r
.
Method
 
!=
 m 
{

                http
.
Error
(
w
,
 http
.
StatusText
(
http
.
StatusBadRequest
),
 http
.
StatusBadRequest
)

                
return

            
}



            f
.
ServeHTTP
(
w
,
 r
)

        
})

    
}

}

然后在我们的路由器中进行引用：

func 
RegisterRoutes
(
r 
*
mux
.
Router
)
 
{

    r
.
Use
(
middleware
.
Logging
())
// 全局应用

    indexRouter 
:=
 r
.
PathPrefix
(
"/index"
).
Subrouter
()

    indexRouter
.
Handle
(
"/"
,
 
&
handler
.
HelloHandler
{})



    userRouter 
:=
 r
.
PathPrefix
(
"/user"
).
Subrouter
()

    userRouter
.
HandleFunc
(
"/names/{name}/countries/{country}"
,
 handler
.
ShowVisitorInfo
)

    userRouter
.
Use
(
middleware
.
Method
(
"GET"
))
//给子路由器应用

}

再次编译启动运行程序后访问

http
:
//localhost:8080/user/names/James/countries/NewZealand

从控制台里可以看到，记录了这个请求的处理时长：

2020
/
02
/
08
 
09
:
29
:
50
 
Starting
 HTTP server
...

2020
/
02
/
08
 
09
:
55
:
20
 
/
user
/
names
/
James
/
countries
/
NewZealan
 
51.157
µ
s

到这里我们探究完了编写Web中间件的过程和原理，在实际开发中只需要根据自己的需求按照我们给的中间件代码模板编写中间件即可，在编写中间件的时候也要注意他们的职责范围，不要所有逻辑都往里放。
