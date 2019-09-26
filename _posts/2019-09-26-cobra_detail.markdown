---
title: cobra detail
layout: post
category: golang
author: 夏泽民
---
Cobra是一个库，其提供简单的接口来创建强大现代的CLI接口，类似于git或者go工具。同时，它也是一个应用，用来生成个人应用框架，从而开发以Cobra为基础的应用。Docker源码中使用了Cobra。

概念
Cobra基于三个基本概念commands,arguments和flags。其中commands代表行为，arguments代表数值，flags代表对行为的改变。

基本模型如下：

APPNAME VERB NOUN --ADJECTIVE或者APPNAME COMMAND ARG --FLAG

例如：

# server是commands，port是flag
hugo server --port=1313

# clone是commands，URL是arguments，brae是flags
git clone URL --bare
Commands
Commands是应用的中心点，同样commands可以有子命令(children commands)，其分别包含不同的行为。

Commands的结构体如下：

type Command struct {
    Use string // The one-line usage message.
    Short string // The short description shown in the 'help' output.
    Long string // The long message shown in the 'help <this-command>' output.
    Run func(cmd *Command, args []string) // Run runs the command.
}
Flags
Flags用来改变commands的行为。其完全支持POSIX命令行模式和Go的flag包。这里的flag使用的是spf13/pflag包，具体可以参考Golang之使用Flag和Pflag.
<!-- more -->
安装与导入
安装
go get -u github.com/spf13/cobra/cobra
导入
import "github.com/spf13/cobra"
Cobra文件结构
cjapp的基本结构
  ▾ cjapp/
    ▾ cmd/
        add.go
        your.go
        commands.go
        here.go
      main.go
main.go
其目的很简单，就是初始化Cobra。其内容基本如下：

package main

import (
  "fmt"
  "os"

  "{pathToYourApp}/cmd"
)

func main() {
  if err := cmd.RootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
使用cobra生成器
cobra init
命令cobra init [yourApp]将会创建初始化应用，同时提供正确的文件结构。同时，其非常智能，你只需给它一个绝对路径，或者一个简单的路径。

cobra add
cobra add serve
<<'COMMENT'
serve created at /home/chenjian/gofile/src/cjapp/cmd/serve.go
COMMENT

cobra add config
<<'COMMENT'
config created at /home/chenjian/gofile/src/cjapp/cmd/config.go
COMMENT

cobra add create -p 'configCmd'
<<'COMMENT'
create created at /home/chenjian/gofile/src/cjapp/cmd/create.go
COMMENT

cobra生成器配置
Cobra生成器通过~/.cjapp.yaml(Linux下)或者$HOME/.cjapp.yaml(windows)来生成LICENSE。

一个.cjapp.yaml格式例子如下：

author: Chen Jian <chenjian158978@gmail.com>
license: MIT
或者可以自定义LICENSE:

license:
  header: This file is part of {\{ .appName }\}.
  text: |
    {\{ .copyright }\}

    This is my license. There are many like it, but this one is mine.
    My license is my best friend. It is my life. I must master it as I must
    master my life.
人工构建Cobra应用
人工构建需要自己创建main.go文件和RootCmd文件。

	Short:   "call me jack",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://o-my-chenjian.com`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("OK")
	},
}

var cfgFile, projectBase, userLicense string

func init() {
	cobra.OnInitialize(initConfig)

	// 在此可以定义自己的flag或者config设置，Cobra支持持久标签(persistent flag)，它对于整个应用为全局
	// 在StringVarP中需要填写`shorthand`，详细见pflag文档
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (defalut in $HOME/.cobra.yaml)")
	RootCmd.PersistentFlags().StringVarP(&projectBase, "projectbase", "b", "", "base project directory eg. github.com/spf13/")
	RootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "Author name for copyright attribution")
	RootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "Name of license for the project (can provide `licensetext` in config)")
	RootCmd.PersistentFlags().Bool("viper", true, "Use Viper for configuration")

	// Cobra同样支持局部标签(local flag)，并只在直接调用它时运行
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 使用viper可以绑定flag
	viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("projectbase", RootCmd.PersistentFlags().Lookup("projectbase"))
	viper.BindPFlag("useViper", RootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	viper.SetDefault("license", "apache")
}

func Execute()  {
	RootCmd.Execute()
}

func initConfig() {
	// 勿忘读取config文件，无论是从cfgFile还是从home文件
	if cfgFile != "" {
		viper.SetConfigName(cfgFile)
	} else {
		// 找到home文件
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 在home文件夹中搜索以“.cobra”为名称的config
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}
	// 读取符合的环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can not read config:", viper.ConfigFileUsed())
	}
}

main.go
main.go的目的就是初始化Cobra

代码下载： cjappmanu_cmd_main.go

package main

import (
	"fmt"
	"os"

	"cjappmanu/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

附加命令
附加命令可以在/cmd/文件夹中写，例如一个版本信息文件，可以创建/cmd/version.go

代码下载： version.go

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ChenJian",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Chen Jian Version: v1.0 -- HEAD")
	},
}

同时，可以将命令添加到父项中，这个例子中RootCmd便是父项。只需要添加：

RootCmd.AddCommand(versionCmd)
处理Flags
Persistent Flags
persistent意思是说这个flag能任何命令下均可使用，适合全局flag：

RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
Local Flags
Cobra同样支持局部标签(local flag)，并只在直接调用它时运行

RootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
Bind flag with Config
使用viper可以绑定flag

var author string

func init() {
  RootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
  viper.BindPFlag("author", RootCmd.PersistentFlags().Lookup("author"))
}
Positional and Custom Arguments
Positional Arguments
Leagacy arg validation有以下几类：

NoArgs: 如果包含任何位置参数，命令报错
ArbitraryArgs: 命令接受任何参数
OnlyValidArgs: 如果有位置参数不在ValidArgs中，命令报错
MinimumArgs(init): 如果参数数目少于N个后，命令行报错
MaximumArgs(init): 如果参数数目多余N个后，命令行报错
ExactArgs(init): 如果参数数目不是N个话，命令行报错
RangeArgs(min, max): 如果参数数目不在范围(min, max)中，命令行报错
Custom Arguments
var cmd = &cobra.Command{
  Short: "hello",
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) < 1 {
      return errors.New("requires at least one arg")
    }
    if myapp.IsValidColor(args[0]) {
      return nil
    }
    return fmt.Errorf("invalid color specified: %s", args[0])
  },
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hello, World!")
  },
}
实例
将root.go修改为以下：

代码下载： example_root.go

package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var echoTimes int

var RootCmd = &cobra.Command{
	Use: "app",
}

var cmdPrint = &cobra.Command{
	Use:   "print [string to print]",
	Short: "Print anything to the screen",
	Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
	},
}

var cmdEcho = &cobra.Command{
	Use:   "echo [string to echo]",
	Short: "Echo anything to the screen",
	Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Print: " + strings.Join(args, " "))
	},
}

var cmdTimes = &cobra.Command{
	Use:   "times [# times] [string to echo]",
	Short: "Echo anything to the screen more times",
	Long: `echo things multiple times back to the user by providing
a count and a string.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < echoTimes; i++ {
			fmt.Println("Echo: " + strings.Join(args, " "))
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	// 两个顶层的命令，和一个cmdEcho命令下的子命令cmdTimes
	RootCmd.AddCommand(cmdPrint, cmdEcho)
	cmdEcho.AddCommand(cmdTimes)
}

func Execute() {
	RootCmd.Execute()
}

func initConfig() {
	// 勿忘读取config文件，无论是从cfgFile还是从home文件
	if cfgFile != "" {
		viper.SetConfigName(cfgFile)
	} else {
		// 找到home文件
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 在home文件夹中搜索以“.cobra”为名称的config
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}
	// 读取符合的环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can not read config:", viper.ConfigFileUsed())
	}
}

操作如下：

go run main.go

<<'COMMENT'
Usage:
  app [command]

Available Commands:
  echo        Echo anything to the screen
  help        Help about any command
  print       Print anything to the screen
  version     Print the version number of ChenJian

Flags:
  -h, --help   help for app

Use "app [command] --help" for more information about a command.
COMMENT


go run main.go echo -h

<<'COMMENT'
echo is for echoing anything back.
Echo works a lot like print, except it has a child command.

Usage:
  app echo [string to echo] [flags]
  app echo [command]

Available Commands:
  times       Echo anything to the screen more times

Flags:
  -h, --help   help for echo

Use "app echo [command] --help" for more information about a command.
COMMENT

go run main.go echo times -h

<<'COMMENT'
echo things multiple times back to the user by providing
a count and a string.

Usage:
  app echo times [# times] [string to echo] [flags]

Flags:
  -h, --help        help for times
  -t, --times int   times to echo the input (default 1)
COMMENT

go run main.go print HERE I AM
<<'COMMENT'
Print: HERE I AM
COMMENT

go run main.go version
<<'COMMENT'
Chen Jian Version: v1.0 -- HEAD
COMMENT

go run main.go echo times WOW -t 3
<<'COMMENT'
Echo: WOW
Echo: WOW
Echo: WOW
COMMENT
自定义help和usage
help
默认的help命令如下：

func (c *Command) initHelp() {
  if c.helpCommand == nil {
    c.helpCommand = &Command{
      Use:   "help [command]",
      Short: "Help about any command",
      Long: `Help provides help for any command in the application.
        Simply type ` + c.Name() + ` help [path to command] for full details.`,
      Run: c.HelpFunc(),
    }
  }
  c.AddCommand(c.helpCommand)
}
可以通过以下来自定义help:

command.SetHelpCommand(cmd *Command)
command.SetHelpFunc(f func(*Command, []string))
command.SetHelpTemplate(s string)
usage
默认的help命令如下：

return func(c *Command) error {
  err := tmpl(c.Out(), c.UsageTemplate(), c)
  return err
}
可以通过以下来自定义help:

command.SetUsageFunc(f func(*Command) error)

command.SetUsageTemplate(s string)
先执行与后执行
Run功能的执行先后顺序如下：

PersistentPreRun
PreRun
Run
PostRun
PersistentPostRun
错误处理函数
RunE功能的执行先后顺序如下：

PersistentPreRunE
PreRunE
RunE
PostRunE
PersistentPostRunE
对不明命令的建议
当遇到不明命令，会有提出一定的建，其采用最小编辑距离算法(Levenshtein distance)。例如：

hugo srever

<<'COMMENT'
Error: unknown command "srever" for "hugo"

Did you mean this?
        server

Run 'hugo --help' for usage.
COMMENT
如果你想关闭智能提示，可以：

command.DisableSuggestions = true

// 或者

command.SuggestionsMinimumDistance = 1
或者使用SuggestFor属性来自定义一些建议，例如：

kubectl remove
<<'COMMENT'
Error: unknown command "remove" for "kubectl"

Did you mean this?
        delete

Run 'kubectl help' for usage.
COMMENT


main 调用 cmd.Execute()，那我们找到这个地方，cmd/root.go 文件：
var RootCmd = &cobra.Command{
    Use:   "cobra_exp1",
    Short: "A brief description of your application",
    Long: 
}

func Execute() {
    if err := RootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }
}


我们看到 Execute() 函数中调用 RootCmd.Execute()，RootCmd 是开始讲组成 Command 结构的一个实例。

在你项目的目录下，运行下面这些命令：
cobra add serve
cobra add config
cobra add create -p 'configCmd'
这样以后，你就可以运行上面那些 app serve 之类的命令了。项目目录如下：
▾ app/
  ▾ cmd/
      serve.go
      config.go
      create.go
    main.go    
    
现在我们有了三个子命令，并且都可以使用，然后只要添加命令逻辑就能真正用了。

Flag
cobra 有两种 flag，一个是全局变量，一个是局部变量。全局什么意思呢，就是所以子命令都可以用。局部的只有自己能用。先看全局的

RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra_exp1.yaml)")
在看局部的：

RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
区别就在 RootCmd 后面的是 Flags 还是 PersistentFlags。

使用标志
标志提供修饰符控制动作命令如何操作

给命令分配一个标志
当标志定义好了，我们需要定义一个变量来关联标志

var Verbose bool
var Source string
持久标志
'持久'表示每个在那个命令下的命令都将能分配到这个标志。对于全局标志，'持久'的标志绑定在root上。

局部标志
Cobra默认只在目标命令上解析标志，父命令忽略任何局部标志。通过打开Command.TraverseChildren Cobra将会在执行任意目标命令前解析标志

command := cobra.Command{
  Use: "print [OPTIONS] [COMMANDS]",
  TraverseChildren: true,
}
绑定标志与配置
你同样可以通过viper绑定标志：

var author string

func init() {
  rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
  viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
}
在这个例子中，永久的标记 author 被viper绑定, 注意, 当用户没有给--author提供值， author不会被赋值。

必须的标记
标记默认是可选的，如果你希望当一个标记没有设置时，命令行报错，你可以标记它为必须的

rootCmd.Flags().StringVarP(&Region, "region", "r", "", "AWS region (required)")
rootCmd.MarkFlagRequired("region")
位置和自定义参数
验证位置参数可以通过 Command的Args字段。

内置下列验证方法

NoArgs - 如果有任何参数，命令行将会报错。
ArbitraryArgs - 命令行将会接收任何参数.
OnlyValidArgs - 如果有如何参数不属于Command的ValidArgs字段，命令行将会报错。
MinimumNArgs(int) - 如果参数个数少于N个，命令行将会报错。
MaximumNArgs(int) - 如果参数个数多余N个，命令行将会报错。
ExactArgs(int) - 如果参数个数不能等于N个，命令行将会报错。
RangeArgs(min, max) - 如果参数个数不在min和max之间, 命令行将会报错.
一个设置自定义验证的例子

var cmd = &cobra.Command{
  Short: "hello",
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) < 1 {
      return errors.New("requires at least one arg")
    }
    if myapp.IsValidColor(args[0]) {
      return nil
    }
    return fmt.Errorf("invalid color specified: %s", args[0])
  },
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hello, World!")
  },
}


运行前和运行后钩子
在命令运行前或运行后，再运行方法非常容易。PersistentPreRun和PreRun方法将会在Run之前执行。PersistentPostRun和PostRun方法将会在Run之后执行。Persistent*Run方法会被子命令继承，如果它们自己没有定义的话。这些方法将按照下面的属性执行：

PersistentPreRun
PreRun
Run
PostRun
PersistentPostRun


生成命令的文档
Cobra 可以基于子命令，标记，等生成文档。以以下格式：

Markdown
ReStructured Text
Man Page
生成bash-completion
Cobra 可以生成一个bash-completion文件。如果你给命令添加更多信息，这些completions可以非常强大和灵活。


https://github.com/spf13/cobra/blob/master/bash_completions.md

You can also configure the bash aliases for the commands and they will also support completions.

alias aliasname=origcommand
complete -o default -F __start_origcommand aliasname

# and now when you run `aliasname` completion will make
# suggestions as it did for `origcommand`.

$) aliasname <tab><tab>
completion     firstcommand   secondcommand

cobra程序只能在GOPATH之下使用，所以首先你需要进入到GOPATH的src目录之下，在该目录下，输入:(否则要输入包名)

$GOPATH/src/$ cobra init demo
在你的当前目录下，应该已经生成了一个demo文件夹:

demo
├── cmd
│   └── root.go
├── LICENSE
└── main.go

$ cobra add test
1
执行完成后，现在我们的demo结构应该是:

.
├── cmd
│   ├── root.go
│   └── test.go
├── LICENSE
└── main.go
可以看到，在cmd目录下，已经生成了一个与我们命令同名的go文件

在init中有一句 RootCmd.AddCommand(testCmd) 这个RootCmd是什么？打开root.go，你会发现RootCmd其实就是我们的根命令。我相信机智的同学已经猜出来我们添加子命令的子命令的方法了。

添加参数
我相信从init函数中的注释中，你已经得到了足够多的信息来自己操作添加flag，但我还是想要啰嗦两句。首先是persistent参数，当你的参数作为persistent flag存在时，如注释所言，在其所有的子命令之下该参数都是可见的。而local flag则只能在该命令调用时执行。可以做一个简单的测试，在test.go的init函数中，添加如下内容:

testCmd.PersistentFlags().String("foo", "", "A help for foo")
testCmd.Flags().String("foolocal", "", "A help for foo")
现在在命令行 go run main.go test -h 得到如下结果:

获取参数值
在知道了如何设置参数后，我们的下一步当然便是需要在运行时获取该参数的值

我们应该在Run这里来获取参数并执行我们的命令功能。获取参数其实也并不复杂。以testCmd.Flags().StringP("aaa", "a", "", "test")此为例，我们可以在Run函数里添加：

str := testCmd.Flags().GetString("aaa")
这样便可以获取到该参数的值了，其余类型参数获取也是同理。


生成文挡
https://github.com/spf13/cobra/blob/master/doc/md_docs.md



生成文档
在root.go中增加
func Get()*cobra.Command{
  return rootCmd
}
main.go 中import
"github.com/spf13/cobra/doc"

在main中增加
  err := doc.GenMarkdownTree(cmd.Get(), "./")
  if err != nil {
    log.Fatal(err)
  }

$ go build main.go
../../../spf13/cobra/doc/man_docs.go:27:2: cannot find package "github.com/cpuguy83/go-md2man/md2man" in any of:
        /usr/local/go/src/github.com/cpuguy83/go-md2man/md2man (from $GOROOT)
        /Users/didi/goLang/src/github.com/cpuguy83/go-md2man/md2man (from $GOPATH)
        /Users/didi/PhpstormProjects/go/src/github.com/cpuguy83/go-md2man/md2man


$ go get github.com/cpuguy83/go-md2man/md2man


$ ./main

生成3个文件
├── gen.md
├── gen_serve.md
├── gen_serve_create.md
├── gen_serve_create_grandchild.md





