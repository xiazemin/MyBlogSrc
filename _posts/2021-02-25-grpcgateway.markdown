---
title: grpc gateway
layout: post
category: golang
author: 夏泽民
---
https://github.com/grpc-ecosystem/grpc-gateway

一、安装

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go

https://github.com/protocolbuffers/protobuf

syntax = "proto3";
package gateway;

import "google/api/annotations.proto";

message StringMessage {
    string value = 1;
}

service Gateway {
   rpc Echo(StringMessage) returns (StringMessage) {
       option (google.api.http) = {
           post: "/v1/example/echo"
           body: "*"
       };
   }
}

<!-- more -->

执行 protoc 编译，生成两个 go 文件，一个是提供 service 的，一个是 gateway 的：

protoc --proto_path=../ -I/usr/local/include -I. -I/home/go-plugin/src -I/home/go-plugin/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. gateway.proto
protoc --proto_path=../ -I/usr/local/include -I. -I/home/go-plugin/src -I/home/go-plugin/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. gateway.proto


https://www.cnblogs.com/linguoguo/p/10148467.html


https://github.com/grpc-ecosystem/grpc-gateway
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

https://www.jianshu.com/p/b78478c1710b

http://dockone.io/article/2836

https://github.com/grpc-ecosystem/grpc-gateway/issues/574    
    
