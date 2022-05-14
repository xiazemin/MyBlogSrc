---
title: middleWare
layout: post
category: algorithm
author: 夏泽民
---
https://www.cnblogs.com/FireworksEasyCool/p/12750339.html
gRPC中TLS认证和自定义方法认证，最后还简单介绍了gRPC拦截器的使用。gRPC自身只能设置一个拦截器，所有逻辑都写一起会比较乱。本篇简单介绍go-grpc-middleware的使用，包括grpc_zap、grpc_auth和grpc_recovery。

go-grpc-middleware简介
go-grpc-middleware封装了认证（auth）, 日志（ logging）, 消息（message）, 验证（validation）, 重试（retries） 和监控（retries）等拦截器。

安装 go get github.com/grpc-ecosystem/go-grpc-middleware
<!-- more -->
使用
import "github.com/grpc-ecosystem/go-grpc-middleware"
myServer := grpc.NewServer(
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        grpc_ctxtags.StreamServerInterceptor(),
        grpc_opentracing.StreamServerInterceptor(),
        grpc_prometheus.StreamServerInterceptor,
        grpc_zap.StreamServerInterceptor(zapLogger),
        grpc_auth.StreamServerInterceptor(myAuthFunction),
        grpc_recovery.StreamServerInterceptor(),
    )),
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        grpc_ctxtags.UnaryServerInterceptor(),
        grpc_opentracing.UnaryServerInterceptor(),
        grpc_prometheus.UnaryServerInterceptor,
        grpc_zap.UnaryServerInterceptor(zapLogger),
        grpc_auth.UnaryServerInterceptor(myAuthFunction),
        grpc_recovery.UnaryServerInterceptor(),
    )),
)
grpc.StreamInterceptor中添加流式RPC的拦截器。
grpc.UnaryInterceptor中添加简单RPC的拦截器。

grpc_zap日志记录
1.创建zap.Logger实例

func ZapInterceptor() *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
	grpc_zap.ReplaceGrpcLogger(logger)
	return logger
}
2.把zap拦截器添加到服务端

grpcServer := grpc.NewServer(
	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
		)),
	)
3.日志分析


各个字段代表的意思如下：

{
	  "level": "info",						// string  zap log levels
	  "msg": "finished unary call",					// string  log message

	  "grpc.code": "OK",						// string  grpc status code
	  "grpc.method": "Ping",					/ string  method name
	  "grpc.service": "mwitkow.testproto.TestService",              // string  full name of the called service
	  "grpc.start_time": "2006-01-02T15:04:05Z07:00",               // string  RFC3339 representation of the start time
	  "grpc.request.deadline": "2006-01-02T15:04:05Z07:00",         // string  RFC3339 deadline of the current request if supplied
	  "grpc.request.value": "something",				// string  value on the request
	  "grpc.time_ms": 1.345,					// float32 run time of the call in ms

	  "peer.address": {
	    "IP": "127.0.0.1",						// string  IP address of calling party
	    "Port": 60216,						// int     port call is coming in on
	    "Zone": ""							// string  peer zone for caller
	  },
	  "span.kind": "server",					// string  client | server
	  "system": "grpc",						// string

	  "custom_field": "custom_value",				// string  user defined field
	  "custom_tags.int": 1337,					// int     user defined tag on the ctx
	  "custom_tags.string": "something"				// string  user defined tag on the ctx
}
4.把日志写到文件中

上面日志是在控制台输出的，现在我们把日志写到文件中，修改ZapInterceptor方法。

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapInterceptor 返回zap.logger实例(把日志写到文件中)
func ZapInterceptor() *zap.Logger {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:  "log/debug.log",
		MaxSize:   1024, //MB
		LocalTime: true,
	})

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		w,
		zap.NewAtomicLevel(),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	grpc_zap.ReplaceGrpcLogger(logger)
	return logger
}
grpc_auth认证
go-grpc-middleware中的grpc_auth默认使用authorization认证方式，以authorization为头部，包括basic, bearer形式等。下面介绍bearer token认证。bearer允许使用access key（如JSON Web Token (JWT)）进行访问。

1.新建grpc_auth服务端拦截器

// TokenInfo 用户信息
type TokenInfo struct {
	ID    string
	Roles []string
}

// AuthInterceptor 认证拦截器，对以authorization为头部，形式为`bearer token`的Token进行验证
func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, " %v", err)
	}
	//使用context.WithValue添加了值后，可以用Value(key)方法获取值
	newCtx := context.WithValue(ctx, tokenInfo.ID, tokenInfo)
	//log.Println(newCtx.Value(tokenInfo.ID))
	return newCtx, nil
}

//解析token，并进行验证
func parseToken(token string) (TokenInfo, error) {
	var tokenInfo TokenInfo
	if token == "grpc.auth.token" {
		tokenInfo.ID = "1"
		tokenInfo.Roles = []string{"admin"}
		return tokenInfo, nil
	}
	return tokenInfo, errors.New("Token无效: bearer " + token)
}

//从token中获取用户唯一标识
func userClaimFromToken(tokenInfo TokenInfo) string {
	return tokenInfo.ID
}
代码中的对token进行简单验证并返回模拟数据。

2.客户端请求添加bearer token

实现和上篇的自定义认证方法大同小异。gRPC 中默认定义了 PerRPCCredentials，是提供用于自定义认证的接口，它的作用是将所需的安全认证信息添加到每个RPC方法的上下文中。其包含 2 个方法：

GetRequestMetadata：获取当前请求认证所需的元数据
RequireTransportSecurity：是否需要基于 TLS 认证进行安全传输
接下来我们实现这两个方法

// Token token认证
type Token struct {
	Value string
}

const headerAuthorize string = "authorization"

// GetRequestMetadata 获取当前请求认证所需的元数据
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{headerAuthorize: t.Value}, nil
}

// RequireTransportSecurity 是否需要基于 TLS 认证进行安全传输
func (t *Token) RequireTransportSecurity() bool {
	return true
}
注意：这里要以authorization为头部，和服务端对应。

发送请求时添加token

//从输入的证书文件中为客户端构造TLS凭证
	creds, err := credentials.NewClientTLSFromFile("../tls/server.pem", "go-grpc-example")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	//构建Token
	token := auth.Token{
		Value: "bearer grpc.auth.token",
	}
	// 连接服务器
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&token))
注意：Token中的Value的形式要以bearer token值形式。因为我们服务端使用了bearer token验证方式。

3.把grpc_auth拦截器添加到服务端

grpcServer := grpc.NewServer(cred.TLSInterceptor(),
	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	        grpc_auth.StreamServerInterceptor(auth.AuthInterceptor),
			grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		    grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
			grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
		)),
	)
写到这里，服务端都会拦截请求并进行bearer token验证，使用bearer token是规范了与HTTP请求的对接，毕竟gRPC也可以同时支持HTTP请求。

grpc_recovery恢复
把gRPC中的panic转成error，从而恢复程序。

1.直接把grpc_recovery拦截器添加到服务端

最简单使用方式

grpcServer := grpc.NewServer(cred.TLSInterceptor(),
	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	        grpc_auth.StreamServerInterceptor(auth.AuthInterceptor),
			grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
			grpc_recovery.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		    grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
			grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
            grpc_recovery.UnaryServerInterceptor(),
		)),
	)
2.自定义错误返回

当panic时候，自定义错误码并返回。

// RecoveryInterceptor panic时返回Unknown错误吗
func RecoveryInterceptor() grpc_recovery.Option {
	return grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		return grpc.Errorf(codes.Unknown, "panic triggered: %v", p)
	})
}
添加grpc_recovery拦截器到服务端

grpcServer := grpc.NewServer(cred.TLSInterceptor(),
	grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	        grpc_auth.StreamServerInterceptor(auth.AuthInterceptor),
			grpc_zap.StreamServerInterceptor(zap.ZapInterceptor()),
			grpc_recovery.StreamServerInterceptor(recovery.RecoveryInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		    grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
			grpc_zap.UnaryServerInterceptor(zap.ZapInterceptor()),
            grpc_recovery.UnaryServerInterceptor(recovery.RecoveryInterceptor()),
		)),
	)
总结
本篇介绍了go-grpc-middleware中的grpc_zap、grpc_auth和grpc_recovery拦截器的使用。go-grpc-middleware中其他拦截器可参考GitHub学习使用。

教程源码地址：https://github.com/Bingjian-Zhu/go-grpc-example
