---
title: protoc-gen-validate
layout: post
category: algorithm
author: 夏泽民
---
https://github.com/envoyproxy/protoc-gen-validate

https://github.com/mwitkow/go-proto-validators


<!-- more -->
This project is currently in alpha. The API should be considered unstable and likely to change

PGV is a protoc plugin to generate polyglot message validators. While protocol buffers effectively guarantee the types of structured data, they cannot enforce semantic rules for values. This plugin adds support to protoc-generated code to validate such constraints.

Developers import the PGV extension and annotate the messages and fields in their proto files with constraint rules:

syntax = "proto3";

package examplepb;

import "validate/validate.proto";

message Person {
  uint64 id    = 1 [(validate.rules).uint64.gt    = 999];

  string email = 2 [(validate.rules).string.email = true];

  string name  = 3 [(validate.rules).string = {
                      pattern:   "^[^[0-9]A-Za-z]+( [^[0-9]A-Za-z]+)*$",
                      max_bytes: 256,
                   }];

  Location home = 4 [(validate.rules).message.required = true];

  message Location {
    double lat = 1 [(validate.rules).double = { gte: -90,  lte: 90 }];
    double lng = 2 [(validate.rules).double = { gte: -180, lte: 180 }];
  }
}
Executing protoc with PGV and the target language's default plugin will create Validate methods on the generated types:

p := new(Person)

err := p.Validate() // err: Id must be greater than 999
p.Id = 1000

err = p.Validate() // err: Email must be a valid email address
p.Email = "example@lyft.com"

err = p.Validate() // err: Name must match pattern '^[^\d\s]+( [^\d\s]+)*$'
p.Name = "Protocol Buffer"

err = p.Validate() // err: Home is required
p.Home = &Location{37.7, 999}

err = p.Validate() // err: Home.Lng must be within [-180, 180]
p.Home.Lng = -122.4

err = p.Validate() // err: nil
 

This project is currently in alpha. The API should be considered unstable and likely to change

PGV is a protoc plugin to generate polyglot message validators. While protocol buffers effectively guarantee the types of structured data, they cannot enforce semantic rules for values. This plugin adds support to protoc-generated code to validate such constraints.

Developers import the PGV extension and annotate the messages and fields in their proto files with constraint rules:

syntax = "proto3";

package examplepb;

import "validate/validate.proto";

message Person {
  uint64 id    = 1 [(validate.rules).uint64.gt    = 999];

  string email = 2 [(validate.rules).string.email = true];

  string name  = 3 [(validate.rules).string = {
                      pattern:   "^[^[0-9]A-Za-z]+( [^[0-9]A-Za-z]+)*$",
                      max_bytes: 256,
                   }];

  Location home = 4 [(validate.rules).message.required = true];

  message Location {
    double lat = 1 [(validate.rules).double = { gte: -90,  lte: 90 }];
    double lng = 2 [(validate.rules).double = { gte: -180, lte: 180 }];
  }
}
Executing protoc with PGV and the target language's default plugin will create Validate methods on the generated types:

p := new(Person)

err := p.Validate() // err: Id must be greater than 999
p.Id = 1000

err = p.Validate() // err: Email must be a valid email address
p.Email = "example@lyft.com"

err = p.Validate() // err: Name must match pattern '^[^\d\s]+( [^\d\s]+)*$'
p.Name = "Protocol Buffer"

err = p.Validate() // err: Home is required
p.Home = &Location{37.7, 999}

err = p.Validate() // err: Home.Lng must be within [-180, 180]
p.Home.Lng = -122.4

err = p.Validate() // err: nil

1. 安装Go
1.1 下载Go
wget https://studygolang.com/dl/golang/go1.13.4.linux-amd64.tar.gz
# 解压
tar -zxvf go1.13.4.linux-amd64.tar.gz
1
2
3
1.2 配置go环境
编辑 /etc/profile 文件

vim ~/.bashrc
1
将下面内容加入到末尾（GOPAT是我Windows中的GOPATH）

export GOROOT=/usr/local/go
export GOPATH=/home/pibigstar/goWork
export PATH=$GOPATH/bin:$PATH
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
1
2
3
4
2. 工具安装
2.1 安装protoc
去这个网址下载：https://github.com/protocolbuffers/protobuf/releases，根据自己的系统，下载对应的文件
解压：

unzip protoc.zip 
1
将bin文件夹下的 protoc 复制到Linux 中的 /bin目录下

sudo cp protoc/bin/protoc /bin/protoc
1
执行 protoc -verson 如果输出版本信息则证明配置成功

2.2 安装protoc-gen-go
# 下载
git clone https://github.com/golang/protobuf.git
# 进入目录
cd protobuf/protoc-gen-go
# 编译
go install
1
2
3
4
5
6
2.3 安装protoc-gen-validate
这个是用来生成pb的校验规则文件，也就是*.pb.validate.go

go get -u github.com/envoyproxy/protoc-gen-validate
1
2.4 安装protoc-gen-grpc-gateway
这个是用来生成让grpc支持http调用的

git clone https://github.com/grpc-ecosystem/grpc-gateway.git
# 安装 protoc-gen-grpc-gateway
cd grpc-gateway/protoc-gen-grpc-gateway
go install .
# 安装protoc-gen-swagger
cd grpc-gateway/protoc-gen-swagger
go install .
1
2
3
4
5
6
7
2.5 安装protoc-gen-doc
这个是用来生成 proto 的文档文件，会生成一个 html 格式的文档， 下载地址：https://github.com/pseudomuto/protoc-gen-doc/releases

wget https://github.com/pseudomuto/protoc-gen-doc/releases/download/v1.3.0/protoc-gen-doc-1.3.0.windows-amd64.go1.11.2.tar.gz

# 解压
tar -zxvf protoc-gen-doc-1.3.0.windows-amd64.go1.11.2.tar.gz
1
2
3
4
2.6 安装proto-gen-java
可在这个地址下载 protoc-gen-java工具，https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/ ，记得把文件名改为protoc-gen-grpc-java.exe

3. 编译
3.1 hello.proto
hello.proto

syntax="proto3";

package main;

message Hello {
	string value = 1;
}
1
2
3
4
5
6
7
3.2 编译
3.2.1 编译为Go代码（protoc-gen-go)
protoc --go_out=plugins=grpc,paths=source_relative:. --validate_out="lang=go,paths=source_relative:." hello.proto
1
注意

paths参数
使用 source_relative 则不会使用option go_package中指定的路径
使用 import 则是使用option go_package 中指定的路径
3.2.2 编译为Java代码 (protoc-gen-java)
可在这个地址下载 protoc-gen-java工具，https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/1.0.1/ ，记得把文件名改为protoc-gen-grpc-java.exe

protoc --java_out=. --grpc-java_out=. hello.proto
1
3.2.3 生成proto文档(proto-gen-doc)
protoc --doc_out=. --doc_opt=html,index.html:Ignore* hello.proto user.proto
1
3.3 复杂点的proto
syntax="proto3";

//生成的pb文件中package为admin
package admin;

//生成go文件的路径
option go_package          = "pb/admin";
//关闭Java多文件
option java_multiple_files = false;
//生成的Java文件的package路径
option java_package        = "pb.admin";

service UserService {
	rpc Login(LoginReq) returns (LoginResp);
}

message LoginReq {
	string username = 1;
	string password = 2;
}

message LoginResp {
	int32 code = 1;
    string msg = 2;
}
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
3.4 生成gw.pb文件
proto 导入 google/api/annotations.proto

syntax = "proto3";

package pb;
option go_package  = "pb/echo";

import "google/api/annotations.proto";

service Echo {
    rpc UnaryEcho (EchoRequest) returns (EchoResponse){
        option (google.api.http) = {
            post: "/v1/example/echo"
            body: "*"
        };
    }
}

// EchoRequest is the request for echo.
message EchoRequest {
    string message = 1;
}

// EchoResponse is the response for echo.
message EchoResponse {
    string message = 1;
}
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
21
22
23
24
25
下载地址： https://github.com/googleapis/googleapis/tree/master/google
在proto下面将 google的api 放进去
在这里插入图片描述
编译时用 -I 命令 引入 extra/src 下面的文件

protoc -I=extra/src:. \
--grpc-gateway_out=logtostderr=true:pb \
--go_out=plugins=grpc,paths=import:pb echo.proto
1
2
3
3.5 编译脚本
#!/usr/bin/env bash

TARGET="../"

if [ -n "$1" ]; then
    TARGET=$1
fi

# 排除掉 extra/src 目录
for file in `find . -path ./extra/src -prune -o -name '*.proto' -print`;
do
	echo $file
	protoc -I=extra/src:. --grpc-gateway_out=$TARGET\
	--go_out=plugins=grpc,paths=import:$TARGET \
	--validate_out="lang=go,paths=import:$TARGET" $file
done
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
3.6 代码注册
package main

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net/http"

	pb "data/proto/demo/echo"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterEchoHandlerFromEndpoint(ctx, mux, "127.0.0.1:6000", opts)
	if err != nil {
		return err
	}
	http.ListenAndServe(":8080", mux)
}

