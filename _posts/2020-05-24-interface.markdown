---
title: golang 父类调用子类方法 继承多态的实现方式
layout: post
category: golang
author: 夏泽民
---
go 语言中，当子类调用父类方法时，“作用域”将进入父类的作用域，看不见子类的方法存在
我们可以通过参数将子类传递给父类，实现在父类中调用子类方法。

总结下有三种
<!-- more -->
一、 基于接口
定义接口，父子类都实现接口，父类方法接收接口类型参数

特点：

结构简单，思路清晰。
基于接口，轻松应对多级继承的情况。

func (a *A) Func3(c C)  参数是接口类型，父类通过接口调用，也可以进行一次类型推断，没有必要

二、 基于反射
父类方法接收子类对象，通过反射调用子类方法
func (self A) sayReal(child interface{}) {
    ref := reflect.ValueOf(child)
    method := ref.MethodByName("Name")
    if (method.IsValid()) {
        r := method.Call(make([]reflect.Value, 0))
        fmt.Println(r[0].String())
    } else {
        // 错误处理
    }
}


三、父类定义方法的成员变量，子类set这个成员变量
b.A.func1=b.func1

完整代码
package main

import "fmt"

func main() {
	b := B{
		A{}}
	b.A.func1=b.func1
	b.Func3(b)
	/**
		A::func3
	    A:: function2
	    panic: runtime error: invalid memory address or nil pointer
	*/
}

type C interface {
	//func1()
	/*
		./main.go:9:9: cannot use b (type B) as type C in argument to b.A.Func3:
			B does not implement C (func1 method has pointer receiver)
	*/
	func3()
}

type A struct{
 func1 func()
}

func (a *A) func2() {
	fmt.Println("A:: function2")
}

func (a *A) Func3(c C) {
	fmt.Println("A::func3")
	a.func2()
	a.func1()
	c.func3()
	if b,ok:=c.(B);ok{ //类型推断没有必要，注意接口的方法接受者不是指针
		b.func2()
	}
}

type B struct {
	A
}

func (b *B) func1() {
	fmt.Println("B::func1")
}

func (b *B) func2() {
	fmt.Println("B::func2")
}

func (b B)func3()  {
	 fmt.Println("B:func3")
}

https://github.com/xiazemin/object

注意b.A.func1=b.func1 这种写法，因为是私有变量，成员变量对B不可以见，所以b.func1 取的是b的函数而不是成员变量，否则取的是继承过来的成员变量。 方便起见可以用大小写区分
