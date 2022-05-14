---
title: typeAssert
layout: post
category: golang
author: 夏泽民
---
https://ieevee.com/tech/2017/07/29/go-type.html
如果某个函数的入参是interface{}，有下面几种方式可以获取入参的方法：

1 fmt:

import "fmt"
func main() {
    v := "hello world"
    fmt.Println(typeof(v))
}
func typeof(v interface{}) string {
    return fmt.Sprintf("%T", v)
}
2 反射：

import (
    "reflect"
    "fmt"
)
func main() {
    v := "hello world"
    fmt.Println(typeof(v))
}
func typeof(v interface{}) string {
    return reflect.TypeOf(v).String()
}
3 类型断言：

func main() {
    v := "hello world"
    fmt.Println(typeof(v))
}
func typeof(v interface{}) string {
    switch t := v.(type) {
    case int:
        return "int"
    case float64:
        return "float64"
    //... etc
    default:
        _ = t
        return "unknown"
    }
}
其实前两个都是用了反射，fmt.Printf(“%T”)里最终调用的还是reflect.TypeOf()。

func (p *pp) printArg(arg interface{}, verb rune) {
    ...
	// Special processing considerations.
	// %T (the value's type) and %p (its address) are special; we always do them first.
	switch verb {
	case 'T':
		p.fmt.fmt_s(reflect.TypeOf(arg).String())
		return
	case 'p':
		p.fmtPointer(reflect.ValueOf(arg), 'p')
		return
	}
r
<!-- more -->
reflect.TypeOf()的参数是v interface{}，golang的反射是怎么做到的呢？

在golang中，interface也是一个结构体，记录了2个指针：

指针1，指向该变量的类型
指针2，指向该变量的value
如下，空接口的结构体就是上述2个指针，第一个指针的类型是type rtype struct；非空接口由于需要携带的信息更多(例如该接口实现了哪些方法)，所以第一个指针的类型是itab，在itab中记录了该变量的动态类型: typ *rtype。

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}

// nonEmptyInterface is the header for a interface value with methods.
type nonEmptyInterface struct {
	// see ../runtime/iface.go:/Itab
	itab *struct {
		ityp   *rtype // static interface type
		typ    *rtype // dynamic concrete type
		link   unsafe.Pointer
		bad    int32
		unused int32
		fun    [100000]unsafe.Pointer // method table
	}
	word unsafe.Pointer
}
我们来看看reflect.TypeOf():

// TypeOf returns the reflection Type that represents the dynamic type of i.
// If i is a nil interface value, TypeOf returns nil.
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}
TypeOf看到的是空接口interface{}，它将变量的地址转换为空接口，然后将将得到的rtype转为Type接口返回。需要注意，当调用reflect.TypeOf的之前，已经发生了一次隐式的类型转换，即将具体类型的向空接口转换。这个过程比较简单，只要拷贝typ *rtype和word unsafe.Pointer就可以了。

例如w := os.Stdout，该变量的接口值在内存里是这样的：

A *os.File interface value

那么对于第三种，类型断言是怎么判断是不是某个接口呢？回到最初，在golang中，接口是一个松耦合的概念，一个类型是不是实现了某个接口，就是看该类型是否实现了该接口要求的所有函数，所以，类型断言判断的方法就是检查该类型是否实现了接口要求的所有函数。

走读k8s代码的时候，可以看到比较多的类型断言的用法：

func LeastRequestedPriorityMap(pod *api.Pod, meta interface{}, nodeInfo *schedulercache.NodeInfo) (schedulerapi.HostPriority, error) {
	var nonZeroRequest *schedulercache.Resource
	if priorityMeta, ok := meta.(*priorityMetadata); ok {
		nonZeroRequest = priorityMeta.nonZeroRequest
	} else {
		// We couldn't parse metadata - fallback to computing it.
		nonZeroRequest = getNonZeroRequests(pod)
	}
	return calculateUnusedPriority(pod, nonZeroRequest, nodeInfo)
}
类型断言的实现在src/runtime/iface.go里(?)，不过这块代码没看懂，等以后再更新吧。


func assertI2I2(inter *interfacetype, i iface) (r iface, b bool) {
	tab := i.tab
	if tab == nil {
		return
	}
	if tab.inter != inter {
		tab = getitab(inter, tab._type, true)
		if tab == nil {
			return
		}
	}
	r.tab = tab
	r.data = i.data
	b = true
	return
}

func assertE2I2(inter *interfacetype, e eface) (r iface, b bool) {
	t := e._type
	if t == nil {
		return
	}
	tab := getitab(inter, t, true)
	if tab == nil {
		return
	}
	r.tab = tab
	r.data = e.data
	b = true
	return
}
Ref:

the go programming language
go internal
how to find a type of a object in golang

https://stackoverflow.com/questions/20170275/how-to-find-the-type-of-an-object-in-go


在Golang类型断言中可以使用反射数组类型吗
import (
  "reflect"
)
func AnythingToSlice(a interface{}) []interface{} {
  rt := reflect.TypeOf(a)
  switch rt.Kind() {
  case reflect.Slice:
    slice, ok := a.([]interface{})
    if ok {
      return slice
    }
    // it works

  case reflect.Array:
    // return a[:]
    // it doesn't work: cannot slice a (type interface {})   
    //
    array, ok := a.([reflect.ValueOf(a).Len()]interface{})
    // :-((( non-constant array bound reflect.ValueOf(a).Len()
    if ok {
       return array[:]
    }

  }
  return []interface{}(a)
}

An explicit type is required in a type assertion. The type cannot be constructed through reflection.

Unless the argument is a []interface{}, the slice or array must be copied to produce a []interface{}.

Try this:

func AnythingToSlice(a interface{}) []interface{} {
    v := reflect.ValueOf(a)
    switch v.Kind() {
    case reflect.Slice, reflect.Array:
        result := make([]interface{}, v.Len())
        for i := 0; i < v.Len(); i++ {
            result[i] = v.Index(i).Interface()
        }
        return result
    default:
        panic("not supported")
    }
}

1.类型断言
类型断言就是将接口类型的值(x)，装换成类型(T)。格式为：

x.(T)
v:=x.(T)
v,ok:=x.(T)
类型断言的必要条件就是x是接口类型，非接口类型的x不能做类型断言：

var i int=10
v:=i.(int) //错误 i不是接口类型
T可以是非接口类型，如果想断言合法，则T应该实现x的接口

T也可以是接口，则x的动态类型也应该实现接口T

var x interface{}=7  //x的动态类型为int,值为7
i:=x.(int)           // i的类型为int ,值为7

type I interface {m()}
var y I

s:=y.(string)      //非法: string 没有实现接口 I (missing method m)
r:=y.(io.Reader)   //y如果实现了接口io.Reader和I的情况下，  r的类型则为io.Reader
类型断言如果非法，运行时就会出现错误，为了避免这种错误，可以使用一下语法：

v,ok:=x.(T)
ok代表类型断言是否合法，如果非法,ok则为false,这样就不会出现panic了

2.类型切换 type switch
类型切换用来比较类型而不是对值进行比较
type switch它用于检测的是值x的类型T是否匹配某个类型.

格式如下，类似类型断言，但是括号内的不是某个具体的类型，而是单词type:

 switch x.(type){
  
 }
type switch语句中可以有一个简写的变量声明，这种情况下，等价于这个变量声明在每个case clause 隐式代码块的开始位置。如果case clause只列出了一个类型，则变量的类型就是这个类型，否则就是原始值的类型

假设下面的例子中的x的类型为x interface{}:

switch i := x.(type) {
case nil:
  printString("x is nil") // i的类型是 x的类型 (interface{})
case int:
  printInt(i) // i的类型 int
case float64:
  printFloat64(i) // i的类型是 float64
case func(int) float64:
  printFunction(i) // i的类型是 func(int) float64
case bool, string:
  printString("type is bool or string") // i的类型是 x (interface{})
default:
  printString("don't know the type") // i的类型是 x的类型 (interface{})
}


https://studygolang.com/articles/28259?fr=sidebar

golang中的类型
golang中的有两种类型，静态类型（static type）和动态类型（dynamic type）

静态类型：静态类型是在声明变量的时候指定的类型，一个变量的静态类型是永恒不变的，所以才被称为静态类型，也可以简称为类型，普通变量只有静态类型。

package main

import "fmt"

func main()  {
    // 变量i和变量g只有静态类型，变量i的类型为int，i的类型为main.Goose，
    var i int
    var g Goose
    fmt.Printf("%T\n", i)
    fmt.Printf("%T\n", g)
}

type Goose struct {
    age  int
    name string
}

执行结果：
int
main.Sparrow
动态类型：接口类型变量除了静态类型之外还有动态类型，动态类型是由给接口类型变量赋值的具体值的类型来决定的，除了动态类型之外还有动态值，动态类型和动态值是对应的，动态值就是接口类型变量赋值的具体值，之所以被称为动态类型，是因为接口类型的动态类型是会变化的，由被赋予的值来决定。

package main
import "fmt"
func main()  {
    // 动态类型
    var b Bird
    // b 是main.Sparrow类型
    b = Sparrow{}
    fmt.Printf("%T\n", b)
    // b 是main.Parrot类型
    b = Parrot{}
    fmt.Printf("%T\n", b)
}
type Bird interface {
    fly()
    sing()
}
// Goose implement Bird interface
type Sparrow struct {
    age  int
    name string
}
func (s Sparrow) fly()  {
    fmt.Println("I am flying.")
}
func (s Sparrow) sing()  {
    fmt.Println("I can sing.")
}
type Parrot struct {
    age  int
    kind int
    name string
}
func (p Parrot) fly()  {
    fmt.Println("I am flying.")
}
func (p Parrot) sing()  {
    fmt.Println("I can sing.")
}

执行结果：
main.Sparrow
main.Parrot
但是我们在变量b赋值之前直接获取它的类型会发现返回的结果是nil，这看起来很奇怪，interface在Golang里也是一种类型，那它声明的变量的类型为什么是nil呢？

package main
import "fmt"
func main()  {
    // 动态类型
    var b Bird
    fmt.Printf("interface类型变量的类型：%T\n", b)
}
type Bird interface {
    fly()
    sing()
}

执行结果：
interface类型变量的类型：<nil>
首先我们需要明确，一个接口类型变量在没有被赋值之前，它的动态类型和动态值都是 nil 。在使用 fmt.Printf("%T\n") 获取一个变量的类型时，其实是调用了reflect包的方法进行获取的， reflect.TypeOf 获取的是接口变量的动态类型， reflect.valueOf() 获取的是接口变量的动态值。所以 fmt.Printf("%T\n",b) 展示的是 reflect.TypeOf 的结果，由于接口变量 b 还没有被赋值，所以它的动态类型是 nil ，动态值也会是 nil 。

对比来看，为什么只是经过了声明未赋值的变量的类型不是 nil 呢？就像在静态类型部分中所展示的那样。原因如下： 我们先来看一下 reflect.TypeOf 函数的定义，func TypeOf(i interface{}) Type{} ，函数的参数是一个 interface 类型的变量，在调用 TypeOf 时，在接口变量 b 没有赋值之前，它的静态类型与参数类型一致，不需要做转换，因为 b 的动态类型为 nil，所以 TypeOf 返回的结果为 nil 。那为什么变量 i 和变量 g 的类型不为 nil 呢？当变量 i 调用 TypeOf 时，会进行类型的转换，将int型变量i转换为 interface 型，在这个过程中会将变量 i 的类型作为 b 的动态类型，变量 i 的值（在这里是变量 i 的零值0）作为 b 的动态值。因为 TypeOf() 获取的是变量 b 的动态类型，所以这个时候展示出的类型为 int。

golang中的类型断言
因为接口变量的动态类型是变化的，有时我们需要知道一个接口变量的动态类型究竟是什么，这就需要使用类型断言，断言就是对接口变量的类型进行检查，其语法结构如下：

value, ok := x.(T)
x表示要断言的接口变量；
T表示要断言的目标类型；
value表示断言成功之后目标类型变量；
ok表示断言的结果，是一个bool型变量，true表示断言成功，false表示失败，如果失败value的值为nil。
代码示例如下：

package main

import (
    "fmt"
)

func main() {
    // 动态类型
    var b Bird
    // b 是main.Sparrow类型
    b = Sparrow{}

    Listen(b)

    // b 是main.Parrot类型
    b = Parrot{}

    Listen(b)
}


func Listen(b Bird) {
    if _, ok := b.(Sparrow); ok {
        fmt.Println("Listren sparrow sing.")
    }

    if _, ok := b.(Parrot); ok {
        fmt.Println("Listren parrot sing.")
    }

    b.sing()
}

type Bird interface {
    fly()
    sing()
}

type Sparrow struct {
    age  int
    name string
}

func (s Sparrow) fly() {
    fmt.Println("I am sparrow, i can fly.")
}
func (s Sparrow) sing() {
    fmt.Println("I am sparrow, i can sing.")
}

type Parrot struct {
    age  int
    kind int
    name string
}

func (p Parrot) fly() {
    fmt.Println("I am parrot, i can fly.")
}
func (p Parrot) sing() {
    fmt.Println("I am parrot, i can sing.")
}

执行结果如下：
Listren sparrow sing.
I am sparrow, i can sing.
Listren parrot sing.
I am parrot, i can sing.
有时结合 switch 使用更方便

func Listen(b Bird) {
    switch b.(type) {
    case Sparrow:
        fmt.Println("Listren sparrow sing.")
    case Parrot:
        fmt.Println("Listren parrot sing.")
    default:
        fmt.Println("Who are you? I don't know you.")
    }
    b.sing()
}

https://studygolang.com/articles/20041

package labs01

import "testing"

type InterfaceA interface {
    AA()
}

type InterfaceB interface {
    BB()
}

type A struct {
    v int
}

func (a *A) AA() {
    a.v += 1
}

type B struct {
    v int
}

func (b *B) BB() {
    b.v += 1
}

func TypeSwitch(v interface{}) {
    switch v.(type) {
    case InterfaceA:
        v.(InterfaceA).AA()
    case InterfaceB:
        v.(InterfaceB).BB()
    }
}

func NormalSwitch(a *A) {
    a.AA()
}

func InterfaceSwitch(v interface{}) {
    v.(InterfaceA).AA()
}

func Benchmark_TypeSwitch(b *testing.B) {
    var a = new(A)

    for i := 0; i < b.N; i++ {
        TypeSwitch(a)
    }
}

func Benchmark_NormalSwitch(b *testing.B) {
    var a = new(A)

    for i := 0; i < b.N; i++ {
        NormalSwitch(a)
    }
}

func Benchmark_InterfaceSwitch(b *testing.B) {
    var a = new(A)

    for i := 0; i < b.N; i++ {
        InterfaceSwitch(a)
    }
}
