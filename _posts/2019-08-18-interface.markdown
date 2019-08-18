---
title: interface
layout: post
category: golang
author: 夏泽民
---
golang中的接口分为带方法的接口和空接口。 带方法的接口在底层用iface表示，空接口的底层则是eface表示。
<!-- more -->
//runtime/runtime2.go

//非空接口
type iface struct {
	tab  *itab
	data unsafe.Pointer
}
type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	hash   uint32 // copy of _type.hash. Used for type switches.
	bad    bool   // type does not implement interface
	inhash bool   // has this itab been added to hash?
	unused [2]byte
	fun    [1]uintptr // variable sized
}

//******************************

//空接口
type eface struct {
	_type *_type
	data  unsafe.Pointer
}

//========================
//这两个接口共同的字段_type
//========================

//runtime/type.go
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	alg        *typeAlg
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}
//_type这个结构体是golang定义数据类型要用的，讲到反射文章的时候在具体讲解这个_type。
复制代码1.iface
1.1 变量类型是如何转换成接口类型的？
看下方代码:
package main
type Person interface {
   run()
}

type xitehip struct {
   age uint8
}
func (o xitehip)run() {
}

func main()  {
   var xh Person = xitehip{age:18}
   xh.run()
}

复制代码xh变量是Person接口类型，那xitehip的struct类型是如何转换成接口类型的呢？
看一下生成的汇编代码：
0x001d 00029 (main.go:13)	PCDATA	$2, $0
0x001d 00029 (main.go:13)	PCDATA	$0, $0
0x001d 00029 (main.go:13)	MOVB	$0, ""..autotmp_1+39(SP)
0x0022 00034 (main.go:13)	MOVB	$18, ""..autotmp_1+39(SP)
0x0027 00039 (main.go:13)	PCDATA	$2, $1
0x0027 00039 (main.go:13)	LEAQ	go.itab."".xitehip,"".Person(SB), AX
0x002e 00046 (main.go:13)	PCDATA	$2, $0
0x002e 00046 (main.go:13)	MOVQ	AX, (SP)
0x0032 00050 (main.go:13)	PCDATA	$2, $1
0x0032 00050 (main.go:13)	LEAQ	""..autotmp_1+39(SP), AX
0x0037 00055 (main.go:13)	PCDATA	$2, $0
0x0037 00055 (main.go:13)	MOVQ	AX, 8(SP)
0x003c 00060 (main.go:13)	CALL	runtime.convT2Inoptr(SB)
0x0041 00065 (main.go:13)	MOVQ	16(SP), AX
0x0046 00070 (main.go:13)	PCDATA	$2, $2
0x0046 00070 (main.go:13)	MOVQ	24(SP), CX
复制代码从汇编发现有个转换函数：
runtime.convT2Inoptr(SB)
我们去看一下这个函数的实现：
func convT2Inoptr(tab *itab, elem unsafe.Pointer) (i iface) {
        t := tab._type
        if raceenabled {
                raceReadObjectPC(t, elem, getcallerpc(), funcPC(convT2Inoptr))
        }
        if msanenabled {
                msanread(elem, t.size)
        }
        x := mallocgc(t.size, t, false)//为elem申请内存
        memmove(x, elem, t.size)//将elem所指向的数据赋值到新的内存中
        i.tab = tab //设置iface的tab
        i.data = x //设置iface的data
        return
}
复制代码从以上实现我们发现编译器生成的struct原始数据会复制一份，然后将新的数据地址赋值给iface.data从而生成了完整的iface，这样如下原始代码中的xh就转换成了Person接口类型。
var xh Person = xitehip{age:18}
用gdb实际运行看一下

convT2Inoptr函数传进来的参数是*itab和源码中的 *xitehip。itab的类型原型和内存中的数据发现itab确实是runtime中源码里的字段。总共占了32个字节。（[4]uint8 不占字节）

是elem的数据他是个名为xitehip的结构体类型里面存放的是age=18。 内存中的0x12正好是age=18。注意此时的地址是:0xc000032777。

xh变量的数据类型和其中data字段的数据。发现xh确实是iface类型了且xh.data的地址不是上面提到的0xc000032777 而是0xc000014098，证明是复制了一份xitehip类型的struct。

1.2 指针变量类型是如何转换成接口类型的呢？
还是上面的例子只是将
var xh Person = xitehip{age:18}
复制代码换成了
var xh Person = &xitehip{age:18}
复制代码那指针类型的变量是如何转换成接口类型的呢？
见下方汇编代码：
0x001d 00029 (main.go:13)	PCDATA	$2, $1
0x001d 00029 (main.go:13)	PCDATA	$0, $0
0x001d 00029 (main.go:13)	LEAQ	type."".xitehip(SB), AX
0x0024 00036 (main.go:13)	PCDATA	$2, $0
0x0024 00036 (main.go:13)	MOVQ	AX, (SP)
0x0028 00040 (main.go:13)	CALL	runtime.newobject(SB)
0x002d 00045 (main.go:13)	PCDATA	$2, $1
0x002d 00045 (main.go:13)	MOVQ	8(SP), AX
0x0032 00050 (main.go:13)	MOVB	$18, (AX)
复制代码发现了这个函数：
runtime.newobject(SB)
复制代码去看一下具体实现：
// implementation of new builtin
// compiler (both frontend and SSA backend) knows the signature
// of this function
func newobject(typ *_type) unsafe.Pointer {
        return mallocgc(typ.size, typ, true)
}
复制代码编译器自动生成了iface并将&xitehip{age:18}创建的对象的地址（通过newobject）赋值给iface.data。就是xitehip这个结构体没有被复制。
用gdb看一下

1.3 那xh是如何找到run方法的呢？
1.4 接口调用规则
把上面的例子添加一个eat()接口方法并实现它（注意这个接口方法的实现的接受者是指针）。
package main
type Person interface {
	run()
	eat(string)
}
type xitehip struct {
	age uint8
}
func (o xitehip)run() { // //接收方o是值
}
func (o *xitehip)eat(food string) { //接收方o是指针
}
func main()  {
	var xh Person = &xitehip{age:18} //xh是指针
	xh.eat("ma la xiao long xia!")
	xh.run()
}
复制代码这个例子的xh变量的实际类型是个指针，那它是如何调用非指针方法run的呢？
继续gdb跟踪一下

直接跟踪xh.tab.fun的内存数据发现eat方法确实在0x44f940。上面已经说了fun这个数组大小只为1那run方法应该在eat的后面，但是gdb没有提示哪个地方是run的起始位置。为了验证run就在eat的后面，我直接往下debug看eat的入口地址在哪里

总结，指针类型的对象调用非指针类型的接收方的方法，编译器自动将接收方转换为指针类型；调用方通过xh.tab.fun这个数组找到对应的方法指令列表。
那xh是值类型的接口，而接口实现的方法的接收方是指针类型，那调用方可以调用这个指针方法吗，答案是不仅不能连编译都编译不过去，

见下表总结：



调用方
接收方
能否编译




值
值
true


值
指针
false


指针
值
true


指针
指针
true


指针
指针和值
true


值
指针和值
false



从上表可以得出如下结论：

调用方是值时，只要接收方有指针方法那编译器不允许通过编译。

2 eface
空接口相对于非空接口没有了方法列表。
type eface struct {
	_type *_type
	data  unsafe.Pointer
}
复制代码第一个属性由itab换成了_type,这个结构体是golang中的变量类型的基础，所以空接口可以指定任意变量类型。
2.1 示例：
cpackage main

import "fmt"

type xitehip struct {
}
func main()  {
	var a interface{} = xitehip{}
	var b interface{} = &xitehip{}
	fmt.Println(a)
	fmt.Println(b)
}
复制代码gdb跟一下

2.2断言
判断变量数据类型
   s, ok := i.(TypeName)
    if ok {
        fmt.Println(s)
    }
复制代码如果没有ok的话类型不正确的话会引起panic。
也可以用switch形式：
    switch v := v.(type) {
      case TypeName:
    ...
    }
复制代码3 检查接口
3.1 利用编译器检查接口实现
var _ InterfaceName = (*TypeName)(nil)
3.2 nil和nil interface
3.2.1 nil
func main() {
    var i interface{}
    if i == nil {
        println(“The interface is nil.“)
    }
}
(gdb) info locals;
i = {_type = 0x0, data = 0x0}
复制代码3.2.2 如果接口内部data值为nil，但tab不为空时，此时接口为nil interface。
// go:noinline
func main() {
    var o *int = nil
    var i interface{} = o

    if i == nil {
        println("Nil")
    }
    println(i)
}

(gdb) info locals;
i = {_type = 0x21432f8 <type.*+36723>, data = 0x0}
o = 0x0
复制代码3.2.3 利用反射检查
  v := reflect.ValueOf(a)
    if v.Isvalid() {
        println(v.IsNil()) // true, This is nil interface
}

 类型方法 reflect.TypeOf(interface{})
示例1代码如下图：

输出I
变量x的类型是I，那将x传入TypeOf()函数之后 Name()函数是如何获取到变量x的类型信息的呢？
接下来我们一步一步分析，第12行代码的Name()函数是如何获取到类型I的。
看一下TypeOf(interface)函数的实现：
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}
复制代码我们发现TypeOf的参数是接口类型，就是说变量x的副本被包装成了runtime/runtime2.go中定义的eface（空接口）。然后将eface强制转换成了emptyInterface,如下是reflect和runtime包下定义两个空接口：
//reflect/type.go
type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}

//runtime/runtime2.go
type eface struct {
	_type *_type
	data  unsafe.Pointer
}
复制代码发现和runtime包中的空接口很像，emptyInterface.word,runtime.eface字段类型是相同的。那就看看rtype和_type是否相同呢？
//reflect/type.go
type rtype struct {
	size       uintptr
	ptrdata    uintptr  // number of bytes in the type that can contain pointers
	hash       uint32   // hash of type; avoids computation in hash tables
	tflag      tflag    // extra type information flags
	align      uint8    // alignment of variable with this type
	fieldAlign uint8    // alignment of struct field with this type
	kind       uint8    // enumeration for C
	alg        *typeAlg // algorithm table
	gcdata     *byte    // garbage collection data
	str        nameOff  // string form
	ptrToThis  typeOff  // type for pointer to this type, may be zero
}

//runtime/type.go
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	alg        *typeAlg
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}
复制代码完全一样所以就可以毫无顾虑转换了。
也就是说emptyInterface.rtype结构体里已经有x的类型信息了。接下来继续看Name()函数是如何获取到类型的字符串信息的：
Type(interface{})函数里有个toType()函数，去看一下：
//reflect/type.go
func toType(t *rtype) Type {
	if t == nil {
		return nil
	}
	return t
}
复制代码上面代码是将*rtype直接转换成了Type类型了，那Type类型是啥?
//reflect/type.go
type Type interface {
......
    Name() string
......
}
复制代码其实Type是个接口类型。
那*rtype肯定实现了此接口中的方法,其中就包括Name()方法。找到了Name()的实现函数如下。如果不先看Name()的实现,其实也能猜到：就是从*rtype类型中定位数据获取数据并返回给调用者的过程，因为*rtype里面有包含值变量类型等信息。
func (t *rtype) Name() string {
	if t.tflag&tflagNamed == 0 {
		return ""
	}
	s := t.String()
	i := len(s) - 1
	for i >= 0 {
		if s[i] == '.' {
			break
		}
		i--
	}
	return s[i+1:]
}
复制代码重点看一下t.String()
func (t *rtype) String() string {
	s := t.nameOff(t.str).name()
	if t.tflag&tflagExtraStar != 0 {
		return s[1:]
	}
	return s
}
复制代码再重点看一下nameOff()：
func (t *rtype) nameOff(off nameOff) name {
	return name{(*byte)(resolveNameOff(unsafe.Pointer(t), int32(off)))}
}
复制代码从名字可以猜测出Off是Offset的缩写（这个函数里面的具体逻辑就探究了）进行偏移从而得到对应内存地址的值。
String()函数中的name()函数如下：

func (n name) name() (s string) {
	if n.bytes == nil {
		return
	}
	b := (*[4]byte)(unsafe.Pointer(n.bytes))

	hdr := (*stringHeader)(unsafe.Pointer(&s))
	hdr.Data = unsafe.Pointer(&b[3])
	hdr.Len = int(b[1])<<8 | int(b[2])
	return s
}
复制代码name()函数的逻辑是根据nameOff()返回的*byte(就是类型信息的首地址)计算出字符串的Data和Len位置，然后通过返回值&s包装出stringHeader（字符串原型）并将Data，Len赋值给字符串原型，从而将返回值s赋值。

总结 ：
普通的变量  => 反射中Type类型 => 获取变量类型信息 。

1，变量副本包装成空接口runtime.eface。
2，将runtime.eface转换成reflat.emptyInterface(结构都一样)。
3，将*emptyInterface.rtype 转换成 reflect.Type接口类型（包装成runtime.iface结构体类型）。
4，接口类型变量根据runtime.iface.tab.fun找到reflat.Name()函数。
5，reflect.Name()根据*rtype结构体str(nameoff类型)找到偏移量。
6，根据偏移量和基地址(基地址没有在*rtype中，这块先略过)。找到类型内存块。
7，包装成stringHeader类型返回给调用者。
其实核心就是将runtime包中的eface结构体数据复制到reflect包中的emptyInterface中然后在从里面获取相应的值类型信息。
refact.Type接口里面的其他方法就不在在这里说了，核心思想就是围绕reflat.emptyInterface中的数据进行查找等操作。
2 值方法 reflect.ValueOf(interface{})
package main
import (
	"reflect"
	"fmt"
)
func main() {
	var a = 3
	v := reflect.ValueOf(a)
	i := v.Interface()
	z := i.(int)
	fmt.Println(z)
}
复制代码看一下reflect.ValueOf()实现：
func ValueOf(i interface{}) Value {
	....
	return unpackEface(i)
}
复制代码返回值是Value类型：
type Value struct {
	typ *rtype
	ptr unsafe.Pointer
	flag //先忽略
}
复制代码Value是个结构体类型，包含着值变量的类型和数据指针。

func unpackEface(i interface{}) Value {
	e := (*emptyInterface)(unsafe.Pointer(&i))

	t := e.typ
	if t == nil {
		return Value{}
	}
	f := flag(t.Kind())
	if ifaceIndir(t) {
		f |= flagIndir
	}
	return Value{t, e.word, f}
}
复制代码具体实现是在unpackEface(interface{})中：
    e := (*emptyInterface)(unsafe.Pointer(&i))
复制代码和上面一样从*runtime.eface转换成*reflect.emptyInterface了。
最后包装成Value：
    return Value{t, e.word, f}
复制代码继续看一下示例代码：
    i := v.Interface()
复制代码的实现：
func (v Value) Interface() (i interface{}) {
	return valueInterface(v, true)
}

func valueInterface(v Value, safe bool) interface{} {
	......
	return packEface(v)
}

func packEface(v Value) interface{} {
	t := v.typ
	var i interface{}
	e := (*emptyInterface)(unsafe.Pointer(&i))
	switch {
	case ifaceIndir(t):
		if v.flag&flagIndir == 0 {
			panic("bad indir")
		}
               //将值的数据信息指针赋值给ptr
		ptr := v.ptr
		if v.flag&flagAddr != 0 {
			c := unsafe_New(t)
			typedmemmove(t, c, ptr)
			ptr = c
		}
                //为空接口赋值
		e.word = ptr 
	case v.flag&flagIndir != 0:
		e.word = *(*unsafe.Pointer)(v.ptr)
	default:
		e.word = v.ptr
	}
        //为空接口赋值
	e.typ = t
	return i
}
复制代码最终调用了packEface()函数，从函数名字面意思理解是打包成空接口。
逻辑是：从value.typ信息包装出reflect.emptyInterface结构体信息，然后将reflect.eface写入i变量中，又因为i是interface{}类型，编译器又会将i转换成runtime.eface类型。
z := i.(int)
复制代码根据字面量int编译器会从runtime.eface._type中查找int的值是否匹配，如果不匹配panic，匹配i的值赋值给z。

总结：从值变量 => value反射变量 => 接口变量：

1，包装成value类型。
2，从value类型中获取rtype包装成reflect.emptyInterface类型。
3，reflect.eface编译器转换成runtime.eface类型。
4，根据程序z :=i(int) 从runtime.eface._type中查找是否匹配。
5，匹配将值赋值给变量z。
总结：Value反射类型转interface{}类型核心还是reflet.emptyInterface与runtime.eface的相互转换。
