---
title: subline text3 自动生成测试代码
layout: post
category: golang
author: 夏泽民
---
方式1:命令行
#!/bin/bash
cd shell/
~/goLang/bin/gotests -all -w ./
#生成测试文件,添加测试用例
go test -coverprofile=coverage.out
#生成coverage.out
go tool cover -html=coverage.out
#弹出页面
#file:///var/folders/r9/35q9g3d56_d9g0v59w9x2l9w0000gn/T/cover915348153/coverage.html#file0

方式2:subline 插件
subline text 3 注册码
----- BEGIN LICENSE -----
sgbteam
Single User License
EA7E-1153259
8891CBB9 F1513E4F 1A3405C1 A865D53F
115F202E 7B91AB2D 0D2A40ED 352B269B
76E84F0B CD69BFC7 59F2DFEF E267328F
215652A3 E88F9D8F 4C38E3BA 5B2DAAE4
969624E7 DC9CD4D5 717FB40C 1B9738CF
20B3C4F1 E917B5B3 87C38D9C ACCE7DD8
5F7EF854 86B9743C FADC04AA FB0DA5C0
F913BE58 42FEA319 F954EFDD AE881E0B
------ END LICENSE ------

1，下载安装gotest
https://github.com/cweill/gotests
2，配置subline text 3
https://github.com/cweill/GoTests-Sublime
具体安装步骤如下： 
（1）安装Package Control，这个如果已经安装过的朋友可以直接跳过。关于Package Control的安装可以参考《Sublime text 2/3 中 Package Control 的安装与使用方法》

（2）打开Sublime Text 3 ，按住Ctrl+Shift+p ，弹出如下输入窗口，在其中输入install package，并选中红框内的列表。

3，安装gotest
Run the Package Control: Install Package command
Find and install GoTests
Restart Sublime Text (if required)

<img src="{{site.url}}{{site.baseurl}}/img/gotests.gif"/>
4，安装GoSublime
GoSublime 是一个交互式的go build 工具，使用起来也是很方便，主要配合Golang build使用。

点击 Preferences > Package control 菜单(MAC快捷键 shift + command + p)
在弹出的输入框输入install 选择Package control:install package
然后输入GoSublime 选择 GoSublime 安装

5，安装mockery
https://github.com/vektra/mockery
go get github.com/vektra/mockery/.../

<!-- more -->
1,$go get -u github.com/cweill/gotests/...
 # cd /Users/didi/goLang/src/golang.org/x/tools; git pull --ff-only
fatal: unable to access 'https://go.googlesource.com/tools/': Failed to connect to go.googlesource.com port 443: Operation timed out
package golang.org/x/tools/imports: exit status 1

解决方案
mkdir －p  golang.org/x/
https://github.com/golang/tools
cd tools/imports
go install

https://github.com/cweill/gotests
cd gotests 
go install

2,GoTests error: [Errno 2] No such file or directory: 'gotests'.
gopath配置不正确
// GoTests.sublime-settings
{
	// Add your GOPATH here.
	"GOPATH": "/Users/didi/goLang",
}
或者cweill/gotests 没有安装成功

3，Error trying to parse settings: Expected value in Packages/User/GoTests.sublime-settings:1:1
配置文件格式不正确

mock是单元测试中常用的一种测试手法，mock对象被定义，并能够替换掉真实的对象被测试的函数所调用。
而mock对象可以被开发人员很灵活的指定传入参数，调用次数，返回值和执行动作，来满足测试的各种情景假设。
使用场景

依赖的服务返回不确定的结果，如获取当前时间。
依赖的服务返回状态中有的难以重建或复现，比如模拟网络错误。
依赖的服务搭建环境代价高，速度慢，需要一定的成本，比如数据库，web服务
依赖的服务行为多变。
为了保证测试的轻量以及开发人员对测试数据的掌控，采用mock来斩断被测试代码中的依赖不失为一种好方法。

每种编程语言根据语言特点其所采用的mock实现有所不同。

使用
Run: mockery -name=Stringer生成的mock名称 and the following will be output to mocks/Stringer.go
1、生成mock文件，默认在./mocks目录下
mockery -name=Mocker接口名


2、生成mock输出到控制台
mockery -name=接口名 -print 



代码处理：
1.业务接口

1.定义被mock的接口
2.定义类（实现接口）
3.定义个方法New（），返回该类对象


// 1.定义接口
type Driver interface {
    Add(*ImportSet) (int64, error)
    }
// 2.定义类，实现该接口
type driverImp struct{}

//3.定义方法，创建该类对象
func NewDriver() Driver {
    return &driverImp{}
}




2.生成mock方法
mockery -name=接口名 -print
生成要被mock的方法
3.service层调用处
// 用创建对象的方法，获取对象，调用接口方法
NewDriver().QueryImportDetails(req.ImportUuid)


4.单元测试

func TestFunc(t *testing.T) {
//保存原来的对象helper.ListDrv就是newDrv()
    oldFd := helper.ListDrv
    //最后还原对象
    defer func() {
        helper.ListDrv = oldFd
    }()
    
    
    type args struct {
        ctx context.Context
        req 
    }
    tests := []struct {
        name    string
        args    args
        want    *****
        wantErr bool
        //在表驱动中，写mock方法
        mock    func()
    }{
        {
            name: "error",
            args: args{
                ctx: context.Background(),
                req: **,
            },
            wantErr: false,
            want:    ***

            mock: func() {
                mfd := &listmock.Driver{}
                mfd.On("QueryCount").Return(int64(0), nil)
                mfd.On("QueryList", mock.Anything, mock.Anything).Return([]list.ImportList{\{UUID: "xxx", Description: "ddddd"}\}, nil)
                //还原对象
                helper.ListDrv = mfd
            },
        },
        },

        // 测试用例 
    for _, tt := range tests {
    // 执行mock方法
        tt.mock()
        
        t.Run(tt.name, func(t *testing.T) {
            里面的单元测试内容不变
}

常用命令
-v是显示出详细的测试结果,
-cover 显示出执行的测试用例的测试覆盖率
1 测试单个文件，一定要带上被测试的原文件
go test -v -cover  file1_test.go file.go

2 测试单个函数方法
go test -v -cover  -run TestFuncName

3 测试整个api包
在包统计目录下
go test -v -cover ./api/...


4、自动生成测试用例，为指定的函数生成单元测试，输出到控制台
gotests -only "函数名称" file(源文件名称).go


5、自动生成测试用例，为指定文件生成单元测试，输出到文件中
gotests -all -w origin.go, origin_test.go

6、生成mock
 //输出到控制台
 mockery -name=接口名 -print
 
 //输出到./mocks/接口名.go文件中
 mockery -name=接口名
