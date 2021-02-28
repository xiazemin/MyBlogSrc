---
title: cobertura gover 测试覆盖率
layout: post
category: golang
author: 夏泽民
---
https://github.com/t-yuki/gocover-cobertura
This is a simple helper tool for generating XML output in Cobertura format for CIs like Jenkins and others from go tool cover output.

$ go get code.google.com/p/go.tools/cmd/cover
$ go get github.com/t-yuki/gocover-cobertura

https://github.com/sozorogami/gover
<!-- more -->

 此处使用命令go test -c -covermode=count -ldflags "-X main._VERSION_=$VERSION.${reversion}"  -coverpkg  ./gopath/src/mvdsp/module/,./gopath/src/mvdsp/extractor/,./gopath/src/mvdsp/mvutil/,./gopath/src/mvdsp/protocol/ -o dsp_server.test

① -c 表示 生成测试二进制文件

② -covermode=count 表示 生成的二进制中包含覆盖率计数信息

③ -ldflags 用来将版本信息写入二进制文件,而不使用额外的version文件

④ -coverpkg 后面是要统计覆盖率的文件源码

⑤ -o 后面是输出的二进制文件名

⑥ ......可能还有更多的可用参数，我就不知道了

3 执行命令，生成一个可执行的二进制文件，拷贝到部署目录下

4 启动服务，在启动命令后加参数： -systemTest -test.coverprofile coverage/coverage.cov

① -systemTest 用来启动前面说过的main test

②  -test.coverprofile 用来指定覆盖率信息写入到哪个文件

三  统计覆盖率
1 执行自动化

2 执行如下命令,生成覆盖率文件coverage.cov

1
2
pid_server=`ps -ef | grep "my_server -systemTest" | grep -v "grep" | awk '{print $2}'`
kill $pid_server
覆盖率产生的条件是程序需要正常退出/结束, 因此当自动化运行完毕后，我们需要给程序发送消息表示结束才可以得到覆盖率文件

对此：go是需要return,即可以kill server_pid(若服务使用supervisor启动,还会自己拉起来，不会影响后续调用)

           c/c++ 是需要exit？，发送p __gcov_flush()，见http://www.cnblogs.com/zhaoxd07/p/5608177.html

           python需要触发atexit模块注册一个回调函数，可使用CTRL+C或者kill -2

3 生成测试报告 

go tool cover -html=./coverage/coverage.cov -o /data/reports/coverage.html
./coverage/gocover-cobertura < coverage.cov > /data/reports/coverage.xml　　
生成xml报告时需要安装一个小插件：go get github.com/t-yuki/gocover-cobertura

 

将bin下的gocover-cobertura放到方便执行的目录下即可使用

4 将测试报告集成到jenkins

https://www.cnblogs.com/zhaoxd07/p/8028847.html


https://github.com/gravityblast/fresh