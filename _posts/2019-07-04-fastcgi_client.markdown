---
title: fastcgi_client
layout: post
category: php
author: 夏泽民
---
CGI全称是“通用网关接口”(Common Gateway Interface)，HTTP服务器与你的或其它机器上的程序进行“交谈”的一种工具，其程序一般运行在网络服务器上。 CGI可以用任何一种语言编写，只要这种语言具有标准输入、输出和环境变量。如php,perl,tcl等。

cgi的弊端
cgi会产生什么问题呢？

当每一个请求进入的时候，cgi都会fork一个新的进程，然后以php为例，每个请求都要耗费相当大的内存，这样一来，并发起来，完全就会GG。
<!-- more -->
为了解决这个问题，于是产生了fastCgi。
	<img src="{{site.url}}{{site.baseurl}}/img/fastcgi1.jpg"/>
		<img src="{{site.url}}{{site.baseurl}}/img/fastcgi2.jpg"/>
			<img src="{{site.url}}{{site.baseurl}}/img/fastcgi3.jpg"/>

从图中可以看出，FastCgi解决的就是这一痛点，当一个新的请求进来的时候，会交给一个已经产生的进程进行处理，而不是fork新的东西出来（php 7还有更多优化，欢迎大家进坑，比如opcached）。

FastCgi协议组成
我目前了解的协议都是由三个部分组成：协议头、正文、结束符（或者说协议尾部），现在我们就来分析下这个协议。

主要参考资料来源于:麻省理工文档 http://www.mit.edu/~yandros/doc/specs/fcgi-spec.html

FastCgi协议头
协议有这几种消息类型

#define FCGI_BEGIN_REQUEST 1
#define FCGI_ABORT_REQUEST 2
#define FCGI_END_REQUEST 3
#define FCGI_PARAMS 4
#define FCGI_STDIN 5
#define FCGI_STDOUT 6
#define FCGI_STDERR 7
#define FCGI_DATA 8
#define FCGI_GET_VALUES 9
#define FCGI_GET_VALUES_RESULT 10
#define FCGI_UNKNOWN_TYPE 11
#define FCGI_MAXTYPE (FCGI_UNKNOWN_TYPE)
很明显，如果我们要发起一个请求，必然要使用消息类型FCGI_BEGIN_REQUEST的进行请求。

看一段通用的包头。

typedef struct {
    unsigned char roleB1;
    unsigned char roleB0;
    unsigned char flags;
    unsigned char reserved[5];
} FCGI_BeginRequestBody;

typedef struct {
    FCGI_Header header;
    FCGI_BeginRequestBody body;
} FCGI_BeginRequestRecord;

typedef struct {
    unsigned char appStatusB3;
    unsigned char appStatusB2;
    unsigned char appStatusB1;
    unsigned char appStatusB0;
    unsigned char protocolStatus;
    unsigned char reserved[3];
} FCGI_EndRequestBody;
我们可以把这个过程看成 begin->(header and content)->end。header这个很有意思，不过联系一下http协议的header,自然一下就了解了，其实在cgi协议里，header应该也属于正文的一种，只是类型比较特殊。

golang的实现
假如我们要实现一个golang fastcgi proxy client的部分，那自然，首先我们要先讲http请求解析成cgi协议的形式。

beginRequest
分析协议得出，web服务器向FastCGI程序发送一个 8 字节 type=FCGI_BEGIN_REQUEST的消息头和一个8字节 FCGI_BeginRequestBody 结构的 消息体，标志一个新请求的开始。

再仔细一想，这不对，服务端是如何知道我们发送的是啥，接受到的数据该如何解析。

从下面开始，我都会用golang来显示代码。

所以，这时候我们需要一个包头。fastCgi通用的包头(此处的header不是http header)为：

type header struct {
	Version       uint8        //协议版本  默认 01
	Type          uint8        // 请求类型
	Id            uint16       // 请求id
	ContentLength uint16  //正文长度
	PaddingLength uint8  //是否有字符补齐
	Reserved      uint8    //预留
}
先产生我们要产生的包体

b := [8]byte{byte(role >> 8), byte(role), flags}
role一般默认用1   1响应器 2验证器 3过滤器
flags 标志是否保持连接 就是http的keeponlive
这个时候，我们就先写包头,然后写包体即可

func (cgi *FCGIClient) writeRecord(recType uint8, reqId uint16, content []byte) (err error) {
	cgi.mutex.Lock()
	defer cgi.mutex.Unlock()
	cgi.buf.Reset()
	cgi.h.init(recType, reqId, len(content))
        //写包头
	if err := binary.Write(&cgi.buf, binary.BigEndian, cgi.h); err != nil {
		return err
	} 
        //写包体
	if _, err := cgi.buf.Write(content); err != nil {
		return err
	}
        //写补位
	if _, err := cgi.buf.Write(pad[:cgi.h.PaddingLength]); err != nil {
		return err
	}
        //将缓冲区写到之前打开的链接  
	_, err = cgi.rwc.Write(cgi.buf.Bytes())
	return err
}
这里应该会很好奇，cgi.rwc是什么。这里我们预留到下一次讲，知道这里是一个建立的socket连接即可。

请求正文
大家都知道 http协议分成 header和body

http request header
cgi协议里有单独对header的处理，即消息类型为FCGI_PARAMS(4)。

我们要从请求里获取到这次的header，很明显,header是一个map[string]string的结构。

func buildEnv(documentRoot string,r *http.Request) (err error,env map[string]string){
	env = make(map[string]string)
	index := "index.php"
	filename := documentRoot+"/"+index
	if r.URL.Path == "/.env" {
		return errors.New("not allow"),env
	} else if r.URL.Path == "/" || r.URL.Path == "" {
		filename = documentRoot + "/" + index
	} else {
		filename = documentRoot + r.URL.Path
	}

	for name,value := range serverEnvironment {
		env[name] = value
	}

	//......其他mapping


	for header, values := range r.Header {
		env["HTTP_" + strings.Replace(strings.ToUpper(header), "-", "_", -1)] = values[0]
	}
	return errors.New("not allow"),env
}
同样，我们先发送包头，再发送包体即可。

http request body
包体就是我们提交的数据,比如Post,delete,put等等操作中，正文包含的数据。

同样，我们先发送包头，再发送包体即可。

但是需要注意的是，一般包体都会很大，但是明显，我们ContentLength只有16位的长度，很大可能是无法一次发送完毕。

因此，我们需要分包进行发送（最大65535）。

请求结束 endRequest
依照上面的，同样，我们发送结束的包头和包体，修改type即可。

获取返回数据
开始我又说一个cgi.rwc是一个socket连接，数据都写往了那里，自然也要从那里读回来。

// recive untill EOF or FCGI_END_REQUEST
	for {
		err1 = rec.read(cgi.rwc)
		if err1 != nil {
			if err1 != io.EOF {
				err = err1
			}
			break
		}
		switch {
		case rec.h.Type == typeStdout:
			retout = append(retout, rec.content()...)
		case rec.h.Type == typeStderr:
			reterr = append(reterr, rec.content()...)
		case rec.h.Type == typeEndRequest:
			fallthrough
		default:
			break
		}
	}
一个简单的demo就完成了。这个时候，我们只需要将接收的数据格式化输出给用户即可。
https://github.com/lwl1989/spinx

http header的构建
golang的http包
golang的http包网上demo实在是太多了，我就不copy了，从golang的core/src来看吧。

net/http

首先进入到golang的安装目录，可以看到src目录，进入到其中，可以找到net/http目录
可以发现，这里几乎涵盖了http所有要用的功能。

那我们要讲请求转发，自然就要看到request.go。

golang http的request
golang 对Request对象（姑且以面向对象的方式解释）定义如下:

type Request struct {                       //忽略这mardown语法影响显示}
	Method           string
	URL              *url.URL
	Proto            string
	ProtoMajor       int
	ProtoMinor       int    
	Header           Header
	Body             io.ReadCloser
	ContentLength    int64
	TransferEncoding []string
	Close            bool
	Host             string
	Form             url.Values
	PostForm         url.Values
	MultipartForm    *multipart.Form
	Trailer          Header
	RemoteAddr       string
	RequestURI       string
	TLS              *tls.ConnectionState
	Cancel <-chan    struct{}
	Response         *Response
	ctx              context.Context
    }
那我们就可以找到大部分我们要的东西。

请求header头，我们可以从Header中获取
请求体，我们可以从Body中获取，如果是表单，我们就要从Form、PostForm或者MultipartForm获取
其他信息，比如host，url之类，我们就可以完美的获取到要的信息。

，当然，不是全部。比如document_root这种
接收数据
 rec := &record{} //建立一个记录缓冲
	var err1 error

	// recive untill EOF or FCGI_END_REQUEST
	for {
		err1 = rec.read(cgi.rwc)   //从cgi建立的链接读取数据
		if err1 != nil {  //判断是否有错误，判断错误是否为终止符
			if err1 != io.EOF {  // keepalive on的时候 不会有终止符  这个要注意
				err = err1
			}
			break
		}
		switch {  //根据返回的类型，将内容写到响应的字节数组中
		case rec.h.Type == typeStdout:   
			retout = append(retout, rec.content()...)
		case rec.h.Type == typeStderr:
			reterr = append(reterr, rec.content()...)
		case rec.h.Type == typeEndRequest:
			fallthrough
		default:
			break
		}
	}
仔细观察type，就是我们在篇一说的消息类型的几种。

我们在假设没有错误出现的情况下，retout这个数组自然是获取所有的返回值。这个时候，golang存储的是字节[22 35 97 ......]，直接将其转化成字符串就可以看到我们熟悉的Http协议的Response了。

如：

头部：
Cache-Control: no-cache, must-revalidate, max-age=0
Connection: keep-alive
Content-Type: text/html; charset=utf-8
Date: Tue, 03 Apr 2018 15:10:58 GMT
Expires: Wed, 11 Jan 1984 05:00:00 GMT
Server: nginx/1.12.2
Transfer-Encoding: chunked
X-Powered-By: PHP/7.1.9

正文(记住正文上面有2个换行符):
//......
这个时候，我们只需要将我们获取到的内容格式化我们想要的内容即可。

发现是PHP有必须对http的content-length进行读，得到长度之后，才会对 body的包进行读取。

所以，我们要加上 header map对content-length的设置

cookies无法设置成功
这个位置我重写了我的赋值位置，不能直接使用获取的 Set-Cookies 的Header进行处理，需要将重写成Cookies，并且 http.SetCookie(write,cookie) 用内置的方法输出。


其他实现：
https://github.com/tomasen/fcgi_client
https://github.com/yookoala/gofast
https://github.com/beberlei/fastcgi-serve
https://github.com/devTransition/job-go-fcgi-proxy