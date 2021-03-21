---
title: golang一个文件里定义多个模板
layout: post
category: golang
author: 夏泽民
---
Go template包下面有两个函数可以创建模板实例
func New(name string) *Template 
func ParseFiles(filenames ...string) (*Template, error)

首先要说的是每一个template内部可以存储多个模板，而且每个模板必须对应一个独立的名字。
两个的不同点在于：

1、使用 New() 在创建时就为其添加一个模板名称，并且执行 t.Execute() 会默认去寻找该名称进行数据融合；

2、使用 ParseFiles() 创建模板可以一次指定多个文件加载多个模板进来，但是就不可以使用 t.Execute() 来执行数据融合；
func (*Template) Execute

func (t *Template) Execute(wr io.Writer, data interface{}) error

Execute方法将解析好的模板应用到data上，并将输出写入wr。如果执行时出现错误，会停止执行，但有可能已经写入wr部分数据。
模板可以安全的并发执行。
但是 ParseFiles() 可以通过

func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error

来进行数据融合，因为该函数可以指定模板名，因此，实例模板就可以知道要去加载自己内部的哪一个模板进行数据融合。

<!-- more -->
func (*Template) ExecuteTemplate

func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error

ExecuteTemplate方法类似Execute，但是使用名为name的t关联的模板产生输出。

因为使用 t.Execute() 无法找到要使用哪个加载过的模板进行数据融合，而只有New()创建时才会指定一个 t.Execute() 执行时默认
加载的模板。
当然无论使用 New() 还是 ParseFiles() 创建模板，都是可以使用 ExecuteTemplate() 来进行数据融合，

但是对于 Execute() 一般与 New() 创建的模板进行配合使用。

html/template 和 text/template
html下的template结构体 实际上是继承了 text 下面的 template结构体

template包下面还有一个 ParseGlob() 方法用于批量解析文件比如在当前目录下有以h开头的模板10个，

使用 template.ParseGlob(“h*”) 即可页将10个模板文件一起解析出来。
注意事项
下面这段代码的输出一定为空

t := template.New("haha")
t, err := t.ParseFiles("header.tmpl")
fmt.Println(err)
t.Execute(os.Stdout, nil)
原因是为什么呢…

首先先记住一个原则 template.New() 和 ParseFiles() 最好不要一起使用，

如果非要一起使用，那么要记住，

New(“TName”) 中的 TName 必须要和 header.tmpl 中定义的{{define name}}中的 name 同名。

但是正常的做法应该是这样的，同样的 ExecuteTemplate() 中输入的 name 也必须和模板中 {{define name}} 相同。

t, _ := template.ParseFiles("header.tmpl")
t.ExecuteTemplate(os.Stdout, "header", nil)
这里要注意下，在这种情况下如果使用 t.Execute() 也是不会输出任何结果的，因为他并不知道你要使用哪个模板。
另外一点要注意的就是

如果模板中没有与填充数据对应的模板语言，那么很有可能panic。

模板中 {{}} 花括号表达式，自动实现了对js代码的过滤，如何不过滤js代码呢，只需要使用 text/template 包下的template，因为html/template包下的模板实现一些针对html的安全操作包括过滤js代码。

Golang 当中支持 Pipeline，一样是使用 |，

Go允许在模板中自定义变量，自定义模板函数。

函数定义必须遵循如下格式：

func FuncName(args ...interface{}) string
通过 template.FuncMap() 强制类型转换为 FuncMap 类型，然后再通过 template实例的 Func(FuncMap) 添加在模板实例中，这样该模板内部在解析时就可以使用该函数。

Go模板包中自定义了一系列内置函数：

var builtins = FuncMap{
    "and":      and,
    "call":     call,
    "html":     HTMLEscaper,
    "index":    index,
    "js":       JSEscaper,
    "len":      length,
    "not":      not,
    "or":       or,
    "print":    fmt.Sprint,
    "printf":   fmt.Sprintf,
    "println":  fmt.Sprintln,
    "urlquery": URLQueryEscaper,
     
     
    // Comparisons
    "eq": eq, // ==
    "ge": ge, // >=
    "gt": gt, // >
    "le": le, // <=
    "lt": lt, // <
    "ne": ne, // !=
}
