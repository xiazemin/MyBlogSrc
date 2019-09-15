---
title: interface
layout: post
category: golang
author: 夏泽民
---
 interface(接口)是golang最重要的特性之一，Interface类型可以定义一组方法，但是这些不需要实现。请注意：此处限定是一组方法，既然是方法，就不能是变量；而且是一组，表明可以有多个方法。再多声明一点，interface本质上是一种类型，确切的说，是指针类型，此处暂且不多表，后文中自然能体会到。
  interface是为实现多态功能，多态是指代码可以根据类型的具体实现采取不同行为的能力。如果一个类型实现了某个接口，所有使用这个接口的地方，都可以支持这种类型的值。

type  接口名称 interface {
    method1(参数列表) 返回值列表
    method2(参数列表) 返回值列表
    ...
    methodn(参数列表) 返回值列表
}
<!-- more -->
如果接口没有任何方法声明，那么就是一个空接口（interface{}），它的用途类似面向对象里的根类型Object，可被赋值为任何类型的对象。接口变量默认值是nil。如果实现接口的类型支持，可做相等运算。

简单的说：

interface是方法的集合
interface是一种类型，并且是指针类型
interface的更重要的作用在于多态实现

接口的使用不仅仅针对结构体，自定义类型、变量等等都可以实现接口。
如果一个接口没有任何方法，我们称为空接口，由于空接口没有方法，所以任何类型都实现了空接口。
要实现一个接口，必须实现该接口里面的所有方法。

Under the covers, interfaces are implemented as two elements, a type and a value. The value, called the interface's dynamic value, is an arbitrary concrete value and the type is that of the value. For the int value 3, an interface value contains, schematically, (int, 3).

interface在Go底层，被表示为一个值和值对应的类型的集合体
go 允许不带任何方法的 interface ，这种类型的 interface 叫empty interface。

如果一个类型实现了一个 interface 中所有方法，我们说类型实现了该 interface，所以所有类型都实现了 empty interface，因为任何一种类型至少实现了 0 个方法。go 没有显式的关键字用来实现 interface，只需要实现 interface 包含的方法即可。

如果定义一个函数参数是 interface{} 类型，这个函数应该可以接受任何类型作为它的参数。

既然空的 interface 可以接受任何类型的参数，那么一个 interface{}类型的 slice 是不是就可以接受任何类型的 slice ?

func printAll(vals []interface{}) { 
names := []string{"stanley", "david", "oscar"}
printAll(names)

执行之后竟然会报cannot use names (type []string) as type []interface {} in argument to printAll 错误，why？

这个错误说明 go 没有帮助我们自动把 slice 转换成 interface{} 类型的 slice，所以出错了。go 不会对 类型是interface{} 的 slice 进行转换。为什么 go 不帮我们自动转换，一开始我也很好奇，最后终于在 go 的 wiki 中找到了答案https://github.com/golang/go/wiki/InterfaceSlice大意是 interface{} 会占用两个字长的存储空间，一个是自身的 methods 数据，一个是指向其存储值的指针，也就是 interface 变量存储的值，因而 slice []interface{} 其长度是固定的N*2，但是 []T 的长度是N*sizeof(T)，两种 slice 实际存储值的大小是有区别的(文中只介绍两种 slice 的不同，至于为什么不能转换猜测可能是 runtime 转换代价比较大)。

但是我们可以手动进行转换来达到我们的目的。

var dataSlice []int = foo()

var interfaceSlice []interface{} = make([]interface{}, len(dataSlice))

for i, d := range dataSlice {

        interfaceSlice[i] = d
        
        
go 中函数都是按值传递即 passed by value。

如果是按 pointer 调用，go 会自动进行转换，因为有了指针总是能得到指针指向的值是什么，如果是 value 调用，go 将无从得知 value 的原始值是什么，因为 value 是份拷贝。go 会把指针进行隐式转换得到 value，但反过来则不行。

对于 receiver 是 value 的 method，任何在 method 内部对 value 做出的改变都不影响调用者看到的 value，这就是按值传递。

https://research.swtch.com/interfaces

在Go语言中interface是一个非常重要的概念，也是与其它语言相比存在很大特色的地方。interface也是一个Go语言中的一种类型，是一种比较特殊的类型，存在两种interface，一种是带有方法的interface，一种是不带方法的interface。Go语言中的所有变量都可以赋值给空interface变量，实现了interface中定义方法的变量可以赋值给带方法的interface变量，并且可以通过interface直接调用对应的方法，实现了其它面向对象语言的多态的概念。

两种不同的interface在Go语言内部被定义成如下的两种结构体

// 没有方法的interface
type eface struct {
    _type *_type
    data  unsafe.Pointer
}

// 记录着Go语言中某个数据类型的基本特征
type _type struct {
    size       uintptr
    ptrdata    uintptr
    hash       uint32
    tflag      tflag
    align      uint8
    fieldalign uint8
    kind       uint8
    alg        *typeAlg
    gcdata    *byte
    str       nameOff
    ptrToThis typeOff
}

// 有方法的interface
type iface struct {
    tab  *itab
    data unsafe.Pointer
}

type itab struct {
    inter  *interfacetype
    _type  *_type
    link   *itab
    hash   uint32
    bad    bool
    inhash bool
    unused [2]byte
    fun    [1]uintptr
}

// interface数据类型对应的type
type interfacetype struct {
    typ     _type
    pkgpath name
    mhdr    []imethod
}

可以看到两种类型的interface在内部实现时都是定义成了一个2个字段的结构体，所以任何一个interface变量都是占用16个byte的内存空间。
 在Go语言中_type这个结构体非常重要，记录着某种数据类型的一些基本特征，比如这个数据类型占用的内存大小（size字段），数据类型的名称（nameOff字段）等等。每种数据类型都存在一个与之对应的_type结构体（Go语言原生的各种数据类型，用户自定义的结构体，用户自定义的interface等等）。如果是一些比较特殊的数据类型，可能还会对_type结构体进行扩展，记录更多的信息，比如interface类型，就会存在一个interfacetype结构体，除了通用的_type外，还包含了另外两个字段pkgpath和mhdr，后文在对这两个字段的作用进行解析。除此之外还有其它类型的数据结构对应的结构体，比如structtype，chantype，slicetype，有兴趣的可以在$GOROOT/src/runtime/type.go文件中查看。
 img src="{{site.url}}{{site.baseurl}}/img/eface.webp"/>
 
 赋值
 存在对没有方法的interface变量和有方法的interface变量赋值这两种不同的情况。分别详解这两种不同的赋值过程。

没有方法的interface变量赋值
 对没有方法的interface变量赋值时编译器做了什么工作？
 
  赋值
 存在对没有方法的interface变量和有方法的interface变量赋值这两种不同的情况。分别详解这两种不同的赋值过程。

没有方法的interface变量赋值
 对没有方法的interface变量赋值时编译器做了什么工作？创建一个eface.go文件，代码如下：
     1  package main
     2
     3  type Struct1 struct {
     4      A int64
     5      B int64
     6  }
     7
     8  func main() {
     9      s := new(Struct1)
    10      var i interface{}
    11      i = a
    12
    13      _ = i
    14  }
 输入命令go build -gcflags '-l -N' eface.go，go tool objdump -s "main.main" eface，查看汇编代码。

     1  TEXT main.main(SB) /Users/didi/Source/Go/src/ppt/eface.go
     2    eface.go:8        0x104f360       4883ec38        SUBQ $0x38, SP
     3    eface.go:8        0x104f364       48896c2430      MOVQ BP, 0x30(SP)
     4    eface.go:8        0x104f369       488d6c2430      LEAQ 0x30(SP), BP
     5    eface.go:9        0x104f36e       48c7042400000000    MOVQ $0x0, 0(SP)
     6    eface.go:9        0x104f376       48c744240800000000  MOVQ $0x0, 0x8(SP)
     7    eface.go:9        0x104f37f       488d0424        LEAQ 0(SP), AX
     8    eface.go:9        0x104f383       4889442410      MOVQ AX, 0x10(SP)
     9    eface.go:10       0x104f388       48c744242000000000  MOVQ $0x0, 0x20(SP)
    10    eface.go:10       0x104f391       48c744242800000000  MOVQ $0x0, 0x28(SP)
    11    eface.go:11       0x104f39a       488b442410      MOVQ 0x10(SP), AX
    12    eface.go:11       0x104f39f       4889442418      MOVQ AX, 0x18(SP)
    13    eface.go:11       0x104f3a4       488d0dd5670000      LEAQ 0x67d5(IP), CX
    14    eface.go:11       0x104f3ab       48894c2420      MOVQ CX, 0x20(SP)
    15    eface.go:11       0x104f3b0       4889442428      MOVQ AX, 0x28(SP)
    16    eface.go:14       0x104f3b5       488b6c2430      MOVQ 0x30(SP), BP
    17    eface.go:14       0x104f3ba       4883c438        ADDQ $0x38, SP
 汇编代码第5~6行给结构体Struct1分配了空间SP+0x0和SP+0x8，第7~8行把这个结构体的地址放在存入了SP+0x10地址，这个地址就是变量s，第9~10行给interface类型的变量i分配了SP+0x20和SP+0x28，第13~14行把结构体A对应的_type的地址赋值到SP+0x20，然后把a变量赋值到了SP+0x28。这就是对没有方法的interface进行赋值的过程。赋值完以后的内存分配如下图：
 <img src="{{site.url}}{{site.baseurl}}/img/eface_no.webp"/>
 有方法的interface变量赋值
 如下一段代码在内存的分布
     1  package main
     2
     3  type I interface {
     4      Add()
     5      Del()
     6  }
     7
     8  type Struct1 struct {
     9      A int64
    10      B int64
    11  }
    12
    13  func (a *Struct1) Add() {
    14      a.A = a.A + 1
    15      a.B = a.B + 1
    16  }
    17
    18  func (a *Struct1) Del() {
    19      a.A = a.A - 1
    20      a.B = a.B - 1
    21  }
    22
    23  func main() {
    24      a := new(Struct1)
    25      var i I
    26      i = a
    27
    28      i.Add()
    29      i.Del()
    30  }
 <img src="{{site.url}}{{site.baseurl}}/img/eface_has.webp"/>
 
 这些内存地址都可以使用gdb调试时得到

(gdb) p i
$11 = {tab = 0x10a70e0 <Struct1,main.I>, data = 0xc42001a0c0}
(gdb) p a
$12 = (struct main.Struct1 *) 0xc42001a0c0
(gdb) p i.tab
$13 = (runtime.itab *) 0x10a70e0 <Struct1,main.I>
(gdb) p i.tab.inter
$14 = (runtime.interfacetype *) 0x105dc60 <type.*+59232>
(gdb) p i.tab._type
$15 = (runtime._type *) 0x105d200 <type.*+56576>
 通过对内存地址的打印，可以很清晰的看出在对有方法的interface变量进行赋值时的内存分布。Struct1类型和interface I类型都存在内存记录着各自的_type结构体信息，在将Struct1类型的变量赋值给interface I类型时，会有一个itab类型的结构体将Struct1类型和interface I类型关联起来。
 上面的例子都是将一个指针赋值给interface变量，如果是将一个值赋值给interface变量。会先对分配一块空间保存该值的副本，然后将该interface变量的data字段指向这个新分配的空间。将一个值赋值给interface变量时，操作的都是该值的一个副本。

2.3 方法的调用
 上面对有方法的interface进行赋值后，是如何实现通过接口变量实现了函数调用呢？参考下面的汇编代码

     1  TEXT main.main(SB) /Users/didi/Source/Go/src/ppt/iface.go
     2    iface.go:23       0x104f3e0       65488b0c25a0080000  MOVQ GS:0x8a0, CX
     3    iface.go:23       0x104f3e9       483b6110        CMPQ 0x10(CX), SP
     4    iface.go:23       0x104f3ed       0f8687000000        JBE 0x104f47a
     5    iface.go:23       0x104f3f3       4883ec38        SUBQ $0x38, SP
     6    iface.go:23       0x104f3f7       48896c2430      MOVQ BP, 0x30(SP)
     7    iface.go:23       0x104f3fc       488d6c2430      LEAQ 0x30(SP), BP
     8    iface.go:23       0x104f401       488d0578ff0000      LEAQ 0xff78(IP), AX
     9    iface.go:24       0x104f408       48890424        MOVQ AX, 0(SP)
    10    iface.go:24       0x104f40c       e86fcefbff      CALL runtime.newobject(SB)
    11    iface.go:24       0x104f411       488b442408      MOVQ 0x8(SP), AX
    12    iface.go:24       0x104f416       4889442410      MOVQ AX, 0x10(SP)
    13    iface.go:25       0x104f41b       48c744242000000000  MOVQ $0x0, 0x20(SP)
    14    iface.go:25       0x104f424       48c744242800000000  MOVQ $0x0, 0x28(SP)
    15    iface.go:26       0x104f42d       488b442410      MOVQ 0x10(SP), AX
    16    iface.go:26       0x104f432       4889442418      MOVQ AX, 0x18(SP)
    17    iface.go:26       0x104f437       488d0da27c0500      LEAQ 0x57ca2(IP), CX
    18    iface.go:26       0x104f43e       48894c2420      MOVQ CX, 0x20(SP)
    19    iface.go:26       0x104f443       4889442428      MOVQ AX, 0x28(SP)
    20    iface.go:28       0x104f448       488b442420      MOVQ 0x20(SP), AX
    21    iface.go:28       0x104f44d       488b4020        MOVQ 0x20(AX), AX
    22    iface.go:28       0x104f451       488b4c2428      MOVQ 0x28(SP), CX
    23    iface.go:28       0x104f456       48890c24        MOVQ CX, 0(SP)
    24    iface.go:28       0x104f45a       ffd0            CALL AX
    25    iface.go:29       0x104f45c       488b442420      MOVQ 0x20(SP), AX
    26    iface.go:29       0x104f461       488b4028        MOVQ 0x28(AX), AX
    27    iface.go:29       0x104f465       488b4c2428      MOVQ 0x28(SP), CX
    28    iface.go:29       0x104f46a       48890c24        MOVQ CX, 0(SP)
    29    iface.go:29       0x104f46e       ffd0            CALL AX
    30    iface.go:30       0x104f470       488b6c2430      MOVQ 0x30(SP), BP
    31    iface.go:30       0x104f475       4883c438        ADDQ $0x38, SP
    32    iface.go:30       0x104f479       c3          RET
    33    iface.go:23       0x104f47a       e8f182ffff      CALL runtime.morestack_noctxt(SB)
    34    iface.go:23       0x104f47f       e95cffffff      JMP main.main(SB)
 汇编代码的第17行和18行，将itab的地址加载到SP+0x20地址处，第20，21行，24行将SP+0x20的值加载到AX寄存器，然后将AX+0x20地址的值加载到AX寄存器，CALL AX就实现了add方法的调用，其中第22行和23行的作用是将interface里面data字段的地址传递给了add方法。
 <img src="{{site.url}}{{site.baseurl}}/img/eface_call.webp"/>
 
 通过对itab结构体进行分析，可以看到偏移0x20处为fun字段，其中0x20处为add函数的入口地址，0x28处就是del函数的入口地址。

2.4 断言的实现
 在Go语言中，经常需要对一个interface变量进行断言

     1  package main
     2
     3  type Struct1 struct {
     4      A int64
     5  }
     6
     7  func main() {
     8      a := new(Struct1)
     9
    10      var i interface{}
    11      i = a
    12
    13      b, ok := i.(Struct1)
    14      if ok {
    15          _ = b
    16      }
    17  }
 生成汇编代码进行分析

     1  TEXT main.main(SB) /Users/didi/Source/Go/src/ppt/assert.go
     2    assert.go:7       0x104f360       4883ec48        SUBQ $0x48, SP
     3    assert.go:7       0x104f364       48896c2440      MOVQ BP, 0x40(SP)
     4    assert.go:7       0x104f369       488d6c2440      LEAQ 0x40(SP), BP
     5    assert.go:8       0x104f36e       48c744241000000000  MOVQ $0x0, 0x10(SP)
     6    assert.go:8       0x104f377       488d442410      LEAQ 0x10(SP), AX
     7    assert.go:8       0x104f37c       4889442420      MOVQ AX, 0x20(SP)
     8    assert.go:10      0x104f381       48c744243000000000  MOVQ $0x0, 0x30(SP)
     9    assert.go:10      0x104f38a       48c744243800000000  MOVQ $0x0, 0x38(SP)
    10    assert.go:11      0x104f393       488b442420      MOVQ 0x20(SP), AX
    11    assert.go:11      0x104f398       4889442428      MOVQ AX, 0x28(SP)
    12    assert.go:11      0x104f39d       488d0d1c680000      LEAQ 0x681c(IP), CX
    13    assert.go:11      0x104f3a4       48894c2430      MOVQ CX, 0x30(SP)
    14    assert.go:11      0x104f3a9       4889442438      MOVQ AX, 0x38(SP)
    15    assert.go:13      0x104f3ae       488b442438      MOVQ 0x38(SP), AX
    16    assert.go:13      0x104f3b3       488b4c2430      MOVQ 0x30(SP), CX
    17    assert.go:13      0x104f3b8       488d1581ed0000      LEAQ 0xed81(IP), DX
    18    assert.go:13      0x104f3bf       4839d1          CMPQ DX, CX
    19    assert.go:13      0x104f3c2       7402            JE 0x104f3c6
    20    assert.go:13      0x104f3c4       eb3f            JMP 0x104f405
    21    assert.go:13      0x104f3c6       488b00          MOVQ 0(AX), AX
    22    assert.go:13      0x104f3c9       b901000000      MOVL $0x1, CX
    23    assert.go:13      0x104f3ce       eb00            JMP 0x104f3d0
 汇编的第12行，17行，18行可以看出，将Struct1对应的_type结构体的地址赋值给interface以后。在进行断言的时候，原理就是将interface变量_type字段的与Struct1对应的_type结构地址进行对比。
 在本例子中，第12行的IP寄存器对应的值是0x104f39d，0x681c(IP)对应的地址为0x1055BB9，第17行的IP寄存器对应的值是0x104f3b8，0xed81(IP)对应的地址为0x105E139，貌似并不相同。可能是对Go的汇编中对IP寄存器的理解存在偏差
 
 3 Go的反射
 反射是一种强大的语言特性，可以“动态”的调用方法，获取结构体运行时的一些特征，很多框架的实现都离不开反射。Go的反射就是通过interface类型来实现的。

3.1 反射获取变量的信息
 Go的反射包主要存在两个重要的结构体。

     1  type Value struct {
     2      typ *rtype
     3      ptr unsafe.Pointer
     4      flag
     5  }
     6
     7  func ValueOf(i interface{}) Value {
     8  }
     9
    10  type Type interface {
    11      Align() int
    12      FieldAlign() int
    13      Method(int) Method
    14      Name() string
    15      //一堆方法
    16      //....
    17  }
    18
    19  func TypeOf(i interface{}) Type {
    20      eface := *(*emptyInterface)(unsafe.Pointer(&i))
    21      return toType(eface.typ)
    22  }
    23
    24  type emptyInterface struct {
    25      typ  *rtype
    26      word unsafe.Pointer
    27  }
 任何一个变量可以通过调用ValueOf来获取到变量的Value结构体，通过TypeOf方法来获取变量的Type接口类型。通过TypeOf方法获取到的Type接口实际上就是该变量对应的_type。
 通过前面的分析，当通过TypeOf方法获取到变量的_type结构体后，很容易获取到该变量的一些基本信息，比如_type结构体中的各种字段都可以直接获取到。

3.2 反射修改变量的值
     1  package main
     2
     3  import (
     4      "reflect"
     5  )
     6
     7  func main() {
     8      var x int64 = 10
     9
    10      reflect.ValueOf(x).SetInt(20)
    11
    12      reflect.ValueOf(&x).SetInt(20)
    13
    14      reflect.ValueOf(&x).Elem().SetInt(20)
    15  }
 上面的例子中，第10行，12行都会报panic，只有第14行能修改变量的值。在使用ValueOf获取到Value结构体以后，flag字段记录着值能否进行修改，这样应该是为了避免误操作，保证api调用者明确了解到是否需要修改值。

3.3 反射修改结构体变量字段的值
 如果需要通过反射修改某结构体里面各个字段的值。

     1  package main
     2
     3  import (
     4      "reflect"
     5      "fmt"
     6  )
     7
     8  type Struct1 struct {
     9      A int64
    10      B int64
    11      C int64
    12  }
    13
    14  func main() {
    15      P := new(Struct1)
    16
    17      V := reflect.ValueOf(P).Elem()
    18      V.FieldByName("A").SetInt(100)
    19      V.FieldByName("B").SetInt(200)
    20      V.FieldByName("C").SetInt(300)
    21
    22      fmt.Printf("%v", P)
    23  }
 上面的代码中，需要根据结构体字段的名称对各个字段的值进行修改，内部是如何实现的呢？
 <img src="{{site.url}}{{site.baseurl}}/img/eface_reflect.webp"/>
 
 每一个自定义的struct类型都存在这一个对应的structType结构体，该结构体记录了每个字段structField。通过对比structField里面的name字段，就可以获取到某个字段的type和偏移量。从而对具体的值进行修改。

3.4 反射动态调用方法
 动态的调用方法是怎么实现的？

     1  package main
     2
     3  import (
     4      "reflect"
     5  )
     6
     7  type Struct1 struct {
     8      A int64
     9      B int64
    10      C int64
    11  }
    12
    13  func (p *Struct1) Set() {
    14      p.A = 200
    15  }
    16
    17  func main() {
    18      P := new(Struct1)
    19      P.A = 100
    20      P.B = 200
    21      P.C = 300
    22
    23      V := reflect.ValueOf(P)
    24
    25      params := make([]reflect.Value, 0)
    26      V.MethodByName("Set").Call(params)
    27  }
 结构体的方法在内存中存在如下的分布
 <img src="{{site.url}}{{site.baseurl}}/img/eface_rc.webp"/>
 在编译过程中，结构体对应方法的相关信息都已经存在于内存中，分配了一块uncommonType的结构体跟在fields字段后面。根据内存的分布，如果需要根据一个结构体的名称获取到方法并且执行，只需要根据uncommonType结构中的moff字段去获取方法相关信息的地址块，然后逐个对比名称是否为想要获取的方法进行调用。
 
 
    
 空interface
空interface(interface{})不包含任何的method，正因为如此，所有的类型都实现了空interface。空interface对于描述起不到任何的作用(因为它不包含任何的method），但是空interface在我们需要存储任意类型的数值的时候相当有用，因为它可以存储任意类型的数值。它有点类似于C语言的void*类型

// 定义a为空接口
var a interface{}
var i int = 5
s := "Hello world"
// a可以存储任意类型的数值
a = i
a = s
一个函数把interface{}作为参数，那么他可以接受任意类型的值作为参数，如果一个函数返回interface{},那么也就可以返回任意类型的值。

 
package main

import (
    "fmt"
    "reflect"
)

func main() {
    a := (*interface{})(nil)
    fmt.Println(reflect.TypeOf(a), reflect.ValueOf(a))
    var b interface{} = (*interface{})(nil)
    fmt.Println(reflect.TypeOf(b), reflect.ValueOf(b))
    fmt.Println(a == nil, b == nil)
}
输出如下：

*interface {} <nil>
*interface {} <nil>
true false

在封面下，接口被实现为两个元素，一个类型和一个值。该值称为接口的动态值，是一个任意的具体值，类型是值的类型。对于int值3，接口值示意性地包含（int，3）。 只有在内部值和类型都未设置时（nil，nil），接口值才为零。特别是一个nil接口将始终保持为零类型。如果我们在接口值中存储一个类型为* int的零指针，则无论指针的值是什么，内部类型都将是* int：（* int，nil）。因此，即使内部指针为零，这样的接口值也不会为零。

a := (*interface{})(nil)与var a *interface{} = nil相等

但是var b interface{} = (*interface{})(nil)，意味着b是类型interface{}，interface{}只有nil当它的类型和值都是变量时才是变量nil，显然类型*interface{}不是nil。