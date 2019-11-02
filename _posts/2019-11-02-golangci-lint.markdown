---
title: golangci-lint
layout: post
category: golang
author: 夏泽民
---
https://github.com/golangci/golangci-lint#macos
https://github.com/alecthomas/gometalinter
GolangCI-Lint是一个lint聚合器，它的速度很快，平均速度是gometalinter的5倍。它易于集成和使用，具有良好的输出并且具有最小数量的误报。而且它还支持go modules。最重要的是免费开源。
<!-- more -->
下面公司或者产品都使用了golangci-lint，例如：Google、Facebook、Red Hat OpenShift、Yahoo、IBM、Xiaomi、Samsung、Arduino、Eclipse Foundation、WooCart、Percona、Serverless、ScyllaDB、NixOS、The New York Times和Istio。

安装
CI安装
大多数安装都是为CI（continuous integration）准备的，强烈推荐安装固定版本的golangci-lint。

// 二进制文件将会被安装在$(go env GOPATH)/bin/golangci-lint目录
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin vX.Y.Z
// 或者安装它到./bin/目录
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s vX.Y.Z
// 在alpine Linux中，curl不是自带的，你需要使用下面命令
wget -O - -q https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s vX.Y.Z

brew install golangci/tap/golangci-lint
brew upgrade golangci/tap/golangci-lint

使用
golangci-lint run [目录]/[文件名]

支持的linter
可以通过命令golangci-lint help linters查看它支持的linters。你可以传入参数-E/--enable来使某个linter可用，也可以使用-D/--disable参数来使某个linter不可用。例如：

golangci-lint run --disable-all -E errcheck

配置文件
GolangCI-Lint 完全可以在没有配置文件的情况下工作，可以通过命令行指定各个配置项

GolangCI-Lint 会从当前目录搜索一下命名的配置文件:

.golangci.yml
.golangci.toml
.golangci.json

